package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/el-j/ts2go/saas/backend/redis"
	"github.com/gin-gonic/gin"
)

// RateLimitConfig defines rate limit configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// RateLimiter provides rate limiting middleware
type RateLimiter struct {
	redis  *redis.Client
	config RateLimitConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		config: config,
	}
}

// Limit applies rate limiting based on IP address
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP or user ID if authenticated)
		identifier := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%v", userID)
		}

		// Check rate limit
		allowed, remaining, resetAt, err := rl.checkRateLimit(c, identifier)
		if err != nil {
			// On error, allow request but log the error
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		if !allowed {
			c.Header("Retry-After", strconv.FormatInt(int64(time.Until(resetAt).Seconds()), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Try again in %d seconds", int(time.Until(resetAt).Seconds())),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LimitByEndpoint applies different rate limits based on endpoint
func (rl *RateLimiter) LimitByEndpoint(endpointLimits map[string]RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()
		config, exists := endpointLimits[endpoint]
		if !exists {
			config = rl.config // Use default
		}

		identifier := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%v", userID)
		}

		key := fmt.Sprintf("ratelimit:%s:%s", endpoint, identifier)
		allowed, remaining, resetAt, err := rl.checkRateLimitWithConfig(c, key, config)
		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests for this endpoint. Try again in %d seconds", int(time.Until(resetAt).Seconds())),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if request is within rate limit
func (rl *RateLimiter) checkRateLimit(c *gin.Context, identifier string) (allowed bool, remaining int, resetAt time.Time, err error) {
	key := fmt.Sprintf("ratelimit:%s", identifier)
	return rl.checkRateLimitWithConfig(c, key, rl.config)
}

// checkRateLimitWithConfig checks rate limit with custom config
func (rl *RateLimiter) checkRateLimitWithConfig(c *gin.Context, key string, config RateLimitConfig) (allowed bool, remaining int, resetAt time.Time, err error) {
	ctx := c.Request.Context()
	now := time.Now()
	window := time.Minute

	// Increment counter
	count, err := rl.redis.IncrementCounter(ctx, key)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	// Set expiry on first request
	if count == 1 {
		if err := rl.redis.SetWithExpiry(ctx, key, count, window); err != nil {
			return false, 0, time.Time{}, err
		}
	}

	resetAt = now.Add(window)
	remaining = config.RequestsPerMinute - int(count)
	if remaining < 0 {
		remaining = 0
	}

	allowed = count <= int64(config.RequestsPerMinute)
	return allowed, remaining, resetAt, nil
}

// DefaultRateLimitConfig returns sensible defaults
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
	}
}

// StrictRateLimitConfig returns strict limits for sensitive endpoints
func StrictRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 10,
		BurstSize:         2,
	}
}

// GenerousRateLimitConfig returns generous limits for authenticated users
func GenerousRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 300,
		BurstSize:         50,
	}
}
