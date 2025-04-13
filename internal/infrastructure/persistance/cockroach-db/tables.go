package cockroachdb

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"gorm.io/gorm"
)

type marketRecord struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type companyRecord struct {
	gorm.Model
	MarketID uint
	Name     string
	ISIN     *string `gorm:"unique"`
}

type brokerageRecord struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type stockRecord struct {
	gorm.Model
	Name      *string
	CompanyID uint   `gorm:"uniqueIndex:idx_ticker_company"`
	Ticker    string `gorm:"uniqueIndex:idx_ticker_company"`
	Price     float64
	Tendency  domain.Tendency
}

type recommendationRecord struct {
	gorm.Model
	StockID     uint
	BrokerageID uint
	RatingTo    domain.Action
	RatingFrom  domain.Action
	TargetTo    float64
	TargetFrom  float64
}

type userRecord struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
	Active   bool
}

type userStockRecord struct {
	gorm.Model
	UserID  uint `gorm:"uniqueIndex:idx_user_stock"`
	StockID uint `gorm:"uniqueIndex:idx_user_stock"`
	Count   uint `gorm:"default:1"`
}
