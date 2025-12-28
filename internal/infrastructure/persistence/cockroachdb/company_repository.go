package cockroachdb

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type CompanyRepository struct{}

func (r CompanyRepository) GetByID(id uint64) (*domain.Company, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r CompanyRepository) Save(company *domain.Company) error {
	return pkg.InternalServerError("not implemented")
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{}
}
