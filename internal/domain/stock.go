package domain

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

/*
Stock
Represents a  stock catalog.
*/
type Stock struct {
	ID        uint64
	CompanyID uint64
	MarketID  uint64
	Ticker    string
	Isin      *string
	Name      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StockUpdates struct {
	Name *string
}

type PopulatedStock struct {
	Stock
	Company      Company
	Market       Market
	LastRegister StockRegister
	IsSaved      *bool
}

type StockCompanySearchParam struct {
	StockTicker string
	CompanyID   uint64
	MarketID    uint64
}

type StockRepository interface {
	/*
		Returns a Stock by its ID. If the ID does not exist, returns nil.

		If userID is not nil, it will return the stock only if the user is the owner of the stock.
	*/
	Get(ctx context.Context, stockID uint64, userID *uint64) (*Stock, error)

	/*
		Returns a map of Stocks by their tickers. If a ticker does not exist, sets the value to nil.
		- nil values returns an error when passed as arguments.
	*/
	GetByStockCompanySearchParams(ctx context.Context, tickers []StockCompanySearchParam) (map[StockCompanySearchParam]*Stock, error)

	/*
		Returns populated stocks by provided stock ID, nil if stock does not exist.
		If userID is not nil, it will return the stock only if the user is the owner of the stock.
	*/
	GetPopulated(ctx context.Context, stockID uint64, userID *uint64) (*PopulatedStock, error)

	/*
		Returns a list of stocks by provided filter.
		If userID is not nil, it will return the stock only if the user has saved the stocks.
	*/
	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter, userID *uint64) (*pkg.PaginationReponse[PopulatedStock], error)

	/*
		Saves a Stock in the database and map missing properties with their default values (if any) including ID.

		- nil values returns an error.
		- invalid constraints returns an error.
		- duplicated ID returns an error.
	*/
	Save(ctx context.Context, stock *Stock) error

	/*
		Saves all Stocks in the database and map missing properties with their default values (if any) including ID.

		- nil values returns an error.
		- invalid constraints returns an error.
		- duplicated ID returns an error.
	*/
	SaveAll(ctx context.Context, stocks []*Stock) error

	/*
		Updates a Stock in the database.
		- if stock does not exist, returns an error.
		- nil values returns an error.
		- invalid constraints returns an error.
	*/

	Update(ctx context.Context, stock StockUpdates) error
}
