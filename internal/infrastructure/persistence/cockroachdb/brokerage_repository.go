package cockroachdb

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type BrokerageRepository struct {
	db sqlx.ExtContext
}

func (r *BrokerageRepository) GetByID(ctx context.Context, brokerageID uint64) (*domain.Brokerage, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *BrokerageRepository) GetByNames(ctx context.Context, names []string) (map[string]*domain.Brokerage, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r *BrokerageRepository) Save(ctx context.Context, brokerage *domain.Brokerage) error {
	return pkg.InternalServerError("not implemented")
}

func (r *BrokerageRepository) SaveAll(ctx context.Context, brokerages []*domain.Brokerage) error {
	return pkg.InternalServerError("not implemented")
}

func NewBrokerageRepository(db sqlx.ExtContext) *BrokerageRepository {
	if db == nil {
		log.Fatalln("required db sqlx.ExtContext passed as nil")
	}

	return &BrokerageRepository{db}
}
