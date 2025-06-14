package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "Layanan REST API Email Sementara berjalan normal",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}
