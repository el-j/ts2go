package analyzer

import (
	"fmt"
	"strings"

	"github.com/el-j/ts2go/internal/transpiler"
)

// ImportType represents the type of import
type ImportType int

const (
	ImportTypePackage ImportType = iota // npm package
	ImportTypeLocal                     // ./file or ../file
	ImportTypeBuiltin                   // fs, path, etc.
)

func (t ImportType) String() string {
	switch t {
	case ImportTypePackage:
		return "package"
	case ImportTypeLocal:
		return "local"
	case ImportTypeBuiltin:
		return "builtin"
	default:
		return "unknown"
	}
}

// Import represents a single import statement
type Import struct {
	Source      string     // 'axios', './user', 'fs'
	Type        ImportType // Package, Local, Builtin
	Symbols     []string   // ['default'], ['get', 'post'], ['*']
	IsNamespace bool       // import * as fs
	IsDefault   bool       // import axios from 'axios'
	Position    int        // Line number for error reporting
	Alias       string     // For namespace imports: import * as name
}

// ImportAnalysis contains all imports found in a file
type ImportAnalysis struct {
	Imports        []Import
	LocalFiles     []string          // List of local files that need to be transpiled
	NpmPackages    []string          // List of npm packages needed
	BuiltinModules []string          // List of Node.js built-ins used
	ImportMap      map[string]Import // Quick lookup by source
}

// AnalyzeImports analyzes all import statements in an AST
func AnalyzeImports(astNode *transpiler.ASTNode) (*ImportAnalysis, error) {
	if astNode == nil {
		return nil, fmt.Errorf("AST node is nil")
	}

	analysis := &ImportAnalysis{
		Imports:        []Import{},
		LocalFiles:     []string{},
		NpmPackages:    []string{},
		BuiltinModules: []string{},
		ImportMap:      make(map[string]Import),
	}

	// Process all statements
	if astNode.Statements != nil {
		for i, stmt := range astNode.Statements {
			if stmt.Kind == "ImportDeclaration" {
				imp, err := parseImportDeclaration(&stmt, i+1)
				if err != nil {
					return nil, fmt.Errorf("failed to parse import at line %d: %w", i+1, err)
				}

				analysis.Imports = append(analysis.Imports, imp)
				analysis.ImportMap[imp.Source] = imp

				// Categorize import
				switch imp.Type {
				case ImportTypeLocal:
					analysis.LocalFiles = append(analysis.LocalFiles, imp.Source)
				case ImportTypePackage:
					analysis.NpmPackages = append(analysis.NpmPackages, imp.Source)
				case ImportTypeBuiltin:
					analysis.BuiltinModules = append(analysis.BuiltinModules, imp.Source)
				}
			}
		}
	}

	return analysis, nil
}

// parseImportDeclaration parses a single import declaration AST node
func parseImportDeclaration(node *transpiler.ASTNode, lineNum int) (Import, error) {
	imp := Import{
		Position: lineNum,
		Symbols:  []string{},
	}

	// Extract module specifier (the 'from' part)
	// Look for StringLiteral in children
	for _, child := range node.Children {
		if child.Kind == "StringLiteral" {
			imp.Source = child.Text
			break
		}
	}

	if imp.Source == "" {
		return imp, fmt.Errorf("import source not found")
	}

	// Determine import type
	imp.Type = classifyImportSource(imp.Source)

	// Extract import clause (what's being imported)
	for _, child := range node.Children {
		if child.Kind == "ImportClause" {
			parseImportClause(&child, &imp)
		}
	}

	return imp, nil
}

