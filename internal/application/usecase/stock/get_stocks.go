package usecase

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type getStocksUseCase struct {
	stockRepository domain.StockRepository
}

func (uc *getStocksUseCase) Execute(ctx context.Context, filters pkg.PaginationFilter) (*pkg.PaginationReponse[domain.PopulatedStock], error) {

	return uc.stockRepository.GetAllPaginated(ctx, filters)
}

func NewGetStocksUseCase(ur domain.StockRepository) *getStocksUseCase {

	if ur == nil {
		panic("stock repository is nil, stopping :b")
	}

	return &getStocksUseCase{stockRepository: ur}
}
