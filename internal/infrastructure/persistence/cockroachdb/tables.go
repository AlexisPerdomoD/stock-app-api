package cockroachdb

import (
	"database/sql"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type userRecord struct {
	ID        uint64    `db:"id"`
	Username  string    `db:"username"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Password  []byte    `db:"password"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r userRecord) ToDomain() *domain.User {
	return &domain.User{
		ID:        r.ID,
		Username:  r.Username,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Password:  r.Password,
		Active:    r.Active,
		CreatedAt: r.CreatedAt,
	}
}

type marketRecord struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *marketRecord) ToDomain() *domain.Market {
	return &domain.Market{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

func (r *marketRecord) MapDomain(dom *domain.Market) {
	dom.ID = r.ID
	dom.Name = r.Name
	dom.CreatedAt = r.CreatedAt
}

type companyRecord struct {
	ID        uint64    `db:"id"`
	MarketID  uint64    `db:"market_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *companyRecord) ToDomain() *domain.Company {
	return &domain.Company{
		ID:        r.ID,
		MarketID:  r.MarketID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

func (r *companyRecord) MapDomain(dom *domain.Company) {
	dom.ID = r.ID
	dom.MarketID = r.MarketID
	dom.Name = r.Name
	dom.CreatedAt = r.CreatedAt
}

type brokerageRecord struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *brokerageRecord) ToDomain() *domain.Brokerage {
	return &domain.Brokerage{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

func (r *brokerageRecord) MapDomain(dom *domain.Brokerage) {
	dom.ID = r.ID
	dom.Name = r.Name
	dom.CreatedAt = r.CreatedAt
}

type stockRecord struct {
	ID        uint64         `db:"id"`
	CompanyID uint64         `db:"company_id"`
	MarketID  uint64         `db:"market_id"`
	Ticker    string         `db:"ticker"`
	Name      sql.NullString `db:"name"`
	Isin      sql.NullString `db:"isin"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (r stockRecord) ToDomain() *domain.Stock {
	var name *string = nil
	var isin *string = nil

	if r.Name.Valid {
		name = &r.Name.String
	}

	if r.Isin.Valid {
		isin = &r.Isin.String
	}

	return &domain.Stock{
		ID:        r.ID,
		CompanyID: r.CompanyID,
		MarketID:  r.MarketID,
		Name:      name,
		Isin:      isin,
		Ticker:    r.Ticker,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type stockRegisterRecord struct {
	ID        uint64          `db:"id"`
	StockID   uint64          `db:"stock_id"`
	Price     float64         `db:"price"`
	Tendency  domain.Tendency `db:"tendency"`
	CreatedAt time.Time       `db:"created_at"`
}

func (r stockRegisterRecord) ToDomain() *domain.StockRegister {
	return &domain.StockRegister{
		ID:        r.ID,
		StockID:   r.StockID,
		Price:     r.Price,
		Tendency:  r.Tendency,
		CreatedAt: r.CreatedAt,
	}
}

// RELATIONSHIPS

type stockUserRecord struct {
	ID        uint64    `db:"id"`
	StockID   uint64    `db:"stock_id"`
	UserID    uint64    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type stockRecommendationRecord struct {
	id              uint64        `db:"id"`
	brokerageID     uint64        `db:"brokerage_id"`
	stockRegisterID uint64        `db:"stock_register_id"`
	ratingFrom      domain.Action `db:"rating_from"`
	ratingTo        domain.Action `db:"rating_to"`
	targetFrom      float64       `db:"target_from"`
	targetTo        float64       `db:"target_to"`
	createdAt       time.Time     `db:"created_at"`
}

func (r stockRecommendationRecord) ToDomain() *domain.Recommendation {
	return &domain.Recommendation{
		ID:              r.id,
		BrokerageID:     r.brokerageID,
		StockRegisterID: r.stockRegisterID,
		RatingFrom:      r.ratingFrom,
		RatingTo:        r.ratingTo,
		TargetFrom:      r.targetFrom,
		TargetTo:        r.targetTo,
		CreatedAt:       r.createdAt,
	}
}

// static queries

type stockTendencyStatRecord struct {
	ID        uint64    `db:"id"`
	StockID   uint64    `db:"stock_id"`
	UpCount   uint64    `db:"up_count"`
	SideCount uint64    `db:"side_count"`
	DownCount uint64    `db:"down_count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *stockTendencyStatRecord) ToDomain() *domain.StockTendencyStat {
	return &domain.StockTendencyStat{
		ID:        r.ID,
		StockID:   r.StockID,
		UpCount:   r.UpCount,
		SideCount: r.SideCount,
		DownCount: r.DownCount,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
