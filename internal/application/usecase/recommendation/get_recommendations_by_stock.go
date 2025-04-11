package usecase

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type getRecommendationsByStockUseCase struct {
	stockRepository          domain.StockRepository
	recommendationRepository domain.RecommendationRepository
}

func (uc *getRecommendationsByStockUseCase) Execute() {
	panic("getRecommendationsUseCase.Execute() not implemented")
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
