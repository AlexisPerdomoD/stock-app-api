package usecases

import (
	"context"
	"log/slog"

	appmodel "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
)

type Login struct {
	ur  domain.UserRepository
	log *slog.Logger
}

func (uc *Login) Execute(ctx context.Context, credentials *appmodel.UserLoginDTO) (*appmodel.UserView, error) {
	password, err := credentials.GetPasswordBytesAndClean()
	if err != nil {
		uc.log.Warn("InvalidStateErr found, credentials probably comsumed before usecase", "err", err)
		return nil, pkg.InvalidStateErr(err.Error())
	}
	defer auth.ZeroBytes(password)

	user, err := uc.ur.GetByUsernameWithPassword(ctx, credentials.GetSafeUsername())
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, pkg.Unauthorized("Invalid credentials")
	}

	defer auth.ZeroBytes(user.Password)

	isValidCredentials, err := auth.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, err
	}

	if !isValidCredentials {
		return nil, pkg.Unauthorized("Invalid credentials")
	}

	return appmodel.NewUserView(user), nil
}

func NewLogin(ur domain.UserRepository, logger *slog.Logger) *Login {
	if ur == nil || logger == nil {
		panic("some dependencies were provided as nil for NewLogin")
	}

	log := logger.With("usercase", "Login")
	return &Login{ur, log}
}
