package cockroachdb

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

/* mapUserInsert assumes that the args are not nil */
func mapUserInsert(args *domain.User) *userRecord {

	return &userRecord{
		UserName: args.UserName,
		Password: string(args.Password),
		Active:   args.Active,
	}
}

/* mapMarketToDomain assumes that the record is not nil */
func mapMarketToDomain(record *marketRecord, target *domain.Market) *domain.Market {
	if target == nil {
		return &domain.Market{
			Name:      record.Name,
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
		}
	}

	target.ID = record.ID
	target.Name = record.Name

	return target
}

/* mapStockToDomain assumes that the record is not nil */
func mapStockToDomain(record *stockRecord, target *domain.Stock) *domain.Stock {

	if target == nil {
		return &domain.Stock{
			ID:        record.ID,
			CompanyID: record.CompanyID,
			Ticker:    record.Ticker,
			Name:      record.Name,
			Tendency:  record.Tendency,
			Price:     record.Price,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}

	}
	target.ID = record.ID
	target.CompanyID = record.CompanyID
	target.Ticker = record.Ticker
	target.Name = record.Name
	target.Tendency = record.Tendency
	target.Price = record.Price
	target.CreatedAt = record.CreatedAt
	target.UpdatedAt = record.UpdatedAt

	return target
}

/* mapPopulatedStockToDomain assumes that the record is not nil */
func mapPopulatedStockToDomain(record *stockRecord, target *domain.PopulatedStock) *domain.PopulatedStock {
	if target == nil {
		return &domain.PopulatedStock{
			Stock:       *mapStockToDomain(record, nil),
			CompanyName: record.Company.Name,
			Market:      *mapMarketToDomain(&record.Company.Market, nil),
		}
	}

	target.Stock = *mapStockToDomain(record, nil)
	target.CompanyName = record.Company.Name
	target.Market = *mapMarketToDomain(&record.Company.Market, nil)
	return target
}

/* mapRecommendationToDomain assumes that the record is not nil */
func mapRecommendationToDomain(record *recommendationRecord, target *domain.Recommendation) *domain.Recommendation {
	if target == nil {
		return &domain.Recommendation{
			ID:          record.ID,
			StockID:     record.StockID,
			BrokerageID: record.BrokerageID,
			RatingTo:    record.RatingTo,
			RatingFrom:  record.RatingFrom,
			TargetTo:    record.TargetTo,
			TargetFrom:  record.TargetFrom,
			CreatedAt:   record.CreatedAt,
		}

	}

	target.ID = record.ID
	target.StockID = record.StockID
	target.BrokerageID = record.BrokerageID
	target.RatingFrom = record.RatingFrom
	target.RatingTo = record.RatingTo
	target.TargetFrom = record.TargetFrom
	target.TargetTo = record.TargetTo
	target.CreatedAt = record.CreatedAt

	return target
}

/* mapPopulatedRecommendationToDomain assumes that the record is not nil */
func mapPopulatedRecommendationToDomain(
	record *recommendationRecord,
	target *domain.PopulatedRecommendation,
) *domain.PopulatedRecommendation {
	if target == nil {
		return &domain.PopulatedRecommendation{
			BrokerageName:  record.Brokerage.Name,
			Recommendation: *mapRecommendationToDomain(record, nil),
		}
	}

	target.ID = record.ID
	target.StockID = record.StockID
	target.BrokerageID = record.BrokerageID
	target.BrokerageName = record.Brokerage.Name
	target.RatingFrom = record.RatingFrom
	target.RatingTo = record.RatingTo
	target.TargetFrom = record.TargetFrom
	target.TargetTo = record.TargetTo
	target.CreatedAt = record.CreatedAt

	return target
}

/* mapUserToDomain assumes that the record is not nil */
func mapUserToDomain(record *userRecord, target *domain.User) *domain.User {

	if target != nil {
		target.ID = record.ID
		target.UserName = record.UserName
		target.Password = []byte(record.Password)
		target.Active = record.Active
		return target
	}

	return &domain.User{
		ID:       record.ID,
		UserName: record.UserName,
		Password: []byte(record.Password),
		Active:   record.Active,
	}

}
