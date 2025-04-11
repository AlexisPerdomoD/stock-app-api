package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
)

type CompanyRepository interface {
	Get(id int) (entity.Company, error)

	GetByMarketIDAndName(marketID int, name string) (entity.Company, error)

	Create(args entity.Company) (entity.Company, error)
}
