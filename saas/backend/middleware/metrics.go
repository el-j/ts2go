package middleware

import (
	"time"

	"github.com/el-j/ts2go/saas/backend/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		metrics.RecordHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			status,
			duration,
		)

		// Record rate limit violations
		if status == 429 {
			metrics.RateLimitExceeded.WithLabelValues(c.FullPath()).Inc()
		}
	}
}

// DatabaseMetricsMiddleware records database metrics
func DatabaseMetricsMiddleware(operation string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		metrics.DatabaseQueryDuration.WithLabelValues(operation).Observe(duration)
	}
}
