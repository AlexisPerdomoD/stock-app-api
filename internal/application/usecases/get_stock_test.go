package usecases_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

func TestGetStock_Execute_OK_WithUser(t *testing.T) {
	ctx := context.Background()

	userID := uint64(10)
	stockID := uint64(1)

	stock := &domain.Stock{
		ID:        stockID,
		MarketID:  2,
		CompanyID: 3,
	}

	market := &domain.Market{ID: 2}
	company := &domain.Company{ID: 3}
	register := &domain.StockRegister{
		StockID:  stockID,
		Price:    100,
		Tendency: domain.Up,
	}

	uc := usecases.NewGetStock(
		&FakeUserRepository{
			HasUserStockFn: func(ctx context.Context, uID, sID uint64) (bool, error) {
				return true, nil
			},
		},
		&FakeMarketRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Market, error) {
				return market, nil
			},
		},
		&FakeCompanyRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Company, error) {
				return company, nil
			},
		},
		&FakeStockRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Stock, error) {
				return stock, nil
			},
		},
		&FakeStockRegisterRepository{
			GetLastByStockIDFn: func(ctx context.Context, id uint64) (*domain.StockRegister, error) {
				return register, nil
			},
		},
	)

	result, err := uc.Execute(ctx, stockID, &userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.IsSaved == nil || *result.IsSaved != true {
		t.Fatalf("expected IsSaved=true, got %v", result.IsSaved)
	}

	if result.ID != stockID {
		t.Fatalf("unexpected stock id")
	}
}

func TestGetStock_Execute_OK_WithoutUser(t *testing.T) {
	ctx := context.Background()
	stockID := uint64(1)

	uc := usecases.NewGetStock(
		&FakeUserRepository{
			HasUserStockFn: func(ctx context.Context, uID, sID uint64) (bool, error) {
				t.Fatal("HasUserStock should not be called")
				return false, nil
			},
		},
		&FakeMarketRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Market, error) {
				return &domain.Market{ID: id}, nil
			},
		},
		&FakeCompanyRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Company, error) {
				return &domain.Company{ID: id}, nil
			},
		},
		&FakeStockRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Stock, error) {
				return &domain.Stock{ID: id, MarketID: 1, CompanyID: 2}, nil
			},
		},
		&FakeStockRegisterRepository{
			GetLastByStockIDFn: func(ctx context.Context, id uint64) (*domain.StockRegister, error) {
				return &domain.StockRegister{StockID: id}, nil
			},
		},
	)

	res, err := uc.Execute(ctx, stockID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsSaved != nil {
		t.Fatal("expected IsSaved to be nil")
	}
}

func TestGetStock_StockNotFound(t *testing.T) {
	uc := usecases.NewGetStock(
		&FakeUserRepository{},
		&FakeMarketRepository{},
		&FakeCompanyRepository{},
		&FakeStockRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Stock, error) {
				return nil, nil
			},
		},
		&FakeStockRegisterRepository{},
	)

	_, err := uc.Execute(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	AssertApiErr(t, err, http.StatusNotFound, "Not Found")
}

func TestGetStock_MarketNil_IsInvalidState(t *testing.T) {
	uc := usecases.NewGetStock(
		&FakeUserRepository{},
		&FakeMarketRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Market, error) {
				return nil, nil
			},
		},
		&FakeCompanyRepository{},
		&FakeStockRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Stock, error) {
				return &domain.Stock{ID: id, MarketID: 1, CompanyID: 2}, nil
			},
		},
		&FakeStockRegisterRepository{},
	)

	_, err := uc.Execute(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetStock_CompanyNil_IsInvalidState(t *testing.T) {
	uc := usecases.NewGetStock(
		&FakeUserRepository{},
		&FakeMarketRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Market, error) {
				return &domain.Market{ID: id}, nil
			},
		},
		&FakeCompanyRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Company, error) {
				return nil, nil
			},
		},
		&FakeStockRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Stock, error) {
				return &domain.Stock{ID: id, MarketID: 1, CompanyID: 2}, nil
			},
		},
		&FakeStockRegisterRepository{},
	)

	_, err := uc.Execute(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetStock_NoRegister_IsInvalidState(t *testing.T) {
	uc := usecases.NewGetStock(
		&FakeUserRepository{},
		&FakeMarketRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Market, error) {
				return &domain.Market{ID: id}, nil
			},
		},
		&FakeCompanyRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Company, error) {
				return &domain.Company{ID: id}, nil
			},
		},
		&FakeStockRepository{
			GetByIDFn: func(ctx context.Context, id uint64) (*domain.Stock, error) {
				return &domain.Stock{ID: id, MarketID: 1, CompanyID: 2}, nil
			},
		},
		&FakeStockRegisterRepository{
			GetLastByStockIDFn: func(ctx context.Context, id uint64) (*domain.StockRegister, error) {
				return nil, nil
			},
		},
	)

	_, err := uc.Execute(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
