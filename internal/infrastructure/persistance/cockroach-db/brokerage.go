package cockroachdb

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"gorm.io/gorm"
)

type brokerageRepository struct {
	db *gorm.DB
}

func (r *brokerageRepository) Get(ctx context.Context, id uint) (*domain.Brokerage, error) {

	panic("Implement me")

}

func (r *brokerageRepository) GetByName(ctx context.Context, name string) (*domain.Brokerage, error) {

	panic("Implement me")

}

func (r *brokerageRepository) Create(ctx context.Context, id uint) error {

	panic("Implement me")

}

func NewBrokerageRepository(db *gorm.DB) *brokerageRepository {
	if db == nil {
		log.Fatalf("bad impl: db not provided for NewBrokerageRepository")
	}

	return &brokerageRepository{db}
}
