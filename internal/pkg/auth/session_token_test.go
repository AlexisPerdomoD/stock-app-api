package auth_test

import (
	"os"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
)

// #nosec: G101 (CWE-798): Potential hardcoded credentials(they're not, but we want to test it)
var INVALID_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiaWF0IjoxNjQ5OTk0OTU0fQ.0x0"
var VALID_SESSION string

func TestGenerateSessionToken(t *testing.T) {
	if err := os.Setenv("SESSION_SECRET", "test-secret"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name:    "fail with nil user",
			user:    nil,
			wantErr: true,
		},
		{
			name:    "fail with invalid user id",
			user:    &domain.User{ID: 0},
			wantErr: true,
		},
		{
			name:    "success with valid user",
			user:    &domain.User{ID: 1},
			wantErr: false,
		},
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
			if got == "" {
				t.Errorf("GenerateSessionToken() returned empty string and error is nil")
			}
			VALID_SESSION = got
		})
	}
}

func TestValidateSessionToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    uint
		wantErr bool
	}{
		{
			name:    "fail with empty token",
			token:   "",
			want:    0,
			wantErr: true,
		},

		{
			name:    "fail with invalid token",
			token:   INVALID_TOKEN,
			want:    0,
			wantErr: true,
		},

		{
			name:    "success with valid token",
			token:   VALID_SESSION,
			want:    1,
			wantErr: false,
		},
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

			if got != tt.want {
				t.Errorf("ValidateSessionToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
