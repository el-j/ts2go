package mapper

import (
	"strings"
	"testing"
)

func TestTransformCall(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	tests := []struct {
		name             string
		code             string
		expectTransform  string
		expectError      bool
		expectConfidence float64
	}{
		{
			name:             "axios.get with mapping",
			code:             "axios.get(url)",
			expectTransform:  "resty.R().Get(url)",
			expectError:      false,
			expectConfidence: 0.9,
		},
		{
			name:             "unknown package",
			code:             "unknown.method()",
			expectError:      true,
			expectConfidence: 0.0,
		},
		{
			name:             "method without API mapping",
			code:             "fs.readFile(path)",
			expectTransform:  "fs.ReadFile(path)",
			expectError:      false,
			expectConfidence: 0.5,
		},
		{
			name:        "invalid call expression",
			code:        "invalid code",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.TransformCall(tt.code)

			if tt.expectError {
				if result.Error == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if result.Error != nil {
				t.Errorf("Unexpected error: %v", result.Error)
				return
			}

			if result.Transformed != tt.expectTransform {
				t.Errorf("Expected transformation %s, got %s", tt.expectTransform, result.Transformed)
			}

			if result.Confidence < tt.expectConfidence-0.1 || result.Confidence > tt.expectConfidence+0.1 {
				t.Errorf("Expected confidence ~%.1f, got %.1f", tt.expectConfidence, result.Confidence)
			}

			t.Logf("Transformed: %s -> %s (confidence: %.1f)", result.Original, result.Transformed, result.Confidence)
		})
	}
}

func TestApplyAPIMapping(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	tests := []struct {
		name     string
		code     string
		goAPI    string
		expected string
	}{
		{
			name:     "simple method",
			code:     "axios.get(url)",
			goAPI:    "resty.Get",
			expected: "resty.Get(url)",
		},
		{
			name:     "chained method",
			code:     "axios.get(url)",
			goAPI:    "resty.R().Get",
			expected: "resty.R().Get(url)",
		},
		{
			name:     "multiple arguments",
			code:     "someFunc(arg1, arg2, arg3)",
			goAPI:    "GoFunc",
			expected: "GoFunc(arg1, arg2, arg3)",
		},
		{
			name:     "empty arguments",
			code:     "someFunc()",
			goAPI:    "GoFunc",
			expected: "GoFunc()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.applyAPIMapping(tt.code, tt.goAPI)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}

			t.Logf("Applied mapping: %s -> %s", tt.code, result)
		})
	}
}

