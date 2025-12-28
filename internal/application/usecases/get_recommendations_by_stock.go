package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type GetRecommendationsByStock struct {
	sr domain.StockRepository
	rr domain.RecommendationRepository
}

func (uc *GetRecommendationsByStock) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	stockID uint64,
) (*pkg.PaginationReponse[domain.PopulatedRecommendation], error) {

	stock, err := uc.sr.Get(ctx, stockID, nil)

	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("Stock not found")
	}

	return uc.rr.GetAllPaginated(ctx, filters, stock.ID)

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
