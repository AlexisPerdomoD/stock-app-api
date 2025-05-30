package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type RegisterUserStockUseCase struct {
	ur domain.UserRepository
}

func (uc *RegisterUserStockUseCase) Execute(ctx context.Context, userID uint, stockID uint) error {

	if err := uc.ur.RegisterUserStock(ctx, userID, stockID); err != nil {
		return pkg.BadRequest("Stock is not valid")
	}
	return nil
}

func NewRegisterUserStockUseCase(ur domain.UserRepository) *RegisterUserStockUseCase {

	if ur == nil {
		log.Fatalln("bad impl: UserRepository was nil for NewRegisterUserStockUseCase")
	}

	return &RegisterUserStockUseCase{ur: ur}
}
