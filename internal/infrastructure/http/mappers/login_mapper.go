package mappers

import (
	appmodels "github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/gin-gonic/gin"
)

func MapUserLoginDTO(c *gin.Context) (*appmodels.UserLoginDTO, error) {
	credentials := &appmodels.UserLoginDTO{}
	if err := c.ShouldBindBodyWithJSON(credentials); err != nil {
		return nil, err
	}

	return credentials, nil
}
