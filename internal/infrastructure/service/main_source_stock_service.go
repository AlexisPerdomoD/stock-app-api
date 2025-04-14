package service

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type MainSourceStockService struct {
	name string
}

func (s *MainSourceStockService) Get(ctx context.Context, limitDate *time.Time) ([]domain.SourceStockData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		panic("Implement mainSourceStockData.Get()")
	}
}

func (s *MainSourceStockService) Name() string {
	return s.name
}

func NewMainSourceStockService() *MainSourceStockService {
	return &MainSourceStockService{name: "main"}
}
