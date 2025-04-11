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

type Action string

const (
	Buy     Action = "buy"
	Sell    Action = "sell"
	Neutral Action = "neutral"
	Hold    Action = "hold"
	Unknown Action = "unknown"
)

/*
Stock
Represents a stock.
*/
type Stock struct {
	ID        int       `json:"id"`
	CompanyID int       `json:"company_id"`
	Ticker    string    `json:"ticker"`
	Name      string    `json:"name,omitempty"`
	Price     float64   `json:"price"`
	Tendency  Tendency  `json:"tendency"`
	CreatedAt time.Time `json:"registered"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PopulatedStock struct {
	Stock
	CompanyName string `json:"company_name"`
	Market      Market `json:"market"`
}

type StockRepository interface {
	Get(ctx context.Context, id int) (*Stock, error)

	GetByTicker(ctx context.Context, marketID int, ticker string) (*Stock, error)

	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[PopulatedStock], error)

	Create(ctx context.Context, args *Stock) error

	Update(ctx context.Context, stockID int, args *Stock) error
}
