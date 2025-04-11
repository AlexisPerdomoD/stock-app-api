package domain

import "time"

/*
Company
Represents the company that owns the stock.
*/
type Company struct {
	ID        int       `json:"id"`
	MarketID  int       `json:"market_id"`
	Name      string    `json:"name"`
	ISIN      string    `json:"isin"`
	CreatedAt time.Time `json:"created_at"`
}

type CompanyRepository interface {
	Get(id int) (Company, error)

	GetByMarketIDAndName(marketID int, name string) (Company, error)

	Create(args Company) (Company, error)
}
