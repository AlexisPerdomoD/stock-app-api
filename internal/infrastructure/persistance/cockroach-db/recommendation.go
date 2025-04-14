package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type recommendationRepository struct {
	db *gorm.DB
}

func (r *recommendationRepository) Get(ctx context.Context, id uint) (*domain.Recommendation, error) {

	record := &recommendationRecord{}

	if err := r.db.WithContext(ctx).First(record, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return mapRecommendationToDomain(record, nil), nil
}

func (r *recommendationRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[domain.Recommendation], error) {

	var records []recommendationRecord
	var total int64

	allowedFilters := map[string]bool{
		"stock_id":     true,
		"created_at":   true,
		"brokerage_id": true,
	}

	allowedSorters := map[string]bool{
		"stock_id":     true,
		"brokerage_id": true,
		"created_at":   true,
	}

	query := r.db.WithContext(ctx).Model(recommendationRecord{})
	query = applyFilters(query, filter.FilterBy, allowedFilters)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := applyPagination(query, &filter, allowedSorters).
		Find(&records).Error; err != nil {
		return nil, err
	}

	recommendations := []domain.Recommendation{}

	for _, record := range records {
		recommendation := domain.Recommendation{}
		_ = mapRecommendationToDomain(&record, &recommendation)
		recommendations = append(recommendations, recommendation)
	}

	page := 1
	if filter.Page > 1 {
		page = filter.Page
	}

	result := &pkg.PaginationReponse[domain.Recommendation]{
		Items:     recommendations,
		Page:      page,
		PageSize:  len(recommendations),
		TotalSize: int(total),
	}

	return result, nil
}

func (r *recommendationRepository) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	if recommendation == nil {
		return pkg.BadRequest("args for recommendation insertion were not provided")
	}

	record := mapRecommendationInsert(recommendation)

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}

	_ = mapRecommendationToDomain(record, recommendation)

	return nil
}

func NewRecommendationRepository(db *gorm.DB) *recommendationRepository {
	if db == nil {
		log.Fatalf("bad impl: db is nil in NewRecommendationRepository")
	}

	return &recommendationRepository{db: db}
}
