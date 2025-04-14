package usecase

import (
	"context"
	"log"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type RegisterStocksUseCase struct {
	sr domain.StockRepository
	rr domain.RecommendationRepository
	cr domain.CompanyRepository
	mr domain.MarketRepository
	br domain.BrokerageRepository
}

func (uc *RegisterStocksUseCase) Execute(ctx context.Context, s domain.SourceStockService, limitDate *time.Time) error {
	if s == nil {
		return pkg.InternalServerError("bad impl: SourceStockService was nil on registerStocksUseCase.Execute()")
	}

	panic("registerStocksUseCase.Execute() not implemented")
}

func NewRegisterStocksUseCase(
	sr domain.StockRepository,
	cr domain.CompanyRepository,
	mr domain.MarketRepository,
	br domain.BrokerageRepository,
	rr domain.RecommendationRepository,
) *RegisterStocksUseCase {

	if sr == nil {
		log.Fatalln("bad impl: StockRepository is nil when creating register stock use case")
	}

	if cr == nil {
		log.Fatalln("bad impl: CompanyRepository is nil when creating register stock use case")
	}

	if mr == nil {
		log.Fatalln("bad impl: MarketRepository is nil when creating register stock use case")
	}

	if br == nil {
		log.Fatalln("bad impl: BrokerageRepository is nil when creating register stock use case")
	}

	if rr == nil {
		log.Fatalln("bad impl:RecommendationRepository is nil when creating register stock use case")
	}

	return &RegisterStocksUseCase{
		sr: sr,
		cr: cr,
		mr: mr,
		br: br,
		rr: rr,
	}
}
