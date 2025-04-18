package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

type RemoveUserStockUseCase struct {
	ur domain.UserRepository
}

func (ruc *RemoveUserStockUseCase) Execute(ctx context.Context, userID uint, stockID uint) error {

	if err:= ruc.ur.RemoveUserStock(ctx, userID, stockID); err !=nil{
		return pkg.BadRequest("Stock is invalid")
	}
	return nil
}

func NewRemoveUserStockUserCase(ur domain.UserRepository) *RemoveUserStockUseCase {
	if ur == nil {
		log.Println("bad impl: UserRepository was passed as nil for NewRemoveUserStockUserCase")
	}

	return &RemoveUserStockUseCase{ur}

}
