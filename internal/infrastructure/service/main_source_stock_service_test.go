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
	mockTime := time.Now()

	mockResponse := service.MainStockSourcePayload{
		Items: []service.MainStockSourceItem{
			{
				Ticker:     "AAPL",
				TargetFrom: "$150.00",
				TargetTo:   "$160.00",
				Company:    "Apple Inc               ",
				Action:     "buy",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Neutral",
				RatingTo:   "buy",
				Time:       mockTime,
			},
		},
		NextPage: nil,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Set env vars
	os.Setenv("MAIN_SOURCE_STOCK_URI", server.URL)
	os.Setenv("MAIN_SOURCE_STOCK_KEY", "test-key")

	svc := service.NewMainSourceStockService()

	data, err := svc.Get(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "aapl", data[0].Stock.Ticker)
	assert.Equal(t, "apple inc", data[0].Company.Name)

	assert.Equal(t, domain.Buy, data[0].Recomendation.RatingTo)
	assert.Equal(t, domain.Neutral, data[0].Recomendation.RatingFrom)
	assert.Equal(t, domain.Up, data[0].Stock.Tendency)
	assert.Equal(t, float64(160), data[0].Stock.Price)
	assert.Equal(t, float64(160), data[0].Recomendation.TargetTo)
	assert.Equal(t, float64(150), data[0].Recomendation.TargetFrom)

	assert.True(t, mockResponse.Items[0].Time.Equal(data[0].Time))

	// add more assertions
}
