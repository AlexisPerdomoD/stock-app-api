package models

import "github.com/alexisPerdomoD/stock-app-api/internal/domain"

type UserLoginDTO struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type RegisterUserDTO struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (dto *RegisterUserDTO) ToDomain() *domain.User {
	return &domain.User{
		UserName: dto.Email,
		Password: []byte(dto.Password),
		Active:   true,
	}
}
