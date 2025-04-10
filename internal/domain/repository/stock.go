package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/aggregate"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type StockRepository interface {
	Get(id int) (entity.Stock, shared.ApiErr)
	GetAllPaginated(filter shared.PaginationFilter) (shared.PaginationReponse[aggregate.PopulatedStock], shared.ApiErr)
	Create(args entity.Stock) (entity.Stock, shared.ApiErr)
	Update(args entity.Stock) (entity.Stock, shared.ApiErr)
}
