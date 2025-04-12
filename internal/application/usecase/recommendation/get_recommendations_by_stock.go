package usecase

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type getRecommendationsByStockUseCase struct {
	stockRepository          domain.StockRepository
	recommendationRepository domain.RecommendationRepository
}

func (uc *getRecommendationsByStockUseCase) Execute(ctx context.Context, stockID uint, filters pkg.PaginationFilter) (*pkg.PaginationReponse[domain.PopulatedRecommendation], error) {

	stock, err := uc.stockRepository.Get(ctx, stockID)

	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("Stock not found")
	}

	filters.FilterBy = []pkg.FilterByItem{{Field: "stock_id", Value: stock.ID}}

	return uc.recommendationRepository.GetAllPaginated(ctx, filters)

}

func NewGetRecommendationsByStockUseCase(
	ur domain.StockRepository,
	rr domain.RecommendationRepository,
) *getRecommendationsByStockUseCase {

	if ur == nil {
		panic("stock repository is nil, stopping :b")
	}

	if rr == nil {
		panic("recommendation repository is nil, stopping :b")
	}

	return &getRecommendationsByStockUseCase{
		stockRepository:          ur,
		recommendationRepository: rr,
	}
}
