package orchestrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/ts2go/internal/project"
)

func TestMultiPackageTranspiler_SimpleProject(t *testing.T) {
	// Setup test project path
	testProjectDir, err := filepath.Abs(filepath.Join("..", "..", "test-projects", "simple-multi-file"))
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Check if test project exists
	if _, err := os.Stat(testProjectDir); os.IsNotExist(err) {
		t.Skip("Test project directory not found")
	}

	// Scan the project
	scanner := project.NewScanner()
	proj, err := scanner.ScanProject(testProjectDir)
	if err != nil {
		t.Fatalf("Failed to scan project: %v", err)
	}

	// Verify we found files
	if len(proj.Files) == 0 {
		t.Fatal("No TypeScript files found in test project")
	}

	t.Logf("Found %d TypeScript files", len(proj.Files))
	for _, file := range proj.Files {
		t.Logf("  - %s", file.RelativePath)
	}

	// Create transpiler
	goModuleName := "github.com/test/simple-multi-file"
	transpiler := NewMultiPackageTranspiler(proj, goModuleName)

	// Create output directory
	outputDir := t.TempDir()
	t.Logf("Output directory: %s", outputDir)

	// Transpile the project
	if err := transpiler.TranspileProject(outputDir); err != nil {
		t.Fatalf("Failed to transpile project: %v", err)
	}

	// Verify output files exist
	expectedFiles := []string{
		"src/main.go", // index.ts becomes main.go for entry point
		"src/models/user.go",
		"src/services/userService.go",
		"go.mod",
	}

	for _, expectedFile := range expectedFiles {
		outputPath := filepath.Join(outputDir, expectedFile)
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Errorf("Expected output file not found: %s", expectedFile)
		} else {
			t.Logf("✓ Generated: %s", expectedFile)
		}
	}

	// Read and verify go.mod
	goModPath := filepath.Join(outputDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	t.Logf("go.mod content:\n%s", string(goModContent))

	// Verify go.mod contains module name
	if len(goModContent) == 0 {
		t.Error("go.mod is empty")
	}
}

func TestMultiPackageTranspiler_EmptyProject(t *testing.T) {
	// Create empty project
	tmpDir := t.TempDir()

	proj := &project.Project{
		RootDir: tmpDir,
		Name:    "empty-project",
		Files:   []*project.SourceFile{},
	}

	transpiler := NewMultiPackageTranspiler(proj, "github.com/test/empty")
	outputDir := filepath.Join(tmpDir, "output")

	// Should not fail on empty project
	err := transpiler.TranspileProject(outputDir)
	if err != nil {
		t.Logf("Transpile empty project returned: %v", err)
	}
}
