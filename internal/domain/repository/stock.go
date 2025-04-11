package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/aggregate"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type StockRepository interface {
	Get(id int) (*entity.Stock, error)

	GetByTicker(marketID int, ticker string) (*entity.Stock, error)

	GetAllPaginated(filter shared.PaginationFilter) (*shared.PaginationReponse[aggregate.PopulatedStock], error)

	Create(args *entity.Stock) error

	Update(stockID int, args *entity.Stock) error
}
