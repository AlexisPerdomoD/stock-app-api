package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type MarketRepository interface {
	Get(id int) (entity.Market, shared.ApiErr)
	GetByName(name string) (entity.Market, shared.ApiErr)
	Create(args entity.Market) (entity.Market, shared.ApiErr)
}
