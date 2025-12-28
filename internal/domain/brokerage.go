package domain

import "time"

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
	GetByID(id uint64) (Brokerage, error)

	Save(brokerage Brokerage) error
}
