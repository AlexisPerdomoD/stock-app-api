package middleware

import (
	"net/http"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg/auth"
	"github.com/gin-gonic/gin"
)

func UserSessionMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {

		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"message": "Session is required"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := auth.ValidateSessionToken(token)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"name": "Unauthorized", "message": "Session is expired or invalid"})
		return
	}

	c.Set("user_id", userID)
	c.Next()
}
