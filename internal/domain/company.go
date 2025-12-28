package domain

import "time"

/*
Company
Represents the company that owns the stock.
*/
type Company struct {
	ID        uint64
	MarketID  uint64
	Name      string
	ISIN      *string
	CreatedAt time.Time
}

/*
CompanyRepository
Repository for the Company entity.
*/
type CompanyRepository interface {
	GetByID(id uint64) (Company, error)

	Save(company Company) error
}
