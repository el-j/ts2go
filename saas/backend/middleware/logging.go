package middleware

import (
	"time"

	"github.com/el-j/ts2go/saas/backend/logger"
	"github.com/gin-gonic/gin"
)

// RequestLogger logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get request ID
		requestID, _ := c.Get("request_id")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user ID if available
		userID, _ := c.Get("user_id")

		// Create log event
		event := logger.Log.Info().
			Str("request_id", requestID.(string)).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency_ms", latency).
			Str("client_ip", c.ClientIP())

		// Add user ID if authenticated
		if userID != nil {
			event = event.Str("user_id", userID.(string))
		}

		// Add error if present
		if len(c.Errors) > 0 {
			event = event.Str("error", c.Errors.String())
		}

		event.Msg("HTTP request")
	}
}
