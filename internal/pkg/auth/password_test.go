package auth_test

import (
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
)

var VALID_PASSWORD = []byte("12345678")
var INVALID_PASSWORD = []byte("87654321")
var INVALID_PASSWORD_TO_LONG = []byte("123456789012345678901234567890123456789012345678901234567890123456789012") // 72 chars
var VALID_HASH []byte

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password []byte
		wantErr  bool
	}{
		{
			name:     "Password is too long",
			password: INVALID_PASSWORD_TO_LONG,
			wantErr:  true,
		},
		{
			name:     "Hash password properly",
			password: VALID_PASSWORD,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPassword() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}

			if got == nil {
				t.Errorf("HashPassword() returned empty string")
			}

			VALID_HASH = got
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name     string
		password []byte
		hash     []byte
		wantErr  bool
	}{
		{
			name:     "Password is too long",
			password: INVALID_PASSWORD_TO_LONG,
			hash:     VALID_HASH,
			wantErr:  true,
		},
		{
			name:     "invalid password",
			password: INVALID_PASSWORD,
			hash:     VALID_HASH,
			wantErr:  true,
		},
		{
			name:     "valid password",
			password: VALID_PASSWORD,
			hash:     VALID_HASH,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := auth.VerifyPassword(tt.password, tt.hash)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("VerifyPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("VerifyPassword() succeeded unexpectedly")
			}
		})
	}
}
