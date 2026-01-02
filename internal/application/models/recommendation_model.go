package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type RecommendationView struct {
	ID          uint64    `json:"id,string" example:"1"`
	StockID     uint64    `json:"stock_id,string" example:"1"`
	BrokerageID uint64    `json:"brokerage_id,string" example:"1"`
	RatingTo    string    `json:"rating_to" example:"1"`
	RatingFrom  string    `json:"rating_from" example:"1"`
	TargetTo    float64   `json:"target_to" example:"1"`
	TargetFrom  float64   `json:"target_from" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2021-01-01T12:00:00Z"`
}

func NewRecommendationView(r domain.Recommendation) RecommendationView {
	return RecommendationView{
		ID:          r.ID,
		StockID:     r.StockRegisterID,
		BrokerageID: r.BrokerageID,
		RatingTo:    r.RatingTo.String(),
		RatingFrom:  r.RatingFrom.String(),
		TargetTo:    r.TargetTo,
		TargetFrom:  r.TargetFrom,
		CreatedAt:   r.CreatedAt,
	}
}

type PopulatedRecommendationView struct {
	RecommendationView
	Brokerage BrokerageView `json:"brokerage"`
}

func NewPopulatedRecommendationView(r domain.PopulatedRecommendation) PopulatedRecommendationView {
	return PopulatedRecommendationView{
		RecommendationView: NewRecommendationView(r.Recommendation),
		Brokerage:          NewBrokerageView(r.Brokerage),
	}
}
