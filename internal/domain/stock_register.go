package domain

import (
	"context"
	"time"
)

type Tendency uint8

const (
	Up   Tendency = 1
	Side Tendency = 2
	Down Tendency = 3
)

type StockRegister struct {
	ID        uint64    `json:"id,string"`
	StockID   uint64    `json:"stock_id,string"`
	Price     float64   `json:"price"`
	Tendency  Tendency  `json:"tendency"`
	CreatedAt time.Time `json:"created_at"`
}

type StockRegisterRepository interface {
	GetByID(ctx context.Context, stockRegisterID uint64) (*StockRegister, error)

	GetLastByStockID(ctx context.Context, stockID uint64) (*StockRegister, error)

	GetRangeByStockID(ctx context.Context, stockID uint64, from, to time.Time) ([]StockRegister, error)

	Save(ctx context.Context, stockRegister StockRegister) error
}
