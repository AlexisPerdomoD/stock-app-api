package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
)

type RegisterUser struct {
	ur domain.UserRepository
}

func (uc RegisterUser) Execute(ctx context.Context, args *models.RegisterUserDTO) (*models.UserView, error) {
	usr := args.ToDomain()

	hashed, err := auth.HashPassword(usr.Password)
	if err != nil {
		return nil, err
	}

	usr.Password = hashed
	defer auth.ZeroBytes(usr.Password)

	if err := uc.ur.Create(ctx, usr); err != nil {
		return nil, err
	}

	return models.NewUserView(usr), nil
}

func NewRegisterUser(ur domain.UserRepository) *RegisterUser {

	if ur == nil {
		log.Fatalln("user repository is nil, stopping :b")
	}

	return &RegisterUser{ur}
}
