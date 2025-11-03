package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/el-j/ts2go/internal/transpiler"
)

// TestEndToEnd runs complete transpilation pipeline tests
func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name      string
		tsFile    string
		wantBuild bool
		wantRun   bool
	}{
		{
			name:      "simple transpilation",
			tsFile:    "../../tests/fixtures/simple.ts",
			wantBuild: true,
			wantRun:   false,
		},
		{
			name:      "advanced features",
			tsFile:    "../../tests/fixtures/advanced.ts",
			wantBuild: true,
			wantRun:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for output
			tmpDir, err := os.MkdirTemp("", "ts2go-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Get absolute paths
			tsPath, err := filepath.Abs(tt.tsFile)
			if err != nil {
				t.Fatalf("failed to resolve TS path: %v", err)
			}

			goFile := filepath.Join(tmpDir, "output.go")

			// Step 1: Transpile
			err = transpiler.Transpile(tsPath, goFile)
			if err != nil {
				t.Fatalf("transpilation failed: %v", err)
			}

			// Verify Go file was created
			if _, err := os.Stat(goFile); os.IsNotExist(err) {
				t.Fatalf("Go file was not created: %s", goFile)
			}

			// Step 2: Try to build (if expected)
			if tt.wantBuild {
				t.Run("build", func(t *testing.T) {
					// Initialize Go module in temp dir
					cmd := exec.Command("go", "mod", "init", "testmodule")
					cmd.Dir = tmpDir
					if err := cmd.Run(); err != nil {
						t.Logf("go mod init failed: %v (might be ok)", err)
					}

					// Try to build
					cmd = exec.Command("go", "build", "-o", "test", goFile)
					cmd.Dir = tmpDir
					output, err := cmd.CombinedOutput()
					if err != nil {
						t.Errorf("build failed: %v\nOutput: %s", err, string(output))
					}
				})
			}

			// Step 3: Try to run (if expected)
			if tt.wantRun {
				t.Run("run", func(t *testing.T) {
					cmd := exec.Command("go", "run", goFile)
					cmd.Dir = tmpDir
					output, err := cmd.CombinedOutput()
					if err != nil {
						t.Errorf("run failed: %v\nOutput: %s", err, string(output))
					}
					t.Logf("Program output: %s", string(output))
				})
			}
		})
	}
}

// TestTranspilerCoverage tests various transpiler features
func TestTranspilerCoverage(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "function declaration",
			input: `
function add(a: number, b: number): number {
	return a + b;
}`,
			want: "func Add(a float64, b float64) float64",
		},
		{
			name: "interface declaration",
			input: `
interface Person {
	name: string;
	age: number;
}`,
			want: "type Person struct",
		},
		{
			name: "const declaration",
			input: `
const PI = 3.14159;
const name = "test";`,
			want: "const PI = 3.14159",
		},
		{
			name: "arrow function",
			input: `
const double = (x: number) => x * 2;`,
			want: "func",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp files
			tmpDir, err := os.MkdirTemp("", "ts2go-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			tsFile := filepath.Join(tmpDir, "test.ts")
			goFile := filepath.Join(tmpDir, "test.go")

			// Write TypeScript
			if err := os.WriteFile(tsFile, []byte(tt.input), 0644); err != nil {
				t.Fatalf("failed to write TS file: %v", err)
			}

			// Transpile
			err = transpiler.Transpile(tsFile, goFile)
			if err != nil {
				t.Fatalf("transpilation failed: %v", err)
			}

			// Read output
			output, err := os.ReadFile(goFile)
			if err != nil {
				t.Fatalf("failed to read Go file: %v", err)
			}

			// Check if output contains expected string
			if !strings.Contains(string(output), tt.want) {
				t.Errorf("output does not contain %q\nGot:\n%s", tt.want, string(output))
			}
		})
	}
}

// TestOptimization tests that optimization is working
func TestOptimization(t *testing.T) {
	input := `
import { unused } from "lib";
import { used } from "lib2";

function main() {
	used();
}

function unusedFunc() {
	return 42;
}
`

	tmpDir, err := os.MkdirTemp("", "ts2go-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tsFile := filepath.Join(tmpDir, "test.ts")
	goFile := filepath.Join(tmpDir, "test.go")

	if err := os.WriteFile(tsFile, []byte(input), 0644); err != nil {
		t.Fatalf("failed to write TS file: %v", err)
	}

	// Transpile with optimization (default)
	err = transpiler.Transpile(tsFile, goFile)
	if err != nil {
		t.Fatalf("transpilation failed: %v", err)
	}

	output, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatalf("failed to read Go file: %v", err)
	}

	outputStr := string(output)

	// Should not contain unused imports/functions
	if strings.Contains(outputStr, "unusedFunc") {
		t.Errorf("optimization failed: unusedFunc still present in output")
	}
}

// BenchmarkTranspilation benchmarks transpilation performance
func BenchmarkTranspilation(b *testing.B) {
	// Create a moderate-sized TypeScript file
	input := `
interface User {
	id: number;
	name: string;
	email: string;
}

function getUser(id: number): User {
	return {
		id: id,
		name: "Test User",
		email: "test@example.com"
	};
}

function processUsers(users: User[]): number {
	return users.length;
}
`

	tmpDir, err := os.MkdirTemp("", "ts2go-bench-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tsFile := filepath.Join(tmpDir, "test.ts")
	if err := os.WriteFile(tsFile, []byte(input), 0644); err != nil {
		b.Fatalf("failed to write TS file: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		goFile := filepath.Join(tmpDir, "test.go")
		err := transpiler.Transpile(tsFile, goFile)
		if err != nil {
			b.Fatalf("transpilation failed: %v", err)
		}
		// Clean up for next iteration
		os.Remove(goFile)
	}
}
