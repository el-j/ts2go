package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/el-j/ts2go/packages/core/internal/analyzer"
	"github.com/el-j/ts2go/packages/core/internal/mapper"
	"github.com/el-j/ts2go/packages/core/internal/project"
	"github.com/el-j/ts2go/packages/core/internal/transpiler"
)

// TranspileCommand implements the 'transpile' subcommand
func TranspileCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ts2go transpile <project-directory> [--out <output-dir>] [--verbose] [--quiet] [--watch]")
	}

	projectDir := args[0]
	outputDir := "output"
	progressLevel := ProgressNormal
	watchMode := false

	// Parse command line options
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--out", "-o":
			if i+1 < len(args) {
				outputDir = args[i+1]
				i++
			}
		case "--verbose", "-v":
			progressLevel = ProgressVerbose
		case "--quiet", "-q":
			progressLevel = ProgressQuiet
		case "--watch", "-w":
			watchMode = true
		}
	}

	// If watch mode is enabled, use the watch implementation
	if watchMode {
		return WatchProject(projectDir, outputDir, progressLevel)
	}

	// Validate input directory
	info, err := os.Stat(projectDir)
	if err != nil {
		return fmt.Errorf("invalid project directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", projectDir)
	}

	// Convert paths to absolute
	projectDir, err = filepath.Abs(projectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	// Create progress reporter
	progress := NewProgressReporter(progressLevel)

	progress.Start("Transpiling TypeScript Project", 5)
	progress.Info("Source: %s", projectDir)
	progress.Info("Output: %s", outputDir)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Scan project
	progress.Step("Scanning project...")
	scanner := project.NewScanner()
	proj, err := scanner.ScanProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to scan project: %w", err)
	}

	progress.Success(fmt.Sprintf("Found %d TypeScript files", len(proj.Files)))

	// Load mappings
	mappingsPath := findMappingsFile()
	db, err := mapper.LoadMappings(mappingsPath)
	if err != nil {
		return fmt.Errorf("failed to load mappings: %w", err)
	}

	classifier := mapper.NewClassifier(db)

	// Analyze imports for all files
	progress.Step("Analyzing imports...")
	for i, file := range proj.Files {
		ast, err := transpiler.ParseTypeScript(file.Path)
		if err != nil {
			progress.Verbose("Warning: Failed to parse %s: %v", file.RelativePath, err)
			continue
		}

		imports, err := analyzer.AnalyzeImports(ast)
		if err != nil {
			progress.Verbose("Warning: Failed to analyze imports in %s: %v", file.RelativePath, err)
			continue
		}

		file.Imports = imports.Imports

		if progressLevel >= ProgressVerbose {
			progress.Verbose("Analyzed %s", file.RelativePath)
		} else if (i+1)%10 == 0 {
			progress.Info("Analyzed %d/%d files", i+1, len(proj.Files))
		}
	}
	progress.Success(fmt.Sprintf("Completed import analysis for %d files", len(proj.Files)))

	// Build dependency graph and resolve build order
	progress.Step("Resolving build order...")
	graph, err := proj.BuildDependencyGraph()
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Check for circular dependencies
	cycles := graph.DetectCircularDependencies()
	if len(cycles) > 0 {
		progress.Warning(fmt.Sprintf("%d circular dependency cycles detected", len(cycles)))
		for i, cycle := range cycles {
			if i < 3 {
				progress.Verbose("Cycle %d: %v", i+1, cycle)
			}
		}
		if len(cycles) > 3 {
			progress.Verbose("... and %d more cycles", len(cycles)-3)
		}
	}

	buildOrder, err := graph.ResolveBuildOrder()
	if err != nil {
		// If circular dependencies prevent build order, just use file order
		progress.Warning(fmt.Sprintf("Could not resolve build order: %v", err))
		progress.Info("Transpiling files in discovery order...")
		buildOrder = make([]string, len(proj.Files))
		for i, file := range proj.Files {
			buildOrder[i] = file.RelativePath
		}
	} else {
		progress.Success(fmt.Sprintf("Build order determined (%d files)", len(buildOrder)))
	}

	// Transpile files in build order
	progress.Step("Transpiling files...")
	successCount := 0
	errorCount := 0

	// Create progress bar in verbose mode
	var bar *ProgressBar
	if progressLevel == ProgressVerbose {
		bar = NewProgressBar(len(buildOrder), 40)
	}

	for i, relPath := range buildOrder {
		file := proj.GetFileByPath(relPath)
		if file == nil {
			continue
		}

		// Determine output path
		outputPath := filepath.Join(outputDir, file.RelativePath)
		outputPath = outputPath[:len(outputPath)-3] + ".go" // Change extension

		// Create output directory
		outputFileDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputFileDir, 0755); err != nil {
			progress.Error(fmt.Sprintf("%s: failed to create directory", file.RelativePath))
			errorCount++
			continue
		}

		// Transpile the file
		err := transpiler.Transpile(file.Path, outputPath)
		if err != nil {
			progress.Error(fmt.Sprintf("%s: %v", file.RelativePath, err))
			errorCount++
		} else {
			successCount++
			progress.Verbose("✓ %s", file.RelativePath)
		}

		// Update progress bar
		if bar != nil {
			bar.Increment()
		} else if (i+1)%5 == 0 || i == len(buildOrder)-1 {
			progress.Info("Progress: %d/%d files (%.0f%%)", successCount, len(buildOrder), float64(successCount)/float64(len(buildOrder))*100)
		}
	}

	// Analyze dependencies and generate go.mod
	progress.Step("Generating go.mod...")

	if proj.PackageJSON != nil {
		allDeps := proj.PackageJSON.GetAllDependencies()
		for _, dep := range allDeps {
			classification, err := classifier.ClassifyPackage(dep)
			if err != nil {
				continue
			}
			proj.Dependencies[dep] = classification
		}
	}

	// Analyze npm package imports from files
	if err := proj.AnalyzeDependencies(classifier); err != nil {
		progress.Warning(fmt.Sprintf("Failed to analyze dependencies: %v", err))
	}

	mg := project.NewModuleGenerator(proj, classifier)
	if err := mg.GenerateGoModule(outputDir); err != nil {
		progress.Warning(fmt.Sprintf("Failed to generate go.mod: %v", err))
	} else {
		progress.Success("go.mod generated")
	}

	// Generate summary
	progress.Complete("Transpilation Complete")

	if progressLevel >= ProgressNormal {
		fmt.Printf("  Total Files: %d\n", len(proj.Files))
		fmt.Printf("  Success: %d\n", successCount)
		fmt.Printf("  Errors: %d\n", errorCount)
		fmt.Printf("  Output: %s\n", outputDir)

		if proj.PackageJSON != nil {
			summary := mg.GetDependencySummary()
			fmt.Printf("\nDependencies:\n")
			fmt.Printf("  NPM Packages: %d\n", summary.TotalNPM)
			fmt.Printf("  Go Packages: %d\n", summary.TotalGo)
			fmt.Printf("  Supported: %d\n", summary.Supported)
			fmt.Printf("  Partial: %d\n", summary.Partial)
			fmt.Printf("  Unsupported: %d\n", summary.Unsupported)

			if summary.Unsupported > 0 && progressLevel >= ProgressVerbose {
				fmt.Println("\n⚠️  Some dependencies are unsupported:")
				for i, pkg := range summary.UnsupportedList {
					if i < 5 { // Show first 5
						fmt.Printf("  - %s\n", pkg)
					}
				}
				if len(summary.UnsupportedList) > 5 {
					fmt.Printf("  ... and %d more\n", len(summary.UnsupportedList)-5)
				}
			}
		}

		fmt.Println("\n✓ Project transpiled successfully")
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  cd %s\n", outputDir)
		fmt.Printf("  go mod tidy\n")
		fmt.Printf("  go build\n")
	}

	return nil
}

func separator() string {
	return "================================================"
}
