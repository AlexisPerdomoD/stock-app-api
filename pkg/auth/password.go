package auth

import (
	"errors"
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

func VerifyPassword(password []byte, hash []byte) (bool, error) {

	if hash == nil || password == nil {
		return false, errors.New("a nil hash or password was provided as argument")
	}

	if len(password) >= 72 {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword(hash, password)

	if err != nil {

		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
