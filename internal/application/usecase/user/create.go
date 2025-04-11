package usecase

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	r "github.com/alexisPerdomoD/stock-app-api/internal/domain/repository"
	s "github.com/alexisPerdomoD/stock-app-api/internal/shared"
	ss "github.com/alexisPerdomoD/stock-app-api/internal/shared/service"
)

type createUserUseCase struct {
	userRepository r.UserRepository
}

func (uc createUserUseCase) Execute(args *entity.User) error {

	hashed, err := ss.HashPassword(args.Password)

	if err != nil {
		return s.InternalServerError("Error hashing password")
	}

	args.Password = hashed

	if err := uc.userRepository.Create(args); err != nil {
		return err
	}

	return nil
}

func NewCreateUserUseCase(ur r.UserRepository) *createUserUseCase {

	if ur == nil {
		panic("user repository is nil, stopping :b")
	}

	return &createUserUseCase{userRepository: ur}
}
