package domain

import (
	"context"
	"time"
)

type SourceStockData struct {
	Market struct {
		Name string `json:"name"`
	} `json:"market"`

	Company struct {
		Name string `json:"name"`
	} `json:"company"`

	Recomendation *struct {
		Action    Action `json:"action"`
		Brokerage struct {
			Name string `json:"name"`
		} `json:"brokerage"`
	} `json:"recomendation"`

	Stock struct {
		Ticker   string   `json:"ticker"`
		Name     string   `json:"name"`
		Price    float64  `json:"price"`
		Tendency Tendency `json:"tendency"`
	} `json:"stock"`

	Time time.Time `json:"time"`
}

type SourceStockService interface {
	Get(ctx context.Context, limitDate *time.Time) ([]SourceStockData, error)
}
