package cockroachdb

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type MarketRepository struct {
	db sqlx.ExtContext
}

func (r *MarketRepository) GetByID(ctx context.Context, marketID uint64) (*domain.Market, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *MarketRepository) GetByNames(ctx context.Context, marketNames []string) (map[string]*domain.Market, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *MarketRepository) Save(ctx context.Context, market *domain.Market) error {
	return pkg.InternalServerError("not implemented")
}

func (r *MarketRepository) SaveAll(ctx context.Context, markets []*domain.Market) error {
	return pkg.InternalServerError("not implemented")
}

func NewMarketRepository(db sqlx.ExtContext) *MarketRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &MarketRepository{db}
}
