package domain

import "time"

/*
Company
Represents the company that owns the stock.
*/
type Company struct {
	ID        uint      `json:"id,string"`
	MarketID  uint      `json:"market_id,string"`
	Name      string    `json:"name"`
	ISIN      *string   `json:"isin"`
	CreatedAt time.Time `json:"created_at"`
}

/*
CompanyRepository
Repository for the Company entity.
*/
type CompanyRepository interface {
	GetByID(id uint) (Company, error)

	Save(company Company) error
}
