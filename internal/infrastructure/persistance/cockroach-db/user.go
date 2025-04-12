package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Get(ctx context.Context, id uint) (*domain.User, error) {
	usr := &userTable{}
	result := r.db.WithContext(ctx).First(usr, id)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return usr.ToDomain(), nil

}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	usr := &userTable{}
	result := r.db.WithContext(ctx).Where("username = ?", username).First(usr)

	if result.Error != nil {
		
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		
		return nil, result.Error
	}

	return usr.ToDomain(), nil

}

func NewUserRepository(db *gorm.DB) *userRepository {

	if db == nil {
		log.Fatalf("db is nil in NewUserRepository")
	}

	return &userRepository{db: db}

}
