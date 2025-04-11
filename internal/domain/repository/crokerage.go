package repository

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"

type BrokerageRepository interface {
	Get(id int) (*entity.Brokerage, error)

	Create(args *entity.Brokerage) error
}
