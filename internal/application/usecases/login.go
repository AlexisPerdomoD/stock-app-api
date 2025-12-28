package usecases

import (
	"context"
	appmodel "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
	"log"
)

type Login struct {
	ur domain.UserRepository
}

func (uc *Login) Execute(ctx context.Context, credentials *appmodel.UserLoginDTO) (*appmodel.UserView, error) {
	password := credentials.GetPasswordBytesAndClean()
	defer auth.ZeroBytes(password)

	includePassword := true
	user, err := uc.ur.GetByUsername(ctx, credentials.Username, includePassword)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, pkg.Unauthorized("Invalid credentials")
	}
	defer auth.ZeroBytes(user.Password)

	validPassword, err := auth.VerifyPassword(password, user.Password)

	if err != nil {
		return nil, err
	}

	if !validPassword {
		return nil, pkg.Unauthorized("Invalid credentials")
	}

	return appmodel.NewUserView(user), nil
}

func NewLogin(ur domain.UserRepository) *Login {

	if ur == nil {
		log.Fatalln("[NewLoginUseCase]: UserRepository was nil")
	}

	return &Login{ur}
}
