package cockroachdb

import (
	"context"
	"errors"
	"log"
	"math"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func (r *stockRepository) Get(ctx context.Context, id uint) (*domain.PopulatedStock, error) {

	record := &stockRecord{}

	if err := r.db.WithContext(ctx).
		Preload("Company.Market").
		First(record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err

	}

	return mapPopulatedStockToDomain(record, nil), nil
}

func (r *stockRepository) GetAllPaginated(
	ctx context.Context,
	filter pkg.PaginationFilter,
	userID *uint,
) (*pkg.PaginationReponse[domain.PopulatedStock], error) {

	allowedFilters := map[string]bool{
		"name":       true,
		"company_id": true,
		"price":      true,
		"ticker":     true,
		"tendency":   true,
	}

	allowedSorters := map[string]bool{
		"tendency":   true,
		"price":      true,
		"updated_at": true,
	}

	var total int64
	var records []stockRecord
	query := r.db.WithContext(ctx).Model(stockRecord{}).Preload("Company.Market")

	if filter.Search != "" {
		query.Joins("JOIN companies ON companies.id = stocks.company_id").
			Where("companies.name LIKE ?", filter.Search)
	}

	if userID != nil {
		query = query.
			Joins("JOIN user_stocks ON user_stocks.stock_record_id = stocks.id").
			Where("user_stocks.user_record_id = ?", *userID)
	}

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
		stocks = append(stocks, *mapPopulatedStockToDomain(&record, nil))
	}

	page := 1
	if filter.Page > 1 {
		page = filter.Page
	}

	result := pkg.PaginationReponse[domain.PopulatedStock]{
		Items:      stocks,
		Page:       page,
		PageSize:   len(stocks),
		TotalSize:  int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Size))),
	}

	return &result, nil
}

func (r *stockRepository) Register(
	ctx context.Context, data []domain.SourceStockData,
) error {

	markets := map[string]marketRecord{}
	companies := map[string]companyRecord{}
	brokerages := map[string]brokerageRecord{}

	tx := r.db.WithContext(ctx).Begin()
	for _, args := range data {

		market, ok := markets[args.Market.Name]
		if !ok {
			market = marketRecord{}
			if err := tx.First(&market, "name = ?", args.Market.Name).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					tx.Rollback()
					return err
				}

				market = marketRecord{Name: args.Market.Name}
				if err := tx.Create(&market).Error; err != nil {
					tx.Rollback()
					return err
				}
				markets[args.Market.Name] = market
			}
		}

		company, ok := companies[args.Company.Name]
		if !ok {
			company = companyRecord{}
			if err := tx.First(&company, "name = ? AND market_id = ?", args.Company.Name, market.ID).
				Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					tx.Rollback()
					return err
				}

				company = companyRecord{
					MarketID: market.ID,
					ISIN:     args.Company.ISIN,
					Name:     args.Company.Name,
				}
				if err := tx.Create(&company).Error; err != nil {
					tx.Rollback()
					return err
				}
				companies[args.Company.Name] = company
			}
		}

		stock := stockRecord{}
		stock.CompanyID = company.ID
		if err := tx.First(&stock, "ticker = ? AND company_id = ?", args.Stock.Ticker, company.ID).
			Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return err
			}
			stock.Ticker = args.Stock.Ticker
			stock.Name = &args.Stock.Name
			stock.Price = args.Stock.Price
			stock.Tendency = args.Stock.Tendency
			stock.CreatedAt = args.Time
			stock.UpdatedAt = args.Time

			if err := tx.Create(&stock).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			updates := &domain.StockUpdates{
				Price:    &args.Stock.Price,
				Tendency: &args.Stock.Tendency,
			}

			if err := tx.Model(&stock).Where("id = ?", stock.ID).Updates(updates).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if args.Recomendation != nil {
			brokerage, ok := brokerages[args.Recomendation.Brokerage.Name]
			if !ok {
				brokerage = brokerageRecord{}
				if err := tx.First(&brokerage, "name = ?", args.Recomendation.Brokerage.Name).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return err
					}

					brokerage = brokerageRecord{
						Name: args.Recomendation.Brokerage.Name,
					}
					if err := tx.Create(&brokerage).Error; err != nil {
						tx.Rollback()
						return err
					}
					brokerages[args.Recomendation.Brokerage.Name] = brokerage
				}
			}

			if err := tx.Create(&recommendationRecord{
				StockID:     stock.ID,
				BrokerageID: brokerage.ID,
				RatingTo:    args.Recomendation.RatingTo,
				RatingFrom:  args.Recomendation.RatingFrom,
				TargetTo:    args.Recomendation.TargetTo,
				TargetFrom:  args.Recomendation.TargetFrom,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	if db == nil {
		log.Fatalf("bad impl: db is nil in NewStockRepository")
	}

	return &stockRepository{db: db}
}
