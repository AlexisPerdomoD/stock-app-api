package usecases

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type GetLastStockRegistersByStock struct {
	stockRegisterRepository domain.StockRegisterRepository
}

func (uc *GetLastStockRegistersByStock) Execute(
	ctx context.Context,
	stockID uint64,
	limit uint16,
) ([]models.StockRegisterView, error) {

	data, err := uc.stockRegisterRepository.GetLastsByStockID(ctx, stockID, limit)
	if err != nil {
		return nil, err
	}

	response := make([]models.StockRegisterView, 0, len(data))
	for _, stockRegister := range data {
		response = append(response, models.StockRegisterView{
			ID:        stockRegister.ID,
			StockID:   stockRegister.StockID,
			Price:     stockRegister.Price,
			Tendency:  stockRegister.Tendency.String(),
			CreatedAt: stockRegister.CreatedAt,
		})
	}

	return response, nil

}

func NewGetLastStockRegistersByStock(stockRegisterRepository domain.StockRegisterRepository) *GetLastStockRegistersByStock {
	if stockRegisterRepository == nil {
		panic("nil arguments for NewGetLastStockRegistersByStock")
	}

	return &GetLastStockRegistersByStock{stockRegisterRepository}
}
