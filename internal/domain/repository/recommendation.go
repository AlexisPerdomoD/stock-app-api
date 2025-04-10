package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/aggregate"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type RecommendationRepository interface {
	Get(id int) (entity.Recommendation, shared.ApiErr)
	GetAllPaginated(filter shared.PaginationFilter) (shared.PaginationReponse[aggregate.PopulatedRecommendation], shared.ApiErr)

	Create(args entity.Recommendation) (entity.Recommendation, shared.ApiErr)
}
