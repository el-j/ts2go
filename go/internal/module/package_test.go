package module

import (
	"path/filepath"
	"testing"
)

func TestPackageMapper_GetPackageName(t *testing.T) {
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "root index file",
			filePath: "/project/index.ts",
			expected: "main",
		},
		{
			name:     "src index file",
			filePath: "/project/src/index.ts",
			expected: "main",
		},
		{
			name:     "nested file in models",
			filePath: "/project/src/models/user.ts",
			expected: "models",
		},
		{
			name:     "nested file in services",
			filePath: "/project/src/services/api.ts",
			expected: "services",
		},
		{
			name:     "deeply nested file",
			filePath: "/project/src/utils/helpers/strings.ts",
			expected: "helpers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GetPackageName(tt.filePath)
			if result != tt.expected {
				t.Errorf("GetPackageName(%s) = %s, want %s", tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestPackageMapper_GetPackagePath(t *testing.T) {
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "root file",
			filePath: "/project/index.ts",
			expected: "github.com/user/project",
		},
		{
			name:     "src file",
			filePath: "/project/src/main.ts",
			expected: "github.com/user/project",
		},
		{
			name:     "nested file in models",
			filePath: "/project/src/models/user.ts",
			expected: "github.com/user/project/models",
		},
		{
			name:     "nested file in services",
			filePath: "/project/src/services/api.ts",
			expected: "github.com/user/project/services",
		},
		{
			name:     "deeply nested file",
			filePath: "/project/src/utils/helpers/strings.ts",
			expected: "github.com/user/project/utils/helpers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GetPackagePath(tt.filePath)
			if result != tt.expected {
				t.Errorf("GetPackagePath(%s) = %s, want %s", tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestPackageMapper_GetOutputPath(t *testing.T) {
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	tests := []struct {
		name      string
		filePath  string
		outputDir string
		expected  string
	}{
		{
			name:      "root index file becomes main.go",
			filePath:  "/project/index.ts",
			outputDir: "/output",
			expected:  "/output/main.go",
		},
		{
			name:      "src index file becomes main.go",
			filePath:  "/project/src/index.ts",
			outputDir: "/output",
			expected:  "/output/src/main.go",
		},
		{
			name:      "nested file preserves structure",
			filePath:  "/project/src/models/user.ts",
			outputDir: "/output",
			expected:  "/output/src/models/user.go",
		},
		{
			name:      "file gets .go extension",
			filePath:  "/project/src/services/api.ts",
			outputDir: "/output",
			expected:  "/output/src/services/api.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GetOutputPath(tt.filePath, tt.outputDir)
			// Normalize paths for comparison
			result = filepath.ToSlash(result)
			expected := filepath.ToSlash(tt.expected)
			if result != expected {
				t.Errorf("GetOutputPath(%s, %s) = %s, want %s", tt.filePath, tt.outputDir, result, expected)
			}
		})
	}
}

func TestPackageMapper_ResolveImportPath(t *testing.T) {
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	tests := []struct {
		name         string
		currentFile  string
		importSource string
		expectedPath string
		expectedOk   bool
	}{
		{
			name:         "relative import same level",
			currentFile:  "/project/src/models/user.ts",
			importSource: "./post",
			expectedPath: "github.com/user/project/models",
			expectedOk:   true,
		},
		{
			name:         "relative import parent level",
			currentFile:  "/project/src/services/api.ts",
			importSource: "../models/user",
			expectedPath: "github.com/user/project/models",
			expectedOk:   true,
		},
		{
			name:         "npm package",
			currentFile:  "/project/src/index.ts",
			importSource: "axios",
			expectedPath: "",
			expectedOk:   false,
		},
		{
			name:         "npm scoped package",
			currentFile:  "/project/src/index.ts",
			importSource: "@types/node",
			expectedPath: "",
			expectedOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, ok := mapper.ResolveImportPath(tt.currentFile, tt.importSource)
			if ok != tt.expectedOk {
				t.Errorf("ResolveImportPath ok = %v, want %v", ok, tt.expectedOk)
			}
			if ok && path != tt.expectedPath {
				t.Errorf("ResolveImportPath path = %s, want %s", path, tt.expectedPath)
			}
		})
	}
}

func TestSanitizePackageName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase letters",
			input:    "models",
			expected: "models",
		},
		{
			name:     "uppercase letters",
			input:    "Models",
			expected: "models",
		},
		{
			name:     "with hyphen",
			input:    "user-models",
			expected: "user_models",
		},
		{
			name:     "with dot",
			input:    "user.models",
			expected: "user_models",
		},
		{
			name:     "starts with number",
			input:    "123models",
			expected: "pkg_123models",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "pkg",
		},
		{
			name:     "special characters",
			input:    "user@models#test",
			expected: "user_models_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizePackageName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizePackageName(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGeneratePackageStructure(t *testing.T) {
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	files := []string{
		"/project/src/index.ts",
		"/project/src/models/user.ts",
		"/project/src/models/post.ts",
		"/project/src/services/api.ts",
	}

	structure, err := GeneratePackageStructure(files, mapper)
	if err != nil {
		t.Fatalf("GeneratePackageStructure failed: %v", err)
	}

	// Check number of packages
	if len(structure.Packages) != 3 {
		t.Errorf("Expected 3 packages, got %d", len(structure.Packages))
	}

	// Check main package
	if structure.MainPackage != "github.com/user/project" {
		t.Errorf("Expected main package github.com/user/project, got %s", structure.MainPackage)
	}

	// Check entry point
	if structure.EntryPoint != "/project/src/index.ts" {
		t.Errorf("Expected entry point /project/src/index.ts, got %s", structure.EntryPoint)
	}

	// Check models package has 2 files
	modelsPackage := "github.com/user/project/models"
	if len(structure.Packages[modelsPackage]) != 2 {
		t.Errorf("Expected models package to have 2 files, got %d", len(structure.Packages[modelsPackage]))
	}
}

func TestGetPackageDeclaration(t *testing.T) {
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	tests := []struct {
		name         string
		filePath     string
		isEntryPoint bool
		expected     string
	}{
		{
			name:         "entry point",
			filePath:     "/project/src/index.ts",
			isEntryPoint: true,
			expected:     "package main",
		},
		{
			name:         "models package",
			filePath:     "/project/src/models/user.ts",
			isEntryPoint: false,
			expected:     "package models",
		},
		{
			name:         "services package",
			filePath:     "/project/src/services/api.ts",
			isEntryPoint: false,
			expected:     "package services",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPackageDeclaration(tt.filePath, mapper, tt.isEntryPoint)
			if result != tt.expected {
				t.Errorf("GetPackageDeclaration = %s, want %s", result, tt.expected)
			}
		})
	}
}
