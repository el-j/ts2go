package middleware

import (
	"fmt"
	"net/http"

	"github.com/el-j/ts2go/saas/backend/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrorHandler middleware handles errors consistently
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add request ID to context
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle AppError
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.StatusCode, errors.ErrorResponse{
					Error: *appErr,
				})
				return
			}

			// Handle generic errors
			fmt.Printf("Request %s error: %v\n", requestID, err)
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
				Error: *errors.ErrInternal,
			})
		}
	}
}

// Recovery middleware recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				fmt.Printf("Panic in request %v: %v\n", requestID, err)

				c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
					Error: *errors.ErrInternal,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
