package department

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/febriandani/material-request-system-backend/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *gin.Context) {
	departments, err := h.service.GetAll()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"failed to fetch departments",
		)
		return
	}

	response.Success(c, http.StatusOK, departments)
}
