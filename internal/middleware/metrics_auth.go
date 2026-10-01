package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func MetricsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Metrics-Key")
		expected := os.Getenv("METRICS_KEY")

		if expected == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "metrics key not configured on server"})
			return
		}

		if key != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}