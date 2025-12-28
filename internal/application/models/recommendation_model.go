package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type RecommendationView struct {
	ID          uint64    `json:"id,string"`
	StockID     uint64    `json:"stock_id,string"`
	BrokerageID uint64    `json:"brokerage_id,string"`
	RatingTo    string    `json:"rating_to"`
	RatingFrom  string    `json:"rating_from"`
	TargetTo    float64   `json:"target_to"`
	TargetFrom  float64   `json:"target_from"`
	CreatedAt   time.Time `json:"created_at"`
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
