package cockroachdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type StockRegisterRepository struct{}

func (r StockRegisterRepository) GetByID(ctx context.Context, stockRegisterID uint64) (*domain.StockRegister, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) GetLastByStockID(ctx context.Context, stockID uint64) (*domain.StockRegister, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) GetRangeByStockID(ctx context.Context, stockID uint64, from, to time.Time) ([]domain.StockRegister, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) Save(ctx context.Context, stockRegister domain.StockRegister) error {
	return pkg.InternalServerError("not implemented")
}

func NewStockRegisterRepository(db *sql.DB) *StockRegisterRepository {
	return &StockRegisterRepository{}
}
