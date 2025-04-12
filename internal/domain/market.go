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
	ID        uint       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type MarketRepository interface {
	Get(ctx context.Context, id uint) (*Market, error)

	GetByName(ctx context.Context, name string) (*Market, error)

	Create(ctx context.Context, args *Market) error
}
