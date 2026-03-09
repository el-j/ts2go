package project

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/el-j/ts2go/internal/analyzer"
	"github.com/el-j/ts2go/internal/mapper"
)

// DependencyGraph represents dependencies between source files
type DependencyGraph struct {
	Nodes       map[string]*GraphNode
	EntryPoints []string
}

// GraphNode represents a file in the dependency graph
type GraphNode struct {
	File         *SourceFile
	Dependencies []*GraphNode // Files this file depends on
	Dependents   []*GraphNode // Files that depend on this file
}

// BuildDependencyGraph builds a dependency graph for the project
// Note: This should be called after files have had their imports analyzed
func (p *Project) BuildDependencyGraph() (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Nodes:       make(map[string]*GraphNode),
		EntryPoints: []string{},
	}

	// Create nodes for all files
	for _, file := range p.Files {
		graph.Nodes[file.RelativePath] = &GraphNode{
			File:         file,
			Dependencies: []*GraphNode{},
			Dependents:   []*GraphNode{},
		}
	}

	// Build edges based on imports
	for _, file := range p.Files {
		// Process each import
		for _, imp := range file.Imports {
			// Only process local imports
			if imp.Type != analyzer.ImportTypeLocal {
				continue
			}

			// Resolve the import path relative to the current file
			importPath := resolveImportPath(p.RootDir, file.Path, imp.Source)
			if importPath == "" {
				continue
			}

			// Find the relative path
			relPath, err := filepath.Rel(p.RootDir, importPath)
			if err != nil {
				continue
			}

			// Add edge in the graph
			sourceNode := graph.Nodes[file.RelativePath]
			targetNode := graph.Nodes[relPath]

			if sourceNode != nil && targetNode != nil {
				sourceNode.Dependencies = append(sourceNode.Dependencies, targetNode)
				targetNode.Dependents = append(targetNode.Dependents, sourceNode)
			}
		}
	}

	// Identify entry points (files with no dependents or marked as entry)
	for _, file := range p.Files {
		node := graph.Nodes[file.RelativePath]
		if len(node.Dependents) == 0 || file.IsEntry {
			graph.EntryPoints = append(graph.EntryPoints, file.RelativePath)
		}
	}

	return graph, nil
}

// ResolveBuildOrder determines the order files should be transpiled
func (g *DependencyGraph) ResolveBuildOrder() ([]string, error) {
	// Use topological sort (Kahn's algorithm)
	inDegree := make(map[string]int)

	// Calculate in-degrees
	for path, node := range g.Nodes {
		inDegree[path] = len(node.Dependencies)
	}

	// Find all nodes with in-degree 0
	queue := []string{}
	for path, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, path)
		}
	}

	// Process nodes in order
	result := []string{}
	for len(queue) > 0 {
		// Pop from queue
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for dependents
		node := g.Nodes[current]
		for _, dependent := range node.Dependents {
			inDegree[dependent.File.RelativePath]--
			if inDegree[dependent.File.RelativePath] == 0 {
				queue = append(queue, dependent.File.RelativePath)
			}
		}
	}

	// Check for cycles
	if len(result) != len(g.Nodes) {
		cycles := []string{}
		for path, degree := range inDegree {
			if degree > 0 {
				cycles = append(cycles, path)
			}
		}
		return nil, fmt.Errorf("circular dependencies detected: %v", cycles)
	}

	return result, nil
}

// DetectCircularDependencies finds circular dependencies in the graph
func (g *DependencyGraph) DetectCircularDependencies() [][]string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	cycles := [][]string{}

	var dfs func(path string, currentPath []string)
	dfs = func(path string, currentPath []string) {
		visited[path] = true
		recStack[path] = true
		currentPath = append(currentPath, path)

		node := g.Nodes[path]
		for _, dep := range node.Dependencies {
			depPath := dep.File.RelativePath

			if !visited[depPath] {
				dfs(depPath, currentPath)
			} else if recStack[depPath] {
				// Found a cycle - extract it
				cycleStart := -1
				for i, p := range currentPath {
					if p == depPath {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := make([]string, len(currentPath)-cycleStart)
					copy(cycle, currentPath[cycleStart:])
					cycles = append(cycles, cycle)
				}
			}
		}

		recStack[path] = false
	}

	for path := range g.Nodes {
		if !visited[path] {
			dfs(path, []string{})
		}
	}

	return cycles
}

