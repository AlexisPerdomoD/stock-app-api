package cockroachdb

import (
	"context"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service/mock"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
)

var r *stockRepository
var getAllPaginatedTest = []struct {
	name             string
	filter           pkg.PaginationFilter
	userID           *uint
	successCondition func(got *pkg.PaginationReponse[domain.PopulatedStock], t *testing.T) bool
	wantErr          bool
}{
	{
		name:   "must resturn all stocks ordered by updated_at desc by default",
		filter: pkg.PaginationFilter{},
		successCondition: func(
			got *pkg.PaginationReponse[domain.PopulatedStock],
			t *testing.T,
		) bool {
			t.Helper()

			if got == nil {
				t.Errorf("GetAllPaginated() got nil, want %v", got)
				return false
			}

			if len(got.Items) == 0 {
				t.Errorf("GetAllPaginated() got empty Items, want %v", got)
				return false
			}

			if len(got.Items) != 20 {
				t.Errorf("GetAllPaginated() got Items len %v, want %v",
					len(got.Items), 20,
				)
				return false
			}

			if got.Page != 1 {
				t.Errorf("GetAllPaginated() got Page %v, want %v", got.Page, 1)
				return false
			}

			for i, item := range got.Items {
				if item.ID == 0 {
					t.Errorf("GetAllPaginated() got Item %v Stock.ID %v, Invalid", i, item.ID)
					return false
				}

				if i == 0 {
					continue
				}

				if got.Items[i-1].UpdatedAt.Before(item.UpdatedAt) {
					t.Errorf("GetAllPaginated() got Item order DESC %v Stock.UpdatedAt %v, must be before %v",
						i, item.UpdatedAt, got.Items[i-1].UpdatedAt)
					return false
				}

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
			t *testing.T,
		) bool {
			t.Helper()

			if got.Page != 2 {
				t.Errorf("GetAllPaginated() got Page %v, want %v", got.Page, 2)
				return false
			}

			if got.PageSize != 5 || len(got.Items) != 5 {
				t.Errorf("GetAllPaginated() got PageSize %v, Items len %v, want %v", got.PageSize, len(got.Items), 5)
				return false
			}

			for i, item := range got.Items {
				if item.ID == 0 {
					t.Errorf("GetAllPaginated() got Item %v Stock.ID %v, Invalid", i, item.ID)
					return false
				}

				if i == 0 {
					continue
				}

				if got.Items[i-1].UpdatedAt.After(item.UpdatedAt) {
					t.Errorf("GetAllPaginated() got Item order ASC %v Stock.UpdatedAt %v, must be after %v", i,
						item.UpdatedAt, got.Items[i-1].UpdatedAt,
					)
					return false
				}

			}

			return true
		},
	},
}

func Test_stockRepository_Register(t *testing.T) {
	r = NewStockRepository(testDB)
	sourceDataService := mock.MockSourceStockService{}

	data, err := sourceDataService.Get(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to get invalid data: %v", err)
	}

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
			if gotErr != nil {
				if !test.wantErr {
					t.Errorf("Register() failed: %v", gotErr)
				}
				return
			}
			if test.wantErr {
				t.Fatal("Register() succeeded unexpectedly")
			}
		})
	}
}

func Test_stockRepository_GetAllPaginated(t *testing.T) {

	for _, test := range getAllPaginatedTest {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := r.GetAllPaginated(context.Background(), test.filter, test.userID)
			if gotErr != nil {
				if !test.wantErr {
					t.Errorf("GetAllPaginated() failed: %v", gotErr)
				}
				return
			}
			if test.wantErr {
				t.Fatal("GetAllPaginated() succeeded unexpectedly")
			}

			if !test.successCondition(got, t) {
				t.Errorf("GetAllPaginated() did not pass success condition")
			}
		})
	}
}

func Test_stockRepository_Get(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		id      uint
		want    *domain.PopulatedStock
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.Get(context.Background(), tt.id, nil)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
