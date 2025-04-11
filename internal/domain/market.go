package domain

import "time"

/*
Market
Represents the market where the stock is traded.
*/
type Market struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type MarketRepository interface {
	Get(id int) (*Market, error)

	GetByName(name string) (*Market, error)

	Create(args *Market) error
}
