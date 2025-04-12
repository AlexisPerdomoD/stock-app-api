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
	ID        uint       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type BrokerageRepository interface {
	Get(ctx context.Context, id uint) (*Brokerage, error)

	GetByName(ctx context.Context, name string) (*Brokerage, error)

	Create(ctx context.Context, args *Brokerage) error
}
