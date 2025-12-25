package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"data": data,
		"meta": gin.H{
			"message":   "success",
			"timestamp": time.Now().UTC(),
		},
	})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC(),
		},
	})
}
