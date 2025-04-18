package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Get(ctx context.Context, id uint) (*domain.User, error) {
	record := &userRecord{}
	result := r.db.WithContext(ctx).First(record, id)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return mapUserToDomain(record, nil), nil

}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	record := &userRecord{}
	result := r.db.WithContext(ctx).Where("user_name = ?", username).First(record)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return mapUserToDomain(record, nil), nil

}

func (r *userRepository) Create(ctx context.Context, usr *domain.User) error {

	if usr == nil {
		return pkg.BadRequest("User args provided were nil")
	}

	record := mapUserInsert(usr)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return pkg.DataBaseErr(err.Error(), 400)
	}

	_ = mapUserToDomain(record, usr)

	return nil
}

func (r *userRepository) RegisterUserStock(ctx context.Context, userID uint, stockID uint) error {

	panic("RegisterUserStock not implemented with cockroachdb")

}

func (r *userRepository) RemoveUserStock(ctx context.Context, userID uint, stockID uint) error {

	panic("RemoveUserStock not implemented with cockroachdb")

}

func NewUserRepository(db *gorm.DB) *userRepository {

	if db == nil {
		log.Fatalf("db is nil in NewUserRepository")
	}

	return &userRepository{db: db}

}
