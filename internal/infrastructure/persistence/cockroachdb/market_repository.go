package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_MARKET_QUERY = `SELECT id, name, created_at FROM markets`
const INSERT_MARKET_QUERY = `INSERT INTO markets(name) VALUES ($1) RETURNING id, name, created_at`
const INSERT_MARKET_NAMED_QUERY = `
	INSERT INTO markets(
		name, 
		batch_index
	) VALUES (
		:name, 
		:batch_index
	) 
	RETURNING 
		id, 
		name, 
		created_at, 
		batch_index`

type MarketRepository struct {
	db sqlx.ExtContext
}

func (r *MarketRepository) GetByID(ctx context.Context, marketID uint64) (*domain.Market, error) {
	row := marketRecord{}
	query := GET_MARKET_QUERY + " WHERE id=$1"
	if err := sqlx.GetContext(ctx, r.db, &row, query, marketID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return row.ToDomain(), nil
}

func (r *MarketRepository) GetByNames(ctx context.Context, marketNames []string) (map[string]*domain.Market, error) {
	result := make(map[string]*domain.Market)
	for _, name := range marketNames {
		result[name] = nil
	}

	rows := []marketRecord{}
	query := GET_MARKET_QUERY + " WHERE name=ANY($1) ORDER BY name"
	if err := sqlx.SelectContext(ctx, r.db, &rows, query, marketNames); err != nil {
		return nil, err
	}

	for _, row := range rows {
		val := row.ToDomain()
		result[val.Name] = val
	}

	return result, nil
}

func (r *MarketRepository) Save(ctx context.Context, market *domain.Market) error {
	if market == nil {
		return pkg.InvalidStateErr("nil pointer passed on market")
	}

	record := &marketRecord{Name: market.Name}
	q := INSERT_MARKET_QUERY
	if err := r.db.QueryRowxContext(ctx, q, record.Name).
		StructScan(record); err != nil {
		return err
	}

	record.MapDomain(market)
	return nil
}

func (r *MarketRepository) SaveAll(ctx context.Context, markets []*domain.Market) error {
	if len(markets) == 0 {
		return nil
	}

	args := make([]marketRecord, 0, len(markets))
	for i, market := range markets {
		if market == nil {
			return pkg.InvalidStateErr("nil pointer passed on markets slice")
		}

		args = append(args, marketRecord{Name: market.Name, BatchIndex: i})
	}
	q := INSERT_STOCK_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := marketRecord{}
		if err := rows.StructScan(&record); err != nil {
			return err
		}

		record.MapDomain(markets[record.BatchIndex])
	}

	return rows.Err()
}

func NewMarketRepository(db sqlx.ExtContext) *MarketRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &MarketRepository{db}
}
