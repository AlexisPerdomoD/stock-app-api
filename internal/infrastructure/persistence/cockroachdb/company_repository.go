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

const GET_COMPANY_QUERY = `SELECT id, market_id, name, created_at FROM companies`
const INSERT_COMPANY_QUERY = `INSERT INTO companies(market_id, name) VALUES ($1, $2) RETURNING id, market_id, name, created_at`
const INSERT_COMPANY_NAMED_QUERY = `
	INSERT INTO companies(
		market_id, 
		name, 
		batch_index
	) VALUES (
		:market_id, 
		:name, 
		:batch_index
	) 
	RETURNING 
		id, 
		market_id, 
		name, 
		created_at, 
		batch_index`

type CompanyRepository struct {
	db sqlx.ExtContext
}

func (r *CompanyRepository) GetByID(ctx context.Context, companyID uint64) (*domain.Company, error) {
	record := &companyRecord{}
	q := GET_COMPANY_QUERY + " WHERE id=$1"
	if err := r.db.QueryRowxContext(ctx, q, companyID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *CompanyRepository) GetByMarketCompanySearch(
	ctx context.Context,
	searchParams []domain.MarketCompanySearchParam,
) (map[domain.MarketCompanySearchParam]*domain.Company, error) {
	results := make(map[domain.MarketCompanySearchParam]*domain.Company)
	for _, key := range searchParams {
		results[key] = nil
	}

	if len(searchParams) == 0 {
		return results, nil
	}
	args := make([]any, 0, len(results)*2)
	where := strings.Builder{}
	where.WriteString(" WHERE (market_id, name) IN ( ")

	i := 0
	first := true
	for key := range results {
		if !first {
			where.WriteString(", ")
		}

		fmt.Fprintf(&where, "($%d, $%d)", i+1, i+2)
		args = append(args, key.MarketID, key.Name)
		i += 2
		first = false
	}
	where.WriteString(" )")

	q := GET_COMPANY_QUERY + where.String()
	rows, err := r.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &companyRecord{}
		if err := rows.StructScan(record); err != nil {
			return nil, err
		}

		key := domain.MarketCompanySearchParam{
			MarketID: record.MarketID,
			Name:     record.Name,
		}
		results[key] = record.ToDomain()
	}

	return results, rows.Err()
}

func (r *CompanyRepository) Save(ctx context.Context, company *domain.Company) error {
	if company == nil {
		return pkg.InvalidStateErr("nil pointer passed on company")
	}

	record := &companyRecord{
		MarketID: company.MarketID,
		Name:     company.Name,
	}

	q := INSERT_COMPANY_QUERY
	if err := r.db.QueryRowxContext(ctx, q, record.MarketID, record.Name).
		StructScan(record); err != nil {
		return err
	}

	record.MapDomain(company)
	return nil
}

func (r *CompanyRepository) SaveAll(ctx context.Context, companies []*domain.Company) error {
	if len(companies) == 0 {
		return nil // no-op
	}

	args := make([]companyRecord, 0, len(companies))
	for i, company := range companies {
		if company == nil {
			return pkg.InvalidStateErr("nil pointer passed on companies slice")
		}

		args = append(args, companyRecord{
			MarketID:   company.MarketID,
			Name:       company.Name,
			BatchIndex: i,
		})
	}

	q := INSERT_COMPANY_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &companyRecord{}
		if err := rows.StructScan(record); err != nil {
			return err
		}

		record.MapDomain(companies[record.BatchIndex])
	}

	return rows.Err()
}

func NewCompanyRepository(db sqlx.ExtContext) *CompanyRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &CompanyRepository{db}
}
