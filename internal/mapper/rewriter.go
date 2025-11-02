package mapper

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yourusername/ts2go/internal/analyzer"
)

// ImportRewriter rewrites TypeScript imports to Go imports
type ImportRewriter struct {
	db         *MappingDatabase
	classifier *Classifier
	projectPkg string // Base package name for the project
}

// NewImportRewriter creates a new import rewriter
func NewImportRewriter(db *MappingDatabase, projectPackage string) *ImportRewriter {
	return &ImportRewriter{
		db:         db,
		classifier: NewClassifier(db),
		projectPkg: projectPackage,
	}
}

// RewriteResult contains the result of rewriting an import
type RewriteResult struct {
	OriginalSource string            // Original import source (e.g., 'axios')
	GoImport       string            // Go import path (e.g., 'github.com/go-resty/resty/v2')
	PackageAlias   string            // Suggested package alias (e.g., 'resty')
	SymbolMappings map[string]string // Symbol name mappings (TS -> Go)
	IsLocal        bool              // Whether this is a local file import
	NeedsTranspile bool              // Whether the source needs to be transpiled
	Error          error             // Error if rewrite failed
}

// RewriteImport rewrites a single import statement
func (r *ImportRewriter) RewriteImport(imp *analyzer.Import) *RewriteResult {
	result := &RewriteResult{
		OriginalSource: imp.Source,
		SymbolMappings: make(map[string]string),
	}

	// Handle local imports
	if imp.Type == analyzer.ImportTypeLocal {
		return r.rewriteLocalImport(imp)
	}

	// Handle built-in and package imports
	mapping, err := r.db.GetMapping(imp.Source)
	if err != nil {
		// Package not in database
		result.Error = fmt.Errorf("no mapping found for package: %s", imp.Source)
		result.NeedsTranspile = true
		return result
	}

	// Check if package is supported
	if !mapping.IsSupported() {
		result.Error = fmt.Errorf("package not supported: %s (reason: %s)", imp.Source, mapping.Reason)
		return result
	}

	result.GoImport = mapping.Go
	result.PackageAlias = r.determineAlias(mapping)
	result.SymbolMappings = r.mapSymbols(imp, mapping)

	return result
}

// RewriteImports rewrites multiple imports
func (r *ImportRewriter) RewriteImports(imports []*analyzer.Import) []*RewriteResult {
	results := make([]*RewriteResult, 0, len(imports))

	for _, imp := range imports {
		result := r.RewriteImport(imp)
		results = append(results, result)
	}

	return results
}

// GenerateGoImports generates the Go import block from rewrite results
func (r *ImportRewriter) GenerateGoImports(results []*RewriteResult) string {
	if len(results) == 0 {
		return ""
	}

	var imports []string
	var errors []string

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("// Error: %s", result.Error.Error()))
			continue
		}

		if result.GoImport != "" {
			imports = append(imports, fmt.Sprintf("\t\"%s\"", result.GoImport))
		}
	}

	if len(imports) == 0 && len(errors) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("import (\n")

	for _, imp := range imports {
		builder.WriteString(imp)
		builder.WriteString("\n")
	}

	if len(errors) > 0 {
		builder.WriteString("\n")
		for _, err := range errors {
			builder.WriteString("\t")
			builder.WriteString(err)
			builder.WriteString("\n")
		}
	}

	builder.WriteString(")")

	return builder.String()
}

// GetRequiredPackages returns a list of unique Go packages needed
func (r *ImportRewriter) GetRequiredPackages(results []*RewriteResult) []string {
	seen := make(map[string]bool)
	packages := make([]string, 0)

	for _, result := range results {
		if result.Error != nil {
			continue
		}
		if result.GoImport != "" && !seen[result.GoImport] {
			seen[result.GoImport] = true
			packages = append(packages, result.GoImport)
		}
	}

	return packages
}

// Helper methods

func (r *ImportRewriter) rewriteLocalImport(imp *analyzer.Import) *RewriteResult {
	result := &RewriteResult{
		OriginalSource: imp.Source,
		IsLocal:        true,
		NeedsTranspile: true,
		SymbolMappings: make(map[string]string),
	}

	// Convert relative path to Go package path
	// ./models/user -> projectpkg/models
	// ../utils/helper -> projectpkg/utils
	cleanPath := strings.TrimSuffix(imp.Source, ".ts")
	cleanPath = strings.TrimSuffix(cleanPath, ".js")

	if strings.HasPrefix(cleanPath, "./") {
		// Same directory
		relativePath := strings.TrimPrefix(cleanPath, "./")
		result.GoImport = filepath.Join(r.projectPkg, filepath.Dir(relativePath))
	} else if strings.HasPrefix(cleanPath, "../") {
		// Parent directory - simplified handling
		relativePath := strings.TrimPrefix(cleanPath, "../")
		result.GoImport = filepath.Join(r.projectPkg, filepath.Dir(relativePath))
	} else {
		// Absolute path or other format
		result.GoImport = filepath.Join(r.projectPkg, filepath.Dir(cleanPath))
	}

	// For local imports, symbols usually map directly (unless we add transformation)
	for _, symbol := range imp.Symbols {
		if symbol != "default" {
			result.SymbolMappings[symbol] = symbol
		}
	}

	result.PackageAlias = filepath.Base(result.GoImport)

	return result
}

