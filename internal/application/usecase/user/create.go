package usecase

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	pkgService "github.com/alexisPerdomoD/stock-app-api/internal/pkg/service"
)

type createUserUseCase struct {
	userRepository domain.UserRepository
}

func (uc createUserUseCase) Execute(args *domain.User) error {

	hashed, err := pkgService.HashPassword(args.Password)

	if err != nil {
		return pkg.InternalServerError("Error hashing password")
	}

	args.Password = hashed

	if err := uc.userRepository.Create(args); err != nil {
		return err
	}

	return nil
}

func NewCreateUserUseCase(ur domain.UserRepository) *createUserUseCase {

	if ur == nil {
		panic("user repository is nil, stopping :b")
	}

	return &createUserUseCase{userRepository: ur}
}
