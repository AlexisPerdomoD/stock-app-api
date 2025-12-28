package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type GetStocks struct {
	sr domain.StockRepository
}

func (uc *GetStocks) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	userID *uint64,
) (*pkg.PaginationReponse[domain.PopulatedStock], error) {

	return uc.sr.GetAllPaginated(ctx, filters, userID)
}

func NewGetStocks(sr domain.StockRepository) *GetStocks {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetStocksUseCase")
	}

	return &GetStocks{sr}
}
