package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/collection"
)

type GetStocks struct {
	sr domain.StockRepository
}

func (uc *GetStocks) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	userID *uint64,
) (*pkg.PaginationReponse[models.PopulatedStockView], error) {

	data, err := uc.sr.GetAllPaginated(ctx, filters, userID)
	if err != nil {
		return nil, err
	}

	response := &pkg.PaginationReponse[models.PopulatedStockView]{
		Page:       data.Page,
		PageSize:   data.PageSize,
		TotalSize:  data.TotalSize,
		TotalPages: data.TotalPages,
		Items:      collection.Map(data.Items, models.NewPopulatedStockView),
	}

	return response, nil
}

func NewGetStocks(sr domain.StockRepository) *GetStocks {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetStocksUseCase")
	}

	return &GetStocks{sr}
}
