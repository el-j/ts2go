package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/el-j/ts2go/core/domain"
	"github.com/el-j/ts2go/core/ports"
)

// TranspilationServiceImpl implements the TranspilationService port
// This is the CORE BUSINESS LOGIC with NO infrastructure dependencies
type TranspilationServiceImpl struct {
	fs       ports.FileSystem
	compiler ports.GoCompiler
	analyzer AnalyzerService
	mapper   MapperService
	codegen  CodeGenService
}

// NewTranspilationService creates a new transpilation service
// Dependencies are injected through ports (interfaces)
func NewTranspilationService(
	fs ports.FileSystem,
	compiler ports.GoCompiler,
	analyzer AnalyzerService,
	mapper MapperService,
	codegen CodeGenService,
) ports.TranspilationService {
	return &TranspilationServiceImpl{
		fs:       fs,
		compiler: compiler,
		analyzer: analyzer,
		mapper:   mapper,
		codegen:  codegen,
	}
}

// TranspileProject transpiles an entire TypeScript project to Go
func (s *TranspilationServiceImpl) TranspileProject(projectPath string, options ports.TranspileOptions) (*domain.TranspilationResult, error) {
	if projectPath == "" {
		return nil, domain.NewValidationError("project path cannot be empty")
	}

	result := domain.NewTranspilationResult(projectPath, options.OutputPath)

	// Scan directory for TypeScript files
	files, err := s.fs.ScanDirectory(projectPath, options.IncludePatterns, options.ExcludePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	// Transpile each file
	for _, file := range files {
		if !file.IsTypeScript() {
			continue // Skip non-TypeScript files
		}

		fileResult, err := s.TranspileFile(file.Path, options)
		if err != nil {
			result.AddFileResult(domain.FileTranspilationResult{
				SourcePath: file.Path,
				Success:    false,
				Error:      err.Error(),
			})
			continue
		}

		result.AddFileResult(*fileResult)
	}

	result.Complete()
	return result, nil
}

// TranspileFile transpiles a single TypeScript file to Go
func (s *TranspilationServiceImpl) TranspileFile(filePath string, options ports.TranspileOptions) (*domain.FileTranspilationResult, error) {
	if filePath == "" {
		return nil, domain.NewValidationError("file path cannot be empty")
	}

	// Read source content
	_, err := s.fs.ReadFile(filePath)
	if err != nil {
		return &domain.FileTranspilationResult{
			SourcePath: filePath,
			Success:    false,
			Error:      fmt.Sprintf("read error: %v", err),
		}, err
	}

	// Parse TypeScript to AST
	ast, err := s.codegen.ParseTypeScript(filePath)
	if err != nil {
		return &domain.FileTranspilationResult{
			SourcePath: filePath,
			Success:    false,
			Error:      fmt.Sprintf("parse error: %v", err),
		}, err
	}

	// Generate Go code
	goCode, err := s.codegen.GenerateGoCode(ast)
	if err != nil {
		return &domain.FileTranspilationResult{
			SourcePath: filePath,
			Success:    false,
			Error:      fmt.Sprintf("codegen error: %v", err),
		}, err
	}

	// Determine output path
	outputPath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".go"
	if options.OutputPath != "" {
		// Use custom output path
		rel, _ := filepath.Rel(filepath.Dir(filePath), filePath)
		outputPath = filepath.Join(options.OutputPath, rel)
		outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".go"
	}

	return &domain.FileTranspilationResult{
		SourcePath: filePath,
		TargetPath: outputPath,
		GoCode:     goCode,
		Success:    true,
	}, nil
}

// AnalyzeProject analyzes a TypeScript project without transpiling
func (s *TranspilationServiceImpl) AnalyzeProject(projectPath string) (*domain.AnalysisReport, error) {
	if projectPath == "" {
		return nil, domain.NewValidationError("project path cannot be empty")
	}

	report := domain.NewAnalysisReport(projectPath)

	// Scan directory for files
	files, err := s.fs.ScanDirectory(projectPath, []string{"**/*.ts", "**/*.tsx"}, []string{"node_modules/**"})
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	report.TotalFiles = len(files)

	// Analyze each TypeScript file
	for _, file := range files {
		if !file.IsTypeScript() {
			continue
		}

		report.TypeScriptFiles++

		// Read and parse file
		_, err := s.fs.ReadFile(file.Path)
		if err != nil {
			report.Warnings = append(report.Warnings, fmt.Sprintf("%s: read error: %v", file.Path, err))
			continue
		}

		ast, err := s.codegen.ParseTypeScript(file.Path)
		if err != nil {
			report.Warnings = append(report.Warnings, fmt.Sprintf("%s: parse error: %v", file.Path, err))
			continue
		}

		// Analyze imports
		analysis, err := s.analyzer.AnalyzeImports(ast)
		if err != nil {
			report.Warnings = append(report.Warnings, fmt.Sprintf("%s: analysis error: %v", file.Path, err))
			continue
		}

		// Collect dependencies
		for _, pkg := range analysis.NpmPackages {
			report.Dependencies = append(report.Dependencies, pkg)
		}
	}

	return report, nil
}

// AnalyzerService defines the interface for analyzing TypeScript code
type AnalyzerService interface {
	AnalyzeImports(ast interface{}) (*ImportAnalysis, error)
}

// MapperService defines the interface for mapping npm packages to Go
type MapperService interface {
	FindGoPackage(npmPackage string) (string, error)
	TransformAPICall(code string) (string, error)
}

// CodeGenService defines the interface for code generation
type CodeGenService interface {
	ParseTypeScript(filePath string) (interface{}, error)
	GenerateGoCode(ast interface{}) (string, error)
}

// ImportAnalysis contains import analysis results
type ImportAnalysis struct {
	LocalImports []string
	NpmPackages  []string
	NodeBuiltins []string
}
