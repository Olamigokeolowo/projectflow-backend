package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // handler runs here

		duration := time.Since(start)
		requestID, _ := c.Get("request_id")

		log.Printf(
			"[%s] %s %s | status=%d | duration=%s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}