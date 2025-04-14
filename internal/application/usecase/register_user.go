package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	pkgService "github.com/alexisPerdomoD/stock-app-api/internal/pkg/service"
)

type registerUserUseCase struct {
	userRepository domain.UserRepository
}

func (uc registerUserUseCase) Execute(ctx context.Context, args *domain.User) error {

	hashed, err := pkgService.HashPassword(args.Password)

	if err != nil {
		return pkg.InternalServerError("Error hashing password")
	}

	args.Password = hashed

	if err := uc.userRepository.Create(ctx, args); err != nil {
		return err
	}

	return nil
}

func NewRegisterUserUseCase(ur domain.UserRepository) *registerUserUseCase {

	if ur == nil {
		log.Fatalln("user repository is nil, stopping :b")
	}

	return &registerUserUseCase{userRepository: ur}
}
