package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePackageJSON(t *testing.T) {
	// Create temporary package.json for testing
	tmpDir := t.TempDir()
	packagePath := filepath.Join(tmpDir, "package.json")

	testPackageJSON := `{
  "name": "test-project",
  "version": "1.0.0",
  "description": "A test project",
  "main": "index.js",
  "type": "module",
  "scripts": {
    "start": "node index.js",
    "test": "jest"
  },
  "dependencies": {
    "axios": "^1.6.0",
    "lodash": "~4.17.21",
    "express": ">=4.18.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/node": "^20.0.0"
  },
  "author": "Test Author",
  "license": "MIT"
}`

	if err := os.WriteFile(packagePath, []byte(testPackageJSON), 0644); err != nil {
		t.Fatalf("Failed to create test package.json: %v", err)
	}

	// Test parsing
	pkg, err := ParsePackageJSON(packagePath)
	if err != nil {
		t.Fatalf("Failed to parse package.json: %v", err)
	}

	// Verify basic fields
	if pkg.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", pkg.Name)
	}

	if pkg.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", pkg.Version)
	}

	if pkg.Type != "module" {
		t.Errorf("Expected type 'module', got '%s'", pkg.Type)
	}

	// Verify dependencies
	if len(pkg.Dependencies) != 3 {
		t.Errorf("Expected 3 dependencies, got %d", len(pkg.Dependencies))
	}

	if pkg.Dependencies["axios"] != "^1.6.0" {
		t.Errorf("Expected axios version '^1.6.0', got '%s'", pkg.Dependencies["axios"])
	}

	// Verify devDependencies
	if len(pkg.DevDependencies) != 2 {
		t.Errorf("Expected 2 devDependencies, got %d", len(pkg.DevDependencies))
	}
}

func TestParsePackageJSONDirectory(t *testing.T) {
	// Create temporary directory with package.json
	tmpDir := t.TempDir()
	packagePath := filepath.Join(tmpDir, "package.json")

	testPackageJSON := `{
  "name": "dir-test",
  "version": "1.0.0"
}`

	if err := os.WriteFile(packagePath, []byte(testPackageJSON), 0644); err != nil {
		t.Fatalf("Failed to create test package.json: %v", err)
	}

	// Test parsing with directory path
	pkg, err := ParsePackageJSON(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse package.json from directory: %v", err)
	}

	if pkg.Name != "dir-test" {
		t.Errorf("Expected name 'dir-test', got '%s'", pkg.Name)
	}
}

func TestGetAllDependencies(t *testing.T) {
	pkg := &PackageInfo{
		Dependencies: map[string]string{
			"axios":  "^1.6.0",
			"lodash": "^4.17.21",
		},
		DevDependencies: map[string]string{
			"typescript": "^5.0.0",
			"jest":       "^29.0.0",
		},
		PeerDependencies: map[string]string{
			"react": "^18.0.0",
		},
	}

	all := pkg.GetAllDependencies()
	if len(all) != 5 {
		t.Errorf("Expected 5 total dependencies, got %d", len(all))
	}

	if all["axios"] != "^1.6.0" {
		t.Errorf("Expected axios in all dependencies")
	}

	if all["typescript"] != "^5.0.0" {
		t.Errorf("Expected typescript in all dependencies")
	}

	if all["react"] != "^18.0.0" {
		t.Errorf("Expected react in all dependencies")
	}
}

func TestGetProductionDependencies(t *testing.T) {
	pkg := &PackageInfo{
		Dependencies: map[string]string{
			"axios":  "^1.6.0",
			"lodash": "^4.17.21",
		},
		DevDependencies: map[string]string{
			"typescript": "^5.0.0",
		},
	}

	prod := pkg.GetProductionDependencies()
	if len(prod) != 2 {
		t.Errorf("Expected 2 production dependencies, got %d", len(prod))
	}

	if _, exists := prod["typescript"]; exists {
		t.Errorf("Dev dependency typescript should not be in production dependencies")
	}
}

func TestIsScopedPackage(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		expected bool
	}{
		{"scoped package", "@babel/core", true},
		{"scoped package 2", "@types/node", true},
		{"regular package", "axios", false},
		{"regular package 2", "lodash", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsScopedPackage(tt.pkg)
			if result != tt.expected {
				t.Errorf("IsScopedPackage('%s') = %v, expected %v", tt.pkg, result, tt.expected)
			}
		})
	}
}

func TestExtractScope(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		expected string
	}{
		{"babel scope", "@babel/core", "babel"},
		{"types scope", "@types/node", "types"},
		{"regular package", "axios", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractScope(tt.pkg)
			if result != tt.expected {
				t.Errorf("ExtractScope('%s') = '%s', expected '%s'", tt.pkg, result, tt.expected)
			}
		})
	}
}

func TestExtractPackageName(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		expected string
	}{
		{"scoped package", "@babel/core", "core"},
		{"scoped package 2", "@types/node", "node"},
		{"regular package", "axios", "axios"},
		{"regular package 2", "lodash", "lodash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPackageName(tt.pkg)
			if result != tt.expected {
				t.Errorf("ExtractPackageName('%s') = '%s', expected '%s'", tt.pkg, result, tt.expected)
			}
		})
	}
}

func TestIsBuiltinModule(t *testing.T) {
	tests := []struct {
		name     string
		module   string
		expected bool
	}{
		{"fs", "fs", true},
		{"path", "path", true},
		{"http", "http", true},
		{"axios", "axios", false},
		{"lodash", "lodash", false},
		{"node:fs", "node:fs", true},
		{"node:path", "node:path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBuiltinModule(tt.module)
			if result != tt.expected {
				t.Errorf("IsBuiltinModule('%s') = %v, expected %v", tt.module, result, tt.expected)
			}
		})
	}
}
