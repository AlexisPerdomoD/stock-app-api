package domain

import (
	"context"
	"time"
)

/*
Company
Represents the company that owns the stock.
*/
type Company struct {
	ID        uint       `json:"id"`
	MarketID  uint       `json:"market_id"`
	Name      string    `json:"name"`
	ISIN      string    `json:"isin"`
	CreatedAt time.Time `json:"created_at"`
}

type CompanyRepository interface {
	Get(ctx context.Context, id uint) (Company, error)

	GetByMarketIDAndName(ctx context.Context, marketID int, name string) (Company, error)

	Create(ctx context.Context, args Company) (Company, error)
}
