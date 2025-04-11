package usecase

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	pkgService "github.com/alexisPerdomoD/stock-app-api/internal/pkg/service"
)

type loginUseCase struct {
	userRepository domain.UserRepository
}

func (uc *loginUseCase) Execute(username, password string) (session string, err error) {
	user, err := uc.userRepository.GetByUsername(username)

	if err != nil {
		return "", err
	}

	if err = pkgService.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}

	session, err = pkgService.GenerateSessionToken(user)

	if err != nil {
		return "", err
	}

	return session, nil
}

func NewLoginUseCase(ur domain.UserRepository) *loginUseCase {

	if ur == nil {
		panic("user repository is nil, stopping :b")
	}

	return &loginUseCase{userRepository: ur}
}
