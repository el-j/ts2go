package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/ts2go/internal/analyzer"
	"github.com/yourusername/ts2go/internal/mapper"
	"github.com/yourusername/ts2go/internal/project"
	"github.com/yourusername/ts2go/internal/transpiler"
)

// AnalyzeCommand implements the 'analyze' subcommand
func AnalyzeCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ts2go analyze <project-directory>")
	}

	projectDir := args[0]

	// Validate directory
	info, err := os.Stat(projectDir)
	if err != nil {
		return fmt.Errorf("invalid project directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", projectDir)
	}

	// Convert to absolute path
	projectDir, err = filepath.Abs(projectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	fmt.Printf("Analyzing TypeScript project: %s\n\n", projectDir)

	// Scan project
	scanner := project.NewScanner()
	proj, err := scanner.ScanProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to scan project: %w", err)
	}

	fmt.Printf("Project: %s\n", proj.Name)
	fmt.Printf("TypeScript Files: %d\n", len(proj.Files))
	fmt.Printf("Entry Points: %d\n", len(proj.EntryPoints))

	if len(proj.EntryPoints) > 0 {
		fmt.Println("\nEntry Points:")
		for _, entry := range proj.EntryPoints {
			fmt.Printf("  - %s\n", entry)
		}
	}

	// Load mapping database
	mappingsPath := findMappingsFile()
	db, err := mapper.LoadMappings(mappingsPath)
	if err != nil {
		return fmt.Errorf("failed to load mappings: %w", err)
	}

	classifier := mapper.NewClassifier(db)

	// Analyze package.json dependencies if available
	if proj.PackageJSON != nil {
		fmt.Printf("\nPackage: %s@%s\n", proj.PackageJSON.Name, proj.PackageJSON.Version)

		allDeps := proj.PackageJSON.GetAllDependencies()
		if len(allDeps) > 0 {
			fmt.Printf("NPM Dependencies: %d\n", len(allDeps))

			// Classify each dependency
			supported := 0
			partial := 0
			unsupported := 0

			for _, dep := range allDeps {
				classification, err := classifier.ClassifyPackage(dep)
				if err != nil {
					continue
				}
				proj.Dependencies[dep] = classification

				switch classification.Status {
				case "supported":
					supported++
				case "partial":
					partial++
				case "unsupported":
					unsupported++
				}
			}

			fmt.Printf("\nDependency Support:\n")
			fmt.Printf("  Supported: %d (%.1f%%)\n", supported, float64(supported)/float64(len(allDeps))*100)
			fmt.Printf("  Partial: %d (%.1f%%)\n", partial, float64(partial)/float64(len(allDeps))*100)
			fmt.Printf("  Unsupported: %d (%.1f%%)\n", unsupported, float64(unsupported)/float64(len(allDeps))*100)

			// Show unsupported packages
			if unsupported > 0 {
				fmt.Println("\nUnsupported Packages:")
				for pkg, classification := range proj.Dependencies {
					if classification.Status == "unsupported" {
						fmt.Printf("  - %s", pkg)
						if classification.Suggestion != "" {
							fmt.Printf(" (suggestion: %s)", classification.Suggestion)
						}
						fmt.Println()
					}
				}
			}
		}
	}

	// Build dependency graph
	fmt.Println("\nAnalyzing file dependencies...")

	// Parse TypeScript files and analyze imports
	for i, file := range proj.Files {
		// Parse the file
		ast, err := transpiler.ParseTypeScript(file.Path)
		if err != nil {
			fmt.Printf("  Warning: Failed to parse %s: %v\n", file.RelativePath, err)
			continue
		}

		// Analyze imports
		imports, err := analyzer.AnalyzeImports(ast)
		if err != nil {
			fmt.Printf("  Warning: Failed to analyze imports in %s: %v\n", file.RelativePath, err)
			continue
		}

		file.Imports = imports.Imports

		// Show progress for large projects
		if (i+1)%10 == 0 {
			fmt.Printf("  Analyzed %d/%d files\n", i+1, len(proj.Files))
		}
	}

	// Build dependency graph
	graph, err := proj.BuildDependencyGraph()
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	stats := graph.GetDependencyStats()
	fmt.Printf("\nDependency Graph:\n")
	fmt.Printf("  Total Files: %d\n", stats["total_files"])
	fmt.Printf("  Entry Points: %d\n", stats["entry_points"])
	fmt.Printf("  Total Dependencies: %d\n", stats["total_deps"])
	fmt.Printf("  Average Dependencies: %.1f\n", stats["avg_deps"])

	if stats["max_deps"].(int) > 0 {
		fmt.Printf("  Most Dependencies: %s (%d)\n", stats["max_deps_file"], stats["max_deps"])
	}

	// Check for circular dependencies
	cycles := graph.DetectCircularDependencies()
	if len(cycles) > 0 {
		fmt.Printf("\n⚠️  Warning: %d circular dependency cycles detected:\n", len(cycles))
		for i, cycle := range cycles {
			if i < 3 { // Show first 3 cycles
				fmt.Printf("  %d. %v\n", i+1, cycle)
			}
		}
		if len(cycles) > 3 {
			fmt.Printf("  ... and %d more\n", len(cycles)-3)
		}
	} else {
		fmt.Println("\n✓ No circular dependencies detected")
	}

	// Generate module info
	if proj.PackageJSON != nil {
		mg := project.NewModuleGenerator(proj, classifier)
		summary := mg.GetDependencySummary()

		fmt.Printf("\nGo Module Info:\n")
		fmt.Printf("  Go Packages Needed: %d\n", summary.TotalGo)
		fmt.Printf("  - Standard Library: %d\n", len(summary.StdlibPackages))
		fmt.Printf("  - External: %d\n", len(summary.ExternalPackages))
		fmt.Printf("  - Runtime: %d\n", len(summary.RuntimePackages))

		if len(summary.ExternalPackages) > 0 && len(summary.ExternalPackages) <= 10 {
			fmt.Println("\nExternal Go Dependencies:")
			for _, pkg := range summary.ExternalPackages {
				fmt.Printf("  - %s\n", pkg)
			}
		}
	}

	fmt.Println("\n✓ Analysis complete")
	return nil
}

// findMappingsFile locates the npm-to-go mappings file
func findMappingsFile() string {
	paths := []string{
		"mappings/npm-to-go.yaml",
		"../../mappings/npm-to-go.yaml",
		"../../../mappings/npm-to-go.yaml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "mappings/npm-to-go.yaml" // Default
}

// findParserScript locates the TypeScript parser script
func findParserScript() string {
	paths := []string{
		"internal/transpiler/parser/parser.js",
		"../../internal/transpiler/parser/parser.js",
		"../../../internal/transpiler/parser/parser.js",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "internal/transpiler/parser/parser.js" // Default
}
