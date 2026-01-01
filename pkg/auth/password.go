package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooLong = errors.New("password exceeds bcrypt 72-byte limit")
var ErrNilHashOrPassword = errors.New("a nil hash or password was provided as argument")

func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func HashPassword(password []byte) ([]byte, error) {

	if len(password) > 72 {
		return nil, ErrPasswordTooLong
	}

	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hashed, nil
}

func VerifyPassword(password []byte, hash []byte) (bool, error) {

	if len(hash) == 0 || len(password) == 0 {
		return false, ErrNilHashOrPassword
	}

	if len(password) > 72 {
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
