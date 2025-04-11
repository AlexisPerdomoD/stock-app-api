package shared_service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	if len(password) >= 72 {
		return "", errors.New("password is too long")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func VerifyPassword(password, hash string) error {
	if len(password) >= 72 {
		return errors.New("password is too long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return err
	}

	return nil
}
