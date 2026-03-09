package mapper

import (
	"fmt"
	"strings"
)

// PackageClass represents the category of an npm package
type PackageClass int

const (
	PackageClassBuiltin      PackageClass = iota // Node.js built-in module
	PackageClassRuntime                          // Custom runtime wrapper
	PackageClassEquivalent                       // Direct Go equivalent exists
	PackageClassFramework                        // Web framework or similar
	PackageClassTranspilable                     // Pure TS, can transpile locally
	PackageClassUnsupported                      // Can't transpile (native modules, browser-only)
)

// String returns the string representation of PackageClass
func (pc PackageClass) String() string {
	switch pc {
	case PackageClassBuiltin:
		return "builtin"
	case PackageClassRuntime:
		return "runtime"
	case PackageClassEquivalent:
		return "equivalent"
	case PackageClassFramework:
		return "framework"
	case PackageClassTranspilable:
		return "transpilable"
	case PackageClassUnsupported:
		return "unsupported"
	default:
		return "unknown"
	}
}

// Classification represents the result of classifying an npm package
type Classification struct {
	Package     string       // npm package name
	Class       PackageClass // Classification category
	GoPackage   string       // Mapped Go package (if available)
	Confidence  float64      // Confidence score 0.0 - 1.0
	Status      string       // "supported", "partial", "unsupported"
	Complexity  string       // "simple", "medium", "high"
	Notes       string       // Additional notes about the classification
	Suggestion  string       // Suggestion for unsupported packages
	IsSupported bool         // Quick check if package is supported
}

// Classifier provides dependency classification functionality
type Classifier struct {
	db *MappingDatabase
}

// NewClassifier creates a new classifier with the given mapping database
func NewClassifier(db *MappingDatabase) *Classifier {
	return &Classifier{db: db}
}

// ClassifyPackage classifies an npm package and returns detailed information
func (c *Classifier) ClassifyPackage(npmPackage string) (*Classification, error) {
	// Try to find in mapping database
	mapping, err := c.db.GetMapping(npmPackage)
	if err != nil {
		// Package not in database - check if it's a local file
		if isLocalImport(npmPackage) {
			return &Classification{
				Package:     npmPackage,
				Class:       PackageClassTranspilable,
				GoPackage:   "", // Will be determined by project structure
				Confidence:  1.0,
				Status:      "transpilable",
				Complexity:  "simple",
				Notes:       "Local file import, will be transpiled",
				IsSupported: true,
			}, nil
		}

		// Unknown package - make educated guess
		return c.classifyUnknownPackage(npmPackage)
	}

	// Package found in database
	class := c.mappingTypeToClass(mapping.Type)
	confidence := c.calculateConfidence(mapping)

	return &Classification{
		Package:     npmPackage,
		Class:       class,
		GoPackage:   mapping.Go,
		Confidence:  confidence,
		Status:      string(mapping.Status),
		Complexity:  string(mapping.Complexity),
		Notes:       mapping.Notes,
		Suggestion:  mapping.Suggestion,
		IsSupported: mapping.IsSupported(),
	}, nil
}

// ClassifyPackages classifies multiple packages at once
func (c *Classifier) ClassifyPackages(packages []string) ([]*Classification, error) {
	results := make([]*Classification, 0, len(packages))

	for _, pkg := range packages {
		classification, err := c.ClassifyPackage(pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to classify %s: %w", pkg, err)
		}
		results = append(results, classification)
	}

	return results, nil
}

// GetSupportedPackages returns only packages that are supported or partially supported
func (c *Classifier) GetSupportedPackages(packages []string) ([]*Classification, error) {
	all, err := c.ClassifyPackages(packages)
	if err != nil {
		return nil, err
	}

	supported := make([]*Classification, 0)
	for _, classification := range all {
		if classification.IsSupported {
			supported = append(supported, classification)
		}
	}

	return supported, nil
}

// GetUnsupportedPackages returns only packages that are not supported
func (c *Classifier) GetUnsupportedPackages(packages []string) ([]*Classification, error) {
	all, err := c.ClassifyPackages(packages)
	if err != nil {
		return nil, err
	}

	unsupported := make([]*Classification, 0)
	for _, classification := range all {
		if !classification.IsSupported {
			unsupported = append(unsupported, classification)
		}
	}

	return unsupported, nil
}

// AnalyzeDependencies provides a summary of all dependencies in a project
func (c *Classifier) AnalyzeDependencies(dependencies map[string]string) (*DependencyAnalysis, error) {
	packages := make([]string, 0, len(dependencies))
	for pkg := range dependencies {
		packages = append(packages, pkg)
	}

	classifications, err := c.ClassifyPackages(packages)
	if err != nil {
		return nil, err
	}

	analysis := &DependencyAnalysis{
		Total:           len(classifications),
		Classifications: classifications,
		ByClass:         make(map[PackageClass]int),
		Supported:       0,
		Partial:         0,
		Unsupported:     0,
	}

	for _, classification := range classifications {
		analysis.ByClass[classification.Class]++

		switch classification.Status {
		case "supported":
			analysis.Supported++
		case "partial":
			analysis.Partial++
		case "unsupported":
			analysis.Unsupported++
		}
	}

	return analysis, nil
}

