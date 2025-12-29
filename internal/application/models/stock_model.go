package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type StockView struct {
	ID        uint64    `json:"id,string"`
	CompanyID uint64    `json:"company_id,string"`
	MarketID  uint64    `json:"market_id,string"`
	Ticker    string    `json:"ticker"`
	Isin      *string   `json:"isin,omitempty"`
	Name      *string   `json:"name,omitempty"`
	CreatedAt time.Time `json:"registered"`
	UpdatedAt time.Time `json:"updated_at"`
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
	ID        uint64    `json:"id,string"`
	StockID   uint64    `json:"stock_id,string"`
	Price     float64   `json:"price"`
	Tendency  string    `json:"tendency"`
	CreatedAt time.Time `json:"created_at"`
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
	IsSaved      *bool             `json:"is_saved,omitempty"`
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
