package http

import (
	"strings"

	pkg "github.com/alexisPerdomoD/stock-app-api/internal/pkg/service"
	"github.com/gin-gonic/gin"
)

func UserSessionMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {

		c.AbortWithStatusJSON(401, gin.H{"message": "Session is required"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if err := pkg.ValidateSessionToken(token); err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Session is expired or invalid"})
		return
	}

	c.Next()
}
