package usecases

import (
	"context"
	"log"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type RegisterStocks struct {
	sr domain.StockRepository
}

func (uc *RegisterStocks) Execute(ctx context.Context, s domain.SourceStockService, limitDate *time.Time) (int, error) {
	if s == nil {
		return 0, pkg.InternalServerError("bad impl: SourceStockService was nil on registerStocksUseCase.Execute()")
	}

	data, err := s.Get(ctx, limitDate)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, nil
	}

	if err := uc.sr.Register(ctx, data); err != nil {
		return 0, err
	}

	return len(data), nil
}

func NewRegisterStocks(sr domain.StockRepository) *RegisterStocks {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository is nil when creating register stock use case")
	}

	return &RegisterStocks{sr}
}
