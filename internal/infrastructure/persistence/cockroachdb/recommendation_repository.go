package cockroachdb

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type RecommendationRepository struct{}

func (r RecommendationRepository) GetAllPaginated(
	ctx context.Context,
	filter pkg.PaginationFilter,
	stockID uint64,
) (*pkg.PaginationReponse[domain.PopulatedRecommendation], error) {
	return nil, pkg.InternalServerError("not implemented")
}

func NewRecommendationRepository(db *sqlx.DB) *RecommendationRepository {
	return &RecommendationRepository{}
}
