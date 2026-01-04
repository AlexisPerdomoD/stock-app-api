package mock

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

var tickers []string = []string{
	"apple",
	"amaz",
	"tijua",
	"donn",
	"sony",
	"nint",
	"macd",
	"nike",
	"redb",
	"expo",
	"nasa",
	"guns",
}

func randomNumber(min, max float64) float64 {
	// #nosec G404 -- rand is fine in test/mock context
	return min + rand.Float64()*(max-min)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	b.Grow(length)
	for range length {
		// #nosec G404 -- rand is fine in test/mock context
		b.WriteByte(charset[rand.Intn(len(charset))])
	}
	return b.String()
}

func nextRating(
	prev domain.Action,
	tendency domain.Tendency,
	price float64,
	targetFrom float64,
	targetTo float64,
) domain.Action {

	// 1. Reglas duras (prioridad absoluta)
	if price >= targetTo {
		return domain.Sell
	}
	if price <= targetFrom {
		return domain.Buy
	}

	// 2. Inercia: probabilidad de mantener rating
	// #nosec G404 -- rand is fine in test/mock context
	r := rand.Float64()

	switch tendency {
	case domain.Up:
		if r < 0.6 {
			return prev
		}
		if r < 0.9 {
			return upgrade(prev)
		}
		return downgrade(prev)

	case domain.Down:
		if r < 0.6 {
			return prev
		}
		if r < 0.9 {
			return downgrade(prev)
		}
		return upgrade(prev)

	default: // Side
		if r < 0.7 {
			return prev
		}
		if r < 0.85 {
			return upgrade(prev)
		}
		return downgrade(prev)
	}
}

func upgrade(a domain.Action) domain.Action {
	switch a {
	case domain.Sell:
		return domain.Hold
	case domain.Hold:
		return domain.Buy
	default:
		return domain.Buy
	}
}

func downgrade(a domain.Action) domain.Action {
	switch a {
	case domain.Buy:
		return domain.Hold
	case domain.Hold:
		return domain.Sell
	default:
		return domain.Sell
	}
}

// mockSourceStockService es una implementación de SourceStockService para pruebas
type mockSourceStockService struct{}

func (m *mockSourceStockService) Name() string {
	return "MockSourceStockService"
}

func (m *mockSourceStockService) Get(ctx context.Context, limitDate *time.Time) ([]services.DataSourceResponse, error) {
	now := time.Now()
	result := make([]services.DataSourceResponse, 0, len(tickers)*6)

	for _, ticker := range tickers {

		var prevPrice float64
		var prevAction domain.Action
		// precio base por ticker (más realista)
		basePrice := randomNumber(50, 500)

		for i := range 6 {
			stockTime := now.Add(-time.Duration(i) * time.Minute * 5)

			if limitDate != nil && stockTime.After(*limitDate) {
				break
			}

			var price, targetFrom, targetTo float64
			var tendency domain.Tendency
			var ratingFrom, ratingTo domain.Action

			if i == 0 {
				price = basePrice
				tendency = domain.Side
				targetFrom = price
				targetTo = price * randomNumber(1.0, 1.10)
				ratingFrom = domain.Neutral

			} else {
				// variación pequeña respecto al precio anterior
				delta := randomNumber(-5, 5)
				price = prevPrice + delta
				targetFrom = prevPrice
				targetTo = price * randomNumber(1.0, 1.10)
				ratingFrom = prevAction

				switch {
				case price > prevPrice:
					tendency = domain.Up
				case price < prevPrice:
					tendency = domain.Down
				default:
					tendency = domain.Side
				}
			}

			ratingTo = nextRating(ratingFrom, tendency, price, targetFrom, targetTo)

			stock := services.DataSourceResponse{
				Market: services.MarketData{
					Name: "mock market",
				},
				Company: services.CompanyData{
					Name: fmt.Sprintf("mock-company-%s", ticker),
				},

				Recomendation: &services.RecommendationData{
					RatingFrom: ratingFrom,
					RatingTo:   ratingTo,
					TargetFrom: targetFrom,
					TargetTo:   targetTo,
					Brokerage: services.BrokerageData{
						Name: "mock-brokerage-" + randomString(8),
					},
				},
				Stock: services.StockRegisterData{
					Ticker:   ticker,
					Name:     ticker,
					Price:    price,
					Tendency: tendency,
				},
				Time: stockTime,
			}

			result = append(result, stock)
			prevPrice = price
			prevAction = ratingTo
		}
	}

	slices.SortFunc(result, func(a, b services.DataSourceResponse) int {
		return b.Time.Compare(a.Time)
	})

	return result, nil
}

func NewMockSourceStockService() services.DataSourceService {
	return &mockSourceStockService{}
}
