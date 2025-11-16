package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/el-j/ts2go/pkg/adapters/driven/filesystem"
	"github.com/el-j/ts2go/pkg/adapters/driven/gocompiler"
	"github.com/el-j/ts2go/pkg/adapters/driven/persistence"
	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports"
	"github.com/el-j/ts2go/pkg/core/services"
)

// Application is the composition root for the CLI application
// It wires together all dependencies using dependency injection
type Application struct {
	// Core services (use cases)
	transpilationService ports.TranspilationService
	runtimeService       ports.GoRuntimeService
	stateService         ports.StateService

	// Infrastructure adapters
	fileSystem   ports.FileSystem
	compiler     ports.GoCompiler
	stateRepo    ports.StateRepository
	settingsRepo ports.SettingsRepository
}

// NewApplication creates a new CLI application with all dependencies wired
func NewApplication() (*Application, error) {
	// Create infrastructure adapters (driven ports)
	fs := filesystem.NewOSFileSystem()

	compiler, err := gocompiler.NewSystemGoCompiler()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Go compiler: %w", err)
	}

	// Get user home directory for state storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	stateDir := filepath.Join(homeDir, ".ts2go", "state")
	stateRepo, err := persistence.NewJSONStateRepository(stateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize state repository: %w", err)
	}

	settingsPath := filepath.Join(homeDir, ".ts2go", "settings.json")
	settingsRepo, err := persistence.NewJSONSettingsRepository(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize settings repository: %w", err)
	}

	// Create placeholder services for analyzer, mapper, codegen
	// TODO: Implement these services in Phase 1 continuation
	analyzer := &placeholderAnalyzer{}
	mapper := &placeholderMapper{}
	codegen := &placeholderCodeGen{}

	// Create core services (business logic)
	transpilationService := services.NewTranspilationService(fs, compiler, analyzer, mapper, codegen)
	runtimeService := services.NewGoRuntimeService(compiler, fs)
	stateService := services.NewStateService(stateRepo, settingsRepo)

	return &Application{
		transpilationService: transpilationService,
		runtimeService:       runtimeService,
		stateService:         stateService,
		fileSystem:           fs,
		compiler:             compiler,
		stateRepo:            stateRepo,
		settingsRepo:         settingsRepo,
	}, nil
}

