package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/ts2go/internal/analyzer"
	"github.com/yourusername/ts2go/internal/mapper"
)

func TestDetermineModuleName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		pkgJSON     *analyzer.PackageInfo
		expected    string
	}{
		{
			"with package.json",
			"my-project",
			&analyzer.PackageInfo{Name: "my-package"},
			"github.com/yourusername/my-package",
		},
		{
			"scoped package",
			"my-project",
			&analyzer.PackageInfo{Name: "@scope/package"},
			"github.com/yourusername/scope-package",
		},
		{
			"without package.json",
			"my-project",
			nil,
			"github.com/yourusername/my-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := &Project{
				Name:        tt.projectName,
				PackageJSON: tt.pkgJSON,
			}
			mg := NewModuleGenerator(project, nil)
			result := mg.determineModuleName()
			if result != tt.expected {
				t.Errorf("determineModuleName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDetermineGoVersion(t *testing.T) {
	project := &Project{Name: "test"}
	mg := NewModuleGenerator(project, nil)

	version := mg.determineGoVersion()
	if version == "" {
		t.Error("determineGoVersion() returned empty string")
	}

	// Should be a valid version format
	if !strings.HasPrefix(version, "1.") {
		t.Errorf("determineGoVersion() = %q, expected to start with '1.'", version)
	}
}

func TestCollectDependencies(t *testing.T) {
	project := &Project{
		Name: "test-project",
		Dependencies: map[string]*mapper.Classification{
			"axios": {
				Package:   "axios",
				Class:     mapper.PackageClassEquivalent,
				GoPackage: "github.com/go-resty/resty/v2",
			},
			"fs": {
				Package:   "fs",
				Class:     mapper.PackageClassBuiltin,
				GoPackage: "os",
			},
			"lodash": {
				Package:   "lodash",
				Class:     mapper.PackageClassUnsupported,
				GoPackage: "",
			},
		},
	}

	mg := NewModuleGenerator(project, nil)
	deps, err := mg.collectDependencies()
	if err != nil {
		t.Fatalf("collectDependencies() failed: %v", err)
	}

	// Should include axios but not fs (stdlib) or lodash (unsupported)
	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(deps))
	}

	if _, ok := deps["github.com/go-resty/resty/v2"]; !ok {
		t.Error("Expected axios mapping to be included")
	}

	if _, ok := deps["os"]; ok {
		t.Error("Standard library packages should not be included")
	}
}

func TestIsStdlibPackage(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		expected bool
	}{
		{"stdlib", "fmt", true},
		{"stdlib", "os", true},
		{"external", "github.com/foo/bar", false},
		{"external", "go.uber.org/zap", false},
		{"stdlib multi", "encoding/json", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isStdlibPackage(tt.pkg)
			if result != tt.expected {
				t.Errorf("isStdlibPackage(%q) = %v, want %v", tt.pkg, result, tt.expected)
			}
		})
	}
}

func TestFormatGoMod(t *testing.T) {
	project := &Project{Name: "test"}
	mg := NewModuleGenerator(project, nil)

	deps := map[string]string{
		"github.com/go-resty/resty/v2": "v2.11.0",
		"github.com/gin-gonic/gin":     "v1.9.1",
	}

	content := mg.formatGoMod("github.com/test/project", "1.21", deps)

	// Verify content
	if !strings.Contains(content, "module github.com/test/project") {
		t.Error("go.mod should contain module declaration")
	}

	if !strings.Contains(content, "go 1.21") {
		t.Error("go.mod should contain go version")
	}

	if !strings.Contains(content, "github.com/go-resty/resty/v2") {
		t.Error("go.mod should contain resty dependency")
	}

	if !strings.Contains(content, "github.com/gin-gonic/gin") {
		t.Error("go.mod should contain gin dependency")
	}

	if !strings.Contains(content, "require (") {
		t.Error("go.mod should have require block")
	}
}

