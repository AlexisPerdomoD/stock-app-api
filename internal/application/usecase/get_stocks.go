package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type GetStocksUseCase struct {
	sr domain.StockRepository
}

func (uc *GetStocksUseCase) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	userID *uint,
) (*pkg.PaginationReponse[domain.PopulatedStock], error) {

	return uc.sr.GetAllPaginated(ctx, filters, userID)
}

func NewGetStocksUseCase(sr domain.StockRepository) *GetStocksUseCase {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetStocksUseCase")
	}

	return &GetStocksUseCase{sr}
}
