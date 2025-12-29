package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "login") || strings.Contains(c.Request.URL.Path, "register") {
			c.Next()
			return
		}
		authHeader := c.GetHeader("Seer-user-name")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User name header required"})
			c.Abort()
			return
		}

		// Example: Bearer token validation
		orgId := c.GetHeader("Seer-org-id")
		if orgId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ORG ID header required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
