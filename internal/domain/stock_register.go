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

func (t Tendency) String() string {
	switch t {
	case Up:
		return "up"
	case Side:
		return "side"
	case Down:
		return "down"
	default:
		return "unknown"
	}
}

type StockRegister struct {
	ID        uint64
	StockID   uint64
	Price     float64
	Tendency  Tendency
	CreatedAt time.Time
}

type StockRegisterRepository interface {

	/*
	   returns a stock register by its id
	*/
	GetByID(ctx context.Context, stockRegisterID uint64) (*StockRegister, error)

	/*
	   returns last stock register by its stock id
	*/
	GetLastByStockID(ctx context.Context, stockID uint64) (*StockRegister, error)

	/*
	   returns lasts stock register by its stock id sorted by date desc and limited by limit count
	*/
	GetLastsByStockID(ctx context.Context, stockID uint64, limit uint16) ([]StockRegister, error)

	/*
	   returns stock registers within a date range ordered by date desc
	*/
	GetRangeByStockID(ctx context.Context, stockID uint64, from, to time.Time) ([]StockRegister, error)

	/*
		saves a stock register in the repository

		- returns error if nil arguments are passed

		- returns error if conflict occurs with arguments provided
	*/
	Save(ctx context.Context, stockRegister *StockRegister) error
	/*
		saves stock registers in the repository

		- returns error if nil arguments are passed

		- returns error if conflict occurs with arguments provided
	*/
	SaveAll(ctx context.Context, stockRegister []*StockRegister) error
}
