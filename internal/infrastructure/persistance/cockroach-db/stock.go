package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func (r *stockRepository) Get(ctx context.Context, id uint) (*domain.Stock, error) {

	record := &stockRecord{}

	if err := r.db.WithContext(ctx).First(record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err

	}

	return mapStockToDomain(record, nil), nil
}

func (r *stockRepository) GetByTicker(ctx context.Context, marketID uint, ticker string) (*domain.Stock, error) {

	record := &stockRecord{}

	if err := r.db.
		WithContext(ctx).
		Joins("JOIN companies on companies.id = stock.companies_id").
		Where("companies.market_id = ? AND stock.ticker = ?", marketID, ticker).
		First(record).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err

	}

	return mapStockToDomain(record, nil), nil
}

func (r *stockRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[domain.PopulatedStock], error) {

	allowedFilters := map[string]bool{
		"name":       true,
		"company_id": true,
		"price":      true,
		"ticker":     true,
		"tendency":   true,
		"user_id":    true,
	}

	allowedSorters := map[string]bool{
		"tendency": true,
		"price":    true,
	}
	var total int64
	var records []stockRecord

	query := r.db.WithContext(ctx).Model(stockRecord{})
	query = applyFilters(query, filter.FilterBy, allowedFilters)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := applyPagination(query, &filter, allowedSorters).
		Find(&records).Error; err != nil {
		return nil, err
	}

	stocks := []domain.PopulatedStock{}

	for _, record := range records {
		stock := domain.Stock{}
		_ = mapStockToDomain(&record, &stock)
		populated := domain.PopulatedStock{
			Stock: stock,
			// TODO: RESOLVE THIS
			CompanyName: "",
			Market:      domain.Market{},
		}

		stocks = append(stocks, populated)
	}

	page := 1
	if filter.Page > 1 {
		page = filter.Page
	}

	result := pkg.PaginationReponse[domain.PopulatedStock]{
		Items:     stocks,
		Page:      page,
		PageSize:  len(stocks),
		TotalSize: int(total),
	}

	return &result, nil
}

func (r *stockRepository) Register(ctx context.Context, data []domain.SourceStockData) error {
	panic("Implement me")
}

func (r *stockRepository) create(ctx context.Context, stock *domain.Stock) error {

	if stock == nil {
		return pkg.BadRequest("args to stock insertion were not provided")
	}

	record := mapStockInsert(stock)

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}
	_ = mapStockToDomain(record, stock)

	return nil
}

func (r *stockRepository) update(ctx context.Context, stockID uint, updates *domain.StockUpdates) error {
	if updates == nil {
		return pkg.BadRequest("args for stock update were not provided")
	}

	if updates.Name == nil && updates.Price == nil && updates.Tendency == nil {
		return pkg.BadRequest("no fields to update")
	}

	return r.db.WithContext(ctx).
		Model(&stockRecord{}).
		Where("id = ?", stockID).
		Updates(updates).Error
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	if db == nil {
		log.Fatalf("bad impl: db is nil in NewStockRepository")
	}

	return &stockRepository{db: db}
}
