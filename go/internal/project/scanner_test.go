package project

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanner_IsTypeScriptFile(t *testing.T) {
	scanner := NewScanner()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"TypeScript file", "test.ts", true},
		{"Declaration file", "test.d.ts", false},
		{"JavaScript file", "test.js", false},
		{"JSON file", "test.json", false},
		{"No extension", "test", false},
		{"TypeScript in path", "test.ts/file.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.isTypeScriptFile(tt.path)
			if result != tt.expected {
				t.Errorf("isTypeScriptFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestScanner_ShouldIgnore(t *testing.T) {
	scanner := NewScanner()

	tests := []struct {
		name     string
		relPath  string
		isDir    bool
		expected bool
	}{
		{"node_modules dir", "node_modules", true, true},
		{"file in node_modules", "node_modules/package/index.ts", false, true},
		{".git dir", ".git", true, true},
		{"dist dir", "dist", true, true},
		{"test file", "test.test.ts", false, true},
		{"spec file", "app.spec.ts", false, true},
		{"normal file", "src/app.ts", false, false},
		{"src dir", "src", true, false},
		{"nested in ignored", "node_modules/pkg/src/index.ts", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &mockFileInfo{isDir: tt.isDir}
			result := scanner.shouldIgnore(tt.relPath, info)
			if result != tt.expected {
				t.Errorf("shouldIgnore(%q) = %v, want %v", tt.relPath, result, tt.expected)
			}
		})
	}
}

func TestScanner_AddIgnorePattern(t *testing.T) {
	scanner := NewScanner()
	initialLen := len(scanner.ignorePatterns)

	scanner.AddIgnorePattern("*.backup.ts")

	if len(scanner.ignorePatterns) != initialLen+1 {
		t.Errorf("AddIgnorePattern did not add pattern")
	}

	if scanner.ignorePatterns[len(scanner.ignorePatterns)-1] != "*.backup.ts" {
		t.Errorf("AddIgnorePattern added wrong pattern")
	}
}

func TestScanner_ScanProject(t *testing.T) {
	// Create a temporary test project
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"src",
		"src/utils",
		"tests",
		"node_modules",
		"node_modules/package",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create test files
	files := map[string]string{
		"package.json":                  `{"name": "test-project", "version": "1.0.0", "dependencies": {"axios": "^1.0.0"}}`,
		"index.ts":                      `import { helper } from './src/utils/helper';`,
		"src/app.ts":                    `export class App {}`,
		"src/utils/helper.ts":           `export function helper() {}`,
		"tests/app.test.ts":             `import { App } from '../src/app';`,
		"node_modules/package/index.ts": `export const pkg = true;`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Scan the project
	scanner := NewScanner()
	project, err := scanner.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	// Verify project name
	if project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", project.Name)
	}

	// Verify file count (should exclude test files and node_modules)
	expectedFiles := 3 // index.ts, src/app.ts, src/utils/helper.ts
	if len(project.Files) != expectedFiles {
		t.Errorf("Expected %d files, got %d", expectedFiles, len(project.Files))
		for _, f := range project.Files {
			t.Logf("  Found: %s", f.RelativePath)
		}
	}

	// Verify package.json was loaded
	if project.PackageJSON == nil {
		t.Error("Expected package.json to be loaded")
	} else {
		if project.PackageJSON.Name != "test-project" {
			t.Errorf("Expected package name 'test-project', got '%s'", project.PackageJSON.Name)
		}
		if len(project.PackageJSON.Dependencies) != 1 {
			t.Errorf("Expected 1 dependency, got %d", len(project.PackageJSON.Dependencies))
		}
	}

	// Verify entry points
	if len(project.EntryPoints) == 0 {
		t.Error("Expected at least one entry point")
	}

	// Verify index.ts is an entry point
	indexFile := project.GetFileByPath("index.ts")
	if indexFile == nil {
		t.Error("Expected index.ts to be found")
	} else if !indexFile.IsEntry {
		t.Error("Expected index.ts to be marked as entry point")
	}
}

func TestScanner_ScanProjectWithoutPackageJSON(t *testing.T) {
	// Create a temporary test project without package.json
	tmpDir := t.TempDir()

	// Create a simple TypeScript file
	err := os.WriteFile(filepath.Join(tmpDir, "main.ts"), []byte("console.log('hello');"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Scan the project
	scanner := NewScanner()
	project, err := scanner.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	// Verify project name is derived from directory
	expectedName := filepath.Base(tmpDir)
	if project.Name != expectedName {
		t.Errorf("Expected project name '%s', got '%s'", expectedName, project.Name)
	}

	// Verify package.json is nil
	if project.PackageJSON != nil {
		t.Error("Expected package.json to be nil")
	}
}

func TestScanner_ScanProjectInvalidDirectory(t *testing.T) {
	scanner := NewScanner()

	// Try to scan non-existent directory
	_, err := scanner.ScanProject("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestScanner_IdentifyEntryPoints(t *testing.T) {
	// Create a temporary test project
	tmpDir := t.TempDir()

	// Create directory structure
	os.MkdirAll(filepath.Join(tmpDir, "src"), 0755)

	// Create test files
	files := []string{
		"index.ts",
		"main.ts",
		"src/helper.ts",
		"src/server.ts",
	}

	for _, file := range files {
		fullPath := filepath.Join(tmpDir, file)
		err := os.WriteFile(fullPath, []byte("export {};"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	// Scan the project
	scanner := NewScanner()
	project, err := scanner.ScanProject(tmpDir)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	// Verify entry points (index.ts, main.ts are in root, server.ts matches pattern)
	expectedEntries := map[string]bool{
		"index.ts":      true,
		"main.ts":       true,
		"src/server.ts": true,
	}

	if len(project.EntryPoints) != len(expectedEntries) {
		t.Errorf("Expected %d entry points, got %d", len(expectedEntries), len(project.EntryPoints))
	}

	for _, entry := range project.EntryPoints {
		if !expectedEntries[entry] {
			t.Errorf("Unexpected entry point: %s", entry)
		}
	}

	// Verify all entry files are marked correctly
	for entryPath := range expectedEntries {
		file := project.GetFileByPath(entryPath)
		if file == nil {
			t.Errorf("Entry point file not found: %s", entryPath)
		} else if !file.IsEntry {
			t.Errorf("File not marked as entry: %s", entryPath)
		}
	}
}

func TestProject_GetFileByPath(t *testing.T) {
	project := &Project{
		Files: []*SourceFile{
			{RelativePath: "index.ts", Path: "/path/to/index.ts"},
			{RelativePath: "src/app.ts", Path: "/path/to/src/app.ts"},
		},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", "index.ts", true},
		{"existing nested file", "src/app.ts", true},
		{"non-existent file", "missing.ts", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := project.GetFileByPath(tt.path)
			if (result != nil) != tt.expected {
				t.Errorf("GetFileByPath(%q) found=%v, want=%v", tt.path, result != nil, tt.expected)
			}
		})
	}
}

func TestProject_Summary(t *testing.T) {
	project := &Project{
		Name:    "test-project",
		RootDir: "/path/to/project",
		Files: []*SourceFile{
			{RelativePath: "index.ts"},
			{RelativePath: "src/app.ts"},
		},
		EntryPoints: []string{"index.ts"},
	}

	summary := project.Summary()

	// Check that summary contains expected information
	expectedStrings := []string{
		"test-project",
		"/path/to/project",
		"TypeScript Files: 2",
		"Entry Points: 1",
	}

	for _, expected := range expectedStrings {
		if !contains(summary, expected) {
			t.Errorf("Summary missing expected string: %s", expected)
		}
	}
}

// Helper types and functions

type mockFileInfo struct {
	isDir bool
}

func (m *mockFileInfo) Name() string       { return "" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
