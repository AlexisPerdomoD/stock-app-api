package auth

import "testing"

var validPassword = []byte("12345678")
var invalidPassword = []byte("87654321")
var invalidPasswordToLong = []byte("1234567890123456789012345678901234567890123456789012345678901234567890123") // 73 chars

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password []byte
		wantErr  error
	}{
		{
			name:     "password too long",
			password: invalidPasswordToLong,
			wantErr:  ErrPasswordTooLong,
		},
		{
			name:     "valid password",
			password: validPassword,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(hash) == 0 {
				t.Fatal("hash is empty")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	hash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password []byte
		hash     []byte
		wantOK   bool
		wantErr  error
	}{
		{
			name:     "password too long",
			password: invalidPasswordToLong,
			hash:     hash,
			wantOK:   false,
			wantErr:  nil,
		},
		{
			name:     "invalid password",
			password: invalidPassword,
			hash:     hash,
			wantOK:   false,
			wantErr:  nil,
		},
		{
			name:     "valid password",
			password: validPassword,
			hash:     hash,
			wantOK:   true,
			wantErr:  nil,
		},
		{
			name:     "nil password",
			password: nil,
			hash:     hash,
			wantOK:   false,
			wantErr:  ErrNilHashOrPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := VerifyPassword(tt.password, tt.hash)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ok != tt.wantOK {
				t.Fatalf("expected %v, got %v", tt.wantOK, ok)
			}
		})
	}
}
