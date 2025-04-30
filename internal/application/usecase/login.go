package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
)

type LoginUseCase struct {
	ur domain.UserRepository
}

func (uc *LoginUseCase) Execute(ctx context.Context, username string, password []byte) (*domain.User, error) {
	user, err := uc.ur.GetByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, pkg.Unauthorized("Invalid credentials")
	}

	if err = auth.VerifyPassword(password, user.Password); err != nil {
		return nil, pkg.Unauthorized("Invalid credentials")
	}

	return user, nil
}

func NewLoginUseCase(ur domain.UserRepository) *LoginUseCase {

	if ur == nil {
		log.Fatalln("[NewLoginUseCase]: UserRepository was nil")
	}

	return &LoginUseCase{ur}
}
