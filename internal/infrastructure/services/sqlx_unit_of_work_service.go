package services

import (
	"context"
	"log"

	appservices "github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/jmoiron/sqlx"
)

type SqlxUnitOfWork struct {
	tx *sqlx.Tx
}

func (u *SqlxUnitOfWork) StockRepository() domain.StockRepository {
	return cockroachdb.NewStockRepository(u.tx)
}

func (u *SqlxUnitOfWork) StockRegisterRepository() domain.StockRegisterRepository {
	return cockroachdb.NewStockRegisterRepository(u.tx)
}

func (u *SqlxUnitOfWork) UserRepository() domain.UserRepository {
	return cockroachdb.NewUserRepository(u.tx)
}

func (u *SqlxUnitOfWork) RecommendationRepository() domain.RecommendationRepository {
	return cockroachdb.NewRecommendationRepository(u.tx)
}

func (u *SqlxUnitOfWork) BrokerageRepository() domain.BrokerageRepository {
	return cockroachdb.NewBrokerageRepository(u.tx)
}

func (u *SqlxUnitOfWork) CompanyRepository() domain.CompanyRepository {
	return cockroachdb.NewCompanyRepository(u.tx)
}

func newSqlxUnitOfWork(db *sqlx.Tx) *SqlxUnitOfWork {
	return &SqlxUnitOfWork{db}
}

type SqlxUnitOfWorkFactory struct {
	db *sqlx.DB
}

func (f *SqlxUnitOfWorkFactory) Do(ctx context.Context, cb func(txCtx context.Context, uow appservices.UnitOfWork) error) error {

	tx, err := f.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err = tx.Rollback(); err != nil {
			log.Printf("[sqlx unit of work] failed to rollback transaction: %v", err)
		}
	}()

	uow := newSqlxUnitOfWork(tx)

	err = cb(ctx, uow)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func NewSqlxUnitOfWorkFactory(db *sqlx.DB) *SqlxUnitOfWorkFactory {
	if db == nil {
		panic("db *sqlx.DB is nil")
	}

	return &SqlxUnitOfWorkFactory{db}
}
