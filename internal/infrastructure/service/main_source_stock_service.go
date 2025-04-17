package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type MainStockSourceItem struct {
	Ticker     string    `json:"ticker"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

type MainStockSourcePayload struct {
	Items    []MainStockSourceItem `json:"items"`
	NextPage *string               `json:"next_page,omitempty"`
}

type MainSourceStockService struct {
	name string
	cl   *http.Client
	uri  string
	key  string
	verbose bool
}

func (s *MainSourceStockService) Name() string {
	return s.name
}

func (s *MainSourceStockService) doRequest(ctx context.Context, nextPage string) (*MainStockSourcePayload, error) {
	payload := &MainStockSourcePayload{}

	req, err := http.NewRequestWithContext(ctx, "GET", s.uri+"?next_page="+nextPage, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.key))
	req.Header.Add("Accept", "application/json")

	res, err := s.cl.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(payload); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, pkg.InternalServerError(fmt.Sprintf("unexpected status code %d", res.StatusCode))
	}

	defer res.Body.Close()
	return payload, nil
}

func (s *MainSourceStockService) getPrice(price string) (float64, error) {
	price = strings.ReplaceAll(price, ",", "")
	price = strings.ReplaceAll(price, "$", "")
	return strconv.ParseFloat(price, 64)
}

func (s *MainSourceStockService) getTendency(current, previus float64) domain.Tendency {
	switch {
	case current > previus:
		return domain.Up
	case current < previus:
		return domain.Down
	default:
		return domain.Side
	}
}

func (s *MainSourceStockService) setRating(rating string) domain.Action {
	switch strings.ToLower(rating) {
	case "buy", "outperform", "overweight":
		return domain.Buy
	case "underweight", "sell":
		return domain.Sell
	default:
		return domain.Neutral
	}
}
func (s *MainSourceStockService) Get(ctx context.Context, limitDate *time.Time) ([]domain.SourceStockData, error) {

	doUntil := limitDate
	response := []domain.SourceStockData{}
	nextPage := ""

	if doUntil != nil && doUntil.After(time.Now()) {
		yesterday := time.Now().AddDate(0, 0, -1)
		doUntil = &yesterday
	}

	if doUntil == nil && s.verbose {
		log.Printf("[main source stock] no limit date provided")
	}

	for {
		if s.verbose {
			log.Printf("[main source stock]: start getting stocks from main source stock from page %s", nextPage)
		}
		payload, err := s.doRequest(ctx, nextPage)
		if err != nil {
			return nil, err
		}

		mkt := domain.MarketArgs{Name: "main source stock"}

		for i, item := range payload.Items {

			/* Stop if the date is before the limit date */
			if doUntil != nil && item.Time.Before(*doUntil) {
				break
			}

			companyName := strings.ToLower(strings.TrimSpace(item.Company))
			if companyName == "" {
				return nil, pkg.InternalServerError(fmt.Sprintf("company is empty in payload at index %d", i))
			}

			brokerageName := strings.ToLower(strings.TrimSpace(item.Brokerage))
			if brokerageName == "" {
				return nil, pkg.InternalServerError(fmt.Sprintf("brokerage is empty in payload at index %d", i))
			}

			ticker := strings.ToLower(strings.TrimSpace(item.Ticker))
			if ticker == "" {
				return nil, pkg.InternalServerError(fmt.Sprintf("ticker is empty in payload at index %d", i))
			}

			previusPrice, err := s.getPrice(item.TargetFrom)
			if err != nil {
				return nil, err
			}
			currentPrice, err := s.getPrice(item.TargetTo)
			if err != nil {
				return nil, err
			}

			tendency := s.getTendency(currentPrice, previusPrice)
			ratingFrom := s.setRating(item.RatingFrom)
			ratingTo := s.setRating(item.RatingTo)

			args := domain.SourceStockData{Time: item.Time}

			args.Market = mkt

			args.Company = domain.CompanyArgs{Name: companyName}

			args.Recomendation = &domain.RecommendationArgs{
				RatingTo:   ratingTo,
				RatingFrom: ratingFrom,
				TargetTo:   currentPrice,
				TargetFrom: previusPrice,
				Brokerage: domain.BrokerageArgs{
					Name: brokerageName,
				},
			}

			args.Stock = domain.StockArgs{
				Ticker:   ticker,
				Price:    currentPrice,
				Tendency: tendency,
			}
			response = append(response, args)
		}

		if payload.NextPage == nil || *payload.NextPage == "" {
			if s.verbose {
				log.Printf("[main source stock] no next page")
				log.Printf("[main source stock] got %v stocks total", len(response))
			}

			break
		}

		nextPage = *payload.NextPage

		if s.verbose {
			log.Printf("[main source stock] next page: %s", nextPage)
		}
	}

	slices.SortFunc(response, func(a, b domain.SourceStockData) int {
		if a.Time.After(b.Time) {
			return 1
		}
		if a.Time.Before(b.Time) {
			return -1
		}
		return 0
	})

	return response, nil
}

func NewMainSourceStockService(verbose bool) *MainSourceStockService {
	uri := os.Getenv("MAIN_SOURCE_STOCK_URI")
	key := os.Getenv("MAIN_SOURCE_STOCK_KEY")

	if uri == "" || key == "" {
		log.Fatalln("please set MAIN_SOURCE_STOCK_URI and MAIN_SOURCE_STOCK_KEY environment variables")
	}

	name := "main source stock"

	cl := &http.Client{
		Timeout: time.Second * 10,
	}

	return &MainSourceStockService{name, cl, uri, key, verbose}
}
