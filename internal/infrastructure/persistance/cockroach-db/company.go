package cockroachdb

import (
	"context"
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func (r *companyRepository) Get(ctx context.Context, id uint) (*domain.Company, error) {
	record := &companyRecord{}
	err := r.db.WithContext(ctx).First(record, id).Error

	if err.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return mapCompanyToDomain(record, nil), nil
}

func (r *companyRepository) GetByMarketIdAndName(ctx context.Context, marketID uint, name string) (*domain.Company, error) {

	record := &companyRecord{}
	err := r.db.WithContext(ctx).First(record).Where("'Name'= ? AND 'MarketID'= ?", name, marketID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return mapCompanyToDomain(record, nil), nil
}

func (r *companyRepository) Create(ctx context.Context, company *domain.Company) error {
	if company == nil {
		return pkg.BadRequest("args for company insertion were not provited")
	}

	record := mapCompanyInsert(company)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}

	_ = mapCompanyToDomain(record, company)

	return nil

}

func NewCompanyRepository(db *gorm.DB) *companyRepository {
	if db == nil {
		log.Fatalf("bad impl: db is nil in NewCompanyRepository")
	}

	return &companyRepository{db: db}
}
