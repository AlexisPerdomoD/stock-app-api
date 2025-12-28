package services

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

/*
UnitOfWork
Represents a Unit of Work meant to be used in repositories operations in a atomic way.

It is meant to be used single time by operations
*/
type UnitOfWork interface {
	/*
		StockRepository returns a transactional stock repository.
	*/
	StockRepository() domain.StockRepository

	/*
		StockRegisterRepository returns a transactional stock repository.
	*/
	StockRegisterRepository() domain.StockRegisterRepository

	/*
		UserRepository returns a transactional user repository.
	*/
	UserRepository() domain.UserRepository

	/*
		RecommendationRepository returns a transactional recommendation repository.
	*/
	RecommendationRepository() domain.RecommendationRepository

	/*
		BrokerageRepository returns a transactional brokerage repository.
	*/
	BrokerageRepository() domain.BrokerageRepository

	/*
		CompanyRepository returns a transactional company repository.
	*/
	CompanyRepository() domain.CompanyRepository
}

/*
UnitOfWorkFactory
Factory for UnitOfWork.
*/
type UnitOfWorkFactory interface {
	/*
		Do executes a transaction.
		- Automatically handles transactions and rollbacks in case of errors.
		- Automatically commits transactions in case of success.
	*/
	Do(ctx context.Context, tx func(txCtx context.Context, uow UnitOfWork) error) error
}
