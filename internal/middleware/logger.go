package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Olamigokeolowo/projectflow-backend/internal/metrics"
)

type logEntry struct {
	RequestID  string  `json:"request_id"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Status     int     `json:"status"`
	DurationMs float64 `json:"duration_ms"`
}

func Logger(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		requestID, _ := c.Get("request_id")
		status := c.Writer.Status()

		m.Record(status, duration)

		entry := logEntry{
			RequestID:  requestID.(string),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Status:     status,
			DurationMs: float64(duration.Microseconds()) / 1000.0,
		}

		data, _ := json.Marshal(entry)
		log.Println(string(data))
	}
}