package middleware

import (
	"net/http"
	"strings"

	"github.com/el-j/ts2go/saas/backend/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Middleware provides authentication middleware
type Middleware struct {
	service    *auth.Service
	repository *auth.Repository
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(service *auth.Service, repository *auth.Repository) *Middleware {
	return &Middleware{
		service:    service,
		repository: repository,
	}
}

// RequireAuth validates JWT token from Authorization header
func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.service.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// RequireAPIKey validates API key from X-API-Key header
func (m *Middleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Hash the provided key and look it up
		keyHash, err := m.service.HashAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process API key"})
			c.Abort()
			return
		}

		key, err := m.repository.GetAPIKeyByHash(c.Request.Context(), keyHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Update last used timestamp (async)
		go m.repository.UpdateAPIKeyLastUsed(c.Request.Context(), key.ID)

		// Store key info in context
		c.Set("api_key_id", key.ID)
		c.Set("user_id", key.UserID)
		c.Set("api_key_scopes", key.Scopes)

		c.Next()
	}
}

// OptionalAuth tries to authenticate but doesn't fail if no auth provided
func (m *Middleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := m.service.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetEmail retrieves email from context
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetUsername retrieves username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	u, ok := username.(string)
	return u, ok
}
