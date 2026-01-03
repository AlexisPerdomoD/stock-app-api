package usecases_test

import (
	"context"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// BROKERAGE REPOSITORY ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeBrokerageRepository struct {
	GetByIDFn    func(ctx context.Context, id uint64) (*domain.Brokerage, error)
	GetByNamesFn func(ctx context.Context, names []string) (map[string]*domain.Brokerage, error)
	SaveFn       func(ctx context.Context, brokerage *domain.Brokerage) error
	SaveAllFn    func(ctx context.Context, brokerages []*domain.Brokerage) error
}

func (f *FakeBrokerageRepository) GetByID(ctx context.Context, id uint64) (*domain.Brokerage, error) {
	if f.GetByIDFn == nil {
		panic("FakeBrokerageRepository.GetByIDFn not set")
	}
	return f.GetByIDFn(ctx, id)
}

func (f *FakeBrokerageRepository) GetByNames(ctx context.Context, names []string) (map[string]*domain.Brokerage, error) {
	if f.GetByNamesFn == nil {
		panic("FakeBrokerageRepository.GetByNamesFn not set")
	}
	return f.GetByNamesFn(ctx, names)
}

func (f *FakeBrokerageRepository) Save(ctx context.Context, brokerage *domain.Brokerage) error {
	if f.SaveFn == nil {
		panic("FakeBrokerageRepository.SaveFn not set")
	}
	return f.SaveFn(ctx, brokerage)
}

func (f *FakeBrokerageRepository) SaveAll(ctx context.Context, brokerages []*domain.Brokerage) error {
	if f.SaveAllFn == nil {
		panic("FakeBrokerageRepository.SaveAllFn not set")
	}
	return f.SaveAllFn(ctx, brokerages)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// USER REPOSITORY ////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeUserRepository struct {
	GetByIDFn                   func(ctx context.Context, id uint64) (*domain.User, error)
	GetByIDWithPasswordFn       func(ctx context.Context, id uint64) (*domain.User, error)
	GetByUsernameFn             func(ctx context.Context, username string) (*domain.User, error)
	GetByUsernameWithPasswordFn func(ctx context.Context, username string) (*domain.User, error)
	SaveFn                      func(ctx context.Context, user *domain.User) error
	HasUserStockFn              func(ctx context.Context, userID, stockID uint64) (bool, error)
	RegisterUserStockFn         func(ctx context.Context, userID, stockID uint64) error
	RemoveUserStockFn           func(ctx context.Context, userID, stockID uint64) error
}

func (f *FakeUserRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	if f.GetByIDFn == nil {
		panic("FakeUserRepository.GetByIDFn not set")
	}
	return f.GetByIDFn(ctx, id)
}

func (f *FakeUserRepository) GetByIDWithPassword(ctx context.Context, id uint64) (*domain.User, error) {
	if f.GetByIDWithPasswordFn == nil {
		panic("FakeUserRepository.GetByIDWithPasswordFn not set")
	}
	return f.GetByIDWithPasswordFn(ctx, id)
}

func (f *FakeUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if f.GetByUsernameFn == nil {
		panic("FakeUserRepository.GetByUsernameFn not set")
	}
	return f.GetByUsernameFn(ctx, username)
}

func (f *FakeUserRepository) GetByUsernameWithPassword(ctx context.Context, username string) (*domain.User, error) {
	if f.GetByUsernameWithPasswordFn == nil {
		panic("FakeUserRepository.GetByUsernameWithPasswordFn not set")
	}
	return f.GetByUsernameWithPasswordFn(ctx, username)
}

func (f *FakeUserRepository) Save(ctx context.Context, user *domain.User) error {
	if f.SaveFn == nil {
		panic("FakeUserRepository.SaveFn not set")
	}
	return f.SaveFn(ctx, user)
}

func (f *FakeUserRepository) HasUserStock(ctx context.Context, userID, stockID uint64) (bool, error) {
	if f.HasUserStockFn == nil {
		panic("FakeUserRepository.HasUserStockFn not set")
	}
	return f.HasUserStockFn(ctx, userID, stockID)
}

func (f *FakeUserRepository) RegisterUserStock(ctx context.Context, userID, stockID uint64) error {
	if f.RegisterUserStockFn == nil {
		panic("FakeUserRepository.RegisterUserStockFn not set")
	}
	return f.RegisterUserStockFn(ctx, userID, stockID)
}

func (f *FakeUserRepository) RemoveUserStock(ctx context.Context, userID, stockID uint64) error {
	if f.RemoveUserStockFn == nil {
		panic("FakeUserRepository.RemoveUserStockFn not set")
	}
	return f.RemoveUserStockFn(ctx, userID, stockID)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// MARKET REPOSITORY ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeMarketRepository struct {
	GetByIDFn    func(ctx context.Context, id uint64) (*domain.Market, error)
	GetByNamesFn func(ctx context.Context, names []string) (map[string]*domain.Market, error)
	SaveFn       func(ctx context.Context, market *domain.Market) error
	SaveAllFn    func(ctx context.Context, markets []*domain.Market) error
}