// TranspileCommand implements the 'transpile' subcommand using hexagonal architecture
func (app *Application) TranspileCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ts2go transpile <project-directory> [--out <output-dir>] [--verbose] [--quiet]")
	}

	projectDir := args[0]
	outputDir := "output"
	verbose := false

	// Parse command line options
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--out", "-o":
			if i+1 < len(args) {
				outputDir = args[i+1]
				i++
			}
		case "--verbose", "-v":
			verbose = true
		case "--quiet", "-q":
			verbose = false
		}
	}

	// Validate and resolve paths
	projectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	exists, err := app.fileSystem.DirectoryExists(projectDir)
	if err != nil {
		return fmt.Errorf("failed to check project directory: %w", err)
	}
	if !exists {
		return fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	if verbose {
		fmt.Printf("Transpiling TypeScript Project\n")
		fmt.Printf("Source: %s\n", projectDir)
		fmt.Printf("Output: %s\n\n", outputDir)
	}

	// Create output directory
	if err := app.fileSystem.CreateDirectory(outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare transpilation options
	options := ports.TranspileOptions{
		OutputPath:      outputDir,
		IncludePatterns: []string{"**/*.ts", "**/*.tsx"},
		ExcludePatterns: []string{"node_modules/**", "**/*.test.ts", "**/*.spec.ts"},
		FormatOutput:    true,
		Optimize:        true,
	}

	// Execute transpilation through the service
	result, err := app.transpilationService.TranspileProject(projectDir, options)
	if err != nil {
		return fmt.Errorf("transpilation failed: %w", err)
	}

	// Display results
	fmt.Printf("\n✓ Transpilation complete\n")
	fmt.Printf("  Files processed: %d\n", result.FilesProcessed)
	fmt.Printf("  Files succeeded: %d\n", result.FilesSuccess)
	fmt.Printf("  Files failed: %d\n", result.FilesFailed)
	fmt.Printf("  Duration: %s\n", result.Duration)

	if len(result.Errors) > 0 && verbose {
		fmt.Printf("\nErrors:\n")
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	if len(result.Warnings) > 0 && verbose {
		fmt.Printf("\nWarnings:\n")
		for _, warn := range result.Warnings {
			fmt.Printf("  - %s\n", warn)
		}
	}

	return nil
}

// AnalyzeCommand implements the 'analyze' subcommand using hexagonal architecture
func (app *Application) AnalyzeCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ts2go analyze <project-directory>")
	}

	projectDir := args[0]

	// Validate and resolve path
	projectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	exists, err := app.fileSystem.DirectoryExists(projectDir)
	if err != nil {
		return fmt.Errorf("failed to check project directory: %w", err)
	}
	if !exists {
		return fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	fmt.Printf("Analyzing TypeScript Project: %s\n\n", projectDir)

	// Execute analysis through the service
	report, err := app.transpilationService.AnalyzeProject(projectDir)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Display results
	fmt.Printf("Project Analysis Results\n")
	fmt.Printf("========================\n\n")
	fmt.Printf("Total files: %d\n", report.TotalFiles)
	fmt.Printf("TypeScript files: %d\n", report.TypeScriptFiles)

	if len(report.Dependencies) > 0 {
		fmt.Printf("\nDependencies (%d):\n", len(report.Dependencies))
		for _, dep := range report.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}
	}

	if len(report.Warnings) > 0 {
		fmt.Printf("\nWarnings (%d):\n", len(report.Warnings))
		for _, warn := range report.Warnings {
			fmt.Printf("  - %s\n", warn)
		}
	}

	return nil
}

// BuildCommand implements the 'build' subcommand using hexagonal architecture
func (app *Application) BuildCommand(args []string) error {
	var sourceDir string
	var outputPath string

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
		case "--tags", "-t":
			if i+1 < len(args) {
				// Parse tags but they're handled by the compiler adapter
				i++
			}
		}
	}

	if sourceDir == "" {
		return fmt.Errorf("usage: ts2go build --source <source-dir> --output <binary-path> [--tags <build-tags>]")
	}

	if outputPath == "" {
		baseName := filepath.Base(sourceDir)
		outputPath = filepath.Join(sourceDir, baseName)
	}

	// Validate source directory
	exists, err := app.fileSystem.DirectoryExists(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to check source directory: %w", err)
	}
	if !exists {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	fmt.Printf("Building Go project from %s...\n", sourceDir)
	fmt.Printf("Output binary: %s\n", outputPath)

	// Execute build through the runtime service
	result, err := app.runtimeService.BuildProject(sourceDir, outputPath)
	if err != nil {
		if result != nil {
			fmt.Fprintf(os.Stderr, "Build output:\n%s\n", result.Output)
		}
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("\n✓ Build successful\n")
	fmt.Printf("  Output: %s\n", result.OutputPath)
	fmt.Printf("  Duration: %s\n", result.Duration)

	return nil
}

// TestCommand implements the 'test' subcommand using hexagonal architecture
func (app *Application) TestCommand(args []string) error {
	var sourceDir string
	verbose := false
	coverage := false

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
		}
	}

	if sourceDir == "" {
		return fmt.Errorf("usage: ts2go test --source <source-dir> [--verbose] [--coverage]")
	}

	// Validate source directory
	exists, err := app.fileSystem.DirectoryExists(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to check source directory: %w", err)
	}
	if !exists {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	fmt.Printf("Running tests in %s...\n\n", sourceDir)

	// Prepare test options
	options := &ports.TestOptions{
		Verbose:  verbose,
		Coverage: coverage,
	}

	// Execute tests through the runtime service
	result, err := app.runtimeService.TestProject(sourceDir, options)
	if err != nil {
		if result != nil {
			fmt.Printf("Test output:\n%s\n", result.Output)
		}
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Printf("✓ All tests passed\n")
	fmt.Printf("  Total: %d\n", result.TotalTests)
	fmt.Printf("  Passed: %d\n", result.PassedTests)
	if result.Coverage > 0 {
		fmt.Printf("  Coverage: %.1f%%\n", result.Coverage)
	}

	return nil
}

// StateCommand implements the 'state' subcommand for persistent state management
func (app *Application) StateCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ts2go state <get|set|delete> <project-path> [json-data]")
	}

	action := args[0]
	projectPath := args[1]

	switch action {
	case "get":
		state, err := app.stateService.GetProjectState(projectPath)
		if err != nil {
			// Return empty state if not found
			fmt.Printf("{}\n")
			return nil
		}
		// Output as JSON
		data, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal state: %w", err)
		}
		fmt.Printf("%s\n", string(data))
		return nil

	case "set":
		if len(args) < 3 {
			return fmt.Errorf("usage: ts2go state set <project-path> <json-data>")
		}
		// Parse JSON data from command line
		jsonData := args[2]
		var state domain.ProjectState
		if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
			return fmt.Errorf("invalid JSON data: %w", err)
		}
		// Ensure projectPath is set correctly
		state.ProjectPath = projectPath
		if err := app.stateService.SaveProjectState(&state); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		fmt.Printf("✓ State saved for project: %s\n", projectPath)
		return nil

	case "delete":
		if err := app.stateService.DeleteProjectState(projectPath); err != nil {
			return fmt.Errorf("failed to delete state: %w", err)
		}
		fmt.Printf("✓ State deleted for project: %s\n", projectPath)
		return nil

	default:
		return fmt.Errorf("unknown action: %s. Valid actions: get, set, delete", action)
	}
}

