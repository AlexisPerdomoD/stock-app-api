package cockroachdb

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type MarketRepository struct {
	db sqlx.ExtContext
}

func (r MarketRepository) GetByID(id uint64) (*domain.Market, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r MarketRepository) Save(market *domain.Market) error {
	return pkg.InternalServerError("not implemented")
}

func NewMarketRepository(db sqlx.ExtContext) *MarketRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	return &MarketRepository{db}
}
