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
	Isin *string
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

type FilterByStock string

const (
	FilterByStockPrice FilterByStock = "filter_stock_price"
)

func (f FilterByStock) String() string {
	return string(f)
}

func (f FilterByStock) IsValid() bool {
	return f == FilterByStockPrice
}

type SortByStock string

const (
	SortByStockPrice    SortByStock = "sort_stock_price"
	SortByStockTendency SortByStock = "sort_stock_tendency"
	SortByStockTicker   SortByStock = "sort_stock_ticker"
	SortByStockDate     SortByStock = "sort_stock_date"
)

func (s SortByStock) String() string {
	return string(s)
}

func (s SortByStock) IsValid() bool {
	return s == SortByStockPrice ||
		s == SortByStockTendency ||
		s == SortByStockTicker ||
		s == SortByStockDate
}

type StockRepository interface {
	/*
		Returns a Stock by its ID. If the ID does not exist, returns nil.

	*/
	GetByID(ctx context.Context, stockID uint64) (*Stock, error)

	/*
		Returns a map of Stocks by their tickers. If a ticker does not exist, sets the value to nil.
		- nil values returns an error when passed as arguments.
	*/
	GetByStockCompanySearchParams(ctx context.Context, tickers []StockCompanySearchParam) (map[StockCompanySearchParam]*Stock, error)

	/*
		Returns a list of stocks by provided filter.

		Allows filters are:

		- FilterByStockPrice (float64) is expected.

		- filter.search (string) filter by ticker, name or company name (case insensitive).

		Allows sorting are:

		- SortByStockPrice

		- SortByStockTendency

		- SortByStockTicker

		- SortByStockDate

		Any other filters (including not valid values) or sorting  will be ignored.
	*/
	GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[PopulatedStock], error)

	/*
		Returns a list of stocks assosiate with the userID by provided filter.

		Allows filters are:

		- FilterByStockPrice (float64) is expected.

		- filter.search (string) filter by ticker, name or company name (case insensitive).

		Allows sorting are:

		- SortByStockPrice

		- SortByStockTendency

		- SortByStockTicker

		- SortByStockDate

		Any other filters (including not valid values) or sorting  will be ignored.
	*/
	GetAllPaginatedByUser(ctx context.Context, filter pkg.PaginationFilter, userID uint64) (*pkg.PaginationReponse[PopulatedStock], error)

	/*
		Saves a Stock in the database and map missing properties with their default values (if any) including ID.

		- nil stock argument is no-op and returns nil.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	Save(ctx context.Context, stock *Stock) error

	/*
		Saves all Stocks in the database and map missing properties with their default values (if any) including ID.

		- nil stock slice is no-op and returns nil.

		- nil elements inside stocks slice returns err.

		- invalid constraints returns an err.

		- duplicated (MarketID, ticket) inside slice arguments  returns an err.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes).
	*/
	SaveAll(ctx context.Context, stocks []*Stock) error

	/*
		Updates a Stock in the database.
		- if stock does not exist, returns an error.

		- invalid constraints returns an error.
	*/

	Update(ctx context.Context, stock StockUpdates) error
}
