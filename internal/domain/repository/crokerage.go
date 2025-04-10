package repository

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"
	"github.com/alexisPerdomoD/stock-app-api/internal/shared"
)

type BrokerageRepository interface {
	Get(id int) (entity.Brokerage, shared.ApiErr)
	Create(args entity.Brokerage) (entity.Brokerage, shared.ApiErr)
}
