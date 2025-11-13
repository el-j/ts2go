package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// BuildCommand handles the build command - compiles transpiled Go code into a binary
func BuildCommand(args []string) error {
	var sourceDir string
	var outputPath string
	var buildFlags []string

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--source", "-s":
			if i+1 < len(args) {
				sourceDir = args[i+1]
				i++
			}
		case "--output", "-o":
			if i+1 < len(args) {
				outputPath = args[i+1]
				i++
			}
		case "--flags", "-f":
			if i+1 < len(args) {
				// Parse comma-separated build flags
				flags := strings.Split(args[i+1], ",")
				buildFlags = append(buildFlags, flags...)
				i++
			}
		}
	}

	// Validate required flags
	if sourceDir == "" {
		return fmt.Errorf("usage: ts2go build --source <source-dir> --output <binary-path> [--flags <build-flags>]")
	}

	// Default output path if not specified
	if outputPath == "" {
		baseName := filepath.Base(sourceDir)
		outputPath = filepath.Join(sourceDir, baseName)
	}

	// Verify source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	// Check for main package
	hasMain, err := hasMainPackage(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to check for main package: %w", err)
	}
	if !hasMain {
		return fmt.Errorf("no main package found in %s. Make sure your transpiled code has a main package", sourceDir)
	}

	fmt.Printf("Building Go project from %s...\n", sourceDir)
	fmt.Printf("Output binary: %s\n", outputPath)
	if len(buildFlags) > 0 {
		fmt.Printf("Build flags: %s\n", strings.Join(buildFlags, " "))
	}

	// Build the project
	startTime := time.Now()
	
	// Construct build command
	cmdArgs := []string{"build", "-o", outputPath}
	cmdArgs = append(cmdArgs, buildFlags...)
	
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = sourceDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute build
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	duration := time.Since(startTime)

	// Get binary size
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("failed to stat output binary: %w", err)
	}

	fmt.Printf("\n✅ Build successful!\n")
	fmt.Printf("   Binary: %s\n", outputPath)
	fmt.Printf("   Size: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	fmt.Printf("   Duration: %v\n", duration)

	return nil
}

// hasMainPackage checks if the directory contains a main package
func hasMainPackage(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".go") {
			filePath := filepath.Join(dir, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			if strings.Contains(string(content), "package main") {
				return true, nil
			}
		}
	}

	return false, nil
}
