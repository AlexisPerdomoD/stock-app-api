package aggregate

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"

type PopulatedStock struct {
	entity.Stock
	CompanyName string `json:"company_name"`
	Market      entity.Market
}