func TestGenerateGoModule(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	project := &Project{
		Name: "test-project",
		PackageJSON: &analyzer.PackageInfo{
			Name: "my-app",
		},
		Dependencies: map[string]*mapper.Classification{
			"axios": {
				Package:   "axios",
				Class:     mapper.PackageClassEquivalent,
				GoPackage: "github.com/go-resty/resty/v2",
			},
		},
	}

	mg := NewModuleGenerator(project, nil)
	err := mg.GenerateGoModule(tmpDir)
	if err != nil {
		t.Fatalf("GenerateGoModule() failed: %v", err)
	}

	// Verify file was created
	goModPath := filepath.Join(tmpDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Fatal("go.mod file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "module github.com/yourusername/my-app") {
		t.Error("go.mod should contain correct module name")
	}

	if !strings.Contains(contentStr, "github.com/go-resty/resty/v2") {
		t.Error("go.mod should contain dependency")
	}
}

func TestGetDependencySummary(t *testing.T) {
	project := &Project{
		Name: "test-project",
		Dependencies: map[string]*mapper.Classification{
			"axios": {
				Package:   "axios",
				Class:     mapper.PackageClassEquivalent,
				GoPackage: "github.com/go-resty/resty/v2",
			},
			"fs": {
				Package:   "fs",
				Class:     mapper.PackageClassBuiltin,
				GoPackage: "os",
			},
			"express": {
				Package:   "express",
				Class:     mapper.PackageClassFramework,
				GoPackage: "github.com/gin-gonic/gin",
			},
			"lodash": {
				Package:   "lodash",
				Class:     mapper.PackageClassUnsupported,
				GoPackage: "",
			},
		},
	}

	mg := NewModuleGenerator(project, nil)
	summary := mg.GetDependencySummary()

	// Verify counts
	if summary.TotalNPM != 4 {
		t.Errorf("Expected 4 NPM packages, got %d", summary.TotalNPM)
	}

	if summary.Supported != 2 { // axios, fs
		t.Errorf("Expected 2 supported packages, got %d", summary.Supported)
	}

	if summary.Partial != 1 { // express
		t.Errorf("Expected 1 partial package, got %d", summary.Partial)
	}

	if summary.Unsupported != 1 { // lodash
		t.Errorf("Expected 1 unsupported package, got %d", summary.Unsupported)
	}

	// Verify unsupported list
	if len(summary.UnsupportedList) != 1 || summary.UnsupportedList[0] != "lodash" {
		t.Errorf("Expected unsupported list to contain lodash")
	}

	// Verify stdlib packages
	if len(summary.StdlibPackages) != 1 || summary.StdlibPackages[0] != "os" {
		t.Errorf("Expected stdlib packages to contain os")
	}

	// Verify external packages
	if len(summary.ExternalPackages) != 2 {
		t.Errorf("Expected 2 external packages, got %d", len(summary.ExternalPackages))
	}
}

func TestFormatSummary(t *testing.T) {
	summary := &DependencySummary{
		TotalNPM:         4,
		TotalGo:          3,
		Supported:        2,
		Partial:          1,
		Unsupported:      1,
		StdlibPackages:   []string{"os"},
		ExternalPackages: []string{"github.com/go-resty/resty/v2"},
		UnsupportedList:  []string{"lodash"},
	}

	output := summary.FormatSummary()

	// Verify output contains key information
	expectedStrings := []string{
		"Dependency Summary",
		"NPM Packages: 4",
		"Go Packages: 3",
		"Supported: 2",
		"Partial: 1",
		"Unsupported: 1",
		"lodash",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Summary output missing expected string: %s", expected)
		}
	}
}

func TestGetPackageVersion(t *testing.T) {
	project := &Project{Name: "test"}
	mg := NewModuleGenerator(project, nil)

	// Test known package
	version := mg.getPackageVersion("axios", "github.com/go-resty/resty/v2")
	if version == "" || version == "latest" {
		t.Errorf("Expected specific version for resty, got %s", version)
	}

	// Test unknown package
	version = mg.getPackageVersion("unknown", "github.com/unknown/package")
	if version != "latest" {
		t.Errorf("Expected 'latest' for unknown package, got %s", version)
	}
}
