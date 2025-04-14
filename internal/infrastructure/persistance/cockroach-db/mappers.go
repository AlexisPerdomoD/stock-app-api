package cockroachdb

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

/* All mappers insertions assume that the args are not nil */

func mapMarketInsert(args *domain.Market) *marketRecord {
	return &marketRecord{Name: args.Name}
}

/* mapCompanyInsert assumes that the args are not nil */
func mapCompanyInsert(args *domain.Company) *companyRecord {

	return &companyRecord{
		MarketID: args.MarketID,
		ISIN:     args.ISIN,
		Name:     args.Name,
	}
}

/* mapBrokerageInsert assumes that the args are not nil */
func mapBrokerageInsert(args *domain.Brokerage) *brokerageRecord {
	return &brokerageRecord{Name: args.Name}
}

/* mapStockInsert assumes that the args are not nil */
func mapStockInsert(args *domain.Stock) *stockRecord {

	return &stockRecord{
		Name:      args.Name,
		CompanyID: args.CompanyID,
		Ticker:    args.Ticker,
		Price:     args.Price,
		Tendency:  args.Tendency,
	}
}

/* mapRecommendationInsert assumes that the args are not nil */
func mapRecommendationInsert(args *domain.Recommendation) *recommendationRecord {

	return &recommendationRecord{
		StockID:     args.StockID,
		BrokerageID: args.BrokerageID,
		RatingFrom:  args.RatingFrom,
		RatingTo:    args.RatingTo,
		TargetFrom:  args.TargetFrom,
		TargetTo:    args.TargetFrom,
	}
}

/* mapUserInsert assumes that the args are not nil */
func mapUserInsert(args *domain.User) *userRecord {

	return &userRecord{
		UserName: args.UserName,
		Password: args.Password,
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

/* mapCompanyToDomain assumes that the record is not nil */
func mapCompanyToDomain(record *companyRecord, target *domain.Company) *domain.Company {

	if target == nil {
		return &domain.Company{
			ID:        record.ID,
			MarketID:  record.MarketID,
			Name:      record.Name,
			ISIN:      record.ISIN,
			CreatedAt: record.CreatedAt,
		}
	}
	return nil
}

/* mapBrokerageToDomain assumes that the record is not nil */
func mapBrokerageToDomain(record *brokerageRecord, target *domain.Brokerage) *domain.Brokerage {

	if target == nil {
		return &domain.Brokerage{
			ID:        record.ID,
			Name:      record.Name,
			CreatedAt: record.CreatedAt,
		}

	}
	return nil
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

/* mapUserToDomain assumes that the record is not nil */
func mapUserToDomain(record *userRecord, target *domain.User) *domain.User {

	if target != nil {
		target.UserName = record.UserName
		target.Password = record.Password
		target.Active = record.Active
		return target
	}

	return &domain.User{
		ID:       record.ID,
		UserName: record.UserName,
		Password: record.Password,
		Active:   record.Active,
	}

}
