package cockroachdb

import (
	"context"
	"log/slog"

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
	db     *sqlx.DB
	logger *slog.Logger
}

func (f *UnitOfWorkFactory) Do(ctx context.Context, cb func(txCtx context.Context, uow appservices.UnitOfWork) error) error {

	tx, err := f.db.BeginTxx(ctx, nil)
	if err != nil {
		f.logger.WarnContext(ctx, "failed to start transaction", "err", err)
		return err
	}
	defer func() {
		if err == nil {
			return
		}

		if rollbackerr := tx.Rollback(); rollbackerr != nil {
			f.logger.WarnContext(ctx, "failed to rollback transaction", "err", rollbackerr)
		}
	}()

	uow := newSqlxUnitOfWork(tx)
	err = cb(ctx, uow)

	if err != nil {
		f.logger.DebugContext(ctx, "error happened while executing unit of work", "err", err)
		return err
	}

	err = tx.Commit()
	return err
}

func NewUnitOfWorkFactory(db *sqlx.DB, logger *slog.Logger) *UnitOfWorkFactory {
	if db == nil || logger == nil {
		panic("db *sqlx.DB is nil")
	}

	return &UnitOfWorkFactory{db, logger}
}
