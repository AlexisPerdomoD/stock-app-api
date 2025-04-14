package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type RemoveUserStockUseCase struct {
	ur domain.UserRepository
}

func (r *RemoveUserStockUseCase) Execute(ctx context.Context, userID uint, stockID uint) error {

	panic("implement me")
}

func NewRemoveUserStockUserCase(ur domain.UserRepository) *RemoveUserStockUseCase {
	if ur == nil {
		log.Println("bad impl: UserRepository was passed as nil for NewRemoveUserStockUserCase")
	}

	return &RemoveUserStockUseCase{ur}

}
