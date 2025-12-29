package cockroachdb

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type StockTendencyStatRepository struct {
	db sqlx.ExtContext
}

func (r *StockTendencyStatRepository) GetByStockID(ctx context.Context, stockID uint64) (*domain.StockTendencyStat, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockTendencyStatRepository) Increment(ctx context.Context, stockID uint64, delta domain.StockTendencyDelta) error {
	return pkg.InternalServerError("not implemented")
}

func (r *StockTendencyStatRepository) IncrementAll(ctx context.Context, stockDeltas map[uint64]domain.StockTendencyDelta) error {
	return pkg.InternalServerError("not implemented")
}

func NewStockTendencyStatRepository(db sqlx.ExtContext) *StockTendencyStatRepository {
	return &StockTendencyStatRepository{db}
}
