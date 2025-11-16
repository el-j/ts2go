package domain_test

import (
	"testing"

	"github.com/el-j/ts2go/core/domain"
)

func TestNewFile(t *testing.T) {
	path := "/path/to/file.ts"
	content := "const x: number = 42;"

	file := domain.NewFile(path, content)

	if file.Path != path {
		t.Errorf("Expected path %s, got %s", path, file.Path)
	}
	if file.Content != content {
		t.Errorf("Expected content %s, got %s", content, file.Content)
	}
	if file.Name != "file.ts" {
		t.Errorf("Expected name 'file.ts', got %s", file.Name)
	}
	if file.Language != "typescript" {
		t.Errorf("Expected language 'typescript', got %s", file.Language)
	}
}

func TestFileIsTypeScript(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"TypeScript file", "/path/test.ts", true},
		{"TypeScript React file", "/path/component.tsx", true},
		{"Go file", "/path/main.go", false},
		{"JavaScript file", "/path/script.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := domain.NewFile(tt.path, "content")
			if file.IsTypeScript() != tt.expected {
				t.Errorf("Expected IsTypeScript() = %v, got %v", tt.expected, file.IsTypeScript())
			}
		})
	}
}

func TestFileIsGo(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Go file", "/path/main.go", true},
		{"TypeScript file", "/path/test.ts", false},
		{"JavaScript file", "/path/script.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := domain.NewFile(tt.path, "content")
			if file.IsGo() != tt.expected {
				t.Errorf("Expected IsGo() = %v, got %v", tt.expected, file.IsGo())
			}
		})
	}
}

func TestFileIsJavaScript(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"JavaScript file", "/path/script.js", true},
		{"JavaScript JSX file", "/path/component.jsx", true},
		{"TypeScript file", "/path/test.ts", false},
		{"Go file", "/path/main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := domain.NewFile(tt.path, "content")
			if file.IsJavaScript() != tt.expected {
				t.Errorf("Expected IsJavaScript() = %v, got %v", tt.expected, file.IsJavaScript())
			}
		})
	}
}

func TestFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"TypeScript", "/path/test.ts", ".ts"},
		{"TypeScript React", "/path/component.tsx", ".tsx"},
		{"Go", "/path/main.go", ".go"},
		{"JavaScript", "/path/script.js", ".js"},
		{"No extension", "/path/README", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := domain.NewFile(tt.path, "content")
			if file.Extension() != tt.expected {
				t.Errorf("Expected extension %s, got %s", tt.expected, file.Extension())
			}
		})
	}
}
