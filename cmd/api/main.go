package main

import (
	"log"

	"github.com/febriandani/material-request-system-backend/internal/config"
	"github.com/febriandani/material-request-system-backend/internal/database"
	"github.com/febriandani/material-request-system-backend/internal/middleware"
	"github.com/febriandani/material-request-system-backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(middleware.Cors(cfg.Cors))

	routes.Register(r, db, cfg)

	log.Println("server running on port", cfg.App.Port)
	r.Run(":" + cfg.App.Port)
}
