package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu             sync.RWMutex
	totalRequests  int
	statusCounts   map[int]int
	totalDuration  time.Duration
}

func New() *Metrics {
	return &Metrics{
		statusCounts: make(map[int]int),
	}
}

func (m *Metrics) Record(status int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	m.statusCounts[status]++
	m.totalDuration += duration
}

func (m *Metrics) Snapshot() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgDurationMs float64
	if m.totalRequests > 0 {
		avgDurationMs = float64(m.totalDuration.Microseconds()) / 1000.0 / float64(m.totalRequests)
	}

	errorCount := 0
	for status, count := range m.statusCounts {
		if status >= 400 {
			errorCount += count
		}
	}

	var errorRate float64
	if m.totalRequests > 0 {
		errorRate = float64(errorCount) / float64(m.totalRequests) * 100
	}

	return map[string]interface{}{
		"total_requests":   m.totalRequests,
		"status_counts":    m.statusCounts,
		"avg_duration_ms":  avgDurationMs,
		"error_rate_pct":   errorRate,
	}
}