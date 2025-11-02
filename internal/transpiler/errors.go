package transpiler

import (
	"fmt"
	"strings"
)

// TranspilationError represents a detailed error during transpilation
type TranspilationError struct {
	File       string   // Source file path
	Line       int      // Line number (1-based)
	Column     int      // Column number (1-based)
	Message    string   // Error message
	Code       string   // Error code (e.g., "UNSUPPORTED_FEATURE")
	Suggestion string   // Suggested fix
	Context    []string // Lines of code showing the error
}

// Error implements the error interface
func (e *TranspilationError) Error() string {
	var builder strings.Builder

	// Main error line
	builder.WriteString(fmt.Sprintf("\n%s:%d:%d: %s\n", e.File, e.Line, e.Column, e.Message))

	// Error code
	if e.Code != "" {
		builder.WriteString(fmt.Sprintf("  Code: %s\n", e.Code))
	}

	// Context
	if len(e.Context) > 0 {
		builder.WriteString("\n")
		for i, line := range e.Context {
			lineNum := e.Line - len(e.Context)/2 + i
			if lineNum == e.Line {
				// Highlight the error line
				builder.WriteString(fmt.Sprintf("  > %4d | %s\n", lineNum, line))
				// Add caret pointing to column
				if e.Column > 0 {
					builder.WriteString(fmt.Sprintf("  %s^ here\n", strings.Repeat(" ", e.Column+7)))
				}
			} else {
				builder.WriteString(fmt.Sprintf("    %4d | %s\n", lineNum, line))
			}
		}
		builder.WriteString("\n")
	}

	// Suggestion
	if e.Suggestion != "" {
		builder.WriteString(fmt.Sprintf("💡 Suggestion: %s\n", e.Suggestion))
	}

	return builder.String()
}

// NewTranspilationError creates a new transpilation error
func NewTranspilationError(file string, line, column int, message string) *TranspilationError {
	return &TranspilationError{
		File:    file,
		Line:    line,
		Column:  column,
		Message: message,
	}
}

// WithCode adds an error code
func (e *TranspilationError) WithCode(code string) *TranspilationError {
	e.Code = code
	return e
}

// WithSuggestion adds a suggested fix
func (e *TranspilationError) WithSuggestion(suggestion string) *TranspilationError {
	e.Suggestion = suggestion
	return e
}

// WithContext adds code context
func (e *TranspilationError) WithContext(lines []string) *TranspilationError {
	e.Context = lines
	return e
}

// Common error constructors

// UnsupportedFeatureError creates an error for unsupported TypeScript features
func UnsupportedFeatureError(file string, line, column int, feature string) *TranspilationError {
	return NewTranspilationError(file, line, column,
		fmt.Sprintf("Unsupported TypeScript feature: %s", feature)).
		WithCode("UNSUPPORTED_FEATURE").
		WithSuggestion("Check the documentation for supported features or file an issue")
}

// TypeConversionError creates an error for type conversion issues
func TypeConversionError(file string, line, column int, fromType, toType string) *TranspilationError {
	return NewTranspilationError(file, line, column,
		fmt.Sprintf("Cannot convert type '%s' to Go type '%s'", fromType, toType)).
		WithCode("TYPE_CONVERSION_ERROR").
		WithSuggestion("Try using a compatible type or add a custom type mapping")
}

// ParseError creates an error for parsing failures
func ParseError(file string, message string) *TranspilationError {
	return NewTranspilationError(file, 0, 0,
		fmt.Sprintf("Failed to parse TypeScript: %s", message)).
		WithCode("PARSE_ERROR").
		WithSuggestion("Check for syntax errors in your TypeScript code")
}

// ImportResolutionError creates an error for import resolution failures
func ImportResolutionError(file string, line, column int, importPath string) *TranspilationError {
	return NewTranspilationError(file, line, column,
		fmt.Sprintf("Cannot resolve import: %s", importPath)).
		WithCode("IMPORT_RESOLUTION_ERROR").
		WithSuggestion("Ensure the imported module exists and is in the dependency list")
}

// Helper function to extract context lines from source code
func ExtractContext(source string, lineNum int, contextLines int) []string {
	lines := strings.Split(source, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return nil
	}

	start := max(0, lineNum-contextLines-1)
	end := min(len(lines), lineNum+contextLines)

	return lines[start:end]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
