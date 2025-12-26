package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
)

type RegisterUser struct {
	ur domain.UserRepository
}

func (uc RegisterUser) Execute(ctx context.Context, usr *domain.User) error {

	hashed, err := auth.HashPassword(usr.Password)
	if err != nil {
		return err
	}

	usr.Password = hashed
	defer auth.ZeroBytes(usr.Password)

	return uc.ur.Create(ctx, usr)
}

func NewRegisterUser(ur domain.UserRepository) *RegisterUser {

	if ur == nil {
		log.Fatalln("user repository is nil, stopping :b")
	}

	return &RegisterUser{ur}
}
