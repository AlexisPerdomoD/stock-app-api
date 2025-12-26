package mappers

import (
	appmodels "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/gin-gonic/gin"
)

func MapRegisterUserDTO(c *gin.Context) (*appmodels.RegisterUserDTO, error) {
	user := &appmodels.RegisterUserDTO{}

	if err := c.ShouldBindBodyWithJSON(user); err != nil {
		return nil, err
	}

	return user, nil
}
