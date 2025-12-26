package domain

import "time"

/*
Market
Represents the market where the stock is traded.
*/
type Market struct {
	ID        uint64    `json:"id,string"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

/*
MarketRepository
Repository for the Market entity.
*/
type MarketRepository interface {
	GetByID(id uint) (Market, error)

	Save(market Market) error
}
