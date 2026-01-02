package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type StockView struct {
	ID        uint64    `json:"id,string" example:"1"`
	CompanyID uint64    `json:"company_id,string" example:"1"`
	MarketID  uint64    `json:"market_id,string" example:"1"`
	Ticker    string    `json:"ticker" example:"AAPL"`
	Isin      *string   `json:"isin,omitempty" example:"US0378331005"`
	Name      *string   `json:"name,omitempty" example:"apple inc"`
	CreatedAt time.Time `json:"registered" example:"2021-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at"  example:"2021-01-01T12:00:00Z"`
}

func NewStockView(s domain.Stock) StockView {

	return StockView{
		ID:        s.ID,
		CompanyID: s.CompanyID,
		MarketID:  s.MarketID,
		Ticker:    s.Ticker,
		Name:      s.Name,
		Isin:      s.Isin,
		CreatedAt: s.CreatedAt,
	}
}

type StockRegisterView struct {
	ID        uint64    `json:"id,string" example:"1"`
	StockID   uint64    `json:"stock_id,string" example:"1"`
	Price     float64   `json:"price" example:"100,00"`
	Tendency  string    `json:"tendency" example:"up"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T12:00:00Z"`
}

func NewStockRegisterView(r domain.StockRegister) StockRegisterView {
	return StockRegisterView{
		ID:        r.ID,
		StockID:   r.StockID,
		Price:     r.Price,
		Tendency:  r.Tendency.String(),
		CreatedAt: r.CreatedAt,
	}
}

type PopulatedStockView struct {
	StockView
	Company      CompanyView       `json:"company"`
	Market       MarketView        `json:"market"`
	LastRegister StockRegisterView `json:"last_register"`
	IsSaved      *bool             `json:"is_saved,omitempty" example:"true"`
}

func NewPopulatedStockView(s domain.PopulatedStock) PopulatedStockView {

	return PopulatedStockView{
		StockView:    NewStockView(s.Stock),
		Company:      NewCompanyView(s.Company),
		Market:       NewMarketView(s.Market),
		LastRegister: NewStockRegisterView(s.LastRegister),
		IsSaved:      s.IsSaved,
	}
}
