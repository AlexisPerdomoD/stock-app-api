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
	ID              uint64
	StockRegisterID uint64
	BrokerageID     uint64
	RatingTo        Action
	RatingFrom      Action
	TargetTo        float64
	TargetFrom      float64
	CreatedAt       time.Time
}

type PopulatedRecommendation struct {
	Recommendation
	Brokerage Brokerage
}

type RecommendationRepository interface {
	/*
	   returns paginated recommendations based on the specified filters
	*/
	GetAllPaginated(
		ctx context.Context,
		filter pkg.PaginationFilter,
		stockID uint64,
	) (*pkg.PaginationReponse[PopulatedRecommendation], error)

	/*
		Saves a list of recommendations in the repository.
		- returns error if nil arguments are passed
		- returns error if conflict occurs with arguments provided
	*/
	SaveAll(ctx context.Context, recommendations []*Recommendation) error
}
