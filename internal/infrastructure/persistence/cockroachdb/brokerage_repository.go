package cockroachdb

import (
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type BrokerageRepository struct {
	db sqlx.ExtContext
}

func (r BrokerageRepository) GetByID(id uint64) (*domain.Brokerage, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r BrokerageRepository) Save(brokerage *domain.Brokerage) error {
	return pkg.InternalServerError("not implemented")
}

func NewBrokerageRepository(db sqlx.ExtContext) *BrokerageRepository {
	if db == nil {
		log.Fatalln("required db sqlx.ExtContext passed as nil")
	}

	return &BrokerageRepository{db}
}
