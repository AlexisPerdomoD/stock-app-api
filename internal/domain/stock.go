package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type Tendency uint8

const (
	Up   Tendency = 1
	Side Tendency = 2
	Down Tendency = 3
)

/*
Stock
Represents a stock.
*/
type Stock struct {
	ID        uint64    `json:"id,string"`
	CompanyID uint64    `json:"company_id,string"`
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
	Get(ctx context.Context, stockID uint64, userID *uint64) (*PopulatedStock, error)

	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter, userID *uint64) (*pkg.PaginationReponse[PopulatedStock], error)

	Save(ctx context.Context, stock *Stock) error

	Update(ctx context.Context, stock StockUpdates) error
}
