package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

/*
Stock
Represents a  stock catalog.
*/
type Stock struct {
	ID        uint64    `json:"id,string"`
	CompanyID uint64    `json:"company_id,string"`
	Ticker    string    `json:"ticker"`
	Name      *string   `json:"name,omitempty"`
	CreatedAt time.Time `json:"registered"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockUpdates struct {
	Name *string
}

type PopulatedStock struct {
	Stock
	Company      Company       `json:"company"`
	Market       Market        `json:"market"`
	LastRegister StockRegister `json:"last_register"`
	IsSaved      *bool         `json:"is_saved,omitempty"`
}

type StockRepository interface {
	Get(ctx context.Context, stockID uint64, userID *uint64) (*Stock, error)

	GetPopulated(ctx context.Context, stockID uint64, userID *uint64) (*PopulatedStock, error)

	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter, userID *uint64) (*pkg.PaginationReponse[PopulatedStock], error)

	Save(ctx context.Context, stock *Stock) error

	Update(ctx context.Context, stock StockUpdates) error
}
