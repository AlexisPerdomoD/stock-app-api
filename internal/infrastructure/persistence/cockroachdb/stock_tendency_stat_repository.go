package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

const GET_STOCK_TENDENCY_STAT_QUERY = `
	SELECT 
		id, 
		stock_id, 
		up_count, 
		side_count, 
		down_count, 
		created_at,
		updated_at
	FROM stock_tendency_stats`

const UPSERT_STOCK_TENDENCY_STAT_NAMED_QUERY = `
	INSERT INTO stock_tendency_stats (
		stock_id,
		up_count,
		side_count,
		down_count
	) VALUES (
		:stock_id,
		:up_count,
		:side_count,
		:down_count
	)
	ON CONFLICT (stock_id)
	DO UPDATE SET
		up_count    = stock_tendency_stats.up_count   + EXCLUDED.up_count,
		side_count  = stock_tendency_stats.side_count + EXCLUDED.side_count,
		down_count  = stock_tendency_stats.down_count + EXCLUDED.down_count,
		updated_at  = now()`

type StockTendencyStatRepository struct {
	db sqlx.ExtContext
}

func (r *StockTendencyStatRepository) GetByStockID(
	ctx context.Context,
	stockID uint64,
) (*domain.StockTendencyStat, error) {
	record := &stockTendencyStatRecord{}
	q := GET_STOCK_TENDENCY_STAT_QUERY + " WHERE stock_id=$1"
	if err := r.db.QueryRowxContext(ctx, q, stockID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *StockTendencyStatRepository) Increment(ctx context.Context, stockID uint64, delta domain.StockTendencyDelta) error {
	args := stockTendencyStatRecord{
		StockID:   stockID,
		UpCount:   delta.Up,
		SideCount: delta.Side,
		DownCount: delta.Down,
	}
	q := UPSERT_STOCK_TENDENCY_STAT_NAMED_QUERY
	_, err := sqlx.NamedExecContext(ctx, r.db, q, args)
	return err
}

func (r *StockTendencyStatRepository) IncrementAll(ctx context.Context, stockDeltas map[uint64]domain.StockTendencyDelta) error {
	args := make([]stockTendencyStatRecord, 0, len(stockDeltas))
	for stockID, delta := range stockDeltas {
		args = append(args, stockTendencyStatRecord{
			StockID:   stockID,
			UpCount:   delta.Up,
			SideCount: delta.Side,
			DownCount: delta.Down,
		})
	}

	q := UPSERT_STOCK_TENDENCY_STAT_NAMED_QUERY
	_, err := sqlx.NamedExecContext(ctx, r.db, q, args)
	return err
}

func NewStockTendencyStatRepository(db sqlx.ExtContext) *StockTendencyStatRepository {
	if db == nil {
		panic("nil arguments provided to NewStockTendencyStatRepository")
	}
	return &StockTendencyStatRepository{db}
}
