package usecase

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/aggregate"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/repository"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type getStocksUseCase struct {
	stockRepository repository.StockRepository
}

func (uc *getStocksUseCase) Execute(filters shared.PaginationFilter) (*shared.PaginationReponse[aggregate.PopulatedStock], error) {

	panic("getStocksUseCase.Execute() not implemented")
}

func NewGetStocksUseCase(ur repository.StockRepository) *getStocksUseCase {

	if ur == nil {
		panic("stock repository is nil, stopping :b")
	}

	return &getStocksUseCase{stockRepository: ur}
}
