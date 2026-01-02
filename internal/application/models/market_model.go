package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type MarketView struct {
	ID        uint64    `json:"id,string" example:"1"`
	Name      string    `json:"name" example:"Market 1"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T12:00:00Z"`
}

func NewMarketView(m domain.Market) MarketView {
	return MarketView{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}
}
