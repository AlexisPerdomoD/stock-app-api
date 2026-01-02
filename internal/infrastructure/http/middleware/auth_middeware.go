package middleware

import (
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserSessionMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {

		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"message": "Session token is malformed or empty"})
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
