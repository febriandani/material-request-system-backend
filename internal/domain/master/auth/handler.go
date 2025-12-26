package auth

import (
	"fmt"
	"net/http"

	"github.com/febriandani/material-request-system-backend/pkg/response"
	"github.com/febriandani/material-request-system-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service              Service
	jwtSecret            string
	jwtExpiration        int64
	jwtExpirationRefresh int64
}

func NewHandler(service Service, jwtSecret string, jwtExpiration, jwtExpirationRefresh int64) *Handler {
	return &Handler{service: service, jwtSecret: jwtSecret, jwtExpiration: jwtExpiration, jwtExpirationRefresh: jwtExpirationRefresh}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
		return
	}

	user, token, err := h.service.Login(req.Username, req.Password, h.jwtSecret, h.jwtExpiration, h.jwtExpirationRefresh)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
		return
	}

	// 1. Parse & validate refresh token
	claims, err := utils.ParseToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		return
	}

	if claims.Type != "refresh" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token type")
		return
	}

	// 2. Validate refresh token hash in DB
	ok, role, err := h.service.ValidateRefreshToken(
		claims.UserID,
		req.RefreshToken,
	)
	fmt.Println("debug:", ok, role, err, claims.UserID, req.RefreshToken)
	if err != nil || !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "refresh token revoked or expired")
		return
	}

	// 3. Generate new access token
	accessToken, err := utils.GenerateAccessToken(
		claims.UserID,
		role,
		h.jwtSecret,
		h.jwtExpiration,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate token")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request")
		return
	}

	// 1. Parse & validate refresh token
	claims, err := utils.ParseToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		return
	}

	err = h.service.Logout(
		claims.UserID,
		req.RefreshToken,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to logout")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
