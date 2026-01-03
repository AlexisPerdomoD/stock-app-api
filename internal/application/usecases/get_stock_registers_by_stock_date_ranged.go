package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type GetStockRegistersByStockDateRanged struct {
	stockRegisterRepository domain.StockRegisterRepository
	limitDateRange          time.Duration
}

func (uc *GetStockRegistersByStockDateRanged) Execute(
	ctx context.Context,
	stockID uint64,
	from, to time.Time,
) ([]models.StockRegisterView, error) {
	if from.After(to) || from.Equal(to) {
		return nil, pkg.BadRequest("invalid date range, from must be before to")
	}

	if to.Sub(from) > uc.limitDateRange {
		message := fmt.Sprintf(
			"invalid date range, from %s to %s is too big, max is %s",
			from,
			to,
			uc.limitDateRange,
		)
		return nil, pkg.BadRequest(message)
	}

	stockRegisters, err := uc.stockRegisterRepository.GetRangeByStockID(ctx, stockID, from, to)
	if err != nil {
		return nil, err
	}

	response := make([]models.StockRegisterView, 0, len(stockRegisters))
	for _, stockRegister := range stockRegisters {
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

func NewGetStockRegistersByStockDateRanged(
	stockRegisterRepository domain.StockRegisterRepository,
	limitDateRange time.Duration,
) *GetStockRegistersByStockDateRanged {
	if stockRegisterRepository == nil {
		panic("nil argument for NewGetStockRegistersByStockDateRanged")
	}

	return &GetStockRegistersByStockDateRanged{
		stockRegisterRepository,
		limitDateRange,
	}
}
