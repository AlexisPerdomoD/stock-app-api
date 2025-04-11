package domain

import "time"

/*
Brokerage
Represents Analytics Brokerage teams that are responsible for the stocks recommendations.
*/
type Brokerage struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type BrokerageRepository interface {
	Get(id int) (*Brokerage, error)

	GetByName(name string) (*Brokerage, error)

	Create(args *Brokerage) error
}
