package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/el-j/ts2go/pkg/adapters/driven/filesystem"
	"github.com/el-j/ts2go/pkg/adapters/driven/gocompiler"
	"github.com/el-j/ts2go/pkg/adapters/driven/persistence"
	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports"
	"github.com/el-j/ts2go/pkg/core/services"
)

// Server is the web API driving adapter for the hexagonal architecture
// It exposes HTTP endpoints that use the same core services as CLI and Desktop UI
type Server struct {
	// Core services (shared with CLI and Desktop UI)
	transpilationService ports.TranspilationService
	runtimeService       ports.GoRuntimeService
	stateService         ports.StateService

	// Server configuration
	addr string
	mux  *http.ServeMux
}

// NewServer creates a new web API server with dependency injection
func NewServer(addr string) (*Server, error) {
	// Initialize infrastructure adapters (same as CLI)
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
	analyzer := &placeholderAnalyzer{}
	mapper := &placeholderMapper{}
	codegen := &placeholderCodeGen{}

	// Create core services (same as CLI and Desktop UI)
	transpilationService := services.NewTranspilationService(fs, compiler, analyzer, mapper, codegen)
	runtimeService := services.NewGoRuntimeService(compiler, fs)
	stateService := services.NewStateService(stateRepo, settingsRepo)

	server := &Server{
		transpilationService: transpilationService,
		runtimeService:       runtimeService,
		stateService:         stateService,
		addr:                 addr,
		mux:                  http.NewServeMux(),
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes configures HTTP endpoints
func (s *Server) setupRoutes() {
	// Health check
	s.mux.HandleFunc("/health", s.handleHealth)

	// Transpilation endpoints
	s.mux.HandleFunc("/api/transpile", s.handleTranspile)
	s.mux.HandleFunc("/api/analyze", s.handleAnalyze)

	// Runtime endpoints
	s.mux.HandleFunc("/api/build", s.handleBuild)
	s.mux.HandleFunc("/api/test", s.handleTest)

	// State management endpoints
	s.mux.HandleFunc("/api/state", s.handleState)
	s.mux.HandleFunc("/api/settings", s.handleSettings)

	// Serve API documentation
	s.mux.HandleFunc("/", s.handleRoot)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("🚀 TS2Go Web API Server starting on %s\n", s.addr)
	fmt.Printf("📚 API Documentation: http://%s/\n", s.addr)
	fmt.Printf("❤️  Health Check: http://%s/health\n\n", s.addr)
	return http.ListenAndServe(s.addr, s.enableCORS(s.mux))
}

// enableCORS wraps the handler with CORS headers
func (s *Server) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "ts2go-web-api",
		"version": "0.2.0-alpha",
	})
}

