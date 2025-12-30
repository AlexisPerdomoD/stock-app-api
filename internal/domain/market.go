package domain

import (
	"context"
	"time"
)

/*
Market
Represents the market where the stock is traded.
*/
type Market struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
}

/*
MarketRepository
Repository for the Market entity.
*/
type MarketRepository interface {
	/*
		Returns a Market by its ID. If the ID does not exist, returns nil.
	*/
	GetByID(ctx context.Context, id uint64) (*Market, error)

	/*
		Returns a map of Markets by their names. If a name does not exist, sets the value to nil.
	*/
	GetByNames(ctx context.Context, names []string) (map[string]*Market, error)

	/*
		Saves a Market in the database and map missing properties with their default values (if any) including ID.

		- returns nil if argument is nil.

		- invalid constraints returns an error.
	*/
	Save(ctx context.Context, market *Market) error

	/*
		Saves all Markets in the database and map missing properties with their default values (if any) including ID.

		- nil slice argument return nil.

		- nil values inside the slice returns an error.

		- invalid constraints returns an error.
	*/
	SaveAll(ctx context.Context, markets []*Market) error
}
