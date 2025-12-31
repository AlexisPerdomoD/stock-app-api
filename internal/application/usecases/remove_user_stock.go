package usecases

import (
	"context"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type RemoveUserStock struct {
	ur domain.UserRepository
}

func (ruc *RemoveUserStock) Execute(ctx context.Context, userID uint64, stockID uint64) error {
	return ruc.ur.RemoveUserStock(ctx, userID, stockID)
}

func NewRemoveUserStock(ur domain.UserRepository) *RemoveUserStock {
	if ur == nil {
		panic("nil arguments were passes for NewRemoveUserStock")
	}

	return &RemoveUserStock{ur}

}
