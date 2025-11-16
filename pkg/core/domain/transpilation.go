package domain

import "time"

// TranspilationResult represents the result of transpiling a TypeScript project to Go
type TranspilationResult struct {
	Success        bool                      `json:"success"`
	ProjectPath    string                    `json:"projectPath"`
	OutputPath     string                    `json:"outputPath"`
	FilesProcessed int                       `json:"filesProcessed"`
	FilesSuccess   int                       `json:"filesSuccess"`
	FilesFailed    int                       `json:"filesFailed"`
	FileResults    []FileTranspilationResult `json:"fileResults"`
	Duration       time.Duration             `json:"duration"`
	StartedAt      time.Time                 `json:"startedAt"`
	CompletedAt    time.Time                 `json:"completedAt"`
	Errors         []string                  `json:"errors,omitempty"`
	Warnings       []string                  `json:"warnings,omitempty"`
}

// FileTranspilationResult represents the result of transpiling a single file
type FileTranspilationResult struct {
	SourcePath string   `json:"sourcePath"`
	TargetPath string   `json:"targetPath"`
	Success    bool     `json:"success"`
	GoCode     string   `json:"goCode,omitempty"`
	Error      string   `json:"error,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
}

// NewTranspilationResult creates a new TranspilationResult
func NewTranspilationResult(projectPath, outputPath string) *TranspilationResult {
	return &TranspilationResult{
		ProjectPath: projectPath,
		OutputPath:  outputPath,
		StartedAt:   time.Now(),
		FileResults: make([]FileTranspilationResult, 0),
		Errors:      make([]string, 0),
		Warnings:    make([]string, 0),
	}
}

// AddFileResult adds a file transpilation result
func (r *TranspilationResult) AddFileResult(result FileTranspilationResult) {
	r.FileResults = append(r.FileResults, result)
	r.FilesProcessed++

	if result.Success {
		r.FilesSuccess++
	} else {
		r.FilesFailed++
	}

	if result.Error != "" {
		r.Errors = append(r.Errors, result.Error)
	}

	r.Warnings = append(r.Warnings, result.Warnings...)
}

// Complete marks the transpilation as complete
func (r *TranspilationResult) Complete() {
	r.CompletedAt = time.Now()
	r.Duration = r.CompletedAt.Sub(r.StartedAt)
	r.Success = r.FilesFailed == 0 && len(r.Errors) == 0
}

// AnalysisReport represents the result of analyzing a TypeScript project
type AnalysisReport struct {
	ProjectPath     string            `json:"projectPath"`
	TotalFiles      int               `json:"totalFiles"`
	TypeScriptFiles int               `json:"typeScriptFiles"`
	Dependencies    []string          `json:"dependencies"`
	Imports         []ImportStatement `json:"imports"`
	Exports         []ExportStatement `json:"exports"`
	Complexity      ComplexityMetrics `json:"complexity"`
	Warnings        []string          `json:"warnings"`
	GeneratedAt     time.Time         `json:"generatedAt"`
}

// ImportStatement represents an import in TypeScript code
type ImportStatement struct {
	Source    string   `json:"source"`
	Imports   []string `json:"imports"`
	IsDefault bool     `json:"isDefault"`
	FilePath  string   `json:"filePath"`
}

// ExportStatement represents an export in TypeScript code
type ExportStatement struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "function", "class", "variable", "interface", "type"
	FilePath string `json:"filePath"`
}

// ComplexityMetrics represents code complexity metrics
type ComplexityMetrics struct {
	TotalLines           int     `json:"totalLines"`
	CodeLines            int     `json:"codeLines"`
	CommentLines         int     `json:"commentLines"`
	BlankLines           int     `json:"blankLines"`
	AvgFunctionLength    float64 `json:"avgFunctionLength"`
	MaxFunctionLength    int     `json:"maxFunctionLength"`
	CyclomaticComplexity int     `json:"cyclomaticComplexity"`
}

// NewAnalysisReport creates a new AnalysisReport
func NewAnalysisReport(projectPath string) *AnalysisReport {
	return &AnalysisReport{
		ProjectPath:  projectPath,
		GeneratedAt:  time.Now(),
		Dependencies: make([]string, 0),
		Imports:      make([]ImportStatement, 0),
		Exports:      make([]ExportStatement, 0),
		Warnings:     make([]string, 0),
	}
}
