package entity

import "time"

type Tendency string

const (
	Up   Tendency = "up"
	Down Tendency = "down"
	Side Tendency = "side"
)

type Action string

const (
	Buy     Action = "buy"
	Sell    Action = "sell"
	Neutral Action = "neutral"
	Hold    Action = "hold"
	Unknown Action = "unknown"
)

/*
Stock
Represents a stock.
*/
type Stock struct {
	ID        int       `json:"id"`
	CompanyID int       `json:"company_id"`
	Ticker    string    `json:"ticker"`
	Name      string    `json:"name,omitempty"`
	Price     float64   `json:"price"`
	Tendency  Tendency  `json:"tendency"`
	CreatedAt time.Time `json:"registered"`
	UpdatedAt time.Time `json:"updated_at"`
}
