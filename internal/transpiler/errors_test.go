package transpiler

import (
	"strings"
	"testing"
)

func TestTranspilationError_Error(t *testing.T) {
	err := NewTranspilationError("test.ts", 10, 5, "Test error message")
	errStr := err.Error()

	// Should contain file:line:column format
	if !strings.Contains(errStr, "test.ts:10:5") {
		t.Errorf("Error string should contain location, got: %s", errStr)
	}

	// Should contain message
	if !strings.Contains(errStr, "Test error message") {
		t.Errorf("Error string should contain message, got: %s", errStr)
	}
}

func TestTranspilationError_WithCode(t *testing.T) {
	err := NewTranspilationError("test.ts", 10, 5, "Test error").
		WithCode("TEST_ERROR")

	errStr := err.Error()
	if !strings.Contains(errStr, "TEST_ERROR") {
		t.Errorf("Error string should contain error code, got: %s", errStr)
	}
}

func TestTranspilationError_WithSuggestion(t *testing.T) {
	err := NewTranspilationError("test.ts", 10, 5, "Test error").
		WithSuggestion("Try this fix")

	errStr := err.Error()
	if !strings.Contains(errStr, "Suggestion") {
		t.Errorf("Error string should contain suggestion header, got: %s", errStr)
	}
	if !strings.Contains(errStr, "Try this fix") {
		t.Errorf("Error string should contain suggestion text, got: %s", errStr)
	}
}

func TestTranspilationError_WithContext(t *testing.T) {
	context := []string{
		"function test() {",
		"  const x = 1;",
		"  return x;",
		"}",
	}

	err := NewTranspilationError("test.ts", 2, 8, "Test error").
		WithContext(context)

	errStr := err.Error()

	// Should contain context lines
	for _, line := range context {
		if !strings.Contains(errStr, line) {
			t.Errorf("Error string should contain context line '%s', got: %s", line, errStr)
		}
	}

	// Should highlight error line
	if !strings.Contains(errStr, ">") {
		t.Errorf("Error string should highlight error line with >, got: %s", errStr)
	}
}

func TestUnsupportedFeatureError(t *testing.T) {
	err := UnsupportedFeatureError("test.ts", 5, 10, "decorators")

	errStr := err.Error()
	if !strings.Contains(errStr, "decorators") {
		t.Errorf("Should mention the unsupported feature")
	}
	if !strings.Contains(errStr, "UNSUPPORTED_FEATURE") {
		t.Errorf("Should have correct error code")
	}
	if !strings.Contains(errStr, "Suggestion") {
		t.Errorf("Should provide a suggestion")
	}
}

func TestTypeConversionError(t *testing.T) {
	err := TypeConversionError("test.ts", 8, 15, "symbol", "string")

	errStr := err.Error()
	if !strings.Contains(errStr, "symbol") || !strings.Contains(errStr, "string") {
		t.Errorf("Should mention both types")
	}
	if !strings.Contains(errStr, "TYPE_CONVERSION_ERROR") {
		t.Errorf("Should have correct error code")
	}
}

func TestParseError(t *testing.T) {
	err := ParseError("test.ts", "Unexpected token")

	errStr := err.Error()
	if !strings.Contains(errStr, "Unexpected token") {
		t.Errorf("Should contain parse error details")
	}
	if !strings.Contains(errStr, "PARSE_ERROR") {
		t.Errorf("Should have correct error code")
	}
}

func TestImportResolutionError(t *testing.T) {
	err := ImportResolutionError("test.ts", 1, 0, "./missing-module")

	errStr := err.Error()
	if !strings.Contains(errStr, "./missing-module") {
		t.Errorf("Should mention the import path")
	}
	if !strings.Contains(errStr, "IMPORT_RESOLUTION_ERROR") {
		t.Errorf("Should have correct error code")
	}
}

func TestExtractContext(t *testing.T) {
	source := `line 1
line 2
line 3
line 4
line 5`

	tests := []struct {
		name         string
		lineNum      int
		contextLines int
		wantLines    int
	}{
		{
			name:         "middle line",
			lineNum:      3,
			contextLines: 1,
			wantLines:    3, // lines 2, 3, 4
		},
		{
			name:         "first line",
			lineNum:      1,
			contextLines: 1,
			wantLines:    2, // lines 1, 2
		},
		{
			name:         "last line",
			lineNum:      5,
			contextLines: 1,
			wantLines:    2, // lines 4, 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractContext(source, tt.lineNum, tt.contextLines)
			if len(result) != tt.wantLines {
				t.Errorf("ExtractContext() returned %d lines, want %d", len(result), tt.wantLines)
			}
		})
	}
}

func TestExtractContext_OutOfBounds(t *testing.T) {
	source := "line 1\nline 2\nline 3"

	// Test with invalid line number
	result := ExtractContext(source, 10, 1)
	if result != nil {
		t.Errorf("ExtractContext() should return nil for out of bounds line number")
	}

	result = ExtractContext(source, 0, 1)
	if result != nil {
		t.Errorf("ExtractContext() should return nil for line number 0")
	}
}

func TestTranspilationError_FormattedOutput(t *testing.T) {
	// Test complete error formatting
	source := `function greet(name: string) {
  console.log("Hello " + name);
  return name.toUppercase(); // Error: should be toUpperCase
}`

	context := ExtractContext(source, 3, 1)
	err := NewTranspilationError("greet.ts", 3, 17, "Method 'toUppercase' does not exist").
		WithCode("METHOD_NOT_FOUND").
		WithSuggestion("Did you mean 'toUpperCase'?").
		WithContext(context)

	errStr := err.Error()

	// Verify all components are present
	checks := []string{
		"greet.ts:3:17",
		"Method 'toUppercase' does not exist",
		"METHOD_NOT_FOUND",
		"Did you mean 'toUpperCase'?",
		"console.log",
		"toUppercase",
	}

	for _, check := range checks {
		if !strings.Contains(errStr, check) {
			t.Errorf("Formatted error should contain '%s'\nGot: %s", check, errStr)
		}
	}

	t.Logf("Formatted error:\n%s", errStr)
}
