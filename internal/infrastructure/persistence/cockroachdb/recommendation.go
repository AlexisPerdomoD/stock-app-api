package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type RecommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) *RecommendationRepository {
	if db == nil {
		log.Fatalf("[RecommendationRepository]: CR db provided is nil")
	}

	return &RecommendationRepository{db: db}
}

func (r *RecommendationRepository) Get(ctx context.Context, id uint) (*domain.Recommendation, error) {

	record := &recommendationRecord{}

	if err := r.db.WithContext(ctx).First(record, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return mapRecommendationToDomain(record, nil), nil
}

func (r *RecommendationRepository) GetAllPaginated(
	ctx context.Context,
	filter pkg.PaginationFilter,
	stockID uint,
) (*pkg.PaginationReponse[domain.PopulatedRecommendation], error) {

	var records []recommendationRecord
	var total int64

	allowedFilters := map[string]bool{
		"created_at":   true,
		"brokerage_id": true,
	}

	allowedSorters := map[string]bool{
		"target_from": true,
		"target_to":   true,
		"rating_from": true,
		"rating_to":   true,
		"created_at":  true,
	}

	query := r.db.WithContext(ctx).
		Model(recommendationRecord{}).
		Preload("Brokerage").
		Where("recommendations.stock_id = ?", stockID)

	if filter.Search != "" {
		query.Joins("JOIN brokerages ON brokerages.id = recommendations.brokerage_id").
			Where("brokerages.name = ?", filter.Search)
	}

	query = applyFilters(query, filter.FilterBy, allowedFilters)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := applyPagination(query, &filter, allowedSorters).
		Find(&records).Error; err != nil {
		return nil, err
	}

	recommendations := []domain.PopulatedRecommendation{}

	for _, record := range records {
		populated := domain.PopulatedRecommendation{}
		_ = mapPopulatedRecommendationToDomain(&record, &populated)
		recommendations = append(recommendations, populated)
	}

	page := 1
	if filter.Page > 1 {
		page = filter.Page
	}

	result := &pkg.PaginationReponse[domain.PopulatedRecommendation]{
		Items:     recommendations,
		Page:      page,
		PageSize:  len(recommendations),
		TotalSize: int(total),
	}

	return result, nil
}
