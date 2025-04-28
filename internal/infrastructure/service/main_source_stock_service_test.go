package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service"
	"github.com/stretchr/testify/assert"
)

func TestMainSourceStockService_Get(t *testing.T) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)
	threeDaysAgo := today.AddDate(0, 0, -3)

	mockResponse := service.MainStockSourcePayload{
		Items: []service.MainStockSourceItem{
			{
				Ticker:     "AAPL",
				TargetFrom: "$100.00",
				TargetTo:   "$120.00",
				Company:    "Apple Inc               ",
				Action:     "buy",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Neutral",
				RatingTo:   "buy",
				Time:       today,
			},
			{
				Ticker:     "BBLC",
				TargetFrom: "$160.00",
				TargetTo:   "$100.00",
				Company:    "Mac donalds               ",
				Action:     "sell",
				Brokerage:  "Goku sasn",
				RatingFrom: "Neutral",
				RatingTo:   "underweight",
				Time:       yesterday,
			}, {
				Ticker:     "AAPL",
				TargetFrom: "$150.00",
				TargetTo:   "$160.00",
				Company:    "Apple Inc               ",
				Action:     "buy",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Neutral",
				RatingTo:   "buy",
				Time:       threeDaysAgo,
			},
		},
		NextPage: nil,
	}

	serverHappyPath := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer serverHappyPath.Close()

	// Set env vars
	os.Setenv("MAIN_SOURCE_STOCK_URI", serverHappyPath.URL)
	os.Setenv("MAIN_SOURCE_STOCK_KEY", "test-key")

	svc := service.NewMainSourceStockService(false)

	data, err := svc.Get(context.Background(), &twoDaysAgo)
	assert.NoError(t, err)

	assert.Len(t, data, 2)

	assert.Equal(t, "bblc", data[0].Stock.Ticker)
	assert.Equal(t, "mac donalds", data[0].Company.Name)
	assert.Equal(t, domain.Sell, data[0].Recomendation.RatingTo)
	assert.Equal(t, domain.Neutral, data[0].Recomendation.RatingFrom)
	assert.Equal(t, domain.Down, data[0].Stock.Tendency)
	assert.Equal(t, float64(100), data[0].Stock.Price)
	assert.Equal(t, float64(100), data[0].Recomendation.TargetTo)
	assert.Equal(t, float64(160), data[0].Recomendation.TargetFrom)

	assert.Equal(t, "aapl", data[1].Stock.Ticker)
	assert.Equal(t, "apple inc", data[1].Company.Name)
	assert.Equal(t, domain.Buy, data[1].Recomendation.RatingTo)
	assert.Equal(t, domain.Neutral, data[1].Recomendation.RatingFrom)
	assert.Equal(t, domain.Up, data[1].Stock.Tendency)
	assert.Equal(t, float64(120), data[1].Stock.Price)
	assert.Equal(t, float64(120), data[1].Recomendation.TargetTo)
	assert.Equal(t, float64(100), data[1].Recomendation.TargetFrom)

	assert.True(t, yesterday.Equal(data[0].Time))
	assert.True(t, today.Equal(data[1].Time))

	t.Log("Test: Validar error al obtener datos de la API")

	serverErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
	}))
	defer serverErr.Close()

	os.Setenv("MAIN_SOURCE_STOCK_URI", serverErr.URL)
	os.Setenv("MAIN_SOURCE_STOCK_KEY", "test-key")
	svc = service.NewMainSourceStockService(false)
	data, err = svc.Get(context.Background(), &twoDaysAgo)
	assert.Nil(t, data)
	assert.Error(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	data, err = svc.Get(ctx, &twoDaysAgo)

	assert.Nil(t, data)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &context.DeadlineExceeded)
}
