package shared_service

import "github.com/alexisPerdomoD/stock-app-api/internal/domain/entity"

func GenerateSessionToken(user *entity.User) (string, error) {
	panic("implement me")
}

func ValidateSessionToken(token string) error {
	panic("implement me")
}
