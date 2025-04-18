package dto

import (
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserDto struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func MapNewUserForm(c *gin.Context) (*domain.User, error) {
	args := &UserDto{}

	if err := c.ShouldBindBodyWithJSON(args); err != nil {
		return nil, err
	}

	user := &domain.User{
		UserName: args.Email,
		Password: args.Password,
		Active:   true,
	}

	return user, nil
}

func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())

			switch e.Tag() {
			case "required":
				errors[field] = "Este campo es obligatorio"
			case "email":
				errors[field] = "Debe ser un email válido"
			case "min":
				errors[field] = "Debe tener al menos " + e.Param() + " caracteres"
			case "max":
				errors[field] = "No puede tener más de " + e.Param() + " caracteres"
			default:
				errors[field] = "Campo inválido: " + e.Tag()
			}
		}
	} else {
		errors["general"] = err.Error()
	}

	return errors
}
