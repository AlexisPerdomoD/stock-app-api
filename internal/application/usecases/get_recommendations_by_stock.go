package usecases

import (
	"context"
	"log"

	appmodel "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/collection"
)

type GetRecommendationsByStock struct {
	sr domain.StockRepository
	rr domain.RecommendationRepository
}

func (uc *GetRecommendationsByStock) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	stockID uint64,
) (*pkg.PaginationReponse[appmodel.PopulatedRecommendationView], error) {

	stock, err := uc.sr.Get(ctx, stockID)

	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("Stock not found")
	}

	data, err := uc.rr.GetAllPaginated(ctx, filters, stock.ID)
	if err != nil {
		return nil, err
	}

	response := &pkg.PaginationReponse[appmodel.PopulatedRecommendationView]{
		Page:       data.Page,
		PageSize:   data.PageSize,
		TotalSize:  data.TotalSize,
		TotalPages: data.TotalPages,
		Items:      collection.Map(data.Items, appmodel.NewPopulatedRecommendationView),
	}
	return response, nil

}

func NewGetRecommendationsByStock(
	sr domain.StockRepository,
	rr domain.RecommendationRepository,
) *GetRecommendationsByStock {

	if sr == nil {
		log.Fatalln("[GetRecommendationsByStockUseCase]: StockRepository provided was nil")
	}

	if rr == nil {
		log.Fatalln("[GetRecommendationsByStockUseCase]: RecommendationRepository was provided as nil")
	}

	return &GetRecommendationsByStock{sr, rr}
}
