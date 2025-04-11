package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
)

type MarketRepository interface {
	Get(id int) (*entity.Market, error)

	GetByName(name string) (*entity.Market, error)

	Create(args *entity.Market) error
}
