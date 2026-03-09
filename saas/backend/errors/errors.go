package errors

import (
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}

// Common error codes
const (
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeRateLimit         = "RATE_LIMIT_EXCEEDED"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeServiceUnavail    = "SERVICE_UNAVAILABLE"
	ErrCodeInvalidToken      = "INVALID_TOKEN"
	ErrCodeExpiredToken      = "EXPIRED_TOKEN"
	ErrCodeInsufficientQuota = "INSUFFICIENT_QUOTA"
)

// Predefined errors
var (
	ErrValidation = &AppError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		StatusCode: http.StatusBadRequest,
	}

	ErrUnauthorized = &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    "Authentication required",
		StatusCode: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       ErrCodeForbidden,
		Message:    "Access denied",
		StatusCode: http.StatusForbidden,
	}

	ErrNotFound = &AppError{
		Code:       ErrCodeNotFound,
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
	}

	ErrConflict = &AppError{
		Code:       ErrCodeConflict,
		Message:    "Resource already exists",
		StatusCode: http.StatusConflict,
	}

	ErrRateLimit = &AppError{
		Code:       ErrCodeRateLimit,
		Message:    "Rate limit exceeded",
		StatusCode: http.StatusTooManyRequests,
	}

	ErrInternal = &AppError{
		Code:       ErrCodeInternal,
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrBadRequest = &AppError{
		Code:       ErrCodeBadRequest,
		Message:    "Bad request",
		StatusCode: http.StatusBadRequest,
	}

	ErrServiceUnavailable = &AppError{
		Code:       ErrCodeServiceUnavail,
		Message:    "Service temporarily unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}

	ErrInvalidToken = &AppError{
		Code:       ErrCodeInvalidToken,
		Message:    "Invalid authentication token",
		StatusCode: http.StatusUnauthorized,
	}

	ErrExpiredToken = &AppError{
		Code:       ErrCodeExpiredToken,
		Message:    "Authentication token expired",
		StatusCode: http.StatusUnauthorized,
	}

	ErrInsufficientQuota = &AppError{
		Code:       ErrCodeInsufficientQuota,
		Message:    "Insufficient quota",
		StatusCode: http.StatusPaymentRequired,
	}
)

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to an error
func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		Details:    details,
		StatusCode: e.StatusCode,
	}
}

// ErrorResponse represents HTTP error response
type ErrorResponse struct {
	Error AppError `json:"error"`
}
