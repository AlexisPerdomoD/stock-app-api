package cockroachdb_test

import (
	"context"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service/mock"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/stretchr/testify/assert"
)

var stock *domain.PopulatedStock
var getAllPaginatedTest = []struct {
	name             string
	filter           pkg.PaginationFilter
	userID           *uint
	successCondition func(got *pkg.PaginationReponse[domain.PopulatedStock], t *assert.Assertions) bool
	wantErr          bool
}{
	{
		name:   "must resturn all stocks ordered by updated_at desc by default",
		filter: pkg.PaginationFilter{},
		successCondition: func(
			got *pkg.PaginationReponse[domain.PopulatedStock],
			a *assert.Assertions,
		) bool {

			a.NotNil(got)
			a.Greater(len(got.Items), 0)
			a.Condition(func() bool {
				return len(got.Items) == 20 || len(got.Items) == got.TotalSize
			})

			a.Equal(got.Page, 1)

			for i, item := range got.Items {
				a.NotEqual(item.ID, 0)

				if i == 0 {
					continue
				}
				a.False(got.Items[i-1].UpdatedAt.Before(item.UpdatedAt))
			}

			return true
		},
	},

	{
		name: "check page size and page number and sort by ASC UpdatedAt filters",
		filter: pkg.PaginationFilter{
			SortBy: map[string]pkg.SortOrder{
				"updated_at": pkg.SortOrderAsc,
			},
			PaginationPage: pkg.PaginationPage{
				Page: 2,
				Size: 5,
			},
		},
		successCondition: func(
			got *pkg.PaginationReponse[domain.PopulatedStock],
			t *assert.Assertions,
		) bool {

			t.Equal(got.Page, 2)
			t.Condition(func() bool {

				return len(got.Items) == 5 || len(got.Items) == got.TotalSize
			})

			for i, item := range got.Items {
				t.NotEqual(item.ID, 0)

				if i == 0 {
					continue
				}

				t.False(got.Items[i-1].UpdatedAt.After(item.UpdatedAt))

			}

			stock = &got.Items[0]
			return true
		},
	},
}

func Test_stockRepository_Register(t *testing.T) {
	db := cockroachdb.NewDB()
	r := cockroachdb.NewStockRepository(db)
	sourceDataService := mock.MockSourceStockService{}

	data, err := sourceDataService.Get(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to get invalid data: %v", err)
	}
	assert := assert.New(t)

	tests := []struct {
		name    string
		data    []domain.SourceStockData
		wantErr bool
	}{
		{
			name: "must do all insertions without error",
			data: data,
		},
		{
			name:    "provided nil data must return error",
			data:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			gotErr := r.Register(context.Background(), test.data)
			if test.wantErr {
				assert.Error(gotErr)
			} else {
				assert.Nil(gotErr)
			}
		})
	}
}

func Test_stockRepository_GetAllPaginated(t *testing.T) {
	db := cockroachdb.NewDB()
	r := cockroachdb.NewStockRepository(db)
	assert := assert.New(t)
	for _, test := range getAllPaginatedTest {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := r.GetAllPaginated(context.Background(), test.filter, test.userID)
			if test.wantErr {
				assert.Error(gotErr)
			} else {
				assert.Nil(gotErr)
			}
			assert.True(test.successCondition(got, assert))
		})
	}
}

func Test_stockRepository_Get(t *testing.T) {
	db := cockroachdb.NewDB()
	r := cockroachdb.NewStockRepository(db)
	assert := assert.New(t)
	tests := []struct {
		name    string
		id      uint
		want    *domain.PopulatedStock
		wantErr bool
	}{
		{
			name: "must return nil if stock does not exist",
			id:   0,
			want: nil,
		},
		{
			name: "must return stock if stock exists",
			id:   stock.ID,
			want: stock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.Get(context.Background(), tt.id, nil)
			if !tt.wantErr {
				assert.Nil(gotErr)
			} else {
				assert.Error(gotErr)
			}

			if tt.want != nil {
				assert.Equal(tt.want, got)
			} else {
				assert.Nil(got)
			}

		})
	}
}
