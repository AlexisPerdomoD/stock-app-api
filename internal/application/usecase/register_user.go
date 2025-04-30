package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
)

type RegisterUserUseCase struct {
	ur domain.UserRepository
}

func (uc RegisterUserUseCase) Execute(ctx context.Context, usr *domain.User) error {

	hashed, err := auth.HashPassword(usr.Password)
	if err != nil {
		return err
	}

	usr.Password = hashed
	defer auth.ZeroBytes(usr.Password)

	return uc.ur.Create(ctx, usr)
}

func NewRegisterUserUseCase(ur domain.UserRepository) *RegisterUserUseCase {

	if ur == nil {
		log.Fatalln("user repository is nil, stopping :b")
	}

	return &RegisterUserUseCase{ur}
}
