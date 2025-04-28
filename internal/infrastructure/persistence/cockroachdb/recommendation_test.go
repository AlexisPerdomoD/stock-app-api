package cockroachdb_test

import (
	"context"
	"strings"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service/mock"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/stretchr/testify/assert"
)

var recommendationGetAllPaginatedTests = []struct {
	name             string
	filter           pkg.PaginationFilter
	stockID          uint
	successCondition func(
		stockID uint,
		got *pkg.PaginationReponse[domain.PopulatedRecommendation],
		a *assert.Assertions,
	) bool
	wantErr bool
}{
	{
		name:    "must resturn empty list if invalid stock id provided",
		filter:  pkg.PaginationFilter{},
		stockID: 0,
		successCondition: func(
			stockID uint,
			got *pkg.PaginationReponse[domain.PopulatedRecommendation],
			a *assert.Assertions,
		) bool {
			a.NotNil(got)
			a.Equal(0, len(got.Items))
			return true
		},
	},
	{
		name: "must return a list of recommendations for a stock default ordering by created_at desc",
		filter: pkg.PaginationFilter{
			PaginationPage: pkg.PaginationPage{
				Page: 1,
				Size: 10,
			},
		},
		successCondition: func(
			stockID uint,
			got *pkg.PaginationReponse[domain.PopulatedRecommendation],
			a *assert.Assertions,
		) bool {
			a.NotNil(got)
			a.Equal(1, got.Page)
			a.Condition(func() bool {
				hasItems := len(got.Items) > 0
				hasValidItemCount := len(got.Items) == got.PageSize ||
					len(got.Items) == got.TotalSize
				return hasItems && hasValidItemCount
			})

			for i, item := range got.Items {
				a.True(item.StockID == stockID)

				if i == 0 {
					continue
				}
				a.True(item.CreatedAt.Before(got.Items[i-1].CreatedAt))
			}

			return true
		},
	},

	{
		name: "must return only recommendations with brokerage name provided",
		filter: pkg.PaginationFilter{
			Search: strings.ToLower("mock"),
		},
		successCondition: func(
			stockID uint,
			got *pkg.PaginationReponse[domain.PopulatedRecommendation],
			a *assert.Assertions,
		) bool {
			a.NotNil(got)
			a.Equal(1, got.Page)
			a.Greater(got.PageSize, 0)
			a.Greater(got.TotalSize, 0)
			a.Greater(len(got.Items), 0)
			a.Condition(func() bool {
				hasValidItemCount := len(got.Items) == got.PageSize ||
					len(got.Items) == got.TotalSize
				return hasValidItemCount
			})

			for i, item := range got.Items {
				a.True(item.StockID == stockID)

				if i == 0 {
					continue
				}
				a.True(strings.HasPrefix(item.BrokerageName, "mock"))
				a.True(item.CreatedAt.Before(got.Items[i-1].CreatedAt))
			}

			return true
		},
	},
}

func TestRecommendationRepository_GetAllPaginated(t *testing.T) {
	sss := mock.MockSourceStockService{}
	sr := cockroachdb.NewStockRepository(cockroachdb.NewDB())

	data, err := sss.Get(context.Background(), nil)
	if err != nil || len(data) == 0 {
		t.Fatalf("failed to get invalid data: %v", err)
	}

	if err := sr.Register(context.Background(), data); err != nil {
		t.Fatalf("failed to register stocks: %v", err)
	}

	stocks, err := sr.GetAllPaginated(context.Background(), pkg.PaginationFilter{}, nil)
	if err != nil {
		t.Fatalf("failed to get stocks: %v", err)
	}
	if len(stocks.Items) < 2 {
		t.Fatalf("failed to get stocks: %v", err)
	}
	recommendationGetAllPaginatedTests[1].stockID = stocks.Items[0].ID
	recommendationGetAllPaginatedTests[2].stockID = stocks.Items[1].ID

	rr := cockroachdb.NewRecommendationRepository(cockroachdb.NewDB())
	assert := assert.New(t)

	for _, tt := range recommendationGetAllPaginatedTests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := rr.GetAllPaginated(context.Background(), tt.filter, tt.stockID)
			if tt.wantErr {
				assert.Error(gotErr)
			} else {
				assert.Nil(gotErr)
			}

			assert.True(tt.successCondition(tt.stockID, got, assert))
		})
	}
}
