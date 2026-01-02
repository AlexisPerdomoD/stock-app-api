package usecases_test

import (
	"context"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

func TestGetRecommendationsByStock_Execute_OK(t *testing.T) {
	repo := &FakeRecommendationRepository{
		getAllFn: func(ctx context.Context, _ pkg.PaginationFilter) (*pkg.PaginationResponse[domain.PopulatedRecommendation], error) {
			return &pkg.PaginationResponse[domain.PopulatedRecommendation]{
				Page:       1,
				PageSize:   10,
				TotalSize:  2,
				TotalPages: 1,
				Items: []domain.PopulatedRecommendation{
					{Recommendation: domain.Recommendation{ID: 1}},
					{Recommendation: domain.Recommendation{ID: 2}},
				},
			}, nil
		},
	}

	uc := usecases.NewGetRecommendations(repo)

	resp, err := uc.Execute(context.Background(), pkg.PaginationFilter{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	if resp.TotalSize != 2 || len(resp.Items) != 2 {
		t.Fatalf("unexpected pagination result: %+v", resp)
	}
}
