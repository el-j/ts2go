package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/ts2go/internal/analyzer"
	"github.com/yourusername/ts2go/internal/mapper"
)

// Project represents a TypeScript project
type Project struct {
	RootDir      string
	Name         string
	PackageJSON  *analyzer.PackageInfo
	Files        []*SourceFile
	Dependencies map[string]*mapper.Classification
	EntryPoints  []string
	OutputDir    string
}

// SourceFile represents a single TypeScript source file
type SourceFile struct {
	Path         string // Absolute path
	RelativePath string // Relative to project root
	Imports      []analyzer.Import
	IsEntry      bool // Is this an entry point?
}

// Scanner scans TypeScript projects
type Scanner struct {
	ignorePatterns []string
}

// NewScanner creates a new project scanner
func NewScanner() *Scanner {
	return &Scanner{
		ignorePatterns: []string{
			"node_modules",
			".git",
			"dist",
			"build",
			"coverage",
			"*.test.ts",
			"*.spec.ts",
		},
	}
}

// ScanProject scans a TypeScript project directory
func (s *Scanner) ScanProject(rootDir string) (*Project, error) {
	// Validate directory exists
	info, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", rootDir)
	}

	project := &Project{
		RootDir:      rootDir,
		Files:        []*SourceFile{},
		Dependencies: make(map[string]*mapper.Classification),
		EntryPoints:  []string{},
	}

	// Try to load package.json
	packageJSONPath := filepath.Join(rootDir, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		pkgInfo, err := analyzer.ParsePackageJSON(packageJSONPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse package.json: %w", err)
		}
		project.PackageJSON = pkgInfo
		project.Name = pkgInfo.Name
	} else {
		// Use directory name as project name
		project.Name = filepath.Base(rootDir)
	}

	// Find all TypeScript files
	err = s.findTypeScriptFiles(project)
	if err != nil {
		return nil, fmt.Errorf("failed to find TypeScript files: %w", err)
	}

	// Identify entry points
	s.identifyEntryPoints(project)

	return project, nil
}

// findTypeScriptFiles recursively finds all .ts files
func (s *Scanner) findTypeScriptFiles(project *Project) error {
	return filepath.Walk(project.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(project.RootDir, path)
		if err != nil {
			return err
		}

		// Check ignore patterns
		if s.shouldIgnore(relPath, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if it's a TypeScript file
		if !info.IsDir() && s.isTypeScriptFile(path) {
			sourceFile := &SourceFile{
				Path:         path,
				RelativePath: relPath,
			}
			project.Files = append(project.Files, sourceFile)
		}

		return nil
	})
}

// shouldIgnore checks if a path should be ignored
func (s *Scanner) shouldIgnore(relPath string, info os.FileInfo) bool {
	// Check each ignore pattern
	for _, pattern := range s.ignorePatterns {
		// For directories, check if the directory name matches
		if info.IsDir() {
			if strings.HasPrefix(pattern, "*") {
				continue // Skip wildcard patterns for directories
			}
			if filepath.Base(relPath) == pattern || relPath == pattern {
				return true
			}
		} else {
			// For files, check patterns
			matched, _ := filepath.Match(pattern, filepath.Base(relPath))
			if matched {
				return true
			}
			// Also check if file is in an ignored directory
			parts := strings.Split(relPath, string(os.PathSeparator))
			for _, part := range parts {
				if part == pattern {
					return true
				}
			}
		}
	}
	return false
}

// isTypeScriptFile checks if a file is a TypeScript file
func (s *Scanner) isTypeScriptFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".ts" && !strings.HasSuffix(path, ".d.ts")
}

// identifyEntryPoints identifies entry point files
func (s *Scanner) identifyEntryPoints(project *Project) {
	// Common entry point patterns
	entryPatterns := []string{
		"index.ts",
		"main.ts",
		"app.ts",
		"server.ts",
		"cli.ts",
	}

	for _, file := range project.Files {
		basename := filepath.Base(file.Path)

		// Check if it matches an entry pattern
		for _, pattern := range entryPatterns {
			if basename == pattern {
				file.IsEntry = true
				project.EntryPoints = append(project.EntryPoints, file.RelativePath)
				break
			}
		}

		// Check if it's in the root directory (likely an entry point)
		if filepath.Dir(file.RelativePath) == "." {
			if !file.IsEntry {
				file.IsEntry = true
				project.EntryPoints = append(project.EntryPoints, file.RelativePath)
			}
		}
	}

	// If no entry points found, mark all root-level files as entry points
	if len(project.EntryPoints) == 0 {
		for _, file := range project.Files {
			if filepath.Dir(file.RelativePath) == "." {
				file.IsEntry = true
				project.EntryPoints = append(project.EntryPoints, file.RelativePath)
			}
		}
	}
}

// AddIgnorePattern adds a custom ignore pattern
func (s *Scanner) AddIgnorePattern(pattern string) {
	s.ignorePatterns = append(s.ignorePatterns, pattern)
}

// GetFileCount returns the number of TypeScript files found
func (p *Project) GetFileCount() int {
	return len(p.Files)
}

// GetEntryPointCount returns the number of entry points
func (p *Project) GetEntryPointCount() int {
	return len(p.EntryPoints)
}

// GetFileByPath returns a source file by its relative path
func (p *Project) GetFileByPath(relPath string) *SourceFile {
	for _, file := range p.Files {
		if file.RelativePath == relPath {
			return file
		}
	}
	return nil
}

// Summary returns a human-readable summary of the project
func (p *Project) Summary() string {
	depCount := 0
	if p.PackageJSON != nil {
		depCount = len(p.PackageJSON.Dependencies) + len(p.PackageJSON.DevDependencies)
	}

	return fmt.Sprintf(
		"Project: %s\n"+
			"Root: %s\n"+
			"TypeScript Files: %d\n"+
			"Entry Points: %d\n"+
			"Dependencies: %d\n",
		p.Name,
		p.RootDir,
		len(p.Files),
		len(p.EntryPoints),
		depCount,
	)
}
