package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_STOCK_QUERY = `
	SELECT 
		id, 
		market_id, 
		company_id, 
		ticker, 
		name, 
		isin, 
		created_at, 
		updated_at 
	FROM stocks`
const INSERT_STOCK_QUERY = `
	INSERT INTO stocks(
		market_id, 
		company_id, 
		ticker, 
		name, 
		isin
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		market_id,
		company_id,
		ticker,
		name,
		isin,
		created_at,
		updated_at`
const INSERT_STOCK_NAMED_QUERY = `
	INSERT INTO stocks(
		market_id, 
		company_id, 
		ticker, 
		name, 
		isin
	) VALUES (:market_id, :company_id, :ticker, :name, :isin)
	RETURNING
		id,
		market_id,
		company_id,
		ticker,
		name,
		isin,
		created_at,
		updated_at`

type StockRepository struct {
	db sqlx.ExtContext
}

func (r *StockRepository) GetByID(ctx context.Context, stockID uint64) (*domain.Stock, error) {
	record := &stockRecord{}
	q := GET_STOCK_QUERY + " WHERE id=$1"

	if err := r.db.QueryRowxContext(ctx, q, stockID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *StockRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[domain.PopulatedStock], error) {

	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockRepository) GetAllPaginatedByUser(ctx context.Context, filter pkg.PaginationFilter, userID uint64) (*pkg.PaginationReponse[domain.PopulatedStock], error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockRepository) GetByStockCompanySearchParams(
	ctx context.Context,
	searchParams []domain.StockCompanySearchParam,
) (map[domain.StockCompanySearchParam]*domain.Stock, error) {
	results := make(map[domain.StockCompanySearchParam]*domain.Stock)
	for _, key := range searchParams {
		results[key] = nil
	}

	if len(searchParams) == 0 {
		return results, nil
	}

	args := make([]any, 0, len(results)*3)
	where := strings.Builder{}
	where.WriteString(" WHERE (market_id, company_id, ticker) IN  (")

	i := 0
	first := true
	for key := range results {
		if !first {
			where.WriteString(", ")
		}

		fmt.Fprintf(&where, "($%d, $%d, $%d)", i+1, i+2, i+3)
		args = append(args, key.MarketID, key.CompanyID, key.StockTicker)

		i += 3
		first = false
	}

	where.WriteString(")")
	q := GET_STOCK_QUERY + where.String()
	rows, err := r.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &stockRecord{}
		if err := rows.StructScan(record); err != nil {
			return nil, err
		}

		key := domain.StockCompanySearchParam{
			MarketID:    record.MarketID,
			CompanyID:   record.CompanyID,
			StockTicker: record.Ticker,
		}

		results[key] = record.ToDomain()
	}

	return results, rows.Err()
}

func (r *StockRepository) Save(ctx context.Context, stock *domain.Stock) error {
	if stock == nil {
		return nil // no-op
	}

	var isin, name sql.NullString
	if stock.Name != nil {
		name.String = *stock.Name
		name.Valid = true
	}

	if stock.Isin != nil {
		isin.String = *stock.Isin
		isin.Valid = true
	}
	record := &stockRecord{}
	if err := r.db.QueryRowxContext(
		ctx,
		INSERT_STOCK_QUERY,
		stock.MarketID,
		stock.CompanyID,
		stock.Ticker,
		name,
		isin,
	).StructScan(record); err != nil {
		return err
	}

	record.MapDomain(stock)
	return nil
}

func (r *StockRepository) SaveAll(ctx context.Context, stocks []*domain.Stock) error {
	if len(stocks) == 0 {
		return nil //no-op
	}

	args := make([]stockRecord, 0, len(stocks))
	stockMap := make(map[domain.StockCompanySearchParam]*domain.Stock)
	for _, stock := range stocks {
		if stock == nil {
			return pkg.InvalidStateErr("nil pointer passed on stock slice")
		}

		key := domain.StockCompanySearchParam{
			MarketID:    stock.MarketID,
			CompanyID:   stock.CompanyID,
			StockTicker: stock.Ticker,
		}
		_, isDuplicatedArg := stockMap[key]
		if isDuplicatedArg {
			return pkg.InvalidStateErr("invalid argument provided, duplicate unique constrain were found")
		}

		arg := stockRecord{
			MarketID:  stock.MarketID,
			CompanyID: stock.CompanyID,
			Ticker:    stock.Ticker,
			Isin:      sql.NullString{},
			Name:      sql.NullString{},
		}

		if stock.Name != nil {
			arg.Name.String = *stock.Name
			arg.Name.Valid = true
		}

		if stock.Isin != nil {
			arg.Isin.String = *stock.Isin
			arg.Isin.Valid = true
		}

		args = append(args, arg)
		stockMap[key] = stock
	}
	q := INSERT_STOCK_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &stockRecord{}
		if err := rows.StructScan(record); err != nil {
			return err
		}

		key := domain.StockCompanySearchParam{
			MarketID:    record.MarketID,
			CompanyID:   record.CompanyID,
			StockTicker: record.Ticker,
		}
		dom, ok := stockMap[key]
		if !ok {
			return pkg.InvalidStateErr("invalid record provided, non mapped entity")
		}

		record.MapDomain(dom)
	}

	return rows.Err()
}

func (r *StockRepository) Update(ctx context.Context, stock domain.StockUpdates) error {
	return pkg.InternalServerError("not implemented")
}

func NewStockRepository(db sqlx.ExtContext) *StockRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &StockRepository{db}
}
