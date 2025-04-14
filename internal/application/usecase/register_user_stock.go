package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type RegisterUserStockUseCase struct {
	ur domain.UserRepository
}

func (uc *RegisterUserStockUseCase) Execute(ctx context.Context, userID uint, stockID uint) error {

	return uc.ur.RegisterUserStock(ctx, userID, stockID)

}

func NewRegisterUserStockUseCase(ur domain.UserRepository) *RegisterUserStockUseCase {

	if ur == nil {
		log.Fatalln("bad impl: UserRepository was nil for NewRegisterUserStockUseCase")
	}

	return &RegisterUserStockUseCase{ur: ur}
}