// parseImportClause extracts symbols from import clause
func parseImportClause(node *transpiler.ASTNode, imp *Import) {
	for _, child := range node.Children {
		switch child.Kind {
		case "Identifier":
			// Default import: import axios from 'axios'
			imp.IsDefault = true
			imp.Symbols = append(imp.Symbols, child.Text)

		case "NamespaceImport":
			// Namespace import: import * as fs from 'fs'
			imp.IsNamespace = true
			// Find the identifier after 'as'
			for _, nameChild := range child.Children {
				if nameChild.Kind == "Identifier" {
					imp.Alias = nameChild.Text
					imp.Symbols = append(imp.Symbols, "*")
					break
				}
			}

		case "NamedImports":
			// Named imports: import { a, b } from 'module'
			for _, nameChild := range child.Children {
				if nameChild.Kind == "ImportSpecifier" {
					// ImportSpecifier can have name or propertyName
					if nameChild.Name != "" {
						imp.Symbols = append(imp.Symbols, nameChild.Name)
					}
				}
			}
		}
	}
}

// classifyImportSource determines the type of import based on the source
func classifyImportSource(source string) ImportType {
	// Remove any "node:" prefix for built-in checking
	checkSource := source
	if strings.HasPrefix(source, "node:") {
		checkSource = source[5:]
	}

	// Check if it's a built-in module
	if IsBuiltinModule(checkSource) {
		return ImportTypeBuiltin
	}

	// Check if it's a local file (starts with . or /)
	if strings.HasPrefix(source, "./") ||
		strings.HasPrefix(source, "../") ||
		strings.HasPrefix(source, "/") {
		return ImportTypeLocal
	}

	// Otherwise it's an npm package
	return ImportTypePackage
}

// GetImportBySource finds an import by its source
func (a *ImportAnalysis) GetImportBySource(source string) (Import, bool) {
	imp, exists := a.ImportMap[source]
	return imp, exists
}

// HasPackage checks if a specific npm package is imported
func (a *ImportAnalysis) HasPackage(packageName string) bool {
	for _, pkg := range a.NpmPackages {
		if pkg == packageName {
			return true
		}
	}
	return false
}

// HasBuiltin checks if a specific built-in module is imported
func (a *ImportAnalysis) HasBuiltin(moduleName string) bool {
	for _, mod := range a.BuiltinModules {
		if mod == moduleName {
			return true
		}
	}
	return false
}

// GetDependencies returns all external dependencies (npm packages + built-ins)
func (a *ImportAnalysis) GetDependencies() []string {
	deps := make([]string, 0, len(a.NpmPackages)+len(a.BuiltinModules))
	deps = append(deps, a.NpmPackages...)
	deps = append(deps, a.BuiltinModules...)
	return deps
}

// Summary returns a human-readable summary of the analysis
func (a *ImportAnalysis) Summary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total imports: %d\n", len(a.Imports)))
	sb.WriteString(fmt.Sprintf("  - npm packages: %d\n", len(a.NpmPackages)))
	sb.WriteString(fmt.Sprintf("  - built-in modules: %d\n", len(a.BuiltinModules)))
	sb.WriteString(fmt.Sprintf("  - local files: %d\n", len(a.LocalFiles)))

	if len(a.NpmPackages) > 0 {
		sb.WriteString("\nNPM Packages:\n")
		for _, pkg := range a.NpmPackages {
			sb.WriteString(fmt.Sprintf("  - %s\n", pkg))
		}
	}

	if len(a.BuiltinModules) > 0 {
		sb.WriteString("\nBuilt-in Modules:\n")
		for _, mod := range a.BuiltinModules {
			sb.WriteString(fmt.Sprintf("  - %s\n", mod))
		}
	}

	if len(a.LocalFiles) > 0 {
		sb.WriteString("\nLocal Files:\n")
		for _, file := range a.LocalFiles {
			sb.WriteString(fmt.Sprintf("  - %s\n", file))
		}
	}

	return sb.String()
}

// AnalyzeFile is a convenience function to analyze imports from a TypeScript file
func AnalyzeFile(filePath string) (*ImportAnalysis, error) {
	// Parse the TypeScript file using the transpiler's parser
	ast, err := transpiler.ParseTypeScript(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TypeScript: %w", err)
	}

	// Analyze the imports
	return AnalyzeImports(ast)
}
