package usecases

import (
	"context"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"log"
)

type GetStock struct {
	sr domain.StockRepository
}

func (uc *GetStock) Execute(ctx context.Context, stockID uint64, userID *uint64) (*domain.PopulatedStock, error) {
	stock, err := uc.sr.Get(ctx, stockID, userID)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("stock does not exist")
	}

	return stock, nil
}

func NewGetStock(sr domain.StockRepository) *GetStock {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetStocksUseCase")
	}

	return &GetStock{sr}
}
