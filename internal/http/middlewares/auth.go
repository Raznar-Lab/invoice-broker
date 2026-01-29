package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
)

func AuthMiddleware(conf *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get token from Authorization header
		authHeader := c.GetHeader("Authorization")

		// 2. Expect "Bearer <token>" format
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Compare with SERVER_API_TOKEN from environment
		if token == "" || token != conf.Server.ApiToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized access: invalid or missing api token",
			})

			// Crucial: stop the execution of subsequent handlers
			c.Abort()
			return
		}

		// Continue to the next handler if authorized
		c.Next()
	}
}
