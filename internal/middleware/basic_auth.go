package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/febriandani/material-request-system-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header is required",
			})
			return
		}

		// 2. Validasi format Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Authorization header format",
			})
			return
		}

		tokenString := parts[1]

		// 3. Parse & validasi JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token",
			})
			return
		}

		// 4. Ambil claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
			})
			return
		}

		// 5. (Optional tapi recommended) set ke context
		c.Set("user_id", claims["sub"])
		c.Set("role", claims["role"])
		c.Set("claims", claims)

		// 6. Lanjut ke handler
		c.Next()
	}
}
