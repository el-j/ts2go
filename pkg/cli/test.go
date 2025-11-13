package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TestCommand handles the test command - runs Go tests on transpiled code
func TestCommand(args []string) error {
	var sourceDir string
	var testFlags []string
	var verbose bool
	var coverage bool

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--source", "-s":
			if i+1 < len(args) {
				sourceDir = args[i+1]
				i++
			}
		case "--verbose", "-v":
			verbose = true
		case "--coverage", "-c":
			coverage = true
		case "--flags", "-f":
			if i+1 < len(args) {
				// Parse comma-separated test flags
				flags := strings.Split(args[i+1], ",")
				testFlags = append(testFlags, flags...)
				i++
			}
		}
	}

	// Validate required flags
	if sourceDir == "" {
		return fmt.Errorf("usage: ts2go test --source <source-dir> [--verbose] [--coverage] [--flags <test-flags>]")
	}

	// Verify source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	fmt.Printf("Running tests in %s...\n", sourceDir)
	if len(testFlags) > 0 {
		fmt.Printf("Test flags: %s\n", strings.Join(testFlags, " "))
	}

	// Build test command
	cmdArgs := []string{"test"}
	
	if verbose {
		cmdArgs = append(cmdArgs, "-v")
	}
	
	if coverage {
		cmdArgs = append(cmdArgs, "-cover", "-coverprofile=coverage.out")
	}
	
	cmdArgs = append(cmdArgs, testFlags...)
	cmdArgs = append(cmdArgs, "./...")

	startTime := time.Now()
	
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = sourceDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute tests
	err := cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("\n❌ Tests failed! Duration: %v\n", duration)
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Printf("\n✅ All tests passed! Duration: %v\n", duration)

	// Show coverage report if requested
	if coverage {
		coveragePath := filepath.Join(sourceDir, "coverage.out")
		if _, err := os.Stat(coveragePath); err == nil {
			fmt.Println("\nGenerating coverage report...")
			coverCmd := exec.Command("go", "tool", "cover", "-func=coverage.out")
			coverCmd.Dir = sourceDir
			coverCmd.Stdout = os.Stdout
			coverCmd.Stderr = os.Stderr
			_ = coverCmd.Run()
		}
	}

	return nil
}
