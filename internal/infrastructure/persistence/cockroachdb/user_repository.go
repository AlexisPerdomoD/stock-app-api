package cockroachdb

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct{}

func (r UserRepository) Get(ctx context.Context, id uint64, includePassword bool) (*domain.User, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r UserRepository) GetByUsername(ctx context.Context, username string, includePassword bool) (*domain.User, error) {
	return nil, pkg.InternalServerError("not implemented")
}

func (r UserRepository) Create(ctx context.Context, args *domain.User) error {
	return pkg.InternalServerError("not implemented")
}

func (r UserRepository) RegisterUserStock(ctx context.Context, userID uint, stockID uint) error {
	return pkg.InternalServerError("not implemented")
}

func (r UserRepository) RemoveUserStock(ctx context.Context, userID uint, stockID uint) error {
	return pkg.InternalServerError("not implemented")
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{}
}