// SettingsCommand implements the 'settings' subcommand for application settings
func (app *Application) SettingsCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ts2go settings <get|set> [json-data]")
	}

	action := args[0]

	switch action {
	case "get":
		settings, err := app.stateService.GetSettings()
		if err != nil {
			// Return empty settings if not found
			fmt.Printf("{}\n")
			return nil
		}
		// Output as JSON
		data, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal settings: %w", err)
		}
		fmt.Printf("%s\n", string(data))
		return nil

	case "set":
		if len(args) < 2 {
			return fmt.Errorf("usage: ts2go settings set <json-data>")
		}
		// Parse JSON data from command line
		jsonData := args[1]
		var settings domain.Settings
		if err := json.Unmarshal([]byte(jsonData), &settings); err != nil {
			return fmt.Errorf("invalid JSON data: %w", err)
		}
		if err := app.stateService.SaveSettings(&settings); err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}
		fmt.Printf("✓ Settings saved\n")
		return nil

	default:
		return fmt.Errorf("unknown action: %s. Valid actions: get, set", action)
	}
}

// Placeholder implementations for services not yet migrated
type placeholderAnalyzer struct{}

func (p *placeholderAnalyzer) AnalyzeImports(ast interface{}) (*services.ImportAnalysis, error) {
	return &services.ImportAnalysis{
		LocalImports: []string{},
		NpmPackages:  []string{},
		NodeBuiltins: []string{},
	}, nil
}

type placeholderMapper struct{}

func (p *placeholderMapper) FindGoPackage(npmPackage string) (string, error) {
	return "", nil
}

func (p *placeholderMapper) TransformAPICall(code string) (string, error) {
	return code, nil
}

type placeholderCodeGen struct{}

func (p *placeholderCodeGen) ParseTypeScript(filePath string) (interface{}, error) {
	return nil, fmt.Errorf("TypeScript parsing not yet implemented in hexagonal architecture")
}

func (p *placeholderCodeGen) GenerateGoCode(ast interface{}) (string, error) {
	return "", fmt.Errorf("Go code generation not yet implemented in hexagonal architecture")
}
