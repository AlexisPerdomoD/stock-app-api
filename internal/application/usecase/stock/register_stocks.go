package usecase

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

type registerStocksUseCase struct {
	stockRepository     domain.StockRepository
	companyRepository   domain.CompanyRepository
	marketRepository    domain.MarketRepository
	brokerageRepository domain.BrokerageRepository
}

func (uc *registerStocksUseCase) Execute() {
	panic("registerStocksUseCase.Execute() not implemented")
}

func NewRegisterStocksUseCase(
	ur domain.StockRepository,
	cr domain.CompanyRepository,
	mr domain.MarketRepository,
	br domain.BrokerageRepository,
) *registerStocksUseCase {

	if ur == nil {
		panic("stock repository is nil, stopping :b")
	}

	if cr == nil {
		panic("company repository is nil, stopping :b")
	}

	if mr == nil {
		panic("market repository is nil, stopping :b")
	}

	if br == nil {
		panic("brokerage repository is nil, stopping :b")
	}

	return &registerStocksUseCase{
		stockRepository:     ur,
		companyRepository:   cr,
		marketRepository:    mr,
		brokerageRepository: br,
	}
}
