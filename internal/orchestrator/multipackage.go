package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/el-j/ts2go/internal/module"
	"github.com/el-j/ts2go/internal/project"
	"github.com/el-j/ts2go/internal/transpiler"
)

// MultiPackageTranspiler handles transpilation of multi-file TypeScript projects
type MultiPackageTranspiler struct {
	project  *project.Project
	registry *module.ExportRegistry
	mapper   *module.PackageMapper
	parser   *module.Parser
}

// NewMultiPackageTranspiler creates a new multi-package transpiler
func NewMultiPackageTranspiler(proj *project.Project, goModuleName string) *MultiPackageTranspiler {
	parser := module.NewParser()

	// Determine source directory
	srcDir := "src"
	if _, err := os.Stat(filepath.Join(proj.RootDir, "src")); os.IsNotExist(err) {
		srcDir = ""
	}

	mapper := module.NewPackageMapper(proj.RootDir, goModuleName, srcDir)

	return &MultiPackageTranspiler{
		project:  proj,
		registry: parser.GetRegistry(),
		mapper:   mapper,
		parser:   parser,
	}
}

// TranspileProject transpiles an entire TypeScript project to Go
func (t *MultiPackageTranspiler) TranspileProject(outputDir string) error {
	// Step 1: Parse all files and build export registry
	fmt.Println("Parsing TypeScript files and analyzing exports...")
	for _, file := range t.project.Files {
		if _, err := t.parser.ParseFile(file.Path); err != nil {
			return fmt.Errorf("failed to parse %s: %w", file.Path, err)
		}
	}

	// Step 2: Generate package structure
	fmt.Println("Generating Go package structure...")
	filePaths := make([]string, len(t.project.Files))
	for i, file := range t.project.Files {
		filePaths[i] = file.Path
	}

	structure, err := module.GeneratePackageStructure(filePaths, t.mapper)
	if err != nil {
		return fmt.Errorf("failed to generate package structure: %w", err)
	}

	// Step 3: Build dependency graph for transpilation order
	fmt.Println("Building dependency graph...")
	graph, err := t.project.BuildDependencyGraph()
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Check for circular dependencies
	cycles := graph.DetectCircularDependencies()
	if len(cycles) > 0 {
		fmt.Printf("Warning: Found %d circular dependencies (will use fallback ordering)\n", len(cycles))
	}

	// Resolve build order
	buildOrder, err := graph.ResolveBuildOrder()
	if err != nil {
		fmt.Printf("Warning: Could not resolve optimal build order: %v\n", err)
		// Fallback: use file order
		buildOrder = make([]string, len(t.project.Files))
		for i, file := range t.project.Files {
			buildOrder[i] = file.RelativePath
		}
	}

	// Step 4: Transpile each file in dependency order
	fmt.Printf("Transpiling %d files...\n", len(buildOrder))
	successCount := 0

	for i, relPath := range buildOrder {
		// Find the source file
		var sourceFile *project.SourceFile
		for _, file := range t.project.Files {
			if file.RelativePath == relPath {
				sourceFile = file
				break
			}
		}

		if sourceFile == nil {
			fmt.Printf("Warning: Could not find source file for %s\n", relPath)
			continue
		}

		// Transpile the file
		isEntry := sourceFile.Path == structure.EntryPoint
		if err := t.transpileFile(sourceFile, outputDir, structure, isEntry); err != nil {
			fmt.Printf("Error transpiling %s: %v\n", relPath, err)
			continue
		}

		successCount++
		if (i+1)%10 == 0 || i+1 == len(buildOrder) {
			fmt.Printf("  Progress: %d/%d files (%d%%)\n", i+1, len(buildOrder), (i+1)*100/len(buildOrder))
		}
	}

	// Step 5: Generate go.mod file
	fmt.Println("Generating go.mod...")
	if err := t.generateGoMod(outputDir, structure); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	fmt.Printf("\n✓ Successfully transpiled %d/%d files\n", successCount, len(buildOrder))
	return nil
}

// transpileFile transpiles a single TypeScript file to Go
func (t *MultiPackageTranspiler) transpileFile(sourceFile *project.SourceFile, outputDir string, structure *module.PackageStructure, isEntry bool) error {
	// Parse TypeScript to AST
	ast, err := transpiler.ParseTypeScript(sourceFile.Path)
	if err != nil {
		return fmt.Errorf("failed to parse TypeScript: %w", err)
	}

	// Get module info
	mod := t.registry.GetModule(sourceFile.Path)
	if mod == nil {
		return fmt.Errorf("module not found in registry")
	}

	// Update module with package name
	mod.PackageName = t.mapper.GetPackageName(sourceFile.Path)

	// Create import resolver
	resolver := module.NewImportResolver(t.registry, t.mapper)
	resolver.SetCurrentFile(sourceFile.Path)

	// Create symbol visibility handler
	visibility := module.NewSymbolVisibility(t.registry, sourceFile.Path)

	// Generate Go code with module context
	generator := transpiler.NewCodeGeneratorWithModule(mod, resolver, visibility, isEntry)
	goCode, err := generator.Generate(ast)
	if err != nil {
		return fmt.Errorf("failed to generate Go code: %w", err)
	}

	// Determine output path
	outputPath := t.mapper.GetOutputPath(sourceFile.Path, outputDir)

	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, []byte(goCode), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// generateGoMod generates the go.mod file for the project
func (t *MultiPackageTranspiler) generateGoMod(outputDir string, structure *module.PackageStructure) error {
	moduleName := t.mapper.GetPackagePath(structure.EntryPoint)
	if moduleName == "" {
		moduleName = "transpiled-project"
	}

	// Extract base module name (without sub-packages)
	parts := strings.Split(moduleName, "/")
	if len(parts) > 3 {
		// For github.com/user/project/subpkg, use github.com/user/project
		moduleName = strings.Join(parts[:3], "/")
	}

	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/yourusername/ts2go-runtime v0.1.0
)
`, moduleName)

	goModPath := filepath.Join(outputDir, "go.mod")
	return os.WriteFile(goModPath, []byte(content), 0644)
}
