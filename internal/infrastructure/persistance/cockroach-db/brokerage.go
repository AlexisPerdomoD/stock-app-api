package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"gorm.io/gorm"
)

type brokerageRepository struct {
	db *gorm.DB
}

func (r *brokerageRepository) Get(ctx context.Context, id uint) (*domain.Brokerage, error) {
	record := &brokerageRecord{}

	if err := r.db.WithContext(ctx).First(record, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err

	}

	return mapBrokerageToDomain(record, nil), nil

}

func (r *brokerageRepository) GetByName(ctx context.Context, name string) (*domain.Brokerage, error) {

	record := &brokerageRecord{}

	if err := r.db.WithContext(ctx).First(record, "'name' = ?", name).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err

	}

	return mapBrokerageToDomain(record, nil), nil

}

func (r *brokerageRepository) Create(ctx context.Context, br *domain.Brokerage) error {

	record := mapBrokerageInsert(br)

	if err := r.db.WithContext(ctx).Create(record).Error;err != nil{
		return err
	}

	_ = mapBrokerageToDomain(record, br)
	return nil
}

func NewBrokerageRepository(db *gorm.DB) *brokerageRepository {
	if db == nil {
		log.Fatalf("bad impl: db not provided for NewBrokerageRepository")
	}

	return &brokerageRepository{db}
}
