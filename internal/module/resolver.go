package module

import (
	"fmt"
	"sort"
	"strings"
)

// ImportResolver resolves TypeScript imports to Go imports
type ImportResolver struct {
	registry      *ExportRegistry
	mapper        *PackageMapper
	currentFile   string
	currentModule *Module
}

// NewImportResolver creates a new import resolver
func NewImportResolver(registry *ExportRegistry, mapper *PackageMapper) *ImportResolver {
	return &ImportResolver{
		registry: registry,
		mapper:   mapper,
	}
}

// SetCurrentFile sets the current file being processed
func (r *ImportResolver) SetCurrentFile(filePath string) {
	r.currentFile = filePath
	r.currentModule = r.registry.GetModule(filePath)
}

// ResolvedImport represents a resolved Go import
type ResolvedImport struct {
	// PackagePath is the full Go import path (e.g., "github.com/user/project/models")
	PackagePath string
	// PackageAlias is the package alias if needed (e.g., "models2" for conflicts)
	PackageAlias string
	// IsStdLib indicates if this is a standard library import
	IsStdLib bool
	// IsRuntime indicates if this is a ts2go runtime import
	IsRuntime bool
	// Symbols is the list of symbols imported from this package
	Symbols []string
}

// ResolveImports resolves all imports for the current file
func (r *ImportResolver) ResolveImports() ([]ResolvedImport, error) {
	if r.currentModule == nil {
		return nil, fmt.Errorf("no module set")
	}

	// Group imports by package
	importMap := make(map[string]*ResolvedImport)

	for _, imp := range r.currentModule.Imports {
		// Skip type-only imports (they don't need runtime imports)
		if imp.IsType {
			continue
		}

		// Check if it's a relative import
		if isRelativeImport(imp.Source) {
			// Resolve to Go package path
			packagePath, ok := r.mapper.ResolveImportPath(r.currentFile, imp.Source)
			if !ok {
				return nil, fmt.Errorf("failed to resolve import %s", imp.Source)
			}

			// Skip imports from the same package
			currentPackagePath := r.mapper.GetPackagePath(r.currentFile)
			if packagePath == currentPackagePath {
				continue
			}

			// Add to import map
			if _, exists := importMap[packagePath]; !exists {
				importMap[packagePath] = &ResolvedImport{
					PackagePath: packagePath,
					Symbols:     []string{},
				}
			}

			// Add symbol if it's a named import
			if imp.Type == ImportNamed || imp.Type == ImportDefault {
				importMap[packagePath].Symbols = append(importMap[packagePath].Symbols, imp.Name)
			}
		} else {
			// It's an npm package - check if we have a runtime mapping
			runtimePath := r.resolveNpmPackage(imp.Source)
			if runtimePath != "" {
				if _, exists := importMap[runtimePath]; !exists {
					importMap[runtimePath] = &ResolvedImport{
						PackagePath: runtimePath,
						IsRuntime:   strings.HasPrefix(runtimePath, "github.com/yourusername/ts2go-runtime"),
						IsStdLib:    isStdLibPackage(runtimePath),
						Symbols:     []string{},
					}
				}

				if imp.Type == ImportNamed || imp.Type == ImportDefault {
					importMap[runtimePath].Symbols = append(importMap[runtimePath].Symbols, imp.Name)
				}
			}
		}
	}

	// Convert map to slice and sort
	resolved := make([]ResolvedImport, 0, len(importMap))
	for _, imp := range importMap {
		resolved = append(resolved, *imp)
	}

	// Sort: stdlib first, then runtime, then project packages
	sort.Slice(resolved, func(i, j int) bool {
		if resolved[i].IsStdLib != resolved[j].IsStdLib {
			return resolved[i].IsStdLib
		}
		if resolved[i].IsRuntime != resolved[j].IsRuntime {
			return resolved[i].IsRuntime
		}
		return resolved[i].PackagePath < resolved[j].PackagePath
	})

	return resolved, nil
}

// GenerateImportBlock generates the Go import block
func (r *ImportResolver) GenerateImportBlock() (string, error) {
	resolved, err := r.ResolveImports()
	if err != nil {
		return "", err
	}

	if len(resolved) == 0 {
		return "", nil
	}

	var builder strings.Builder
	builder.WriteString("import (\n")

	// Group by type with blank lines between
	hasStdLib := false
	hasRuntime := false

	for _, imp := range resolved {
		// Add blank line before runtime imports if we had stdlib
		if !imp.IsStdLib && hasStdLib && !hasRuntime && imp.IsRuntime {
			builder.WriteString("\n")
			hasRuntime = true
		}

		// Add blank line before project imports if we had runtime or stdlib
		if !imp.IsStdLib && !imp.IsRuntime && (hasStdLib || hasRuntime) {
			builder.WriteString("\n")
		}

		if imp.IsStdLib {
			hasStdLib = true
		}
		if imp.IsRuntime {
			hasRuntime = true
		}

		// Generate import line
		if imp.PackageAlias != "" {
			builder.WriteString(fmt.Sprintf("\t%s \"%s\"\n", imp.PackageAlias, imp.PackagePath))
		} else {
			builder.WriteString(fmt.Sprintf("\t\"%s\"\n", imp.PackagePath))
		}
	}

	builder.WriteString(")")
	return builder.String(), nil
}

// resolveNpmPackage maps npm package names to Go packages
func (r *ImportResolver) resolveNpmPackage(npmPackage string) string {
	// Map common npm packages to Go equivalents
	mapping := map[string]string{
		"fs":      "github.com/yourusername/ts2go-runtime/fs",
		"path":    "github.com/yourusername/ts2go-runtime/path",
		"console": "fmt", // stdlib
		"axios":   "github.com/yourusername/ts2go-runtime/axios",
		"lodash":  "github.com/yourusername/ts2go-runtime/lodash",
	}

	if goPath, exists := mapping[npmPackage]; exists {
		return goPath
	}

	// Unknown package - return empty (will be skipped)
	return ""
}

// isStdLibPackage checks if a package is from Go standard library
func isStdLibPackage(packagePath string) bool {
	stdLibPackages := []string{
		"fmt", "os", "io", "strings", "strconv", "time", "sync",
		"path", "path/filepath", "net", "net/http", "encoding/json",
		"errors", "context", "bytes", "bufio", "crypto", "math",
	}

	for _, pkg := range stdLibPackages {
		if packagePath == pkg || strings.HasPrefix(packagePath, pkg+"/") {
			return true
		}
	}

	return false
}

// GetImportedSymbols returns the list of symbols imported from a specific package
func (r *ImportResolver) GetImportedSymbols(packagePath string) []string {
	if r.currentModule == nil {
		return nil
	}

	symbols := []string{}

	for _, imp := range r.currentModule.Imports {
		if isRelativeImport(imp.Source) {
			resolvedPath, ok := r.mapper.ResolveImportPath(r.currentFile, imp.Source)
			if ok && resolvedPath == packagePath {
				if imp.Type == ImportNamed || imp.Type == ImportDefault {
					symbols = append(symbols, imp.Name)
				}
			}
		}
	}

	return symbols
}

// ShouldImportPackage checks if a package should be imported
func (r *ImportResolver) ShouldImportPackage(packagePath string) bool {
	// Don't import the same package
	currentPackagePath := r.mapper.GetPackagePath(r.currentFile)
	if packagePath == currentPackagePath {
		return false
	}

	// Check if we actually use this package
	symbols := r.GetImportedSymbols(packagePath)
	return len(symbols) > 0
}
