package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_STOCK_REGISTER_QUERY = `SELECT id, stock_id, price, tendency, created_at FROM stock_registers`

const INSERT_STOCK_REGISTER_QUERY = `
	INSERT INTO stock_registers (
		stock_id, 
		price, 
		tendency, 
		created_at
	) VALUES ($1, $2, $3, $4)
	RETURNING
		id, 
		stock_id, 
		price, 
		tendency, 
		created_at`

const INSERT_STOCK_REGISTER_NAMED_QUERY = `
		INSERT INTO stock_registers (
			stock_id, 
			price, 
			tendency, 
			created_at,
			batch_index
		) VALUES (
			:stock_id, 
			:price, 
			:tendency, 
			:created_at,
			:batch_index
		)RETURNING
			id, 
			stock_id, 
			price, 
			tendency, 
			created_at,
			batch_index`

type StockRegisterRepository struct {
	db sqlx.ExtContext
}

func (r *StockRegisterRepository) GetByID(
	ctx context.Context,
	stockRegisterID uint64,
) (*domain.StockRegister, error) {

	record := &stockRegisterRecord{}
	q := GET_STOCK_REGISTER_QUERY + " WHERE id=$1"
	if err := r.db.QueryRowxContext(ctx, q, stockRegisterID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *StockRegisterRepository) GetLastByStockID(ctx context.Context, stockID uint64) (*domain.StockRegister, error) {
	record := &stockRegisterRecord{}
	q := GET_STOCK_REGISTER_QUERY + " WHERE stock_id=$1 ORDER BY created_at DESC LIMIT 1"
	if err := r.db.QueryRowxContext(ctx, q, stockID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *StockRegisterRepository) GetLastsByStockID(ctx context.Context, stockID uint64, limit uint16) ([]domain.StockRegister, error) {
	results := make([]domain.StockRegister, 0, limit)

	q := GET_STOCK_REGISTER_QUERY + " WHERE stock_id=$1 ORDER BY created_at DESC LIMIT $2"

	rows, err := r.db.QueryxContext(ctx, q, stockID, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &stockRegisterRecord{}
		if err := rows.StructScan(record); err != nil {
			return nil, err
		}

		results = append(results, *record.ToDomain())
	}

	return results, rows.Err()
}

func (r *StockRegisterRepository) GetRangeByStockID(ctx context.Context, stockID uint64, from, to time.Time) ([]domain.StockRegister, error) {
	results := make([]domain.StockRegister, 0)
	if from.After(to) {
		return results, nil
	}

	q := GET_STOCK_REGISTER_QUERY + " WHERE stock_id=$1 AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC"

	rows, err := r.db.QueryxContext(ctx, q, stockID, from, to)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &stockRegisterRecord{}
		if err := rows.StructScan(record); err != nil {
			return nil, err
		}

		results = append(results, *record.ToDomain())
	}

	return results, rows.Err()
}

func (r *StockRegisterRepository) Save(ctx context.Context, stockRegister *domain.StockRegister) error {
	if stockRegister == nil {
		return pkg.InvalidStateErr("nil pointer passed on stockRegister")
	}

	record := &stockRegisterRecord{}
	if err := r.db.QueryRowxContext(
		ctx,
		INSERT_STOCK_REGISTER_QUERY,
		stockRegister.StockID,
		stockRegister.Price,
		stockRegister.Tendency,
		stockRegister.CreatedAt,
	).StructScan(record); err != nil {
		return err
	}

	record.MapDomain(stockRegister)
	return nil
}

func (r *StockRegisterRepository) SaveAll(ctx context.Context, stockRegisters []*domain.StockRegister) error {
	if len(stockRegisters) == 0 {
		return nil //no-op
	}

	args := make([]stockRegisterRecord, 0, len(stockRegisters))
	for i, stockRegister := range stockRegisters {
		if stockRegister == nil {
			return pkg.InvalidStateErr("nil pointer passed on stockRegisters slice")
		}

		arg := stockRegisterRecord{
			StockID:    stockRegister.StockID,
			Price:      stockRegister.Price,
			Tendency:   stockRegister.Tendency,
			CreatedAt:  stockRegister.CreatedAt,
			BatchIndex: i,
		}

		args = append(args, arg)
	}
	q := INSERT_STOCK_REGISTER_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	insertedCount := 0
	for rows.Next() {
		record := &stockRegisterRecord{}
		if err := rows.StructScan(record); err != nil {
			return err
		}

		record.MapDomain(stockRegisters[record.BatchIndex])
		insertedCount++
	}

	if insertedCount != len(stockRegisters) {
		return pkg.InvalidStateErr(fmt.Sprintf("inserted count %d != len(stockRegisters) %d", insertedCount, len(stockRegisters)))
	}

	return rows.Err()
}

func NewStockRegisterRepository(db sqlx.ExtContext) *StockRegisterRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &StockRegisterRepository{db}
}
