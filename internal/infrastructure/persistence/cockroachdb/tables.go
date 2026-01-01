package cockroachdb

import (
	"database/sql"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type userRecord struct {
	ID         uint64    `db:"id"`
	Username   string    `db:"username"`
	Firstname  string    `db:"firstname"`
	Lastname   string    `db:"lastname"`
	Password   []byte    `db:"password"`
	Active     bool      `db:"active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	BatchIndex int       `db:"batch_index"`
}

func (r *userRecord) ToDomain() *domain.User {
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

func (r *userRecord) MapDomain(dom *domain.User) {
	dom.ID = r.ID
	dom.Username = r.Username
	dom.Firstname = r.Firstname
	dom.Lastname = r.Lastname
	dom.Password = r.Password
	dom.Active = r.Active
	dom.CreatedAt = r.CreatedAt
}

type marketRecord struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	BatchIndex int       `db:"batch_index"`
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
	ID         uint64    `db:"id"`
	MarketID   uint64    `db:"market_id"`
	Name       string    `db:"name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	BatchIndex int       `db:"batch_index"`
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
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	BatchIndex int       `db:"batch_index"`
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
	ID         uint64         `db:"id"`
	CompanyID  uint64         `db:"company_id"`
	MarketID   uint64         `db:"market_id"`
	Ticker     string         `db:"ticker"`
	Name       sql.NullString `db:"name"`
	Isin       sql.NullString `db:"isin"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
	BatchIndex int            `db:"batch_index"`
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

func (r stockRecord) MapDomain(dom *domain.Stock) {
	var name *string = nil
	var isin *string = nil

	if r.Name.Valid {
		name = &r.Name.String
	}

	if r.Isin.Valid {
		isin = &r.Isin.String
	}

	dom.ID = r.ID
	dom.MarketID = r.MarketID
	dom.CompanyID = r.CompanyID
	dom.Ticker = r.Ticker
	dom.Name = name
	dom.Isin = isin
}

type stockRegisterRecord struct {
	ID         uint64          `db:"id"`
	StockID    uint64          `db:"stock_id"`
	Price      float64         `db:"price"`
	Tendency   domain.Tendency `db:"tendency"`
	CreatedAt  time.Time       `db:"created_at"`
	BatchIndex int             `db:"batch_index"`
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

func (r *stockRegisterRecord) MapDomain(stockRegister *domain.StockRegister) {
	r.ID = stockRegister.ID
	r.StockID = stockRegister.StockID
	r.Price = stockRegister.Price
	r.Tendency = stockRegister.Tendency
	r.CreatedAt = stockRegister.CreatedAt
}

type stockRecommendationRecord struct {
	ID              uint64        `db:"id"`
	BrokerageID     uint64        `db:"brokerage_id"`
	StockRegisterID uint64        `db:"stock_register_id"`
	RatingFrom      domain.Action `db:"rating_from"`
	RatingTo        domain.Action `db:"rating_to"`
	TargetFrom      float64       `db:"target_from"`
	TargetTo        float64       `db:"target_to"`
	CreatedAt       time.Time     `db:"created_at"`
	BatchIndex      int           `db:"batch_index"`
}

func (r *stockRecommendationRecord) ToDomain() *domain.Recommendation {
	return &domain.Recommendation{
		ID:              r.ID,
		BrokerageID:     r.BrokerageID,
		StockRegisterID: r.StockRegisterID,
		RatingFrom:      r.RatingFrom,
		RatingTo:        r.RatingTo,
		TargetFrom:      r.TargetFrom,
		TargetTo:        r.TargetTo,
		CreatedAt:       r.CreatedAt,
	}
}

func (r *stockRecommendationRecord) MapDomain(dom *domain.Recommendation) {
	dom.ID = r.ID
	dom.BrokerageID = r.BrokerageID
	dom.StockRegisterID = r.StockRegisterID
	dom.RatingFrom = r.RatingFrom
	dom.RatingTo = r.RatingTo
	dom.TargetFrom = r.TargetFrom
	dom.TargetTo = r.TargetTo
	dom.CreatedAt = r.CreatedAt
}

// RELATIONSHIPS

type stockUserRecord struct {
	ID        uint64    `db:"id"`
	StockID   uint64    `db:"stock_id"`
	UserID    uint64    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

// STATIC QUERIES

type stockTendencyStatRecord struct {
	ID         uint64    `db:"id"`
	StockID    uint64    `db:"stock_id"`
	UpCount    uint64    `db:"up_count"`
	SideCount  uint64    `db:"side_count"`
	DownCount  uint64    `db:"down_count"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	BatchIndex int       `db:"batch_index"`
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

// POPULATED QUERIES

type populatedStockRecommendationQueryRow struct {
	ID              uint64        `db:"id"`
	BrokerageID     uint64        `db:"brokerage_id"`
	StockRegisterID uint64        `db:"stock_register_id"`
	RatingFrom      domain.Action `db:"rating_from"`
	RatingTo        domain.Action `db:"rating_to"`
	TargetFrom      float64       `db:"target_from"`
	TargetTo        float64       `db:"target_to"`
	CreatedAt       time.Time     `db:"created_at"`

	BrokerageName      string    `db:"brokerage_name"`
	BrokerageCreatedAt time.Time `db:"brokerage_created_at"`
}
