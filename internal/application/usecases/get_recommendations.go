package usecases

import (
	"context"

	appmodel "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/collection"
)

type GetRecommendationsByStock struct {
	recommendationRepository domain.RecommendationRepository
}

func (uc *GetRecommendationsByStock) Execute(ctx context.Context, filters pkg.PaginationFilter) (*pkg.PaginationResponse[appmodel.PopulatedRecommendationView], error) {
	data, err := uc.recommendationRepository.GetAllPaginated(ctx, filters)
	if err != nil {
		return nil, err
	}

	response := &pkg.PaginationResponse[appmodel.PopulatedRecommendationView]{
		Page:       data.Page,
		PageSize:   data.PageSize,
		TotalSize:  data.TotalSize,
		TotalPages: data.TotalPages,
		Items:      collection.Map(data.Items, appmodel.NewPopulatedRecommendationView),
	}

	return response, nil

}

func NewGetRecommendations(
	rr domain.RecommendationRepository,
) *GetRecommendationsByStock {

	if rr == nil {
		panic("[GetRecommendations]: RecommendationRepository was provided as nil")
	}

	return &GetRecommendationsByStock{rr}
}
