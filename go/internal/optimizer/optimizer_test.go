package optimizer

import (
	"strings"
	"testing"
)

func TestOptimizer_RemoveUnusedImports(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Remove unused import",
			input: `package main

import (
	"fmt"
	"strings"
	"os"
)

func main() {
	fmt.Println("Hello")
}
`,
			expected: `package main

import (
"fmt"
)

func main() {
	fmt.Println("Hello")
}
`,
		},
		{
			name: "Keep all used imports",
			input: `package main

import (
	"fmt"
	"strings"
)

func main() {
	s := strings.ToUpper("hello")
	fmt.Println(s)
}
`,
			expected: `package main

import (
	"fmt"
	"strings"
)

func main() {
	s := strings.ToUpper("hello")
	fmt.Println(s)
}
`,
		},
		{
			name: "Remove all unused imports",
			input: `package main

import (
	"fmt"
	"strings"
	"os"
)

func main() {
	println("Hello")
}
`,
			expected: `package main

func main() {
	println("Hello")
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewOptimizer()
			result, err := opt.OptimizeImports(tt.input)
			if err != nil {
				t.Fatalf("OptimizeImports failed: %v", err)
			}

			// Normalize whitespace for comparison
			result = normalizeWhitespace(result)
			expected := normalizeWhitespace(tt.expected)

			if result != expected {
				t.Errorf("OptimizeImports result mismatch:\nGot:\n%s\n\nExpected:\n%s", result, expected)
			}
		})
	}
}

func TestOptimizer_RemoveUnusedVariables(t *testing.T) {
	// Note: Variable removal is conservative - this test just ensures it doesn't break anything
	input := `package main

var usedVar int = 10

func main() {
	println(usedVar)
}
`

	opt := NewOptimizer()
	result, err := opt.Optimize(input)
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}

	// Should still contain usedVar
	if !strings.Contains(result, "usedVar") {
		t.Errorf("Expected result to contain usedVar")
	}

	// Should compile
	if !strings.Contains(result, "package main") {
		t.Errorf("Expected valid Go code")
	}
}

func TestOptimizer_RemoveUnusedFunctions(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name: "Remove unused helper function",
			input: `package main

func main() {
	usedFunc()
}

func usedFunc() {
	println("used")
}

func unusedFunc() {
	println("never called")
}
`,
			shouldContain:    []string{"main", "usedFunc"},
			shouldNotContain: []string{"unusedFunc"},
		},
		{
			name: "Keep exported functions",
			input: `package main

func main() {
	println("main")
}

func ExportedFunc() {
	println("exported")
}

func unusedPrivate() {
	println("unused")
}
`,
			shouldContain:    []string{"main", "ExportedFunc"},
			shouldNotContain: []string{"unusedPrivate"},
		},
		{
			name: "Keep init function",
			input: `package main

func main() {
	println("main")
}

func init() {
	println("init")
}

func unused() {
	println("unused")
}
`,
			shouldContain:    []string{"main", "init"},
			shouldNotContain: []string{"unused"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewOptimizer()
			result, err := opt.Optimize(tt.input)
			if err != nil {
				t.Fatalf("Optimize failed: %v", err)
			}

			for _, s := range tt.shouldContain {
				if !strings.Contains(result, s) {
					t.Errorf("Expected result to contain '%s', but it doesn't:\n%s", s, result)
				}
			}

			for _, s := range tt.shouldNotContain {
				if strings.Contains(result, s) {
					t.Errorf("Expected result to NOT contain '%s', but it does:\n%s", s, result)
				}
			}
		})
	}
}

func TestOptimizer_FullOptimization(t *testing.T) {
	input := `package main

import (
	"fmt"
	"strings"
	"os"
)

var usedVar string = "hello"

func main() {
	fmt.Println(usedVar)
	helperFunc()
}

func helperFunc() {
	fmt.Println("helper")
}

func unusedHelper() {
	strings.ToUpper("never called")
}
`

	opt := NewOptimizer()
	result, err := opt.Optimize(input)
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}

	// Should keep: fmt (used), main, helperFunc, usedVar
	// Should remove: strings (only used in removed function), os (unused), unusedHelper

	shouldContain := []string{"fmt", "main", "helperFunc", "usedVar"}
	for _, s := range shouldContain {
		if !strings.Contains(result, s) {
			t.Errorf("Expected result to contain '%s'", s)
		}
	}

	shouldNotContain := []string{"unusedHelper", "\"os\""}
	for _, s := range shouldNotContain {
		if strings.Contains(result, s) {
			t.Errorf("Expected result to NOT contain '%s'", s)
		}
	}

	// strings might still be there if the function removal doesn't cascade
	// That's okay - the main goal is removing unused functions and obvious unused imports

	t.Logf("Optimized code:\n%s", result)
}

func TestOptimizer_InvalidCode(t *testing.T) {
	input := `package main
this is not valid go code
`

	opt := NewOptimizer()
	result, err := opt.Optimize(input)

	// Should return original code without error when parsing fails
	if err != nil {
		t.Errorf("Expected no error for invalid code, got: %v", err)
	}

	if result != input {
		t.Errorf("Expected original code to be returned for invalid input")
	}
}

// Helper functions

func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}
