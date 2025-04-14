package usecase

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type removeUserStockUseCase struct {
	ur domain.UserRepository
}

func (r *removeUserStockUseCase) Execute(ctx context.Context, userID uint, stockID uint) error {

	panic("implement me")
}

func NewRemoveUserStockUserCase(ur domain.UserRepository) *removeUserStockUseCase {
	if ur == nil {
		log.Println("bad impl: UserRepository was passed as nil for NewRemoveUserStockUserCase")
	}

	return &removeUserStockUseCase{ur}

}
