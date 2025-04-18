package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/golang-jwt/jwt/v5"
)

type SessionClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateSessionToken(user *domain.User) (token string, err error) {

	if user == nil {
		log.Fatalln("[GenerateSessionToken] user is nil")
	}

	secret := os.Getenv("SESSION_SECRET")

	if secret == "" {
		log.Fatalln("[GenerateSessionToken] secret is empty")
	}

	claims := &SessionClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%d", time.Now().Unix()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "stock-app-api",
		},
	}

	session := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return session.SignedString([]byte(secret))
}

func ValidateSessionToken(token string) (userID uint, err error) {
	secret := os.Getenv("SESSION_SECRET")
	session := &SessionClaims{}
	payload, err := jwt.ParseWithClaims(token, session, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if !payload.Valid {
		return 0, pkg.Unauthorized("Session token invalid or expired")
	}

	return session.UserID, nil
}
