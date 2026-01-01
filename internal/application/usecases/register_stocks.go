package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/collection"
)

type RegisterStocks struct {
	unitOfWorkFactory  services.UnitOfWorkFactory
	dataSourceProvider map[string]services.DataSourceService
}

// Execute runs the stock registration use case.
//
// Parameters:
//
// - ctx: request-scoped context
//
// - service: identifier of the data source service to use
//
// - limitDate: optional lower bound date for fetching data
//
// Returns:
//
// - number of stock registers inserted
//
// - error if any step fails
func (uc *RegisterStocks) Execute(ctx context.Context, service string, limitDate *time.Time) (int, error) {

	s, ok := uc.dataSourceProvider[service]
	if !ok || s == nil {
		log.Printf("[register stocks] no data source service found for %s", service)
		return 0, pkg.InternalServerError("dataSourceService was nil for this request")
	}

	data, err := s.Get(ctx, limitDate)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		log.Printf("[register stocks] no stocks found for %s for date range %s", service, limitDate)
		return 0, nil
	}

	insertedCount := 0
	err = uc.unitOfWorkFactory.Do(ctx, func(txCtx context.Context, uow services.UnitOfWork) error {

		brokerages, err := uc.GetBrokerages(txCtx, uow.BrokerageRepository(), data)
		if err != nil {
			return err
		}

		markets, err := uc.GetMarkets(txCtx, uow.MarketRepository(), data)
		if err != nil {
			return err
		}

		companies, err := uc.GetCompanies(txCtx, uow.CompanyRepository(), markets, data)
		if err != nil {
			return err
		}

		stocks, err := uc.GetStocks(txCtx, uow.StockRepository(), markets, companies, data)
		if err != nil {
			return err
		}

		stockRegisters, err := uc.SetStockRegisters(
			txCtx,
			uow.StockRegisterRepository(),
			uow.RecommendationRepository(),
			brokerages,
			markets,
			companies,
			stocks,
			data,
		)
		if err != nil {
			return err
		}

		if err = uc.IncrementTendencyStats(txCtx, uow.StockTendencyStatRepository(), stockRegisters); err != nil {
			return err
		}
		insertedCount = len(stockRegisters)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return insertedCount, nil
}

func (uc *RegisterStocks) GetMarkets(
	ctx context.Context,
	r domain.MarketRepository,
	data []services.DataSourceResponse,
) (map[string]*domain.Market, error) {

	marketNames := collection.NewSet[string]()
	for _, d := range data {
		marketNames.Add(d.Market.Name)
	}

	marketNamesSlice := marketNames.Slice()
	marketMap, err := r.GetByNames(ctx, marketNamesSlice)
	if err != nil {
		return nil, err
	}

	newMarkets := make([]*domain.Market, 0, len(marketNamesSlice))
	for name, market := range marketMap {
		if market != nil {
			continue
		}

		newMarkets = append(newMarkets, &domain.Market{Name: name})
	}

	if len(newMarkets) == 0 {
		return marketMap, nil
	}

	if err = r.SaveAll(ctx, newMarkets); err != nil {
		return nil, err
	}

	for i := range newMarkets {
		m := newMarkets[i]
		marketMap[m.Name] = m
	}

	return marketMap, nil
}

func (uc *RegisterStocks) GetCompanies(
	ctx context.Context,
	r domain.CompanyRepository,
	markets map[string]*domain.Market,
	data []services.DataSourceResponse,
) (map[domain.MarketCompanySearchParam]*domain.Company, error) {
	searchParams := collection.NewSet[domain.MarketCompanySearchParam]()
	for _, d := range data {
		market, ok := markets[d.Market.Name]
		if !ok || market == nil {
			return nil, pkg.InternalServerError("market was nil when getting companies, invalid data state")
		}

		searchParams.Add(domain.MarketCompanySearchParam{MarketID: market.ID, Name: d.Company.Name})
	}

	searchParamsSlice := searchParams.Slice()
	companyMap, err := r.GetByMarketCompanySearch(ctx, searchParamsSlice)
	if err != nil {
		return nil, err
	}

	newCompanies := make([]*domain.Company, 0, len(searchParamsSlice))

	for key, company := range companyMap {
		if company != nil {
			continue
		}

		newCompanies = append(newCompanies, &domain.Company{MarketID: key.MarketID, Name: key.Name})
	}

	if len(newCompanies) == 0 {
		return companyMap, nil
	}

	if err = r.SaveAll(ctx, newCompanies); err != nil {
		return nil, err
	}

	for i := range newCompanies {
		c := newCompanies[i]
		companyMap[domain.MarketCompanySearchParam{MarketID: c.MarketID, Name: c.Name}] = c
	}

	return companyMap, nil
}

func (uc *RegisterStocks) GetStocks(
	ctx context.Context,
	r domain.StockRepository,
	markets map[string]*domain.Market,
	companies map[domain.MarketCompanySearchParam]*domain.Company,
	data []services.DataSourceResponse,
) (map[domain.StockCompanySearchParam]*domain.Stock, error) {

	stockArgMap := make(map[domain.StockCompanySearchParam]struct {
		Name *string
		ISIN *string
	})
	searchKeySet := collection.NewSet[domain.StockCompanySearchParam]()

	for _, d := range data {
		market, ok := markets[d.Market.Name]
		if !ok || market == nil {
			return nil, pkg.InternalServerError("market was nil when getting stocks, invalid data state")
		}

		company, ok := companies[domain.MarketCompanySearchParam{MarketID: market.ID, Name: d.Company.Name}]
		if !ok || company == nil {
			return nil, pkg.InternalServerError("company was nil when getting stocks, invalid data state")
		}

		if company.MarketID != market.ID {
			log.Printf("[register stocks] illegal state: company market does not match stock market, company: %+v, market: %+v", company, market)
			return nil, pkg.InternalServerError("illegal state: company market does not match stock market")
		}

		key := domain.StockCompanySearchParam{
			StockTicker: d.Stock.Ticker,
			CompanyID:   company.ID,
			MarketID:    market.ID,
		}
		stockArgs := struct {
			Name *string
			ISIN *string
		}{}

		if d.Stock.Name != "" {
			stockArgs.Name = &d.Stock.Name
		}

		if d.Stock.ISIN != "" {
			stockArgs.ISIN = &d.Stock.ISIN
		}
		stockArgMap[key] = stockArgs
		searchKeySet.Add(key)
	}

	searchKeySlice := searchKeySet.Slice()
	stockMap, err := r.GetByStockCompanySearchParams(ctx, searchKeySlice)
	if err != nil {
		return nil, err
	}

	newStocks := make([]*domain.Stock, 0, len(searchKeySlice))
	for key, stock := range stockMap {
		if stock != nil {
			continue
		}

		newStocks = append(newStocks, &domain.Stock{
			Ticker:    key.StockTicker,
			CompanyID: key.CompanyID,
			MarketID:  key.MarketID,
			Name:      stockArgMap[key].Name,
			Isin:      stockArgMap[key].ISIN,
		})
	}

	if len(newStocks) == 0 {
		return stockMap, nil
	}

	if err = r.SaveAll(ctx, newStocks); err != nil {
		return nil, err
	}

	for i := range newStocks {
		s := newStocks[i]
		key := domain.StockCompanySearchParam{
			StockTicker: s.Ticker,
			CompanyID:   s.CompanyID,
			MarketID:    s.MarketID,
		}

		stockMap[key] = s
	}

	return stockMap, nil
}

func (uc *RegisterStocks) GetBrokerages(
	ctx context.Context,
	r domain.BrokerageRepository,
	data []services.DataSourceResponse,
) (map[string]*domain.Brokerage, error) {

	brokerageNames := collection.NewSet[string]()

	for _, d := range data {
		if d.Recomendation == nil {
			continue
		}

		brokerageNames.Add(d.Recomendation.Brokerage.Name)
	}

	if len(brokerageNames) == 0 {
		return make(map[string]*domain.Brokerage), nil
	}

	brokerageMap, err := r.GetByNames(ctx, brokerageNames.Slice())
	if err != nil {
		return nil, err
	}

	newBrokerages := make([]*domain.Brokerage, 0, len(brokerageMap))
	for name, val := range brokerageMap {
		if val != nil {
			continue
		}

		newBrokerages = append(newBrokerages, &domain.Brokerage{Name: name})
	}

	if len(newBrokerages) == 0 {
		return brokerageMap, nil
	}

	if err = r.SaveAll(ctx, newBrokerages); err != nil {
		return nil, err
	}

	for _, brokerage := range newBrokerages {
		brokerageMap[brokerage.Name] = brokerage
	}

	return brokerageMap, nil
}

func (uc *RegisterStocks) SetStockRegisters(
	ctx context.Context,
	registerRepository domain.StockRegisterRepository,
	recommendationRepository domain.RecommendationRepository,
	brokerages map[string]*domain.Brokerage,
	markets map[string]*domain.Market,
	companies map[domain.MarketCompanySearchParam]*domain.Company,
	stocks map[domain.StockCompanySearchParam]*domain.Stock,
	data []services.DataSourceResponse,
) ([]*domain.StockRegister, error) {
	newRegisters := make([]*domain.StockRegister, 0, len(data))
	newRecommendationArgs := make([]struct {
		recommendation *domain.Recommendation
		stockRegister  *domain.StockRegister
	}, 0, len(data))

	for _, d := range data {
		market, ok := markets[d.Market.Name]
		if !ok || market == nil {
			return nil, pkg.InternalServerError("market was not properly mapped for stock register setting")
		}

		companyKey := domain.MarketCompanySearchParam{
			MarketID: market.ID,
			Name:     d.Company.Name,
		}
		company, ok := companies[companyKey]
		if !ok || company == nil {
			return nil, pkg.InternalServerError("company was not properly mapped for stock register setting")
		}

		stockKey := domain.StockCompanySearchParam{
			MarketID:    market.ID,
			CompanyID:   company.ID,
			StockTicker: d.Stock.Ticker,
		}
		stock, ok := stocks[stockKey]
		if !ok || stock == nil {
			return nil, pkg.InternalServerError("stock was not properly mapped for stock register setting")
		}

		register := &domain.StockRegister{
			StockID:   stock.ID,
			Price:     d.Stock.Price,
			Tendency:  d.Stock.Tendency,
			CreatedAt: d.Time.UTC(),
		}

		newRegisters = append(newRegisters, register)

		if d.Recomendation == nil {
			continue
		}

		brokerage, ok := brokerages[d.Recomendation.Brokerage.Name]
		if !ok || brokerage == nil {
			return nil, pkg.InternalServerError("brokerage was not properly mapped for stock recommendations setting")
		}

		recommendation := &domain.Recommendation{
			BrokerageID:     brokerage.ID,
			StockRegisterID: 0, // assing later
			TargetTo:        d.Recomendation.TargetTo,
			TargetFrom:      d.Recomendation.TargetFrom,
			RatingTo:        d.Recomendation.RatingTo,
			RatingFrom:      d.Recomendation.RatingFrom,
			CreatedAt:       d.Time.UTC(),
		}

		newRecommendationArgs = append(newRecommendationArgs,
			struct {
				recommendation *domain.Recommendation
				stockRegister  *domain.StockRegister
			}{recommendation, register})
	}

	if err := registerRepository.SaveAll(ctx, newRegisters); err != nil {
		return nil, err
	}

	if len(newRecommendationArgs) > 0 {
		newRecommendations := collection.Map(newRecommendationArgs, func(arg struct {
			recommendation *domain.Recommendation
			stockRegister  *domain.StockRegister
		}) *domain.Recommendation {
			arg.recommendation.StockRegisterID = arg.stockRegister.ID
			return arg.recommendation
		})

		if err := recommendationRepository.SaveAll(ctx, newRecommendations); err != nil {
			return nil, err
		}
	}

	return newRegisters, nil
}

func (uc *RegisterStocks) IncrementTendencyStats(
	ctx context.Context,
	r domain.StockTendencyStatRepository,
	registers []*domain.StockRegister,
) error {
	deltas := make(map[uint64]domain.StockTendencyDelta)

	for _, register := range registers {
		delta := deltas[register.StockID]

		switch register.Tendency {
		case domain.Up:
			delta.Up++
		case domain.Side:
			delta.Side++
		case domain.Down:
			delta.Down++
		default:
			return pkg.InternalServerError(fmt.Sprintf("invalid tendency found from register record %+v", register))
		}

		deltas[register.StockID] = delta
	}

	return r.IncrementAll(ctx, deltas)
}

func NewRegisterStocks(uow services.UnitOfWorkFactory, ds map[string]services.DataSourceService) *RegisterStocks {

	if uow == nil || ds == nil {
		panic(fmt.Sprintf(
			"nil arguments for NewRegisterStocks, unitOfWorkFactory: %+v, dataSources: %+v", uow, ds),
		)
	}

	return &RegisterStocks{uow, ds}
}
