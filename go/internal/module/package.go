package module

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PackageMapper maps TypeScript file paths to Go package names
type PackageMapper struct {
	// projectRoot is the root directory of the project
	projectRoot string
	// moduleName is the Go module name (e.g., "github.com/user/project")
	moduleName string
	// srcDir is the source directory (e.g., "src")
	srcDir string
}

// NewPackageMapper creates a new package mapper
func NewPackageMapper(projectRoot, moduleName, srcDir string) *PackageMapper {
	return &PackageMapper{
		projectRoot: projectRoot,
		moduleName:  moduleName,
		srcDir:      srcDir,
	}
}

// GetPackageName returns the Go package name for a TypeScript file
func (m *PackageMapper) GetPackageName(tsFilePath string) string {
	// Get relative path from project root
	relPath, err := filepath.Rel(m.projectRoot, tsFilePath)
	if err != nil {
		// Fallback to base filename
		return sanitizePackageName(filepath.Base(filepath.Dir(tsFilePath)))
	}

	// Remove file extension
	dir := filepath.Dir(relPath)

	// Handle root-level files (index.ts, main.ts)
	if dir == "." || dir == m.srcDir {
		return "main"
	}

	// Remove src prefix if present
	if strings.HasPrefix(dir, m.srcDir+string(filepath.Separator)) {
		dir = strings.TrimPrefix(dir, m.srcDir+string(filepath.Separator))
	}

	// Get the last directory component for the package name
	parts := strings.Split(dir, string(filepath.Separator))
	packageName := parts[len(parts)-1]

	return sanitizePackageName(packageName)
}

// GetPackagePath returns the full Go import path for a TypeScript file
func (m *PackageMapper) GetPackagePath(tsFilePath string) string {
	// Get relative path from project root
	relPath, err := filepath.Rel(m.projectRoot, tsFilePath)
	if err != nil {
		return m.moduleName
	}

	// Remove file extension
	dir := filepath.Dir(relPath)

	// Handle root-level files
	if dir == "." || dir == m.srcDir {
		return m.moduleName
	}

	// Remove src prefix if present
	if strings.HasPrefix(dir, m.srcDir+string(filepath.Separator)) {
		dir = strings.TrimPrefix(dir, m.srcDir+string(filepath.Separator))
	}

	// Convert path separators to /
	dir = filepath.ToSlash(dir)

	return m.moduleName + "/" + dir
}

// GetOutputPath returns the output file path for a TypeScript file
func (m *PackageMapper) GetOutputPath(tsFilePath, outputDir string) string {
	// Get relative path from project root
	relPath, err := filepath.Rel(m.projectRoot, tsFilePath)
	if err != nil {
		// Fallback to filename
		return filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(tsFilePath), ".ts")+".go")
	}

	// Remove file extension and add .go
	relPath = strings.TrimSuffix(relPath, ".ts") + ".go"

	// Handle index.ts -> main.go for root
	dir := filepath.Dir(relPath)
	base := filepath.Base(relPath)

	if base == "index.go" && (dir == "." || dir == m.srcDir) {
		base = "main.go"
	}

	return filepath.Join(outputDir, dir, base)
}

// ResolveImportPath resolves a TypeScript import to a Go import path
func (m *PackageMapper) ResolveImportPath(currentFile, importSource string) (string, bool) {
	// Check if it's a relative import
	if !isRelativeImport(importSource) {
		// It's an npm package - not a local import
		return "", false
	}

	// Resolve relative import
	currentDir := filepath.Dir(currentFile)

	// Clean the import source (remove ./ and ../ and add .ts extension if missing)
	importPath := filepath.Join(currentDir, importSource)
	if !strings.HasSuffix(importPath, ".ts") {
		importPath += ".ts"
	}

	// Get the package path for the resolved file
	packagePath := m.GetPackagePath(importPath)

	return packagePath, true
}

// sanitizePackageName ensures the package name is valid for Go
func sanitizePackageName(name string) string {
	// Replace invalid characters
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A') // Convert to lowercase
		}
		return '_'
	}, name)

	// Ensure it doesn't start with a number
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		name = "pkg_" + name
	}

	// Ensure it's not empty
	if name == "" {
		name = "pkg"
	}

	return name
}

// PackageStructure represents the Go package structure
type PackageStructure struct {
	// Packages maps package path to list of files
	Packages map[string][]string
	// MainPackage is the main package path
	MainPackage string
	// EntryPoint is the main file path
	EntryPoint string
}

// GeneratePackageStructure generates the Go package structure from TypeScript files
func GeneratePackageStructure(files []string, mapper *PackageMapper) (*PackageStructure, error) {
	structure := &PackageStructure{
		Packages: make(map[string][]string),
	}

	// Group files by package
	for _, file := range files {
		packagePath := mapper.GetPackagePath(file)
		structure.Packages[packagePath] = append(structure.Packages[packagePath], file)

		// Check if this is a main/entry file
		base := filepath.Base(file)
		if base == "index.ts" || base == "main.ts" {
			dir := filepath.Dir(file)
			// Check if it's at root or src level
			relPath, _ := filepath.Rel(mapper.projectRoot, dir)
			if relPath == "." || relPath == mapper.srcDir {
				structure.MainPackage = packagePath
				structure.EntryPoint = file
			}
		}
	}

	// If no entry point found, use the first file
	if structure.EntryPoint == "" && len(files) > 0 {
		structure.EntryPoint = files[0]
		structure.MainPackage = mapper.GetPackagePath(files[0])
	}

	return structure, nil
}

// GetPackageDeclaration returns the package declaration for a TypeScript file
func GetPackageDeclaration(tsFilePath string, mapper *PackageMapper, isEntryPoint bool) string {
	if isEntryPoint {
		return "package main"
	}
	return fmt.Sprintf("package %s", mapper.GetPackageName(tsFilePath))
}
