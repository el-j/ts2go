package domain

import (
	"path/filepath"
	"strings"
)

// File represents a source code file in the project
type File struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	Content      string `json:"content"`
	Language     string `json:"language"` // "typescript", "go", etc.
	RelativePath string `json:"relativePath"`
}

// NewFile creates a new File instance
func NewFile(path, content string) *File {
	name := filepath.Base(path)
	language := detectLanguage(name)

	return &File{
		Path:     path,
		Name:     name,
		Content:  content,
		Language: language,
	}
}

// IsTypeScript returns true if the file is a TypeScript file
func (f *File) IsTypeScript() bool {
	return f.Language == "typescript"
}

// IsGo returns true if the file is a Go file
func (f *File) IsGo() bool {
	return f.Language == "go"
}

// IsJavaScript returns true if the file is a JavaScript file
func (f *File) IsJavaScript() bool {
	return f.Language == "javascript"
}

// Extension returns the file extension
func (f *File) Extension() string {
	return filepath.Ext(f.Name)
}

// detectLanguage determines the programming language based on file extension
func detectLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".ts":
		return "typescript"
	case ".tsx":
		return "typescript"
	case ".js":
		return "javascript"
	case ".jsx":
		return "javascript"
	case ".go":
		return "go"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".md":
		return "markdown"
	default:
		return "unknown"
	}
}
