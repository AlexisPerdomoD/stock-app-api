package mock

import (
	"context"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"math/rand"
	"strings"
	"time"
)

func RandomNumber(min, max float64) float64 {
	// #nosec G404 -- rand is fine in test/mock context
	return min + rand.Float64()*(max-min)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	b.Grow(length)
	for range length {
		// #nosec G404 -- rand is fine in test/mock context
		b.WriteByte(charset[rand.Intn(len(charset))])
	}
	return b.String()
}

func RandomTendency() domain.Tendency {
	// #nosec G404 -- rand is fine in test/mock context
	return domain.Tendency(rand.Intn(3) + 1)
}

func RandomAction() domain.Action {
	// #nosec G404 -- rand is fine in test/mock context
	return domain.Action(rand.Intn(4) + 1)
}

func RandomTicker() string {
	tickers := []string{"apple", "amaz", "tijua", "donn", "sony", "nint", "macd", "nike", "redb", "expo", "nasa", "guns"}
	// #nosec G404 -- rand is fine in test/mock context
	return tickers[rand.Intn(len(tickers))]
}

// MockSourceStockService es una implementación de SourceStockService para pruebas
type MockSourceStockService struct{}

func (m *MockSourceStockService) Name() string {
	return "MockSourceStockService"
}

func (m *MockSourceStockService) Get(ctx context.Context, limitDate *time.Time) ([]services.DataSourceResponse, error) {
	var result []services.DataSourceResponse

	for i := range 500 {
		ticker := RandomTicker()
		stock := services.DataSourceResponse{
			Market: services.MarketData{
				Name: "mock market",
			},
			Company: services.CompanyData{
				Name: ticker,
			},
			Recomendation: &services.RecommendationData{
				RatingTo:   RandomAction(),
				RatingFrom: RandomAction(),
				TargetTo:   RandomNumber(10, 2000),
				TargetFrom: RandomNumber(10, 2000),
				Brokerage:  services.BrokerageData{Name: "mock " + RandomString(10)},
			},
			Stock: services.StockRegisterData{
				Ticker:   ticker,
				Name:     ticker,
				Price:    RandomNumber(10, 2000),
				Tendency: RandomTendency(),
			},
			Time: time.Now().Add(time.Duration(-i) * time.Hour * 6),
		}
		result = append(result, stock)
	}
	return result, nil
}
