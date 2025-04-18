package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type Action uint8

const (
	Buy     Action = 1
	Hold    Action = 2
	Neutral Action = 3
	Sell    Action = 4
)

/*
Brokerage
Represents Analytics Brokerage teams that are responsible for the stocks recommendations.
*/
type Brokerage struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

/*
Recommendation
Represents a recommendation made by a brokerage team.
*/
type Recommendation struct {
	ID          uint      `json:"id"`
	StockID     uint      `json:"stock_id"`
	BrokerageID uint      `json:"brokerage_id"`
	RatingTo    Action    `json:"rating_to"`
	RatingFrom  Action    `json:"rating_from"`
	TargetTo    float64   `json:"target_to"`
	TargetFrom  float64   `json:"target_from"`
	CreatedAt   time.Time `json:"created_at"`
}

type PopulatedRecommendation struct {
	Recommendation
	BrokerageName string `json:"brokerage_name"`
}

type RecommendationRepository interface {
	Get(ctx context.Context, id uint) (*Recommendation, error)

	GetAllPaginated(
		ctx context.Context,
		filter pkg.PaginationFilter,
		stockID uint,
	) (*pkg.PaginationReponse[PopulatedRecommendation], error)
}