// handleRoot serves API documentation
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>TS2Go Web API</title>
    <style>
        body { font-family: system-ui; max-width: 800px; margin: 50px auto; padding: 0 20px; }
        h1 { color: #2563eb; }
        .endpoint { background: #f3f4f6; padding: 15px; margin: 10px 0; border-radius: 8px; }
        .method { display: inline-block; padding: 4px 8px; border-radius: 4px; font-weight: bold; }
        .post { background: #10b981; color: white; }
        .get { background: #3b82f6; color: white; }
        code { background: #e5e7eb; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🚀 TS2Go Web API</h1>
    <p>Hexagonal Architecture - Web Driving Adapter</p>
    
    <h2>Endpoints</h2>
    
    <div class="endpoint">
        <span class="method get">GET</span> <code>/health</code>
        <p>Health check endpoint</p>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/transpile</code>
        <p>Transpile TypeScript project to Go</p>
        <pre>{ "projectPath": "/path/to/project", "outputPath": "/path/to/output" }</pre>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/analyze</code>
        <p>Analyze TypeScript project dependencies</p>
        <pre>{ "projectPath": "/path/to/project" }</pre>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/build</code>
        <p>Build Go code into binary</p>
        <pre>{ "sourcePath": "/path/to/source", "outputPath": "/path/to/binary" }</pre>
    </div>
    
    <div class="endpoint">
        <span class="method post">POST</span> <code>/api/test</code>
        <p>Run Go tests</p>
        <pre>{ "sourcePath": "/path/to/source" }</pre>
    </div>
    
    <div class="endpoint">
        <span class="method get">GET</span> <code>/api/state?project=/path/to/project</code>
        <p>Get project state</p>
    </div>
    
    <div class="endpoint">
        <span class="method get">GET</span> <code>/api/settings</code>
        <p>Get application settings</p>
    </div>
    
    <h2>Architecture</h2>
    <p>This Web API is a <strong>driving adapter</strong> in the hexagonal architecture. It shares the same core business logic with:</p>
    <ul>
        <li>✅ CLI (Command Line Interface)</li>
        <li>✅ Desktop UI (Tauri + Vue)</li>
        <li>✅ Web API (this server)</li>
    </ul>
    <p>All three interfaces use identical <code>TranspilationService</code>, <code>GoRuntimeService</code>, and <code>StateService</code> implementations.</p>
</body>
</html>`))
}

// TranspileRequest represents a transpilation request
type TranspileRequest struct {
	ProjectPath string   `json:"projectPath"`
	OutputPath  string   `json:"outputPath"`
	Optimize    bool     `json:"optimize"`
	Include     []string `json:"include,omitempty"`
	Exclude     []string `json:"exclude,omitempty"`
}

// handleTranspile processes transpilation requests
func (s *Server) handleTranspile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req TranspileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.ProjectPath == "" || req.OutputPath == "" {
		respondError(w, http.StatusBadRequest, "projectPath and outputPath are required")
		return
	}

	// Use default patterns if not provided
	if len(req.Include) == 0 {
		req.Include = []string{"**/*.ts", "**/*.tsx"}
	}
	if len(req.Exclude) == 0 {
		req.Exclude = []string{"node_modules/**", "**/*.test.ts", "**/*.spec.ts"}
	}

	// Call hexagonal core service
	options := ports.TranspileOptions{
		OutputPath:      req.OutputPath,
		IncludePatterns: req.Include,
		ExcludePatterns: req.Exclude,
		FormatOutput:    true,
		Optimize:        req.Optimize,
	}

	result, err := s.transpilationService.TranspileProject(req.ProjectPath, options)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Transpilation failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":        result.Success,
		"filesProcessed": result.FilesProcessed,
		"outputPath":     req.OutputPath,
		"duration":       result.Duration.String(),
	})
}

// AnalyzeRequest represents an analysis request
type AnalyzeRequest struct {
	ProjectPath string `json:"projectPath"`
}

// handleAnalyze processes analysis requests
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.ProjectPath == "" {
		respondError(w, http.StatusBadRequest, "projectPath is required")
		return
	}

	// Call hexagonal core service
	result, err := s.transpilationService.AnalyzeProject(req.ProjectPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Analysis failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"totalFiles":   result.TotalFiles,
		"complexity":   result.Complexity,
		"dependencies": result.Dependencies,
	})
}

// BuildRequest represents a build request
type BuildRequest struct {
	SourcePath string   `json:"sourcePath"`
	OutputPath string   `json:"outputPath"`
	BuildFlags []string `json:"buildFlags,omitempty"`
}

// handleBuild processes build requests
func (s *Server) handleBuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req BuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.SourcePath == "" || req.OutputPath == "" {
		respondError(w, http.StatusBadRequest, "sourcePath and outputPath are required")
		return
	}

	// Call hexagonal core service
	result, err := s.runtimeService.BuildProject(req.SourcePath, req.OutputPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Build failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    result.Success,
		"outputPath": result.OutputPath,
		"duration":   result.Duration.String(),
	})
}

// TestRequest represents a test request
type TestRequest struct {
	SourcePath string   `json:"sourcePath"`
	TestFlags  []string `json:"testFlags,omitempty"`
	Coverage   bool     `json:"coverage,omitempty"`
}

// handleTest processes test requests
func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.SourcePath == "" {
		respondError(w, http.StatusBadRequest, "sourcePath is required")
		return
	}

	// Call hexagonal core service
	options := &ports.TestOptions{
		Coverage: req.Coverage,
		Verbose:  false,
	}

	result, err := s.runtimeService.TestProject(req.SourcePath, options)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Tests failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":     result.Success,
		"totalTests":  result.TotalTests,
		"passedTests": result.PassedTests,
		"failedTests": result.FailedTests,
		"coverage":    result.Coverage,
		"duration":    result.Duration.String(),
	})
}

// handleState processes state get/set requests
func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	projectPath := r.URL.Query().Get("project")
	if projectPath == "" {
		respondError(w, http.StatusBadRequest, "project query parameter is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get project state
		state, err := s.stateService.GetProjectState(projectPath)
		if err != nil {
			respondJSON(w, http.StatusOK, domain.NewProjectState(projectPath))
			return
		}
		respondJSON(w, http.StatusOK, state)

	case http.MethodPost:
		// Set project state
		var state domain.ProjectState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid state data: "+err.Error())
			return
		}
		state.ProjectPath = projectPath
		if err := s.stateService.SaveProjectState(&state); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to save state: "+err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "saved"})

	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleSettings processes settings get/set requests
func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get settings
		settings, err := s.stateService.GetSettings()
		if err != nil {
			respondJSON(w, http.StatusOK, &domain.Settings{})
			return
		}
		respondJSON(w, http.StatusOK, settings)

	case http.MethodPost:
		// Set settings
		var settings domain.Settings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid settings data: "+err.Error())
			return
		}
		if err := s.stateService.SaveSettings(&settings); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to save settings: "+err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "saved"})

	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Placeholder implementations (same as CLI)
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
