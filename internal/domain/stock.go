package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type Tendency uint8

const (
	Up   Tendency = 1
	Side Tendency = 2
	Down Tendency = 3
)

/*
Market
Represents the market where the stock is traded.
*/
type Market struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

/*
Company
Represents the company that owns the stock.
*/
type Company struct {
	ID        uint      `json:"id"`
	MarketID  uint      `json:"market_id"`
	Name      string    `json:"name"`
	ISIN      *string   `json:"isin"`
	CreatedAt time.Time `json:"created_at"`
}

/*
Stock
Represents a stock.
*/
type Stock struct {
	ID        uint      `json:"id"`
	CompanyID uint      `json:"company_id"`
	Ticker    string    `json:"ticker"`
	Name      *string   `json:"name,omitempty"`
	Price     float64   `json:"price"`
	Tendency  Tendency  `json:"tendency"`
	CreatedAt time.Time `json:"registered"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockUpdates struct {
	Name     *string
	Price    *float64
	Tendency *Tendency
}

type PopulatedStock struct {
	Stock
	CompanyName string `json:"company_name"`
	Market      Market `json:"market"`
	IsSaved     *bool  `json:"is_saved,omitempty"`
}

type StockRepository interface {
	Get(ctx context.Context, stockID uint, userID *uint) (*PopulatedStock, error)

	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter, userID *uint) (*pkg.PaginationReponse[PopulatedStock], error)

	Register(ctx context.Context, stock []SourceStockData) error
}
