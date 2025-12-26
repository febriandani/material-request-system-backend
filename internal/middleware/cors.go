package middleware

import (
	"time"

	"github.com/febriandani/material-request-system-backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(cfg config.Cors) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     cfg.Methods,
		AllowHeaders:     cfg.Headers,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
