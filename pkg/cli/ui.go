package cli

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/yourusername/ts2go/internal/analyzer"
	"github.com/yourusername/ts2go/internal/project"
	"github.com/yourusername/ts2go/internal/transpiler"
)

//go:embed ui_templates/*
var uiTemplates embed.FS

// UICommand starts the web-based desktop UI for the transpiler
func UICommand(args []string) error {
	port := 8080
	openBrowser := false

	// Parse command line options
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port", "-p":
			if i+1 < len(args) {
				p, err := strconv.Atoi(args[i+1])
				if err != nil {
					return fmt.Errorf("invalid port number: %s", args[i+1])
				}
				port = p
				i++
			}
		case "--open", "-o":
			openBrowser = true
		}
	}

	server := &UIServer{
		port: port,
	}

	if openBrowser {
		go func() {
			url := fmt.Sprintf("http://localhost:%d", port)
			fmt.Printf("Opening browser at %s\n", url)
			openURL(url)
		}()
	}

	return server.Start()
}

// UIServer handles the web UI server
type UIServer struct {
	port int
}

// Start begins serving the web UI
func (s *UIServer) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/transpile", s.handleTranspile)
	http.HandleFunc("/api/analyze", s.handleAnalyze)
	http.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("🚀 TS2Go Web UI starting on http://localhost:%d\n", s.port)
	fmt.Printf("📝 Open your browser to start transpiling TypeScript to Go\n")
	fmt.Printf("Press Ctrl+C to stop the server\n\n")

	return http.ListenAndServe(addr, nil)
}

// handleIndex serves the main UI page
func (s *UIServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(uiTemplates, "ui_templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"Port": s.port,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// handleHealth returns server health status
func (s *UIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"version": "0.1.0",
	})
}

// TranspileRequest represents a transpilation request
type TranspileRequest struct {
	TypeScript string `json:"typescript"`
	FileName   string `json:"fileName"`
}

// TranspileResponse represents a transpilation response
type TranspileResponse struct {
	Success bool   `json:"success"`
	Go      string `json:"go"`
	Error   string `json:"error,omitempty"`
}

// handleTranspile processes TypeScript code and returns Go code
func (s *UIServer) handleTranspile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TranspileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.TypeScript == "" {
		sendJSONError(w, "TypeScript code is required", http.StatusBadRequest)
		return
	}

	if req.FileName == "" {
		req.FileName = "input.ts"
	}

	// Create temporary directory for transpilation
	tempDir, err := os.MkdirTemp("", "ts2go-ui-*")
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to create temp directory: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// Write TypeScript file
	tsFile := filepath.Join(tempDir, req.FileName)
	if err := os.WriteFile(tsFile, []byte(req.TypeScript), 0644); err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to write TypeScript file: %v", err), http.StatusInternalServerError)
		return
	}

	// Create output directory
	outDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to create output directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Scan the project
	scanner := project.NewScanner()
	proj, err := scanner.ScanProject(tempDir)
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to scan project: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse and transpile each file
	for _, file := range proj.Files {
		ast, err := transpiler.ParseTypeScript(file.Path)
		if err != nil {
			sendJSONError(w, fmt.Sprintf("Failed to parse TypeScript: %v", err), http.StatusInternalServerError)
			return
		}

		// Analyze imports
		imports, err := analyzer.AnalyzeImports(ast)
		if err != nil {
			sendJSONError(w, fmt.Sprintf("Failed to analyze imports: %v", err), http.StatusInternalServerError)
			return
		}
		file.Imports = imports.Imports

		// Transpile to Go
		generator := transpiler.NewCodeGenerator()
		goCode, err := generator.Generate(ast)
		if err != nil {
			sendJSONError(w, fmt.Sprintf("Failed to transpile: %v", err), http.StatusInternalServerError)
			return
		}

		// Write output file
		baseNameNoExt := strings.TrimSuffix(filepath.Base(file.Path), ".ts")
		outFile := filepath.Join(outDir, baseNameNoExt+".go")
		if err := os.WriteFile(outFile, []byte(goCode), 0644); err != nil {
			sendJSONError(w, fmt.Sprintf("Failed to write output: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Read the generated Go code
	goFiles, err := findGoFiles(outDir)
	if err != nil || len(goFiles) == 0 {
		sendJSONError(w, "No Go files generated", http.StatusInternalServerError)
		return
	}

	// Read all generated Go files
	var goCode strings.Builder
	for i, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		if err != nil {
			continue
		}
		if i > 0 {
			goCode.WriteString("\n\n// ===============================\n\n")
		}
		goCode.WriteString(fmt.Sprintf("// File: %s\n", filepath.Base(goFile)))
		goCode.Write(content)
	}

	response := TranspileResponse{
		Success: true,
		Go:      goCode.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AnalyzeRequest represents an analyze request
type AnalyzeRequest struct {
	TypeScript string `json:"typescript"`
	FileName   string `json:"fileName"`
}

// AnalyzeResponse represents an analyze response
type AnalyzeResponse struct {
	Success      bool              `json:"success"`
	Dependencies []string          `json:"dependencies"`
	Imports      []string          `json:"imports"`
	Exports      []string          `json:"exports"`
	Stats        map[string]int    `json:"stats"`
	Warnings     []string          `json:"warnings"`
	Error        string            `json:"error,omitempty"`
}

// handleAnalyze analyzes TypeScript code and returns information
func (s *UIServer) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.TypeScript == "" {
		sendJSONError(w, "TypeScript code is required", http.StatusBadRequest)
		return
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "ts2go-analyze-*")
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to create temp directory: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// Write TypeScript file
	fileName := "input.ts"
	if req.FileName != "" {
		fileName = req.FileName
	}
	tsFile := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(tsFile, []byte(req.TypeScript), 0644); err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to write file: %v", err), http.StatusInternalServerError)
		return
	}

	// Scan project
	scanner := project.NewScanner()
	proj, err := scanner.ScanProject(tempDir)
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to scan: %v", err), http.StatusInternalServerError)
		return
	}

	// Build response
	response := AnalyzeResponse{
		Success:      true,
		Dependencies: []string{},
		Imports:      []string{},
		Exports:      []string{},
		Stats: map[string]int{
			"files":       len(proj.Files),
			"entryPoints": len(proj.EntryPoints),
		},
		Warnings: []string{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func findGoFiles(dir string) ([]string, error) {
	var goFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	return goFiles, err
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, freebsd, openbsd, netbsd
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