func (r *ImportRewriter) determineAlias(mapping *Mapping) string {
	// Use the last component of the Go package path as alias
	parts := strings.Split(mapping.Go, "/")
	if len(parts) == 0 {
		return ""
	}

	lastPart := parts[len(parts)-1]

	// Remove version suffix if present (e.g., v2, v3)
	if strings.HasPrefix(lastPart, "v") && len(lastPart) > 1 {
		if len(parts) > 1 {
			lastPart = parts[len(parts)-2]
		}
	}

	return lastPart
}

func (r *ImportRewriter) mapSymbols(imp *analyzer.Import, mapping *Mapping) map[string]string {
	symbolMappings := make(map[string]string)

	// If it's a namespace import (import * as fs), map to package name
	if imp.IsNamespace {
		alias := r.determineAlias(mapping)
		symbolMappings["*"] = alias
		return symbolMappings
	}

	// For named imports, try to find API mappings
	for _, symbol := range imp.Symbols {
		if symbol == "default" {
			// Default import maps to package alias
			alias := r.determineAlias(mapping)
			symbolMappings["default"] = alias
			continue
		}

		// Try to find API mapping for this symbol
		// Look for patterns like "package.symbol" in API mappings
		apiKey := fmt.Sprintf("%s.%s", imp.Source, symbol)
		if goAPI, ok := mapping.APIMappings[apiKey]; ok {
			// Extract just the function name from patterns like "resty.R().Get"
			symbolMappings[symbol] = extractFunctionName(goAPI)
		} else {
			// No mapping found, capitalize first letter (Go convention)
			symbolMappings[symbol] = capitalizeFirst(symbol)
		}
	}

	return symbolMappings
}

func extractFunctionName(goAPI string) string {
	// Extract the last part of a chain like "resty.R().Get" -> "Get"
	parts := strings.Split(goAPI, ".")
	if len(parts) == 0 {
		return goAPI
	}

	lastPart := parts[len(parts)-1]

	// Remove parentheses if present
	lastPart = strings.TrimSuffix(lastPart, "()")

	return lastPart
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// AnalysisResult combines import analysis with rewrite results
type AnalysisResult struct {
	Imports        []*analyzer.Import
	RewriteResults []*RewriteResult
	GoImports      string
	RequiredPkgs   []string
	Errors         []error
}

// AnalyzeAndRewrite performs complete import analysis and rewriting
func (r *ImportRewriter) AnalyzeAndRewriteFromAnalysis(analysis *analyzer.ImportAnalysis) (*AnalysisResult, error) {
	if analysis == nil {
		return nil, fmt.Errorf("import analysis is nil")
	}

	// Get all imports
	allImports := make([]*analyzer.Import, len(analysis.Imports))
	for i := range analysis.Imports {
		allImports[i] = &analysis.Imports[i]
	}

	// Rewrite each import
	rewriteResults := r.RewriteImports(allImports)

	// Collect errors
	var errors []error
	for _, result := range rewriteResults {
		if result.Error != nil {
			errors = append(errors, result.Error)
		}
	}

	// Generate Go import block
	goImports := r.GenerateGoImports(rewriteResults)

	// Get required packages
	requiredPkgs := r.GetRequiredPackages(rewriteResults)

	return &AnalysisResult{
		Imports:        allImports,
		RewriteResults: rewriteResults,
		GoImports:      goImports,
		RequiredPkgs:   requiredPkgs,
		Errors:         errors,
	}, nil
}

// Summary returns a summary of the analysis
func (ar *AnalysisResult) Summary() string {
	totalImports := len(ar.Imports)
	successful := totalImports - len(ar.Errors)

	return fmt.Sprintf(
		"Import Analysis Summary\n"+
			"Total Imports: %d\n"+
			"Successfully Rewritten: %d\n"+
			"Errors: %d\n"+
			"Required Go Packages: %d\n",
		totalImports,
		successful,
		len(ar.Errors),
		len(ar.RequiredPkgs),
	)
}
