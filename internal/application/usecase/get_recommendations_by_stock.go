package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type GetRecommendationsByStockUseCase struct {
	sr domain.StockRepository
	rr domain.RecommendationRepository
}

func (uc *GetRecommendationsByStockUseCase) Execute(
	ctx context.Context,
	filters pkg.PaginationFilter,
	stockID uint,
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

func NewGetRecommendationsByStockUseCase(
	sr domain.StockRepository,
	rr domain.RecommendationRepository,
) *GetRecommendationsByStockUseCase {

	if sr == nil {
		log.Fatalln("[GetRecommendationsByStockUseCase]: StockRepository provided was nil")
	}

	if rr == nil {
		log.Fatalln("[GetRecommendationsByStockUseCase]: RecommendationRepository was provided as nil")
	}

	return &GetRecommendationsByStockUseCase{sr, rr}
}
