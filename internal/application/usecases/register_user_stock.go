package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type RegisterUserStock struct {
	ur domain.UserRepository
}

func (uc *RegisterUserStock) Execute(ctx context.Context, userID uint, stockID uint) error {

	if err := uc.ur.RegisterUserStock(ctx, userID, stockID); err != nil {
		return pkg.BadRequest("Stock is not valid")
	}
	return nil
}

func NewRegisterUserStock(ur domain.UserRepository) *RegisterUserStock {

	if ur == nil {
		log.Fatalln("bad impl: UserRepository was nil for NewRegisterUserStockUseCase")
	}

	return &RegisterUserStock{ur: ur}
}
