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

type FilterByRecommendation string

const (
	FilterByRecommendationStockID FilterByRecommendation = "filter_recomendation_stock_id"
)

func (f FilterByRecommendation) String() string {
	return string(f)
}

func (f FilterByRecommendation) IsValid() bool {
	return f == FilterByRecommendationStockID
}

type SortByRecommendation string

const (
	SortByRecommendationCreatedAt             SortByRecommendation = "sort_recommendation_created_at"
	SortByRecommendationStockRegisterTendency SortByRecommendation = "sort_recommendation_stock_register_tendency"
)

func (s SortByRecommendation) String() string {
	return string(s)
}

func (s SortByRecommendation) IsValid() bool {
	return s == SortByRecommendationCreatedAt || s == SortByRecommendationStockRegisterTendency
}

type RecommendationRepository interface {
	/*
		returns paginated recommendations based on the specified filters

		Allows filters are:
		- FilterByRecommendationStockID (uint64) is expected.

		- filter.search (string) filter by brokerage name (case insensitive).

		Allows sorting are:
		- SortByRecommendationCreatedAt

		- SortByRecommendationStockRegisterTendency

		Any other filters (including not valid values) or sorting  will be ignored.
	*/
	GetAllPaginated(
		ctx context.Context,
		filter pkg.PaginationFilter,
	) (*pkg.PaginationResponse[PopulatedRecommendation], error)

	/*
		Saves a list of recommendations in the repository.

		- nil recommendations slice is no-op and returns nil.

		- returns error if nil arguments are passed inside recommendations slice.

		- returns error if conflict occurs with arguments provided.
	*/
	SaveAll(ctx context.Context, recommendations []*Recommendation) error
}
