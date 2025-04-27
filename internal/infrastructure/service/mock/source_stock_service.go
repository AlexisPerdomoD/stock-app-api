package mock

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

func RandomNumber(min, max float64) float64 {
	// #nosec G404 -- rand is fine in test/mock context
	return min + rand.Float64()*(max-min)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	b.Grow(length)
	for i := 0; i < length; i++ {
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

func (m *MockSourceStockService) Get(ctx context.Context, limitDate *time.Time) ([]domain.SourceStockData, error) {
	var result []domain.SourceStockData

	for i := 0; i < 500; i++ {
		ticker := RandomTicker()
		stock := domain.SourceStockData{
			Market: domain.MarketArgs{
				Name: "mock-market",
			},
			Company: domain.CompanyArgs{
				Name: ticker,
				ISIN: nil,
			},
			Recomendation: &domain.RecommendationArgs{
				RatingTo:   RandomAction(),
				RatingFrom: RandomAction(),
				TargetTo:   RandomNumber(10, 2000),
				TargetFrom: RandomNumber(10, 2000),
				Brokerage:  domain.BrokerageArgs{Name: RandomString(10)},
			},
			Stock: domain.StockArgs{
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
