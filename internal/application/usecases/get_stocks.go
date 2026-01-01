package usecases

import (
	"context"
	"fmt"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/collection"
)

type GetStocks struct {
	stockRepository domain.StockRepository
}

func (uc *GetStocks) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	userID *uint64,
) (*pkg.PaginationResponse[models.PopulatedStockView], error) {
	var data *pkg.PaginationResponse[domain.PopulatedStock]
	var err error

	if userID != nil {
		data, err = uc.stockRepository.GetAllPaginatedByUser(ctx, filters, *userID)
	} else {
		data, err = uc.stockRepository.GetAllPaginated(ctx, filters)
	}

	if err != nil {
		return nil, err
	}

	response := &pkg.PaginationResponse[models.PopulatedStockView]{
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
		panic(fmt.Sprintf("nil arguments provided for NewGetStocksUseCase stockRepository=%+V", sr))
	}

	return &GetStocks{sr}
}
