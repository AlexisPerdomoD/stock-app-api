package usecase

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/repository"

type registerStocksUseCase struct {
	stockRepository     repository.StockRepository
	companyRepository   repository.CompanyRepository
	marketRepository    repository.MarketRepository
	brokerageRepository repository.BrokerageRepository
}

func (uc *registerStocksUseCase) Execute() {
	panic("registerStocksUseCase.Execute() not implemented")
}

func NewRegisterStocksUseCase(
	ur repository.StockRepository,
	cr repository.CompanyRepository,
	mr repository.MarketRepository,
	br repository.BrokerageRepository,
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
