package auth

import (
	"log"
	"os"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateSessionToken(user *domain.User) (token string, err error) {

	if user == nil {
		log.Fatalln("[GenerateSessionToken] user is nil")
	}

	secret := os.Getenv("SESSION_SECRET")

	if secret == "" {
		log.Fatalln("[GenerateSessionToken] secret is empty")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	session := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return session.SignedString([]byte(secret))
}

func ValidateSessionToken(token string) (userID uint, err error) {
	secret := os.Getenv("SESSION_SECRET")

	session, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := session.Claims.(jwt.MapClaims)
	if !ok {
		return 0, pkg.Unauthorized("[session]: not claims found on session token")
	}


	return 0, nil
}
