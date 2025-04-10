package entity

import "time"

/*
Market
Represents the market where the stock is traded.
*/
type Market struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