func TestExtractArguments(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "single argument",
			code:     "func(arg)",
			expected: "arg",
		},
		{
			name:     "multiple arguments",
			code:     "func(arg1, arg2, arg3)",
			expected: "arg1, arg2, arg3",
		},
		{
			name:     "nested parentheses",
			code:     "func(nested(arg))",
			expected: "nested(arg)",
		},
		{
			name:     "empty arguments",
			code:     "func()",
			expected: "",
		},
		{
			name:     "complex nested",
			code:     "func(a, b(c, d), e)",
			expected: "a, b(c, d), e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.extractArguments(tt.code)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetPackageAlias(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	tests := []struct {
		goPackage string
		expected  string
	}{
		{
			goPackage: "github.com/go-resty/resty/v2",
			expected:  "resty",
		},
		{
			goPackage: "github.com/samber/lo",
			expected:  "lo",
		},
		{
			goPackage: "net/http",
			expected:  "http",
		},
		{
			goPackage: "simple",
			expected:  "simple",
		},
	}

	for _, tt := range tests {
		t.Run(tt.goPackage, func(t *testing.T) {
			result := transformer.getPackageAlias(tt.goPackage)

			if result != tt.expected {
				t.Errorf("Expected alias %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPatternMatcher(t *testing.T) {
	matcher := NewPatternMatcher()

	code := `
		const result = axios.get(url);
		const data = _.map(items, func);
		client.post().then(handler);
		fs.readFile(path);
	`

	calls := matcher.FindAPICalls(code)

	if len(calls) == 0 {
		t.Fatal("Expected to find API calls")
	}

	t.Logf("Found %d API calls:", len(calls))
	for _, call := range calls {
		t.Logf("  - %s", call)
	}

	// Check for specific patterns
	expectedPatterns := []string{"axios.get", "_.map", "fs.readFile"}
	for _, expected := range expectedPatterns {
		found := false
		for _, call := range calls {
			if strings.Contains(call, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find pattern %s", expected)
		}
	}
}

func TestCreateTransformationPlan(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	code := `
		const result = axios.get(url);
		const file = fs.readFile(path);
	`

	plan := transformer.CreateTransformationPlan(code)

	if len(plan.Transformations) == 0 {
		t.Fatal("Expected transformations in plan")
	}

	if plan.SuccessRate == 0.0 {
		t.Error("Expected non-zero success rate")
	}

	t.Logf("Transformation plan:\n%s", plan.Summary)

	for i, trans := range plan.Transformations {
		t.Logf("  %d. %s -> %s", i+1, trans.Original, trans.Transformed)
		if trans.Error != nil {
			t.Logf("     Error: %v", trans.Error)
		}
	}
}

func TestApplyTransformations(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	originalCode := `const result = axios.get(url);`

	plan := transformer.CreateTransformationPlan(originalCode)
	transformedCode := plan.ApplyTransformations()

	if transformedCode == originalCode {
		t.Error("Expected code to be transformed")
	}

	if !strings.Contains(transformedCode, "resty") {
		t.Error("Expected transformed code to contain 'resty'")
	}

	t.Logf("Original:    %s", originalCode)
	t.Logf("Transformed: %s", transformedCode)
}

func TestContextualTransformer(t *testing.T) {
	db := createTestMappingDB()
	ctrans := NewContextualTransformer(db)

	// Set context: client is an axios instance
	ctrans.SetVariableContext("client", "axios")

	code := "client.get(url)"
	result := ctrans.TransformWithContext(code, "client")

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}

	if !strings.Contains(result.Transformed, "resty") {
		t.Errorf("Expected transformed code to contain 'resty', got: %s", result.Transformed)
	}

	t.Logf("Contextual transform: %s -> %s", code, result.Transformed)
}

func TestContextualTransformerNoContext(t *testing.T) {
	db := createTestMappingDB()
	ctrans := NewContextualTransformer(db)

	code := "unknown.get(url)"
	result := ctrans.TransformWithContext(code, "unknown")

	if result.Error == nil {
		t.Error("Expected error for unknown context")
	}

	t.Logf("No context error: %v", result.Error)
}

func TestGenericTransform(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	mapping, _ := db.GetMapping("fs")

	code := "fs.readFile(path)"
	result := transformer.genericTransform(code, "fs", "readFile", mapping)

	// Should capitalize method name
	if !strings.Contains(result, "ReadFile") {
		t.Errorf("Expected capitalized method name, got: %s", result)
	}

	// Should preserve arguments
	if !strings.Contains(result, "path") {
		t.Errorf("Expected arguments to be preserved, got: %s", result)
	}

	t.Logf("Generic transform: %s -> %s", code, result)
}

func TestTransformMultipleCalls(t *testing.T) {
	db := createTestMappingDB()
	transformer := NewCallTransformer(db)

	patterns := []string{
		"axios.get(url)",
		"fs.readFile(path)",
		"axios.post(url, data)",
	}

	results := transformer.TransformMultipleCalls("", patterns)

	if len(results) != len(patterns) {
		t.Errorf("Expected %d results, got %d", len(patterns), len(results))
	}

	successCount := 0
	for pattern, result := range results {
		t.Logf("%s -> %s", pattern, result.Transformed)
		if result.Error == nil {
			successCount++
		}
	}

	if successCount == 0 {
		t.Error("Expected at least one successful transformation")
	}
}