// DependencyAnalysis contains summary information about project dependencies
type DependencyAnalysis struct {
	Total           int
	Classifications []*Classification
	ByClass         map[PackageClass]int
	Supported       int
	Partial         int
	Unsupported     int
}

// Summary returns a human-readable summary of the analysis
func (da *DependencyAnalysis) Summary() string {
	supportedPct := float64(da.Supported) / float64(da.Total) * 100
	partialPct := float64(da.Partial) / float64(da.Total) * 100
	unsupportedPct := float64(da.Unsupported) / float64(da.Total) * 100

	return fmt.Sprintf(
		"Dependency Analysis\n"+
			"Total Dependencies: %d\n"+
			"  Supported: %d (%.1f%%)\n"+
			"  Partial: %d (%.1f%%)\n"+
			"  Unsupported: %d (%.1f%%)\n"+
			"By Class:\n"+
			"  Builtin: %d\n"+
			"  Runtime: %d\n"+
			"  Equivalent: %d\n"+
			"  Framework: %d\n"+
			"  Transpilable: %d\n"+
			"  Unsupported: %d\n",
		da.Total,
		da.Supported, supportedPct,
		da.Partial, partialPct,
		da.Unsupported, unsupportedPct,
		da.ByClass[PackageClassBuiltin],
		da.ByClass[PackageClassRuntime],
		da.ByClass[PackageClassEquivalent],
		da.ByClass[PackageClassFramework],
		da.ByClass[PackageClassTranspilable],
		da.ByClass[PackageClassUnsupported],
	)
}

// GetUnsupportedList returns a list of unsupported package names
func (da *DependencyAnalysis) GetUnsupportedList() []string {
	unsupported := make([]string, 0)
	for _, classification := range da.Classifications {
		if !classification.IsSupported {
			unsupported = append(unsupported, classification.Package)
		}
	}
	return unsupported
}

// Helper functions

func (c *Classifier) mappingTypeToClass(mappingType MappingType) PackageClass {
	switch mappingType {
	case MappingTypeRuntime:
		return PackageClassRuntime
	case MappingTypeStdlib:
		return PackageClassBuiltin
	case MappingTypeEquivalent:
		return PackageClassEquivalent
	case MappingTypeFramework:
		return PackageClassFramework
	case MappingTypeUnsupported:
		return PackageClassUnsupported
	default:
		return PackageClassUnsupported
	}
}

func (c *Classifier) calculateConfidence(mapping *Mapping) float64 {
	// Base confidence on status and complexity
	var confidence float64

	switch mapping.Status {
	case StatusSupported:
		confidence = 1.0
	case StatusPartial:
		confidence = 0.7
	case StatusUnsupported:
		confidence = 0.0
	default:
		confidence = 0.5
	}

	// Adjust based on complexity
	switch mapping.Complexity {
	case ComplexitySimple:
		// No adjustment for simple
	case ComplexityMedium:
		confidence *= 0.9
	case ComplexityHigh:
		confidence *= 0.8
	}

	return confidence
}

func (c *Classifier) classifyUnknownPackage(npmPackage string) (*Classification, error) {
	// Try to make educated guess based on package name patterns

	// Scoped packages (e.g., @types/node, @babel/core)
	if strings.HasPrefix(npmPackage, "@types/") {
		return &Classification{
			Package:     npmPackage,
			Class:       PackageClassTranspilable,
			GoPackage:   "",
			Confidence:  0.5,
			Status:      "unknown",
			Complexity:  "simple",
			Notes:       "TypeScript type definitions, can be ignored",
			IsSupported: true,
		}, nil
	}

	// Common frontend frameworks/tools
	frontendKeywords := []string{"react", "vue", "angular", "svelte", "webpack", "babel", "rollup", "vite", "esbuild"}
	for _, keyword := range frontendKeywords {
		if strings.Contains(npmPackage, keyword) {
			return &Classification{
				Package:     npmPackage,
				Class:       PackageClassUnsupported,
				GoPackage:   "",
				Confidence:  0.8,
				Status:      "unsupported",
				Complexity:  "high",
				Notes:       "Frontend framework or build tool",
				Suggestion:  "Use Go templates or server-side rendering",
				IsSupported: false,
			}, nil
		}
	}

	// Unknown package - default to transpilable with low confidence
	return &Classification{
		Package:     npmPackage,
		Class:       PackageClassTranspilable,
		GoPackage:   "",
		Confidence:  0.3,
		Status:      "unknown",
		Complexity:  "unknown",
		Notes:       "Package not in mapping database, may need manual review",
		Suggestion:  "Add mapping to npm-to-go.yaml or transpile as local code",
		IsSupported: true, // Assume transpilable until proven otherwise
	}, nil
}

func isLocalImport(path string) bool {
	return strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, "/")
}
