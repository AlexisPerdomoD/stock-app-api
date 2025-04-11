package usecase

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type registerUserStockUseCase struct {
	userRepository domain.UserRepository
}

func (uc *registerUserStockUseCase) Execute(ctx context.Context, userID int, stockID int) error {

	return uc.userRepository.RegisterUserStock(ctx, userID, stockID)

}

func NewRegisterUserStockUseCase(ur domain.UserRepository) *registerUserStockUseCase {

	if ur == nil {
		panic("user repository is nil, stopping :b")
	}

	return &registerUserStockUseCase{userRepository: ur}
}
