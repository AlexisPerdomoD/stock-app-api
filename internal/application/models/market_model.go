package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type MarketView struct {
	ID        uint64    `json:"id,string"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMarketView(m domain.Market) MarketView {
	return MarketView{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}
}
