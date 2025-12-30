package cockroachdb

import (
	"context"
	"database/sql"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_BROKERAGE_QUERY = `SELECT id, name, created_at FROM brokerages`
const INSERT_BROKERAGE_QUERY = `INSERT INTO brokerages(name) VALUES ($1) RETURNING id, name, created_at`
const INSERT_BROKERAGE_NAMED_QUERY = `INSERT INTO brokerages(name) VALUES (:name) RETURNING id, name, created_at`

type BrokerageRepository struct {
	db sqlx.ExtContext
}

func (r *BrokerageRepository) GetByID(ctx context.Context, brokerageID uint64) (*domain.Brokerage, error) {
	var record brokerageRecord
	q := GET_BROKERAGE_QUERY + " WHERE id=$1"
	if err := r.db.QueryRowxContext(ctx, q, brokerageID).
		StructScan(&record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *BrokerageRepository) GetByNames(ctx context.Context, names []string) (map[string]*domain.Brokerage, error) {
	results := make(map[string]*domain.Brokerage, len(names))
	for _, name := range names {
		results[name] = nil
	}

	if len(names) == 0 {
		return results, nil
	}

	q := GET_BROKERAGE_QUERY + " WHERE name=ANY($1) ORDER BY name"
	rows, err := r.db.QueryxContext(ctx, q, names)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &brokerageRecord{}
		if err := rows.StructScan(record); err != nil {
			return nil, err
		}

		result := record.ToDomain()
		results[result.Name] = result
	}

	return results, rows.Err()
}

func (r *BrokerageRepository) Save(ctx context.Context, brokerage *domain.Brokerage) error {
	if brokerage == nil {
		return nil
	}

	record := brokerageRecord{Name: brokerage.Name}
	q := INSERT_BROKERAGE_QUERY
	if err := r.db.QueryRowxContext(ctx, q, record.Name).
		StructScan(&record); err != nil {
		return err
	}

	record.MapDomain(brokerage)
	return nil
}

func (r *BrokerageRepository) SaveAll(ctx context.Context, brokerages []*domain.Brokerage) error {
	args := make([]brokerageRecord, 0, len(brokerages))
	brokerageMap := make(map[string]*domain.Brokerage)
	for _, brokerage := range brokerages {
		if brokerage == nil {
			return pkg.InvalidStateErr("nil pointer passed on brokerages slice")
		}

		if _, isDuplicated := brokerageMap[brokerage.Name]; isDuplicated {
			return pkg.InvalidStateErr("invalid argument provided, duplicate unique constrain were found")
		}

		brokerageMap[brokerage.Name] = brokerage
		args = append(args, brokerageRecord{Name: brokerage.Name})
	}

	if len(args) == 0 {
		return nil
	}

	q := INSERT_BROKERAGE_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &brokerageRecord{}
		if err := rows.StructScan(record); err != nil {
			return err
		}

		record.MapDomain(brokerageMap[record.Name])
	}

	return rows.Err()
}

func NewBrokerageRepository(db sqlx.ExtContext) *BrokerageRepository {
	if db == nil {
		log.Fatalln("required db sqlx.ExtContext passed as nil")
	}

	return &BrokerageRepository{db}
}
