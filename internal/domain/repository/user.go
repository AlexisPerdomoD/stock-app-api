package repository

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"

type UserRepository interface {
	Get(id int) (*entity.User, error)

	GetByUsername(username string) (*entity.User, error)

	Create(args *entity.User) error
}
