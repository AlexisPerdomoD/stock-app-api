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
	ID        uint64
	CompanyID uint64
	Ticker    string
	Name      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StockUpdates struct {
	Name *string
}

type PopulatedStock struct {
	Stock
	Company      Company
	Market       Market
	LastRegister StockRegister
	IsSaved      *bool
}

type StockRepository interface {
	Get(ctx context.Context, stockID uint64, userID *uint64) (*Stock, error)

	GetPopulated(ctx context.Context, stockID uint64, userID *uint64) (*PopulatedStock, error)

	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter, userID *uint64) (*pkg.PaginationReponse[PopulatedStock], error)

	Save(ctx context.Context, stock *Stock) error

	Update(ctx context.Context, stock StockUpdates) error
}
