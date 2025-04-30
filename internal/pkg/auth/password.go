package auth

import (
	"errors"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"golang.org/x/crypto/bcrypt"
)

func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func HashPassword(password []byte) ([]byte, error) {

	if len(password) >= 72 {
		return nil, errors.New("password is too long")
	}

	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hashed, nil
}

func VerifyPassword(password []byte, hash []byte) error {
	if hash == nil {
		log.Panicln("[VerifyPassword]: a nil hash was provided")
	}

	defer ZeroBytes(hash)

	if password == nil {
		return pkg.Unauthorized("Invalid credentials")
	}

	defer ZeroBytes(password)

	if len(password) >= 72 {
		return pkg.Unauthorized("Invalid credentials")
	}

	return bcrypt.CompareHashAndPassword(hash, password)
}
