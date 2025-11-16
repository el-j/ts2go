package domain

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements error unwrapping
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Common error codes
const (
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeAlreadyExists       = "ALREADY_EXISTS"
	ErrCodeTranspilationFailed = "TRANSPILATION_FAILED"
	ErrCodeAnalysisFailed      = "ANALYSIS_FAILED"
	ErrCodeBuildFailed         = "BUILD_FAILED"
	ErrCodeRunFailed           = "RUN_FAILED"
	ErrCodeTestFailed          = "TEST_FAILED"
	ErrCodeGoNotFound          = "GO_NOT_FOUND"
	ErrCodeInvalidProject      = "INVALID_PROJECT"
	ErrCodeFileSystemError     = "FILESYSTEM_ERROR"
	ErrCodeStateSaveError      = "STATE_SAVE_ERROR"
	ErrCodeStateLoadError      = "STATE_LOAD_ERROR"
)

// NewDomainError creates a new domain error
func NewDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Error constructors for common scenarios
func ErrInvalidInput(message string) *DomainError {
	return &DomainError{Code: ErrCodeInvalidInput, Message: message}
}

// ErrNotFound is a sentinel error for not found resources
var ErrNotFound = &DomainError{
	Code:    ErrCodeNotFound,
	Message: "resource not found",
}

func ErrNotFoundWithMessage(resource string) *DomainError {
	return &DomainError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func ErrTranspilationFailed(message string, cause error) *DomainError {
	return &DomainError{
		Code:    ErrCodeTranspilationFailed,
		Message: message,
		Cause:   cause,
	}
}

func ErrGoNotFound() *DomainError {
	return &DomainError{
		Code:    ErrCodeGoNotFound,
		Message: "Go compiler not found in system PATH",
	}
}

// NewValidationError creates a validation error
func NewValidationError(format string, args ...interface{}) *DomainError {
	return &DomainError{
		Code:    ErrCodeInvalidInput,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrInvalidProject(path string) *DomainError {
	return &DomainError{
		Code:    ErrCodeInvalidProject,
		Message: fmt.Sprintf("invalid project at path: %s", path),
	}
}
