package cockroachdb

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"time"
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

func (r marketRecord) ToDomain() *domain.Market {
	return &domain.Market{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

type companyRecord struct {
	ID        uint64    `db:"id"`
	MarketID  uint64    `db:"market_id"`
	Name      string    `db:"name"`
	ISIN      *string   `db:"isin"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r companyRecord) ToDomain() *domain.Company {
	return &domain.Company{
		ID:        r.ID,
		MarketID:  r.MarketID,
		Name:      r.Name,
		ISIN:      r.ISIN,
		CreatedAt: r.CreatedAt,
	}
}

type brokerageRecord struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r brokerageRecord) ToDomain() *domain.Brokerage {
	return &domain.Brokerage{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

type stockRecord struct {
	ID        uint64          `db:"id"`
	CompanyID uint64          `db:"company_id"`
	Name      *string         `db:"name"`
	Ticker    string          `db:"ticker"`
	Price     float64         `db:"price"`
	Tendency  domain.Tendency `db:"tendency"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (r stockRecord) ToDomain() *domain.Stock {
	return &domain.Stock{
		ID:        r.ID,
		Name:      r.Name,
		CompanyID: r.CompanyID,
		Ticker:    r.Ticker,
		Price:     r.Price,
		Tendency:  r.Tendency,
		CreatedAt: r.CreatedAt,
	}
}
