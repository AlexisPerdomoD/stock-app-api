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

		- returns error if nil arguments are passed
	*/
	GetByNames(ctx context.Context, names []string) (map[string]*Brokerage, error)

	/*
		Saves a brokerage in the repository

		- returns nil if nil arguments are passed

		- returns error if conflict occurs with arguments provided
	*/
	Save(ctx context.Context, brokerage *Brokerage) error

	/*
		Saves brokerages in the repository

		- returns nil if nil brokerages slice is passed

		- returns error if nil pointers are passed inside brokerages slice

		- returns error if conflict occurs with arguments provided
	*/
	SaveAll(ctx context.Context, brokerages []*Brokerage) error
}
