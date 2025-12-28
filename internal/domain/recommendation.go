package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type Action uint8

const (
	Buy     Action = 1
	Hold    Action = 2
	Neutral Action = 3
	Sell    Action = 4
)

func (a Action) String() string {
	switch a {
	case Buy:
		return "buy"
	case Hold:
		return "hold"
	case Neutral:
		return "neutral"
	case Sell:
		return "sell"
	default:
		return "unknown"
	}
}

/*
Recommendation
Represents a recommendation made by a brokerage team.
*/
type Recommendation struct {
	ID              uint64    `json:"id,string"`
	StockRegisterID uint64    `json:"stock_register_id,string"`
	BrokerageID     uint64    `json:"brokerage_id,string"`
	RatingTo        Action    `json:"rating_to"`
	RatingFrom      Action    `json:"rating_from"`
	TargetTo        float64   `json:"target_to"`
	TargetFrom      float64   `json:"target_from"`
	CreatedAt       time.Time `json:"created_at"`
}

type PopulatedRecommendation struct {
	Recommendation
	Brokerage *Brokerage `json:"brokerage"`
}

type RecommendationRepository interface {
	GetAllPaginated(
		ctx context.Context,
		filter pkg.PaginationFilter,
		stockID uint64,
	) (*pkg.PaginationReponse[PopulatedRecommendation], error)
}
