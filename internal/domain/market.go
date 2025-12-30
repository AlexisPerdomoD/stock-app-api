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

		- nil market is no-op and returns nil.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	Save(ctx context.Context, market *Market) error

	/*
		Saves all Markets in the database and map missing properties with their default values (if any) including ID.

		- nil slice argument is no-op and return nil.

		- nil values inside the markets slice returns an error.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	SaveAll(ctx context.Context, markets []*Market) error
}
