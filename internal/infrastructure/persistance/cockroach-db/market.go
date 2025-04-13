package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type marketRepository struct {
	db *gorm.DB
}

func (r *marketRepository) Get(ctx context.Context, id uint) (*domain.Market, error) {
	record := &marketRecord{}

	err := r.db.WithContext(ctx).First(record, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return mapMarketToDomain(record, nil), nil

}

func (r *marketRepository) GetByName(ctx context.Context, name string) (*domain.Market, error) {

	record := &marketRecord{}

	err := r.db.WithContext(ctx).First(record, "'Name' = ?", name).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return mapMarketToDomain(record, nil), nil

}

func (r *marketRepository) Create(ctx context.Context, market *domain.Market) error {
	if market == nil {
		return pkg.BadRequest("args for Market insertion were not provided")
	}

	record := mapMarketInsert(market)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}

	_ = mapMarketToDomain(record, nil)

	return nil
}

func NewMarketRepository(db *gorm.DB) *marketRepository {

	if db == nil {
		log.Fatalf("bad impl: db not provided for NewMarketRepository call")
	}

	return &marketRepository{db}
}
