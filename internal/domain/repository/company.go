package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type CompanyRepository interface {
	Get(id int) (entity.Company, shared.ApiErr)
	GetByMarketIDAndName(marketID int, name string) (entity.Company, shared.ApiErr)
	Create(args entity.Company) (entity.Company, shared.ApiErr)
}
