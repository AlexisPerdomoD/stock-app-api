package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/aggregate"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type RecommendationRepository interface {
	Get(id int) (*entity.Recommendation, error)

	GetAllPaginated(filter shared.PaginationFilter) (*shared.PaginationReponse[aggregate.PopulatedRecommendation], error)

	Create(args *entity.Recommendation) error
}
