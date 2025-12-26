package usecases

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type RemoveUserStock struct {
	ur domain.UserRepository
}

func (ruc *RemoveUserStock) Execute(ctx context.Context, userID uint, stockID uint) error {

	if err := ruc.ur.RemoveUserStock(ctx, userID, stockID); err != nil {
		return pkg.BadRequest("Stock is invalid")
	}
	return nil
}

func NewRemoveUserStock(ur domain.UserRepository) *RemoveUserStock {
	if ur == nil {
		log.Println("bad impl: UserRepository was passed as nil for NewRemoveUserStockUserCase")
	}

	return &RemoveUserStock{ur}

}
