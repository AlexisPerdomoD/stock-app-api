package domain

import (
	"context"
	"time"
)

/*
Company
Represents the company that owns the stock.
*/
type Company struct {
	ID        uint64
	MarketID  uint64
	Name      string
	CreatedAt time.Time
}

type MarketCompanySearchParam struct {
	MarketID uint64
	Name     string
}

/*
CompanyRepository
Repository for the Company entity.
*/
type CompanyRepository interface {
	/*
		Returns a Company by its ID. If the ID does not exist, returns nil.
	*/
	GetByID(ctx context.Context, id uint64) (*Company, error)

	/*
		Returns Companies by their MarketID and Name. If the ID does not exist, sets the value to nil.
	*/
	GetByMarketCompanySearch(ctx context.Context, searchParams []MarketCompanySearchParam) (map[MarketCompanySearchParam]*Company, error)

	/*
		Saves a Company in the database and map missing properties with their default values (if any) including ID.

		- returns error if nil arguments are passed.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	Save(ctx context.Context, company *Company) error

	/*
		Saves all Companies in the database and map missing properties with their default values (if any) including ID.

		- nil companies slice is no-op and returns nil

		- nil values inside companies slice returns an error.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).

	*/
	SaveAll(ctx context.Context, companies []*Company) error
}
