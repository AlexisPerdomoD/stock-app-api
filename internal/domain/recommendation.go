package domain

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

/*
Recommendation
Represents a recommendation made by a brokerage team.
*/
type Recommendation struct {
	ID          int       `json:"id"`
	StockID     int       `json:"stock_id"`
	BrokerageID int       `json:"brokerage_id"`
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
	Get(id int) (*Recommendation, error)

	GetAllPaginated(filter pkg.PaginationFilter) (*pkg.PaginationReponse[PopulatedRecommendation], error)

	Create(args *Recommendation) error
}
