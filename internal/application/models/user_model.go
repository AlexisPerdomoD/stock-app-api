package models

import (
	"errors"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type UserLoginDTO struct {
	Username string `json:"email" binding:"email,required" example:"alexis@perdomo.com"`
	Password string `json:"password" binding:"required" example:"123456789"`

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
	Username  string `json:"email" binding:"email,required" example:"alexis@perdomo.com"`
	Firstname string `json:"firstname" binding:"min=1" example:"Alexis"`
	Lastname  string `json:"lastname" binding:"min=1" example:"Perdomo"`
	Password  string `json:"password" binding:"required,min=8,max=72" example:"123456789"`

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
	ID        uint64    `json:"id,string" example:"1"`
	Username  string    `json:"username" example:"alexis@perdomo.com"`
	Firstname string    `json:"firstname" example:"Alexis"`
	Lastname  string    `json:"lastname" example:"Perdomo"`
	Active    bool      `json:"active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T12:00:00Z"`
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
