package cockroachdb

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type StockRepository struct {
	db sqlx.ExtContext
}

func (r *StockRepository) Get(ctx context.Context, stockID uint64, userID *uint64) (*domain.Stock, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockRepository) GetPopulated(ctx context.Context, stockID uint64, userID *uint64) (*domain.PopulatedStock, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter, userID *uint64) (*pkg.PaginationReponse[domain.PopulatedStock], error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockRepository) GetByStockCompanySearchParams(ctx context.Context, tickers []domain.StockCompanySearchParam) (map[domain.StockCompanySearchParam]*domain.Stock, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *StockRepository) Save(ctx context.Context, stock *domain.Stock) error {
	return pkg.InternalServerError("not implemented")
}

func (r *StockRepository) SaveAll(ctx context.Context, stocks []*domain.Stock) error {
	return pkg.InternalServerError("not implemented")
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
