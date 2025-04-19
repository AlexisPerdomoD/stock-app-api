package cockroachdb

import (
	"context"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
)

func Test_stockRepository_Get(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		id      uint
		want    *domain.PopulatedStock
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStockRepository(tt.db)
			got, gotErr := r.Get(context.Background(), tt.id)
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


func Test_stockRepository_Register(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		data    []domain.SourceStockData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStockRepository(tt.db)
			gotErr := r.Register(context.Background(), tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Register() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Register() succeeded unexpectedly")
			}
		})
	}
}


func Test_stockRepository_GetAllPaginated(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		filter  pkg.PaginationFilter
		userID  *uint
		want    *pkg.PaginationReponse[domain.PopulatedStock]
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStockRepository(tt.db)
			got, gotErr := r.GetAllPaginated(context.Background(), tt.filter, tt.userID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetAllPaginated() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetAllPaginated() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetAllPaginated() = %v, want %v", got, tt.want)
			}
		})
	}
}

