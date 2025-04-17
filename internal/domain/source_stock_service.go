package domain

import (
	"context"
	"time"
)

type MarketArgs struct {
	Name string
}

type CompanyArgs struct {
	Name string
	ISIN *string
}

type BrokerageArgs struct {
	Name string
}

type RecommendationArgs struct {
	RatingTo   Action
	RatingFrom Action
	TargetTo   float64
	TargetFrom float64
	Brokerage  BrokerageArgs
}

type StockArgs struct {
	Ticker   string
	Name     string
	Price    float64
	Tendency Tendency
}

type SourceStockData struct {
	Market MarketArgs

	Company CompanyArgs

	Recomendation *RecommendationArgs

	Stock StockArgs

	Time time.Time
}

type SourceStockService interface {
	Name() string
	Get(ctx context.Context, limitDate *time.Time) ([]SourceStockData, error)
}
