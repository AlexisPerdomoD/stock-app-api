package domain

import (
	"context"
	"time"
)

type StockTendencyStat struct {
	ID        uint64
	StockID   uint64
	UpCount   uint64
	SideCount uint64
	DownCount uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StockTendencyDelta struct {
	Up   uint64
	Side uint64
	Down uint64
}

// StockTendencyStatRepository defines persistence operations for stock tendency statistics.
type StockTendencyStatRepository interface {

	// GetByStockID returns the tendency statistics associated with the given stock ID.
	//
	// Returns (nil, nil) if no statistics exist for the provided stock ID.
	GetByStockID(ctx context.Context, stockID uint64) (*StockTendencyStat, error)

	// Increment applies an incremental update to the tendency statistics of a stock.
	//
	// If no statistics exist for the given StockID, a new record is created using
	// the provided delta values as initial counters.
	//
	// Returns an error if:
	//   - the delta contains an invalid StockID
	//   - a persistence or concurrency conflict occurs
	Increment(ctx context.Context, stockID uint64, delta StockTendencyDelta) error

	// IncrementAll applies incremental updates to multiple stock tendency statistics
	// in a single operation.
	//
	// The operation should be atomic when supported by the underlying storage.
	//
	// Returns an error if:
	//   - any delta contains an invalid StockID
	//   - a persistence or concurrency conflict occurs
	IncrementAll(ctx context.Context, stockDeltas map[uint64]StockTendencyDelta) error
}
