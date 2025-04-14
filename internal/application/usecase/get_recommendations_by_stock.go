package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type getRecommendationsByStockUseCase struct {
	sr domain.StockRepository
	rr domain.RecommendationRepository
}

func (uc *getRecommendationsByStockUseCase) Execute(
	ctx context.Context,
	stockID uint,
	filters pkg.PaginationFilter,
) (*pkg.PaginationReponse[domain.PopulatedRecommendation], error) {

	stock, err := uc.sr.Get(ctx, stockID)

	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("Stock not found")
	}

	filters.FilterBy = []pkg.FilterByItem{{Field: "stock_id", Value: stock.ID}}

	return uc.rr.GetAllPaginated(ctx, filters)

}

func NewGetRecommendationsByStockUseCase(
	sr domain.StockRepository,
	rr domain.RecommendationRepository,
) *getRecommendationsByStockUseCase {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetRecommendationsByStockUseCase")
	}

	if rr == nil {
		log.Fatalln("bad impl: RecommendationRepository was nil for NewGetRecommendationsByStockUseCase")
	}

	return &getRecommendationsByStockUseCase{sr, rr}
}
