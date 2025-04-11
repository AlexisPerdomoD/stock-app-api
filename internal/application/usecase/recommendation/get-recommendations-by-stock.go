package usecase

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/repository"

type getRecommendationsByStockUseCase struct {
	stockRepository          repository.StockRepository
	recommendationRepository repository.RecommendationRepository
}

func (uc *getRecommendationsByStockUseCase) Execute() {
	panic("getRecommendationsUseCase.Execute() not implemented")
}

func NewGetRecommendationsByStockUseCase(
	ur repository.StockRepository,
	rr repository.RecommendationRepository,
) *getRecommendationsByStockUseCase {

	if ur == nil {
		panic("stock repository is nil, stopping :b")
	}

	if rr == nil {
		panic("recommendation repository is nil, stopping :b")
	}

	return &getRecommendationsByStockUseCase{
		stockRepository: ur, 
		recommendationRepository: rr,
	}
}
