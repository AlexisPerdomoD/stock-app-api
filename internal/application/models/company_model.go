package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type CompanyView struct {
	ID        uint64    `json:"id,string"`
	MarketID  uint64    `json:"market_id,string"`
	Name      string    `json:"name"`
	ISIN      *string   `json:"isin"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCompanyView(c domain.Company) CompanyView {

	return CompanyView{
		ID:        c.ID,
		MarketID:  c.MarketID,
		Name:      c.Name,
		ISIN:      c.ISIN,
		CreatedAt: c.CreatedAt,
	}
}
