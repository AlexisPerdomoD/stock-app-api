package services

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

/*
Represents a response from a data source.
*/
type MarketData struct {
	Name string
}

/*
Represents a company data from a data source.
*/
type CompanyData struct {
	Name string
	ISIN *string
}

/*
Represents a brokerage data from a data source.
*/
type BrokerageData struct {
	Name string
}

/*
Represents a recommendation data from a data source.
*/
type RecommendationData struct {
	RatingTo   domain.Action
	RatingFrom domain.Action
	TargetTo   float64
	TargetFrom float64
	Brokerage  BrokerageData
}

/*
Represents a stock data from a data source.
*/
type StockData struct {
	Ticker   string
	Name     string
	Price    float64
	Tendency domain.Tendency
}

/*
Represents expected response from a data source.
*/
type DataSourceResponse struct {
	Market        MarketData
	Company       CompanyData
	Recomendation *RecommendationData
	Stock         StockData
	Time          time.Time
}

/*
DataSourceService
Service that provides data from different sources of stocks, companies, brokerages, etc.
*/
type DataSourceService interface {
	/*
		Name returns the name of the data source.
	*/
	Name() string

	/*
		Get returns a list of stocks from different sources to an specific limit date.
		If limitDate is nil, should provide last day stocks updates by default.
	*/
	Get(ctx context.Context, limitDate *time.Time) ([]DataSourceResponse, error)
}

type DataSourceServiceProvider interface {
	/*
		GetDataSourceService returns a DataSourceService by name.
	*/
	GetDataSourceService(name string) DataSourceService
	/*
		Exists returns true if a DataSourceService is registered by name.
	*/
	Exists(name string) bool

	/*
		RegisterDataSourceService registers a DataSourceService by name.
	*/
	RegisterDataSourceService(name string, ds DataSourceService)
}
