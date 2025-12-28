package models

import (
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
)

type UserLoginDTO struct {
	Username string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

// GetPasswordBytesAndClean returns the password bytes and clean the password field
// to avoid leaking it in the logs
func (dto *UserLoginDTO) GetPasswordBytesAndClean() []byte {
	pwd := []byte(dto.Password)
	dto.Password = ""
	return pwd
}

type RegisterUserDTO struct {
	Username  string `json:"email" binding:"email,required"`
	Firstname string `json:"firstname" binding:"min=1"`
	Lastname  string `json:"lastname" binding:"min=1"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
}

func (dto *RegisterUserDTO) ToDomain() *domain.User {
	return &domain.User{
		Username:  dto.Username,
		Firstname: dto.Firstname,
		Lastname:  dto.Lastname,
		Active:    true,
	}
}

// GetPasswordBytesAndClean returns the password bytes and clean the password field
func (dto *RegisterUserDTO) GetPasswordBytesAndClean() []byte {
	pwd := []byte(dto.Password)
	dto.Password = ""
	return pwd
}

type UserView struct {
	ID        uint64    `json:"id,string"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

func (u UserView) GetID() uint64 {
	return u.ID
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