func (f *FakeMarketRepository) GetByID(ctx context.Context, id uint64) (*domain.Market, error) {
	if f.GetByIDFn == nil {
		panic("FakeMarketRepository.GetByIDFn not set")
	}
	return f.GetByIDFn(ctx, id)
}

func (f *FakeMarketRepository) GetByNames(ctx context.Context, names []string) (map[string]*domain.Market, error) {
	if f.GetByNamesFn == nil {
		panic("FakeMarketRepository.GetByNamesFn not set")
	}
	return f.GetByNamesFn(ctx, names)
}

func (f *FakeMarketRepository) Save(ctx context.Context, market *domain.Market) error {
	if f.SaveFn == nil {
		panic("FakeMarketRepository.SaveFn not set")
	}
	return f.SaveFn(ctx, market)
}

func (f *FakeMarketRepository) SaveAll(ctx context.Context, markets []*domain.Market) error {
	if f.SaveAllFn == nil {
		panic("FakeMarketRepository.SaveAllFn not set")
	}
	return f.SaveAllFn(ctx, markets)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// COMPANY REPOSITORY ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeCompanyRepository struct {
	GetByIDFn                  func(ctx context.Context, id uint64) (*domain.Company, error)
	GetByMarketCompanySearchFn func(ctx context.Context, params []domain.MarketCompanySearchParam) (map[domain.MarketCompanySearchParam]*domain.Company, error)
	SaveFn                     func(ctx context.Context, company *domain.Company) error
	SaveAllFn                  func(ctx context.Context, companies []*domain.Company) error
}

func (f *FakeCompanyRepository) GetByID(ctx context.Context, id uint64) (*domain.Company, error) {
	if f.GetByIDFn == nil {
		panic("FakeCompanyRepository.GetByIDFn not set")
	}
	return f.GetByIDFn(ctx, id)
}

func (f *FakeCompanyRepository) GetByMarketCompanySearch(ctx context.Context, params []domain.MarketCompanySearchParam) (map[domain.MarketCompanySearchParam]*domain.Company, error) {
	if f.GetByMarketCompanySearchFn == nil {
		panic("FakeCompanyRepository.GetByMarketCompanySearchFn not set")
	}
	return f.GetByMarketCompanySearchFn(ctx, params)
}

func (f *FakeCompanyRepository) Save(ctx context.Context, company *domain.Company) error {
	if f.SaveFn == nil {
		panic("FakeCompanyRepository.SaveFn not set")
	}
	return f.SaveFn(ctx, company)
}

func (f *FakeCompanyRepository) SaveAll(ctx context.Context, companies []*domain.Company) error {
	if f.SaveAllFn == nil {
		panic("FakeCompanyRepository.SaveAllFn not set")
	}
	return f.SaveAllFn(ctx, companies)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// STOCK REPOSITORY ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeStockRepository struct {
	GetByIDFn                       func(ctx context.Context, stockID uint64) (*domain.Stock, error)
	GetByStockCompanySearchParamsFn func(ctx context.Context, params []domain.StockCompanySearchParam) (map[domain.StockCompanySearchParam]*domain.Stock, error)
	GetAllPaginatedFn               func(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationResponse[domain.PopulatedStock], error)
	GetAllPaginatedByUserFn         func(ctx context.Context, filter pkg.PaginationFilter, userID uint64) (*pkg.PaginationResponse[domain.PopulatedStock], error)
	SaveFn                          func(ctx context.Context, stock *domain.Stock) error
	SaveAllFn                       func(ctx context.Context, stocks []*domain.Stock) error
	UpdateFn                        func(ctx context.Context, updates domain.StockUpdates) error
}

func (f *FakeStockRepository) GetByID(ctx context.Context, stockID uint64) (*domain.Stock, error) {
	if f.GetByIDFn == nil {
		panic("FakeStockRepository.GetByIDFn not set")
	}
	return f.GetByIDFn(ctx, stockID)
}

func (f *FakeStockRepository) GetByStockCompanySearchParams(ctx context.Context, params []domain.StockCompanySearchParam) (map[domain.StockCompanySearchParam]*domain.Stock, error) {
	if f.GetByStockCompanySearchParamsFn == nil {
		panic("FakeStockRepository.GetByStockCompanySearchParamsFn not set")
	}
	return f.GetByStockCompanySearchParamsFn(ctx, params)
}

