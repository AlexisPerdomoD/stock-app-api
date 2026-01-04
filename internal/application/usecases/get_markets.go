package usecases

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type GetMarkets struct {
	marketRepository domain.MarketRepository
}

func (uc *GetMarkets) Execute(ctx context.Context) ([]models.MarketView, error) {

	data, err := uc.marketRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	markets := make([]models.MarketView, 0, len(data))
	for _, market := range data {
		markets = append(markets, models.NewMarketView(market))
	}

	return markets, nil
}

func NewGetMarkets(mr domain.MarketRepository) *GetMarkets {
	return &GetMarkets{mr}
}
