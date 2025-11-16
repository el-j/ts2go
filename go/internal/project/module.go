package project

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/el-j/ts2go/internal/mapper"
)

// ModuleGenerator handles go.mod generation
type ModuleGenerator struct {
	project    *Project
	classifier *mapper.Classifier
}

// NewModuleGenerator creates a new module generator
func NewModuleGenerator(project *Project, classifier *mapper.Classifier) *ModuleGenerator {
	return &ModuleGenerator{
		project:    project,
		classifier: classifier,
	}
}

// GenerateGoModule generates a go.mod file for the project
func (mg *ModuleGenerator) GenerateGoModule(outputDir string) error {
	// Determine module name
	moduleName := mg.determineModuleName()

	// Determine Go version
	goVersion := mg.determineGoVersion()

	// Collect Go dependencies
	dependencies, err := mg.collectDependencies()
	if err != nil {
		return fmt.Errorf("failed to collect dependencies: %w", err)
	}

	// Generate go.mod content
	content := mg.formatGoMod(moduleName, goVersion, dependencies)

	// Write to file
	goModPath := filepath.Join(outputDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	return nil
}

// determineModuleName determines the Go module name
func (mg *ModuleGenerator) determineModuleName() string {
	// Use package.json name if available
	if mg.project.PackageJSON != nil && mg.project.PackageJSON.Name != "" {
		name := mg.project.PackageJSON.Name
		// Convert scoped packages
		name = strings.ReplaceAll(name, "@", "")
		name = strings.ReplaceAll(name, "/", "-")
		return "github.com/yourusername/" + name
	}

	// Use project name
	return "github.com/yourusername/" + mg.project.Name
}

// determineGoVersion determines the Go version to use
func (mg *ModuleGenerator) determineGoVersion() string {
	// Use a recent stable version
	return "1.21"
}

// collectDependencies collects all Go package dependencies
func (mg *ModuleGenerator) collectDependencies() (map[string]string, error) {
	dependencies := make(map[string]string)

	// Iterate through all project dependencies
	for npmPkg, classification := range mg.project.Dependencies {
		// Skip unsupported packages
		if classification.Class == mapper.PackageClassUnsupported {
			continue
		}

		// Get Go package name
		goPackage := classification.GoPackage
		if goPackage == "" {
			continue
		}

		// Skip standard library packages
		if isStdlibPackage(goPackage) {
			continue
		}

		// Add to dependencies with version
		version := mg.getPackageVersion(npmPkg, goPackage)
		dependencies[goPackage] = version
	}

	return dependencies, nil
}

// getPackageVersion determines the version for a Go package
func (mg *ModuleGenerator) getPackageVersion(_ /* npmPackage */, goPackage string) string {
	// Default version mapping
	versions := map[string]string{
		"github.com/go-resty/resty/v2":  "v2.11.0",
		"github.com/gin-gonic/gin":      "v1.9.1",
		"github.com/gorilla/mux":        "v1.8.1",
		"github.com/stretchr/testify":   "v1.8.4",
		"github.com/sirupsen/logrus":    "v1.9.3",
		"go.uber.org/zap":               "v1.26.0",
		"github.com/spf13/cobra":        "v1.8.0",
		"github.com/google/uuid":        "v1.5.0",
		"github.com/lib/pq":             "v1.10.9",
		"go.mongodb.org/mongo-driver":   "v1.13.1",
		"github.com/go-redis/redis/v8":  "v8.11.5",
		"gopkg.in/yaml.v3":              "v3.0.1",
		"github.com/el-j/ts2go-runtime": "v0.1.0",
	}

	if version, ok := versions[goPackage]; ok {
		return version
	}

	// Default to latest
	return "latest"
}

// formatGoMod formats the go.mod file content
func (mg *ModuleGenerator) formatGoMod(moduleName, goVersion string, dependencies map[string]string) string {
	var sb strings.Builder

	// Module declaration
	sb.WriteString(fmt.Sprintf("module %s\n\n", moduleName))

	// Go version
	sb.WriteString(fmt.Sprintf("go %s\n", goVersion))

	// Dependencies
	if len(dependencies) > 0 {
		sb.WriteString("\nrequire (\n")

		// Sort dependencies for consistent output
		packages := make([]string, 0, len(dependencies))
		for pkg := range dependencies {
			packages = append(packages, pkg)
		}
		sort.Strings(packages)

		// Write each dependency
		for _, pkg := range packages {
			version := dependencies[pkg]
			sb.WriteString(fmt.Sprintf("\t%s %s\n", pkg, version))
		}

		sb.WriteString(")\n")
	}

	return sb.String()
}

// isStdlibPackage checks if a package is in the Go standard library
func isStdlibPackage(pkg string) bool {
	// Standard library packages don't have a domain
	return !strings.Contains(pkg, ".")
}

// DependencySummary provides a summary of dependencies
type DependencySummary struct {
	TotalNPM         int
	TotalGo          int
	Supported        int
	Partial          int
	Unsupported      int
	StdlibPackages   []string
	ExternalPackages []string
	RuntimePackages  []string
	UnsupportedList  []string
}

// GetDependencySummary returns a summary of the project's dependencies
func (mg *ModuleGenerator) GetDependencySummary() *DependencySummary {
	summary := &DependencySummary{
		StdlibPackages:   []string{},
		ExternalPackages: []string{},
		RuntimePackages:  []string{},
		UnsupportedList:  []string{},
	}

	summary.TotalNPM = len(mg.project.Dependencies)

	for npmPkg, classification := range mg.project.Dependencies {
		goPackage := classification.GoPackage

		switch classification.Class {
		case mapper.PackageClassBuiltin:
			summary.Supported++
			if goPackage != "" && isStdlibPackage(goPackage) {
				summary.StdlibPackages = append(summary.StdlibPackages, goPackage)
			}

		case mapper.PackageClassRuntime:
			summary.Supported++
			summary.RuntimePackages = append(summary.RuntimePackages, goPackage)

		case mapper.PackageClassEquivalent:
			summary.Supported++
			if goPackage != "" && !isStdlibPackage(goPackage) {
				summary.ExternalPackages = append(summary.ExternalPackages, goPackage)
			}

		case mapper.PackageClassFramework:
			summary.Partial++
			if goPackage != "" {
				summary.ExternalPackages = append(summary.ExternalPackages, goPackage)
			}

		case mapper.PackageClassTranspilable:
			summary.Partial++

		case mapper.PackageClassUnsupported:
			summary.Unsupported++
			summary.UnsupportedList = append(summary.UnsupportedList, npmPkg)
		}
	}

	// Calculate total Go packages
	summary.TotalGo = len(summary.StdlibPackages) + len(summary.ExternalPackages) + len(summary.RuntimePackages)

	// Sort lists
	sort.Strings(summary.StdlibPackages)
	sort.Strings(summary.ExternalPackages)
	sort.Strings(summary.RuntimePackages)
	sort.Strings(summary.UnsupportedList)

	return summary
}

// FormatSummary formats the dependency summary as a human-readable string
func (ds *DependencySummary) FormatSummary() string {
	var sb strings.Builder

	sb.WriteString("Dependency Summary\n")
	sb.WriteString("==================\n\n")

	sb.WriteString(fmt.Sprintf("NPM Packages: %d\n", ds.TotalNPM))
	sb.WriteString(fmt.Sprintf("Go Packages: %d\n", ds.TotalGo))
	sb.WriteString(fmt.Sprintf("  - Standard Library: %d\n", len(ds.StdlibPackages)))
	sb.WriteString(fmt.Sprintf("  - External: %d\n", len(ds.ExternalPackages)))
	sb.WriteString(fmt.Sprintf("  - Runtime: %d\n", len(ds.RuntimePackages)))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("Supported: %d\n", ds.Supported))
	sb.WriteString(fmt.Sprintf("Partial: %d\n", ds.Partial))
	sb.WriteString(fmt.Sprintf("Unsupported: %d\n", ds.Unsupported))

	if len(ds.UnsupportedList) > 0 {
		sb.WriteString("\nUnsupported Packages:\n")
		for _, pkg := range ds.UnsupportedList {
			sb.WriteString(fmt.Sprintf("  - %s\n", pkg))
		}
	}

	return sb.String()
}
