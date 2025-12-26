package domain

import "time"

/*
Brokerage
Represents Analytics Brokerage teams that are responsible for the stocks recommendations.
*/
type Brokerage struct {
	ID        uint64    `json:"id,string"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

/*
BrokerageRepository
Repository for the Brokerage entity.
*/
type BrokerageRepository interface {
	GetByID(id uint) (Brokerage, error)

	Save(brokerage Brokerage) error
}
