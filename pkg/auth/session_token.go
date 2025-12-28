package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/golang-jwt/jwt/v5"
)

type SessionClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type UserLike interface {
	GetID() uint64
}

func GenerateSessionToken(user UserLike) (token string, err error) {

	if user == nil || user.GetID() == 0 {
		return "", fmt.Errorf("[GenerateSessionToken] user is invalid or nil %+v", user)
	}

	secret := os.Getenv("SESSION_SECRET")

	if secret == "" {
		return "", fmt.Errorf("[GenerateSessionToken] secret is empty")
	}

	claims := &SessionClaims{
		UserID: user.GetID(),
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

func ValidateSessionToken(token string) (userID uint64, err error) {
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
