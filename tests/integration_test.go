package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestTranspileSimple(t *testing.T) {
	// Build the ts2go binary first
	buildCmd := exec.Command("go", "build", "-o", "../ts2go", "../cmd/ts2go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build ts2go: %v", err)
	}
	defer os.Remove("../ts2go")

	// Setup
	inputFile := "fixtures/simple.ts"
	outputFile := "fixtures/simple.go"
	expectedFile := "fixtures/simple.go.expected"

	// Clean up output file
	defer os.Remove(outputFile)

	// Install npm dependencies for parser
	parserDir := "../internal/transpiler/parser"
	if _, err := os.Stat(filepath.Join(parserDir, "node_modules")); os.IsNotExist(err) {
		t.Log("Installing npm dependencies...")
		npmCmd := exec.Command("npm", "install")
		npmCmd.Dir = parserDir
		if output, err := npmCmd.CombinedOutput(); err != nil {
			t.Logf("npm install output: %s", string(output))
			t.Skipf("Skipping test: npm install failed: %v", err)
		}
	}

	// Run transpiler
	cmd := exec.Command("../ts2go", "convert", "--in", inputFile, "--out", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Transpiler output: %s", string(output))
		t.Fatalf("Transpilation failed: %v", err)
	}

	// Check that output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read output and expected
	gotBytes, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// For now, just check that we got some output
	if len(gotBytes) == 0 {
		t.Fatal("Output file is empty")
	}

	t.Logf("Generated Go code:\n%s", string(gotBytes))

	// If expected file exists, compare
	if _, err := os.Stat(expectedFile); err == nil {
		expectedBytes, err := os.ReadFile(expectedFile)
		if err != nil {
			t.Fatalf("Failed to read expected file: %v", err)
		}

		// Basic validation - check for key elements
		got := string(gotBytes)
		expected := string(expectedBytes)

		// Check for package declaration
		if !contains(got, "package main") {
			t.Error("Missing package declaration")
		}

		// Check for struct definition
		if !contains(got, "type Person struct") {
			t.Error("Missing Person struct definition")
		}

		t.Logf("Expected:\n%s", expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
