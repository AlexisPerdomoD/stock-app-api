package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_MARKET_QUERY = `SELECT id, name, created_at FROM markets`
const INSERT_MARKET_QUERY = `INSERT INTO markets(name) VALUES (:name) RETURNING id, name, created_at`

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
		return nil
	}

	record := marketRecord{Name: market.Name}
	rows, err := sqlx.NamedQueryContext(ctx, r.db, INSERT_MARKET_QUERY, record)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		if err := rows.StructScan(&record); err != nil {
			return err
		}
	}

	record.MapDomain(market)
	return rows.Err()
}

func (r *MarketRepository) SaveAll(ctx context.Context, markets []*domain.Market) error {
	if len(markets) == 0 {
		return nil
	}

	records := make([]marketRecord, 0, len(markets))
	for _, market := range markets {
		if market == nil {
			return pkg.InternalServerError("nil arguments were provided as reference for saving markets record")
		}

		records = append(records, marketRecord{Name: market.Name})
	}

	rows, err := sqlx.NamedQueryContext(ctx, r.db, INSERT_MARKET_QUERY, records)
	if err != nil {
		return err
	}

	defer func() { _ = rows.Close() }()

	i := 0
	for rows.Next() {
		record := marketRecord{}
		if err := rows.StructScan(&record); err != nil {
			return err
		}

		record.MapDomain(markets[i])
		i++
	}

	if i != len(markets) {
		return pkg.InternalServerError("inserted rows count mismatch")
	}

	return rows.Err()
}

func NewMarketRepository(db sqlx.ExtContext) *MarketRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &MarketRepository{db}
}
