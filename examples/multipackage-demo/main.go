package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/yourusername/ts2go/internal/orchestrator"
	"github.com/yourusername/ts2go/internal/project"
)

func main() {
	inputDir := flag.String("in", "", "Input TypeScript project directory")
	outputDir := flag.String("out", "", "Output directory for Go code")
	moduleName := flag.String("module", "github.com/example/project", "Go module name")
	flag.Parse()

	if *inputDir == "" || *outputDir == "" {
		log.Fatal("Usage: multipackage-demo -in <input-dir> -out <output-dir> [-module <module-name>]")
	}

	// Convert to absolute paths
	absInputDir, err := filepath.Abs(*inputDir)
	if err != nil {
		log.Fatalf("Failed to resolve input directory: %v", err)
	}

	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		log.Fatalf("Failed to resolve output directory: %v", err)
	}

	fmt.Printf("Transpiling TypeScript project...\n")
	fmt.Printf("  Input:  %s\n", absInputDir)
	fmt.Printf("  Output: %s\n", absOutputDir)
	fmt.Printf("  Module: %s\n\n", *moduleName)

	// Scan the project
	scanner := project.NewScanner()
	proj, err := scanner.ScanProject(absInputDir)
	if err != nil {
		log.Fatalf("Failed to scan project: %v", err)
	}

	fmt.Printf("Found %d TypeScript files\n\n", len(proj.Files))

	// Create transpiler
	transpiler := orchestrator.NewMultiPackageTranspiler(proj, *moduleName)

	// Transpile
	if err := transpiler.TranspileProject(absOutputDir); err != nil {
		log.Fatalf("Transpilation failed: %v", err)
	}

	fmt.Println("\n✓ Transpilation complete!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", absOutputDir)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go build\n")
}
