package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type GetStockUseCase struct {
	sr domain.StockRepository
}

func (uc *GetStockUseCase) Execute(ctx context.Context, stockID uint) (*domain.PopulatedStock, error) {
	stock, err := uc.sr.Get(ctx, stockID)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("stock does not exist")
	}

	return stock, nil
}

func NewGetStockUseCase(sr domain.StockRepository) *GetStockUseCase {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetStocksUseCase")
	}

	return &GetStockUseCase{sr}
}
