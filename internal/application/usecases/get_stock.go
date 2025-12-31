package usecases

import (
	"context"
	"fmt"

	appmodels "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type GetStock struct {
	userRepository          domain.UserRepository
	marketRepository        domain.MarketRepository
	companyRepository       domain.CompanyRepository
	stockRepository         domain.StockRepository
	stockRegisterRepository domain.StockRegisterRepository
}

func (uc *GetStock) Execute(ctx context.Context, stockID uint64, userID *uint64) (*appmodels.PopulatedStockView, error) {
	stock, err := uc.stockRepository.GetByID(ctx, stockID)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("stock does not exist")
	}

	var isSaved *bool

	if userID != nil {
		hasStock, err := uc.userRepository.HasUserStock(ctx, *userID, stockID)
		if err != nil {
			return nil, err
		}

		isSaved = &hasStock
	}

	market, err := uc.marketRepository.GetByID(ctx, stock.MarketID)
	if err != nil {
		return nil, err
	}
	if market == nil {
		return nil, pkg.InvalidStateErr("stocks must have always associate market")
	}

	company, err := uc.companyRepository.GetByID(ctx, stock.CompanyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, pkg.InvalidStateErr("stocks must have always associate company")
	}

	lastRegister, err := uc.stockRegisterRepository.GetLastByStockID(ctx, stock.ID)
	if err != nil {
		return nil, err
	}
	if lastRegister == nil {
		return nil, pkg.InvalidStateErr("stocks must have always atleast one register")
	}

	populatedStock := domain.PopulatedStock{
		Stock:        *stock,
		Company:      *company,
		Market:       *market,
		LastRegister: *lastRegister,
		IsSaved:      isSaved,
	}
	response := appmodels.NewPopulatedStockView(populatedStock)
	return &response, nil
}

/*
	concurrent version maybe overkill for this case
func (uc *GetStock) Execute(ctx context.Context, stockID uint64, userID *uint64) (*appmodels.PopulatedStockView, error) {
	stock, err := uc.stockRepository.GetByID(ctx, stockID)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, pkg.NotFound("stock does not exist")
	}

	ctx, cancel := context.WithCancel(ctx)

	var (
		wg           sync.WaitGroup
		mu           sync.Mutex
		market       *domain.Market
		company      *domain.Company
		lastRegister *domain.StockRegister
		isSaved      *bool
	)

	errhandler := func(e error) {
		if e == nil || err != nil {
			return
		}

		mu.Lock()
		err = e
		mu.Unlock()
		cancel()
	}

	if userID != nil {
		wg.Go(func() {
			hasStock, err := uc.userRepository.HasUserStock(ctx, *userID, stockID)
			if err != nil {
				errhandler(err)
				return
			}

			isSaved = &hasStock
		})
	}

	wg.Go(func() {
		mrkt, err := uc.marketRepository.GetByID(ctx, stock.MarketID)
		if err != nil {
			errhandler(err)
			return
		}

		market = mrkt
	})

	wg.Go(func() {
		cmpn, err := uc.companyRepository.GetByID(ctx, stock.CompanyID)
		if err != nil {
			errhandler(err)
			return
		}

		company = cmpn
	})

	wg.Go(func() {
		rgtr, err := uc.stockRegisterRepository.GetLastByStockID(ctx, stock.ID)
		if err != nil {
			errhandler(err)
			return
		}

		if rgtr == nil {
			errhandler(pkg.InvalidStateErr("stocks must have always atleast one register"))
			return
		}

		lastRegister = rgtr
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	populatedStock := domain.PopulatedStock{
		Stock:        *stock,
		Company:      *company,
		Market:       *market,
		LastRegister: *lastRegister,
		IsSaved:      isSaved,
	}
	response := appmodels.NewPopulatedStockView(populatedStock)
	return &response, nil
}

*/

func NewGetStock(
	ur domain.UserRepository,
	mr domain.MarketRepository,
	cr domain.CompanyRepository,
	sr domain.StockRepository,
	srr domain.StockRegisterRepository,
) *GetStock {

	if sr == nil || mr == nil || cr == nil || ur == nil || srr == nil {
		panic(
			fmt.Sprintf(
				"nil arguments provided for NewGetStocksUseCase userRepository=%v marketRepository=%v companyRepository=%v  stockRepository=%v stockRegisterRepository=%v",
				ur, mr, cr, sr, srr,
			),
		)
	}

	return &GetStock{ur, mr, cr, sr, srr}
}
