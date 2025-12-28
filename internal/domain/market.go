package domain

import "time"

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
	GetByID(id uint64) (Market, error)

	Save(market Market) error
}
