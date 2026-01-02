package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type CompanyView struct {
	ID        uint64    `json:"id,string" example:"1"`
	MarketID  uint64    `json:"market_id,string" example:"1"`
	Name      string    `json:"name" example:"Company 1"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T12:00:00Z"`
}

func NewCompanyView(c domain.Company) CompanyView {

	return CompanyView{
		ID:        c.ID,
		MarketID:  c.MarketID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}
