package cockroachdb

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type StockRegisterRepository struct {
	db sqlx.ExtContext
}

func (r StockRegisterRepository) GetByID(ctx context.Context, stockRegisterID uint64) (*domain.StockRegister, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) GetLastByStockID(ctx context.Context, stockID uint64) (*domain.StockRegister, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) GetRangeByStockID(ctx context.Context, stockID uint64, from, to time.Time) ([]domain.StockRegister, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) Save(ctx context.Context, stockRegister *domain.StockRegister) error {
	return pkg.InternalServerError("not implemented")
}

func (r StockRegisterRepository) SaveAll(ctx context.Context, stockRegisters []*domain.StockRegister) error {
	return pkg.InternalServerError("not implemented")
}

func NewStockRegisterRepository(db sqlx.ExtContext) *StockRegisterRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &StockRegisterRepository{db}
}
