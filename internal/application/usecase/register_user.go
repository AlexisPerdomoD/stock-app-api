package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
)

type RegisterUserUseCase struct {
	ur domain.UserRepository
}

func (uc RegisterUserUseCase) Execute(ctx context.Context, usr *domain.User) (session string, err error) {

	hashed, err := auth.HashPassword(usr.Password)

	if err != nil {
		return "", pkg.InternalServerError("Error hashing password")
	}

	usr.Password = hashed

	if err := uc.ur.Create(ctx, usr); err != nil {
		return "", err
	}

	return auth.GenerateSessionToken(usr)
}

func NewRegisterUserUseCase(ur domain.UserRepository) *RegisterUserUseCase {

	if ur == nil {
		log.Fatalln("user repository is nil, stopping :b")
	}

	return &RegisterUserUseCase{ur}
}
