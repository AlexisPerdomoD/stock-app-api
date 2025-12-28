package usecases

import (
	"context"
	"log"

	appmodels "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type GetStock struct {
	sr domain.StockRepository
}

func (uc *GetStock) Execute(ctx context.Context, stockID uint64, userID *uint64) (*appmodels.PopulatedStockView, error) {
	stock, err := uc.sr.GetPopulated(ctx, stockID, userID)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("stock does not exist")
	}

	response := appmodels.NewPopulatedStockView(*stock)
	return &response, nil
}

func NewGetStock(sr domain.StockRepository) *GetStock {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository was nil for NewGetStocksUseCase")
	}

	return &GetStock{sr}
}
