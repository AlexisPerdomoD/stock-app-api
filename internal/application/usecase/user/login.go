package usecase

import (
	r "github.com/alexisPerdomoD/stock-app-api/internal/domain/repository"
	ss "github.com/alexisPerdomoD/stock-app-api/internal/shared/service"
)

type loginUseCase struct {
	userRepository r.UserRepository
}

func (uc *loginUseCase) Execute(username, password string) (session string, err error) {
	user, err := uc.userRepository.GetByUsername(username)

	if err != nil {
		return "", err
	}

	if err = ss.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}

	session, err = ss.GenerateSessionToken(user)

	if err != nil {
		return "", err
	}

	return session, nil
}

func NewLoginUseCase(ur r.UserRepository) *loginUseCase {

	if ur == nil {
		panic("user repository is nil, stopping :b")
	}

	return &loginUseCase{userRepository: ur}
}