// GetDependencyStats returns statistics about the dependency graph
func (g *DependencyGraph) GetDependencyStats() map[string]interface{} {
	totalDeps := 0
	maxDeps := 0
	maxDepsFile := ""

	for path, node := range g.Nodes {
		depCount := len(node.Dependencies)
		totalDeps += depCount

		if depCount > maxDeps {
			maxDeps = depCount
			maxDepsFile = path
		}
	}

	avgDeps := 0.0
	if len(g.Nodes) > 0 {
		avgDeps = float64(totalDeps) / float64(len(g.Nodes))
	}

	return map[string]interface{}{
		"total_files":   len(g.Nodes),
		"entry_points":  len(g.EntryPoints),
		"total_deps":    totalDeps,
		"avg_deps":      avgDeps,
		"max_deps":      maxDeps,
		"max_deps_file": maxDepsFile,
	}
}

// AnalyzeDependencies analyzes npm dependencies from imports
func (p *Project) AnalyzeDependencies(classifier *mapper.Classifier) error {
	// Collect all npm imports from all files
	npmPackages := make(map[string]bool)

	for _, file := range p.Files {
		for _, imp := range file.Imports {
			if imp.Type == analyzer.ImportTypePackage {
				// Extract package name (handle scoped packages)
				pkgName := extractPackageName(imp.Source)
				npmPackages[pkgName] = true
			}
		}
	}

	// Classify each package
	for pkg := range npmPackages {
		classification, err := classifier.ClassifyPackage(pkg)
		if err != nil {
			return fmt.Errorf("failed to classify package %s: %w", pkg, err)
		}
		p.Dependencies[pkg] = classification
	}

	return nil
}

// GetSupportedPackages returns packages that have mappings
func (p *Project) GetSupportedPackages() []*mapper.Classification {
	result := []*mapper.Classification{}
	for _, classification := range p.Dependencies {
		if classification.Class == mapper.PackageClassBuiltin ||
			classification.Class == mapper.PackageClassRuntime ||
			classification.Class == mapper.PackageClassEquivalent {
			result = append(result, classification)
		}
	}
	return result
}

// GetUnsupportedPackages returns packages without mappings
func (p *Project) GetUnsupportedPackages() []*mapper.Classification {
	result := []*mapper.Classification{}
	for _, classification := range p.Dependencies {
		if classification.Class == mapper.PackageClassUnsupported {
			result = append(result, classification)
		}
	}
	return result
}

// Helper functions

// resolveImportPath resolves an import path to an absolute file path
func resolveImportPath(rootDir, sourceFile, importPath string) string {
	// Get directory of source file
	sourceDir := filepath.Dir(sourceFile)

	// Resolve relative path
	var resolvedPath string
	if strings.HasPrefix(importPath, ".") {
		resolvedPath = filepath.Join(sourceDir, importPath)
	} else {
		// Assume it's relative to root
		resolvedPath = filepath.Join(rootDir, importPath)
	}

	// Try with .ts extension
	if !strings.HasSuffix(resolvedPath, ".ts") {
		if _, err := filepath.Abs(resolvedPath + ".ts"); err == nil {
			resolvedPath = resolvedPath + ".ts"
		}
	}

	// Clean the path
	resolvedPath = filepath.Clean(resolvedPath)

	return resolvedPath
}

// extractPackageName extracts the package name from an import source
func extractPackageName(source string) string {
	// Handle scoped packages (@scope/package)
	if strings.HasPrefix(source, "@") {
		parts := strings.SplitN(source, "/", 3)
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
	}

	// Regular packages
	parts := strings.Split(source, "/")
	return parts[0]
}
