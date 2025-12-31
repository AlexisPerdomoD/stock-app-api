package usecases

import (
	"context"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type RegisterUserStock struct {
	userRepository domain.UserRepository
}

func (uc *RegisterUserStock) Execute(ctx context.Context, userID uint64, stockID uint64) error {
	hasRegister, err := uc.userRepository.HasUserStock(ctx, userID, stockID)
	if err != nil {
		return err
	}

	if hasRegister {
		return nil
	}

	return uc.userRepository.RegisterUserStock(ctx, userID, stockID)
}

func NewRegisterUserStock(ur domain.UserRepository) *RegisterUserStock {

	if ur == nil {
		panic("nil dependencies for NewRegisterUserStock")
	}

	return &RegisterUserStock{userRepository: ur}
}
