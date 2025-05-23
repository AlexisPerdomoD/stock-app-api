package cockroachdb

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"gorm.io/gorm"
)

type marketRecord struct {
	gorm.Model
	Name      string          `gorm:"unique"`
	Companies []companyRecord `gorm:"foreignKey:MarketID"`
}

type companyRecord struct {
	gorm.Model
	MarketID uint
	Market   marketRecord
	Name     string
	ISIN     *string `gorm:"unique"`
}

type brokerageRecord struct {
	gorm.Model
	Name            string                 `gorm:"unique"`
	Recommendations []recommendationRecord `gorm:"foreignKey:BrokerageID"`
}

type recommendationRecord struct {
	gorm.Model
	StockID     uint
	Stock       stockRecord
	BrokerageID uint
	Brokerage   brokerageRecord
	RatingTo    domain.Action
	RatingFrom  domain.Action
	TargetTo    float64
	TargetFrom  float64
}

type stockRecord struct {
	gorm.Model
	Name      *string
	CompanyID uint `gorm:"uniqueIndex:idx_ticker_company"`
	Company   companyRecord
	Ticker    string `gorm:"uniqueIndex:idx_ticker_company"`
	Price     float64
	Tendency  domain.Tendency
	Users     []userRecord `gorm:"many2many:user_stocks;"`
}

type userRecord struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
	Stocks   []stockRecord `gorm:"many2many:user_stocks;"`
	Active   bool
}

/* gorm struct table naming overide for conviniance */

func (marketRecord) TableName() string {
	return "markets"
}

func (companyRecord) TableName() string {
	return "companies"
}

func (brokerageRecord) TableName() string {
	return "brokerages"
}

func (stockRecord) TableName() string {
	return "stocks"
}

func (recommendationRecord) TableName() string {
	return "recommendations"
}

func (userRecord) TableName() string {
	return "users"
}
