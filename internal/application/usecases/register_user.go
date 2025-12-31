package usecases

import (
	"context"
	"log/slog"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
)

type RegisterUser struct {
	userRepository domain.UserRepository
	log            *slog.Logger
}

func (uc *RegisterUser) Execute(ctx context.Context, dto *models.RegisterUserDTO) (*models.UserView, error) {
	rawPassword, err := dto.GetPasswordBytesAndClean()
	if err != nil {
		uc.log.Warn("InvalidStateErr found, credentials probably comsumed before usecase", "err", err)
		return nil, pkg.InvalidStateErr(err.Error())
	}
	password, err := auth.HashPassword(rawPassword)
	auth.ZeroBytes(rawPassword)
	if err != nil {
		return nil, err
	}

	usr := dto.ToDomain()
	usr.Password = password
	err = uc.userRepository.Save(ctx, usr)
	auth.ZeroBytes(password)
	if err != nil {
		return nil, err
	}

	return models.NewUserView(usr), nil
}

func NewRegisterUser(ur domain.UserRepository, logger *slog.Logger) *RegisterUser {

	if ur == nil || logger == nil {
		panic("nil dependencies where passed for NewRegisterUser")
	}

	log := logger.With("usecase", "RegisterUser")
	return &RegisterUser{ur, log}
}
