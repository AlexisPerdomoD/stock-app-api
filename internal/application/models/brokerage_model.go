package models

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

type BrokerageView struct {
	ID        uint64 `json:"id,string" example:"1"`
	Name      string `json:"name" example:"Brokerage 1"`
	CreatedAt string `json:"created_at" example:"2021-01-01T12:00:00Z"`
}

func NewBrokerageView(b domain.Brokerage) BrokerageView {
	return BrokerageView{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt.String(),
	}
}
