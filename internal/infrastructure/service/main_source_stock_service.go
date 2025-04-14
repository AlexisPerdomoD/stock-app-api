package service

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type mainSourceStockData struct{}

func (s *mainSourceStockData) Get(ctx context.Context, limitDate *time.Time) ([]domain.SourceStockData, error) {
	panic("Implement mainSourceStockData.Get()")
}

func NewMainSourceStockData() *mainSourceStockData {
	return &mainSourceStockData{}
}
