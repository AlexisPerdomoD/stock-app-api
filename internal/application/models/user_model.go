package models

import (
	"errors"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type UserLoginDTO struct {
	Username string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`

	passwordConsumed bool
}

// GetPasswordBytesAndClean returns the password bytes and clean the password field
// this action can only be done once, else error is returned.
func (dto *UserLoginDTO) GetPasswordBytesAndClean() ([]byte, error) {
	if dto.passwordConsumed {
		return nil, errors.New("password field already consummed, ilegal action")
	}

	pwd := []byte(dto.Password)
	dto.Password = ""
	dto.passwordConsumed = true
	return pwd, nil
}

type RegisterUserDTO struct {
	Username  string `json:"email" binding:"email,required"`
	Firstname string `json:"firstname" binding:"min=1"`
	Lastname  string `json:"lastname" binding:"min=1"`
	Password  string `json:"password" binding:"required,min=8,max=72"`

	passwordConsumed bool
}

func (dto *RegisterUserDTO) ToDomain() *domain.User {
	return &domain.User{
		Username:  dto.Username,
		Firstname: dto.Firstname,
		Lastname:  dto.Lastname,
		Active:    true,
	}
}

// GetPasswordBytesAndClean returns the password bytes and clean the password field.
// this action can only be done once, else error is returned.
func (dto *RegisterUserDTO) GetPasswordBytesAndClean() ([]byte, error) {
	if dto.passwordConsumed {
		return nil, errors.New("password field already consummed, ilegal action")
	}

	pwd := []byte(dto.Password)
	dto.Password = ""
	dto.passwordConsumed = true
	return pwd, nil
}

type UserView struct {
	ID        uint64    `json:"id,string"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserView(user *domain.User) *UserView {
	return &UserView{
		ID:        user.ID,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}
}
