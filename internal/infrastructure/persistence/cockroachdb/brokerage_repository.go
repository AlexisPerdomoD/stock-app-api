package cockroachdb

import (
	"database/sql"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type BrokerageRepository struct{}

func (r BrokerageRepository) GetByID(id uint64) (*domain.Brokerage, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r BrokerageRepository) Save(brokerage *domain.Brokerage) error {
	return pkg.InternalServerError("not implemented")
}

func NewBrokerageRepository(db *sql.DB) *BrokerageRepository {
	return &BrokerageRepository{}
}
