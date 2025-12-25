package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/febriandani/material-request-system-backend/pkg/response"
)

func BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, p, ok := c.Request.BasicAuth()
		if !ok || u != username || p != password {
			response.Error(
				c,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				"invalid basic auth credentials",
			)
			c.Abort()
			return
		}
		c.Next()
	}
}
