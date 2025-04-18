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

func (uc *LoginUseCase) Execute(ctx context.Context, username, password string) (session string, err error) {
	user, err := uc.ur.GetByUsername(ctx, username)

	if err != nil {
		return "", err
	}

	if err = auth.VerifyPassword(password, user.Password); err != nil {
		return "", pkg.Unauthorized(err.Error())
	}

	session, err = auth.GenerateSessionToken(user)

	if err != nil {
		return "", err
	}

	return session, nil
}

func NewLoginUseCase(ur domain.UserRepository) *LoginUseCase {

	if ur == nil {
		log.Fatalln("bad impl: UserRepository was nil for NewLoginUseCase")
	}

	return &LoginUseCase{ur}
}
