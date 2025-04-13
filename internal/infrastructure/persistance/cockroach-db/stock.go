package cockroachdb

import (
	"context"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func (r *stockRepository) Get(ctx context.Context, id uint) (*domain.Stock, error) {

	panic("implement me")
}

func (r *stockRepository) GetByTicker(ctx context.Context, marketID uint, ticker string) (*domain.Stock, error) {

	panic("implement me")
}

func (r *stockRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationReponse[domain.Stock], error) {

	panic("implement me")
}

func (r *stockRepository) Create(ctx context.Context, stock *domain.Stock) error {

	panic("implement me")
}

func (r *stockRepository) Update(ctx context.Context, stockID uint, stock *domain.Stock) error {

	panic("implement me")
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	if db == nil {
		log.Fatalf("bad impl: db is nil in NewStockRepository")
	}

	return &stockRepository{db: db}
}
