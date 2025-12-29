package models

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

type BrokerageView struct {
	ID        uint64 `json:"id,string"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func NewBrokerageView(b domain.Brokerage) BrokerageView {
	return BrokerageView{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt.String(),
	}
}
