package domain

import (
	"context"
	"time"
)

/*
Brokerage
Represents Analytics Brokerage teams that are responsible for the stocks recommendations.
*/
type Brokerage struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
}

/*
BrokerageRepository
Repository for the Brokerage entity.
*/
type BrokerageRepository interface {
	/*
	   returns a brokerage by its id
	*/
	GetByID(ctx context.Context, brokerageID uint64) (*Brokerage, error)

	/*
		returns brokerages by their names, if not found any it is set to nil in the final map
	*/
	GetByNames(ctx context.Context, names []string) (map[string]*Brokerage, error)

	/*
		Saves a brokerage in the repository

		- returns error if nil arguments are passed.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	Save(ctx context.Context, brokerage *Brokerage) error

	/*
		Saves brokerages in the repository

		- nil brokerages slice is no-op and returns nil.

		- nil values inside brokerages slice returns an error.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	SaveAll(ctx context.Context, brokerages []*Brokerage) error
}