func (f *FakeStockRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationResponse[domain.PopulatedStock], error) {
	if f.GetAllPaginatedFn == nil {
		panic("FakeStockRepository.GetAllPaginatedFn not set")
	}
	return f.GetAllPaginatedFn(ctx, filter)
}

func (f *FakeStockRepository) GetAllPaginatedByUser(ctx context.Context, filter pkg.PaginationFilter, userID uint64) (*pkg.PaginationResponse[domain.PopulatedStock], error) {
	if f.GetAllPaginatedByUserFn == nil {
		panic("FakeStockRepository.GetAllPaginatedByUserFn not set")
	}
	return f.GetAllPaginatedByUserFn(ctx, filter, userID)
}

func (f *FakeStockRepository) Save(ctx context.Context, stock *domain.Stock) error {
	if f.SaveFn == nil {
		panic("FakeStockRepository.SaveFn not set")
	}
	return f.SaveFn(ctx, stock)
}

func (f *FakeStockRepository) SaveAll(ctx context.Context, stocks []*domain.Stock) error {
	if f.SaveAllFn == nil {
		panic("FakeStockRepository.SaveAllFn not set")
	}
	return f.SaveAllFn(ctx, stocks)
}

func (f *FakeStockRepository) Update(ctx context.Context, updates domain.StockUpdates) error {
	if f.UpdateFn == nil {
		panic("FakeStockRepository.UpdateFn not set")
	}
	return f.UpdateFn(ctx, updates)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// RECOMMENDATION REPOSITORY ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeRecommendationRepository struct {
	getAllFn func(ctx context.Context, filters pkg.PaginationFilter) (*pkg.PaginationResponse[domain.PopulatedRecommendation], error)

	saveAllFn func(ctx context.Context, recommendations []*domain.Recommendation) error
}

func (f *FakeRecommendationRepository) GetAllPaginated(
	ctx context.Context,
	filters pkg.PaginationFilter,
) (*pkg.PaginationResponse[domain.PopulatedRecommendation], error) {
	return f.getAllFn(ctx, filters)
}

func (f *FakeRecommendationRepository) SaveAll(
	ctx context.Context,
	recommendations []*domain.Recommendation,
) error {
	return f.saveAllFn(ctx, recommendations)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// STOCK REGISTER REPOSITORY //////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FakeStockRegisterRepository struct {
	GetByIDFn           func(ctx context.Context, id uint64) (*domain.StockRegister, error)
	GetLastByStockIDFn  func(ctx context.Context, stockID uint64) (*domain.StockRegister, error)
	GetLastsByStockIDFn func(ctx context.Context, stockID uint64, limit uint16) ([]domain.StockRegister, error)
	GetRangeByStockIDFn func(ctx context.Context, stockID uint64, from, to time.Time) ([]domain.StockRegister, error)
	SaveFn              func(ctx context.Context, register *domain.StockRegister) error
	SaveAllFn           func(ctx context.Context, registers []*domain.StockRegister) error
}

func (f *FakeStockRegisterRepository) GetByID(ctx context.Context, id uint64) (*domain.StockRegister, error) {
	if f.GetByIDFn == nil {
		panic("FakeStockRegisterRepository.GetByIDFn not set")
	}
	return f.GetByIDFn(ctx, id)
}

func (f *FakeStockRegisterRepository) GetLastByStockID(ctx context.Context, stockID uint64) (*domain.StockRegister, error) {
	if f.GetLastByStockIDFn == nil {
		panic("FakeStockRegisterRepository.GetLastByStockIDFn not set")
	}
	return f.GetLastByStockIDFn(ctx, stockID)
}

func (f *FakeStockRegisterRepository) GetLastsByStockID(ctx context.Context, stockID uint64, limit uint16) ([]domain.StockRegister, error) {
	if f.GetLastsByStockIDFn == nil {
		panic("FakeStockRegisterRepository.GetLastsByStockIDFn not set")
	}
	return f.GetLastsByStockIDFn(ctx, stockID, limit)
}

func (f *FakeStockRegisterRepository) GetRangeByStockID(ctx context.Context, stockID uint64, from, to time.Time) ([]domain.StockRegister, error) {
	if f.GetRangeByStockIDFn == nil {
		panic("FakeStockRegisterRepository.GetRangeByStockIDFn not set")
	}
	return f.GetRangeByStockIDFn(ctx, stockID, from, to)
}

func (f *FakeStockRegisterRepository) Save(ctx context.Context, register *domain.StockRegister) error {
	if f.SaveFn == nil {
		panic("FakeStockRegisterRepository.SaveFn not set")
	}
	return f.SaveFn(ctx, register)
}

