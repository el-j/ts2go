package ports

import "github.com/el-j/ts2go/core/domain"

// TranspilationService defines the core transpilation use case.
// This is a DRIVING PORT (inbound) - it's the API of our core logic.
type TranspilationService interface {
	// TranspileProject transpiles an entire TypeScript project to Go
	TranspileProject(projectPath string, options TranspileOptions) (*domain.TranspilationResult, error)

	// TranspileFile transpiles a single TypeScript file to Go
	TranspileFile(filePath string, options TranspileOptions) (*domain.FileTranspilationResult, error)

	// AnalyzeProject analyzes a TypeScript project without transpiling
	AnalyzeProject(projectPath string) (*domain.AnalysisReport, error)
}

// TranspileOptions contains options for transpilation
type TranspileOptions struct {
	OutputPath      string
	GoModuleName    string
	IncludePatterns []string
	ExcludePatterns []string
	FormatOutput    bool
	Optimize        bool
}
