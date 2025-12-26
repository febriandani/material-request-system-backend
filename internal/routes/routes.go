package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	config "github.com/febriandani/material-request-system-backend/internal/config"
	auth "github.com/febriandani/material-request-system-backend/internal/domain/master/auth"
	department "github.com/febriandani/material-request-system-backend/internal/domain/master/department"
	user "github.com/febriandani/material-request-system-backend/internal/domain/master/user"
	middleware "github.com/febriandani/material-request-system-backend/internal/middleware"
)

func Register(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {

	r.Use(middleware.Logger())

	api := r.Group("/api/v1")

	// =====================
	// AUTH (NO BASIC AUTH)
	// =====================
	{
		authRepo := auth.NewRepository(db)
		authService := auth.NewService(authRepo)
		authHandler := auth.NewHandler(authService, cfg.JWT.Secret, cfg.JWT.Expiration, cfg.JWT.ExpirationRefresh)

		authGroup := api.Group("/auth")
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		// Protected routes
		authGroup.Use(
			middleware.JWTAuth(
				cfg.JWT.Secret,
			),
		)
		authGroup.POST("/logout", authHandler.Logout)
	}

	// =====================
	// MASTER DATA (BASIC AUTH)
	// =====================
	master := api.Group("/master")
	master.Use(
		middleware.BasicAuth(
			cfg.Auth.Master.Username,
			cfg.Auth.Master.Password,
		),
	)

	{
		deptRepo := department.NewRepository(db)
		deptService := department.NewService(deptRepo)
		deptHandler := department.NewHandler(deptService)

		master.GET("/departments", deptHandler.GetAll)
	}

	// =====================
	// USER DATA (JWT AUTH)
	// =====================
	masterUser := api.Group("/user")
	masterUser.Use(
		middleware.JWTAuth(
			cfg.JWT.Secret,
		),
	)

	{
		userRepo := user.NewRepository(db)
		userService := user.NewService(userRepo)
		userHandler := user.NewHandler(userService)

		masterUser.GET("/", userHandler.GetAll)
		masterUser.GET("/approvers", userHandler.GetApprovers)
	}
}