func (f *FakeStockRegisterRepository) SaveAll(ctx context.Context, registers []*domain.StockRegister) error {
	if f.SaveAllFn == nil {
		panic("FakeStockRegisterRepository.SaveAllFn not set")
	}
	return f.SaveAllFn(ctx, registers)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// STOCK TENDENCY STAT REPOSITORY //////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FakeStockTendencyStatRepository struct {
	GetByStockIDFn func(ctx context.Context, stockID uint64) (*domain.StockTendencyStat, error)
	IncrementFn    func(ctx context.Context, stockID uint64, delta domain.StockTendencyDelta) error
	IncrementAllFn func(ctx context.Context, deltas map[uint64]domain.StockTendencyDelta) error
}

func (f *FakeStockTendencyStatRepository) GetByStockID(ctx context.Context, stockID uint64) (*domain.StockTendencyStat, error) {
	if f.GetByStockIDFn == nil {
		panic("FakeStockTendencyStatRepository.GetByStockIDFn not set")
	}

	return f.GetByStockIDFn(ctx, stockID)
}

func (f *FakeStockTendencyStatRepository) Increment(ctx context.Context, stockID uint64, delta domain.StockTendencyDelta) error {
	if f.IncrementFn == nil {
		panic("FakeStockTendencyStatRepository.IncrementFn not set")
	}

	return f.IncrementFn(ctx, stockID, delta)
}

func (f *FakeStockTendencyStatRepository) IncrementAll(
	ctx context.Context,
	deltas map[uint64]domain.StockTendencyDelta,
) error {
	if f.IncrementAllFn == nil {
		panic("FakeStockTendencyStatRepository.IncrementAllFn not set")
	}

	return f.IncrementAllFn(ctx, deltas)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// DATA SOURCE ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeDataSource struct {
	NameFn func() string
	GetFn  func(ctx context.Context, since *time.Time) ([]services.DataSourceResponse, error)
}

func (f *FakeDataSource) Name() string {
	if f.NameFn != nil {
		return f.NameFn()
	}
	return "unknown"
}

func (f *FakeDataSource) Get(
	ctx context.Context,
	since *time.Time,
) ([]services.DataSourceResponse, error) {
	if f.GetFn == nil {
		panic("FakeDataSource.GetFn not set")
	}

	return f.GetFn(ctx, since)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////// UNIT OF WORK ///////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FakeUnitOfWork struct {
	StockRepo          domain.StockRepository
	StockRegisterRepo  domain.StockRegisterRepository
	StockTendencyRepo  domain.StockTendencyStatRepository
	UserRepo           domain.UserRepository
	RecommendationRepo domain.RecommendationRepository
	BrokerageRepo      domain.BrokerageRepository
	CompanyRepo        domain.CompanyRepository
	MarketRepo         domain.MarketRepository
}

func (u *FakeUnitOfWork) StockRepository() domain.StockRepository {
	return u.StockRepo
}

func (u *FakeUnitOfWork) StockRegisterRepository() domain.StockRegisterRepository {
	return u.StockRegisterRepo
}

func (u *FakeUnitOfWork) StockTendencyStatRepository() domain.StockTendencyStatRepository {
	return u.StockTendencyRepo
}

func (u *FakeUnitOfWork) UserRepository() domain.UserRepository {
	return u.UserRepo
}

func (u *FakeUnitOfWork) RecommendationRepository() domain.RecommendationRepository {
	return u.RecommendationRepo
}

func (u *FakeUnitOfWork) BrokerageRepository() domain.BrokerageRepository {
	return u.BrokerageRepo
}

func (u *FakeUnitOfWork) CompanyRepository() domain.CompanyRepository {
	return u.CompanyRepo
}

func (u *FakeUnitOfWork) MarketRepository() domain.MarketRepository {
	return u.MarketRepo
}

type FakeUnitOfWorkFactory struct {
	UOW  services.UnitOfWork
	DoFn func(ctx context.Context, tx func(context.Context, services.UnitOfWork) error) error
}

func (f *FakeUnitOfWorkFactory) Do(
	ctx context.Context,
	tx func(txCtx context.Context, uow services.UnitOfWork) error,
) error {

	if f.DoFn != nil {
		return f.DoFn(ctx, tx)
	}

	// happy path por defecto
	return tx(ctx, f.UOW)
}

func NewHappyPathUnitOfWork() *FakeUnitOfWork {
	return &FakeUnitOfWork{
		StockRepo:          &FakeStockRepository{},
		StockRegisterRepo:  &FakeStockRegisterRepository{},
		StockTendencyRepo:  &FakeStockTendencyStatRepository{},
		UserRepo:           &FakeUserRepository{},
		RecommendationRepo: &FakeRecommendationRepository{},
		BrokerageRepo:      &FakeBrokerageRepository{},
		CompanyRepo:        &FakeCompanyRepository{},
		MarketRepo:         &FakeMarketRepository{},
	}
}
