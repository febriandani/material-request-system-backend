package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	department "github.com/febriandani/material-request-system-backend/internal/domain/master/department"
)

func Register(r *gin.Engine, db *sqlx.DB) {
	api := r.Group("/api/v1")

	// health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// master data
	master := api.Group("/master")
	{
		deptRepo := department.NewRepository(db)
		deptService := department.NewService(deptRepo)
		deptHandler := department.NewHandler(deptService)

		master.GET("/departments", deptHandler.GetAll)
	}
}
