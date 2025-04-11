package pkg

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

func GenerateSessionToken(user *domain.User) (string, error) {
	panic("implement me")
}

func ValidateSessionToken(token string) error {
	panic("implement me")
}
