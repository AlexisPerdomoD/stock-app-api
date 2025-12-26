package models

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

type UserLoginDTO struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type RegisterUserDTO struct {
	Username  string `json:"email" binding:"email,required"`
	Firstname string `json:"firstname" binding:"min=1"`
	Lastname  string `json:"lastname" binding:"min=1"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
}

func (dto *RegisterUserDTO) ToDomain() *domain.User {
	return &domain.User{
		UserName: dto.Username,
		Password: []byte(dto.Password),
		Active:   true,
	}
}
