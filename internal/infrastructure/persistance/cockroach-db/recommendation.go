package cockroachdb

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type recommendationRepository struct {
	db *gorm.DB
}

func (r *recommendationRepository) Get(ctx context.Context, id uint) (*domain.Recommendation, error) {

	panic("implement me")

}

func (r *recommendationRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[domain.Recommendation], error) {

	panic("implement me")
}

func (r *recommendationRepository) Create(ctx context.Context, recommendation *domain.Recommendation) error {

	panic("implement me")
}

func NewRecommendationRepository(db *gorm.DB) *recommendationRepository {
	if db == nil {
		log.Fatalf("bad impl: db is nil in NewRecommendationRepository")
	}

	return &recommendationRepository{db: db}
}
