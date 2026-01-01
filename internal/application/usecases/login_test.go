package usecases_test

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	appmodel "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
)

func TestLogin_Execute_OK(t *testing.T) {
	ctx := context.Background()

	password := "secret123"
	hashed, err := auth.HashPassword([]byte(password))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user := &domain.User{
		ID:       1,
		Username: "alexis",
		Password: hashed,
	}

	uc := usecases.NewLogin(
		&FakeUserRepository{
			GetByUsernameWithPasswordFn: func(ctx context.Context, username string) (*domain.User, error) {
				return user, nil
			},
		},
		slog.Default(),
	)

	credentials := &appmodel.UserLoginDTO{
		Username: "alexis",
		Password: password,
	}

	result, err := uc.Execute(ctx, credentials)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ID != user.ID {
		t.Fatal("unexpected user id")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	uc := usecases.NewLogin(
		&FakeUserRepository{
			GetByUsernameWithPasswordFn: func(ctx context.Context, username string) (*domain.User, error) {
				return nil, nil
			},
		},
		slog.Default(),
	)

	credentials := &appmodel.UserLoginDTO{
		Username: "ghost",
		Password: "whatever",
	}

	_, err := uc.Execute(context.Background(), credentials)
	if err == nil {
		t.Fatal("expected error")
	}

	AssertApiErr(t, err, http.StatusUnauthorized, "Unauthorized")
}

func TestLogin_InvalidPassword(t *testing.T) {
	hashed, _ := auth.HashPassword([]byte("correct"))

	user := &domain.User{
		Username: "alexis",
		Password: hashed,
	}

	uc := usecases.NewLogin(
		&FakeUserRepository{
			GetByUsernameWithPasswordFn: func(ctx context.Context, username string) (*domain.User, error) {
				return user, nil
			},
		},
		slog.Default(),
	)

	credentials := &appmodel.UserLoginDTO{
		Username: "alexis",
		Password: "wrong",
	}

	_, err := uc.Execute(context.Background(), credentials)
	if err == nil {
		t.Fatal("expected error")
	}

	AssertApiErr(t, err, http.StatusUnauthorized, "Unauthorized")
}

func TestLogin_CredentialsAlreadyConsumed(t *testing.T) {
	uc := usecases.NewLogin(
		&FakeUserRepository{},
		slog.Default(),
	)

	credentials := &appmodel.UserLoginDTO{
		Username: "alexis",
		Password: "password",
	}

	_, err := credentials.GetPasswordBytesAndClean() // simulate consumption
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = uc.Execute(context.Background(), credentials)
	if err == nil {
		t.Fatal("expected error")
	}

	AssertApiErr(t, err, http.StatusInternalServerError, "Invalid State Error")
}
