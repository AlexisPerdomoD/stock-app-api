package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

/*
	PAYLOAD EXAMPLE

[

	{
	    "name": "Nasdaq Crypto Index",
	    "symbol": "NCI.IDX-NAS",
	    "current_price": 4483.12,
	    "prev_close_price": 4472.79,
	    "price_change_from_prev_close": 10.329999999999927,
	    "percent_change_from_prev_close": 0.002309520455912289,

"prev_close_date": "2025-04-25",

a	    "sort_order_index": 0,

	    "last_updated": "2025-04 -28T19:44:00.049000+00:00",
	    "event_timestamp": "2025-04-28T19:44:00.049000+00:00",
	    "pretty_symbol": "NCI.IDX"
	},
	{
	    "name": "Bitcoin",
	    "symbol": "BTCUSD-BITS",
	    "current_price": 94923.0,
	    "prev_close_price": 93753.0,
	    "price_change_from_prev_close": 1170.0,
	    "percent_change_from_prev_close": 0.01247960065277911,
	    "prev_close_date": "2025-04-27",
	    "sort_order_index": 1,
	    "last_updated": "2025-04-28T19:44:02.505000+00:00",
	    "event_timestamp": "2025-04-28T19:44:02.505000+00:00",
	    "pretty_symbol": "BTCUSD"
	},
	{

"name": "Ether",

	    "symbol": "ETHUSD-BITS",
	    "current_price": 1800.6,
	    "prev_close_price": 1793.7,
	    "price_change_from_prev_close": 6.899999999999864,
	    "percent_change_from_prev_close": 0.003846797123264684,
	    "prev_close_date": "2025-04-27",
	    "sort_order_index": 2,
	    "last_updated": "2025-04-28T19:43:56.039000+00:00",
	    "event_timestamp": "2025-04-28T19:43:56.039000+00:00",
	    "pretty_symbol": "ETHUSD"
	},
	{
	    "name": "Litecoin",
	    "symbol": "LTCUSD-BITS",
	    "current_price": 85.69,
	    "prev_close_price": 85.49,
	    "price_change_from_prev_close": 0.20000000000000284,
	    "percent_change_from_prev_close": 0.002339454907006701,
	    "prev_close_date": "2025-04-27",
	    "sort_order_index": 3,
	    "last_updated": "2025-04-28T19:42:42.579000+00:00",
	    "event_timestamp": "2025-04-28T19:42:42.579000+00:00",
	    "pretty_symbol": "LTCUSD"
	},
	{
	    "name": "XRP",
	    "symbol": "XRPUSD-BITS",
	    "current_price": 2.29442,
	    "prev_close_price": 2.25686,
	    "price_change_from_prev_close": 0.03756000000000004,
	"percent_change_from_prev_close": 0.016642591919746923,
	    "prev_close_date": "2025-04-27",
	    "sort_order_index": 4,
	"last_updated": "2025-04-28T19:43:56.758000+00:00",
	    "event_timestamp": "2025-04-28T19:43:56.758000+00:00",
	    "pretty_symbol": "XRPUSD"
	}

]
*/
type CnnStockSourceItem struct {
	CompanyName    string    `json:"name"`
	Ticker         string    `json:"symbol"`
	CurrentPrice   float64   `json:"current_price"`
	PrevClosePrice float64   `json:"prev_close_price"`
	LastUpdated    time.Time `json:"last_updated"`
}

type CnnStockSourceService struct {
	url string
	cl  *http.Client
}

func (s *CnnStockSourceService) Name() string {
	return "CnnStockSourceService"
}

func (s *CnnStockSourceService) Get(ctx context.Context, limitDate *time.Time) ([]domain.SourceStockData, error) {
	payload := []CnnStockSourceItem{}

	req, err := http.NewRequestWithContext(ctx, "GET", s.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	// SIMULATE BROWSER REQUEST
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:137.0) Gecko/20100101 Firefox/137.0")

	res, err := s.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("[cnn source] failed to close response body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, pkg.InternalServerError(fmt.Sprintf("unexpected status code %d", res.StatusCode))
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	data := []domain.SourceStockData{}
	log.Printf("[CnnStockSourceService]: service started and sorcing %d stocks", len(payload))
	for _, item := range payload {
		if limitDate != nil && limitDate.After(item.LastUpdated) {
			continue
		}

		tendency := domain.Side

		if item.PrevClosePrice > item.CurrentPrice {
			tendency = domain.Down
		} else if item.PrevClosePrice < item.CurrentPrice {
			tendency = domain.Up
		}

		dataItem := domain.SourceStockData{
			Time: item.LastUpdated,
			Market: domain.MarketArgs{
				Name: "cnn stock source",
			},
			Company: domain.CompanyArgs{
				Name: strings.ToLower(item.CompanyName),
			},
			Stock: domain.StockArgs{
				Ticker:   strings.ToLower(item.Ticker),
				Price:    item.CurrentPrice,
				Tendency: tendency,
			},
		}

		data = append(data, dataItem)
	}

	slices.SortFunc(data, func(a, b domain.SourceStockData) int {
		if a.Time.After(b.Time) {
			return 1
		}
		if a.Time.Before(b.Time) {
			return -1
		}
		return 0
	})
	log.Printf("[CnnStockSourceService]: sourced %d stocks", len(data))

	return data, nil
}

func NewCnnStockSourceService() *CnnStockSourceService {
	url := "https://production.dataviz.cnn.io/markets/crypto/summary"
	cl := &http.Client{
		Timeout: time.Second * 10,
	}

	return &CnnStockSourceService{url, cl}
}
