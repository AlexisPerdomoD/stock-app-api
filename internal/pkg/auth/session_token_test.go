package auth_test

import (
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
)

func TestGenerateSessionToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		user    *domain.User
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GenerateSessionToken(tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GenerateSessionToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GenerateSessionToken() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GenerateSessionToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateSessionToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		token   string
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateSessionToken(tt.token)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateSessionToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateSessionToken() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ValidateSessionToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
