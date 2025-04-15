package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type Tendency string

const (
	Up   Tendency = "up"
	Down Tendency = "down"
	Side Tendency = "side"
)

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
}

type StockRepository interface {
	Get(ctx context.Context, id uint) (*Stock, error)

	GetByTicker(ctx context.Context, marketID uint, ticker string) (*Stock, error)

	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[PopulatedStock], error)

	Register(ctx context.Context, stock []SourceStockData) error
}
