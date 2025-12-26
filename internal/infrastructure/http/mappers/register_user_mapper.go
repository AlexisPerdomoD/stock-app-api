package mappers

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/models"
	"github.com/gin-gonic/gin"
)

func MapRegisterUserDTO(c *gin.Context) (*models.RegisterUserDTO, error) {
	user := &models.RegisterUserDTO{}

	if err := c.ShouldBindBodyWithJSON(user); err != nil {
		return nil, err
	}

	return user, nil
}
