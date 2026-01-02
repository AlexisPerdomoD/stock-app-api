package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

func TestRegisterStocks_Execute_OK(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	data := []services.DataSourceResponse{
		{
			Time: now,
			Market: services.MarketData{
				Name: "NASDAQ",
			},
			Company: services.CompanyData{
				Name: "Apple",
			},
			Stock: services.StockRegisterData{
				Ticker:   "AAPL",
				Price:    100,
				Tendency: domain.Up,
			},
		},
	}

	uow := NewHappyPathUnitOfWork()

	uow.MarketRepo = &FakeMarketRepository{
		GetByNamesFn: func(ctx context.Context, names []string) (map[string]*domain.Market, error) {
			return map[string]*domain.Market{
				"NASDAQ": &domain.Market{Name: "NASDAQ", ID: 1},
			}, nil
		},
	}

	uow.CompanyRepo = &FakeCompanyRepository{
		GetByMarketCompanySearchFn: func(
			ctx context.Context,
			params []domain.MarketCompanySearchParam,
		) (map[domain.MarketCompanySearchParam]*domain.Company, error) {
			return map[domain.MarketCompanySearchParam]*domain.Company{
				{MarketID: 1, Name: "Apple"}: &domain.Company{Name: "Apple", MarketID: 1, ID: 1},
			}, nil
		},
	}

	uow.StockRepo = &FakeStockRepository{
		GetByStockCompanySearchParamsFn: func(
			ctx context.Context,
			params []domain.StockCompanySearchParam,
		) (map[domain.StockCompanySearchParam]*domain.Stock, error) {
			return map[domain.StockCompanySearchParam]*domain.Stock{
				{StockTicker: "AAPL", CompanyID: 1, MarketID: 1}: &domain.Stock{Ticker: "AAPL", CompanyID: 1, MarketID: 1, ID: 1},
			}, nil
		},
	}

	uow.StockRegisterRepo = &FakeStockRegisterRepository{
		SaveAllFn: func(ctx context.Context, regs []*domain.StockRegister) error {
			// simula IDs asignados
			for i := range regs {
				// #nosec
				regs[i].ID = uint64(i + 1)
			}
			return nil
		},
	}

	uow.StockTendencyRepo = &FakeStockTendencyStatRepository{
		IncrementAllFn: func(ctx context.Context, deltas map[uint64]domain.StockTendencyDelta) error {
			return nil
		},
	}

	uowFactory := &FakeUnitOfWorkFactory{
		UOW: uow,
	}

	ds := map[string]services.DataSourceService{
		"test": &FakeDataSource{
			GetFn: func(ctx context.Context, limit *time.Time) ([]services.DataSourceResponse, error) {
				return data, nil
			},
		},
	}

	uc := usecases.NewRegisterStocks(uowFactory, ds)

	count, err := uc.Execute(ctx, "test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 inserted register, got %d", count)
	}
}
