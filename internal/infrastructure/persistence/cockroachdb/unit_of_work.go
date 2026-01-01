package cockroachdb

import (
	"context"
	"log"

	appservices "github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	tx *sqlx.Tx
}

func (u *UnitOfWork) MarketRepository() domain.MarketRepository {
	return NewMarketRepository(u.tx)
}

func (u *UnitOfWork) StockRepository() domain.StockRepository {
	return NewStockRepository(u.tx)
}

func (u *UnitOfWork) StockRegisterRepository() domain.StockRegisterRepository {
	return NewStockRegisterRepository(u.tx)
}

func (u *UnitOfWork) StockTendencyStatRepository() domain.StockTendencyStatRepository {
	return NewStockTendencyStatRepository(u.tx)
}

func (u *UnitOfWork) UserRepository() domain.UserRepository {
	return NewUserRepository(u.tx)
}

func (u *UnitOfWork) RecommendationRepository() domain.RecommendationRepository {
	return NewRecommendationRepository(u.tx)
}

func (u *UnitOfWork) BrokerageRepository() domain.BrokerageRepository {
	return NewBrokerageRepository(u.tx)
}

func (u *UnitOfWork) CompanyRepository() domain.CompanyRepository {
	return NewCompanyRepository(u.tx)
}

func newSqlxUnitOfWork(db *sqlx.Tx) *UnitOfWork {
	return &UnitOfWork{db}
}

type UnitOfWorkFactory struct {
	db *sqlx.DB
}

func (f *UnitOfWorkFactory) Do(ctx context.Context, cb func(txCtx context.Context, uow appservices.UnitOfWork) error) error {

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

func NewSqlxUnitOfWorkFactory(db *sqlx.DB) *UnitOfWorkFactory {
	if db == nil {
		panic("db *sqlx.DB is nil")
	}

	return &UnitOfWorkFactory{db}
}
