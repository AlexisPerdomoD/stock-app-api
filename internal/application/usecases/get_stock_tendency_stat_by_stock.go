package usecases

import (
	"context"
	"fmt"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type GetStockRegistersStatsByStock struct {
	stockRepository      domain.StockRepository
	stockStatsRepository domain.StockTendencyStatRepository
}

func (uc *GetStockRegistersStatsByStock) Execute(
	ctx context.Context,
	stockID uint64,
) (*models.StockTendencyStatView, error) {
	stock, err := uc.stockRepository.GetByID(ctx, stockID)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound(fmt.Sprintf("stock with id %d not found", stockID))
	}

	stat, err := uc.stockStatsRepository.GetByStockID(ctx, stockID)
	if err != nil {
		return nil, err
	}

	if stat == nil {
		return nil, pkg.InvalidStateErr(fmt.Sprintf("no stats found for stock with id %d", stockID))
	}

	response := models.NewStockTendencyStatView(*stat)
	return &response, nil
}

func NewGetStockRegistersStatsByStock(
	stockRepository domain.StockRepository,
	stockStatsRepository domain.StockTendencyStatRepository,
) *GetStockRegistersStatsByStock {
	if stockRepository == nil || stockStatsRepository == nil {
		panic("nil arguments for NewGetStockRegistersStatsByStock")
	}

	return &GetStockRegistersStatsByStock{
		stockRepository,
		stockStatsRepository,
	}
}
