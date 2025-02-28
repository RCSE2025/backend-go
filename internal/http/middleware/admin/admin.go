package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("user_role") != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
