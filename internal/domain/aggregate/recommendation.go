package aggregate

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"

type PopulatedRecommendation struct {
	entity.Recommendation
	BrokerageName string `json:"brokerage_name"`
}
