package mapper

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// MappingType represents the type of Go package mapping
type MappingType string

const (
	MappingTypeRuntime     MappingType = "runtime"     // Custom runtime package
	MappingTypeStdlib      MappingType = "stdlib"      // Go standard library
	MappingTypeEquivalent  MappingType = "equivalent"  // Third-party equivalent
	MappingTypeFramework   MappingType = "framework"   // Web framework
	MappingTypeUnsupported MappingType = "unsupported" // Not supported
)

// MappingStatus represents support level
type MappingStatus string

const (
	StatusSupported   MappingStatus = "supported"   // Fully supported
	StatusPartial     MappingStatus = "partial"     // Partially supported
	StatusUnsupported MappingStatus = "unsupported" // Not supported
)

// ComplexityLevel represents conversion complexity
type ComplexityLevel string

const (
	ComplexitySimple ComplexityLevel = "simple" // Direct 1:1 mapping
	ComplexityMedium ComplexityLevel = "medium" // Some API differences
	ComplexityHigh   ComplexityLevel = "high"   // Significant differences
)

// Mapping represents a single npm-to-go package mapping
type Mapping struct {
	Npm         string            `yaml:"npm"`
	Go          string            `yaml:"go"`
	Type        MappingType       `yaml:"type"`
	Status      MappingStatus     `yaml:"status"`
	Complexity  ComplexityLevel   `yaml:"complexity"`
	Description string            `yaml:"description"`
	Notes       string            `yaml:"notes"`
	Example     string            `yaml:"example"`
	Reason      string            `yaml:"reason"`
	Suggestion  string            `yaml:"suggestion"`
	APIMappings map[string]string `yaml:"api_mappings"`
}

// MappingDatabase contains all package mappings
type MappingDatabase struct {
	Version  string              `yaml:"version"`
	Mappings []Mapping           `yaml:"mappings"`
	index    map[string]*Mapping // Internal index for fast lookups
}

// LoadMappings loads the npm-to-go mappings from a YAML file
func LoadMappings(filePath string) (*MappingDatabase, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mapping file: %w", err)
	}

	var db MappingDatabase
	if err := yaml.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("failed to parse mapping YAML: %w", err)
	}

	// Build index for fast lookups
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	return &db, nil
}

// LoadDefaultMappings loads mappings from the default location
func LoadDefaultMappings() (*MappingDatabase, error) {
	// Try to find mappings/npm-to-go.yaml relative to the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Search up the directory tree for mappings directory
	dir := cwd
	for {
		mappingPath := filepath.Join(dir, "mappings", "npm-to-go.yaml")
		if _, err := os.Stat(mappingPath); err == nil {
			return LoadMappings(mappingPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}

	return nil, fmt.Errorf("could not find mappings/npm-to-go.yaml in directory tree")
}

// GetMapping retrieves a mapping for an npm package
func (db *MappingDatabase) GetMapping(npmPackage string) (*Mapping, error) {
	if mapping, ok := db.index[npmPackage]; ok {
		return mapping, nil
	}
	return nil, fmt.Errorf("no mapping found for npm package: %s", npmPackage)
}

// HasMapping checks if a mapping exists for an npm package
func (db *MappingDatabase) HasMapping(npmPackage string) bool {
	_, ok := db.index[npmPackage]
	return ok
}

// GetSupportedMappings returns all supported mappings
func (db *MappingDatabase) GetSupportedMappings() []Mapping {
	var result []Mapping
	for _, m := range db.Mappings {
		if m.Status == StatusSupported || m.Status == StatusPartial {
			result = append(result, m)
		}
	}
	return result
}

// GetMappingsByType returns all mappings of a specific type
func (db *MappingDatabase) GetMappingsByType(mappingType MappingType) []Mapping {
	var result []Mapping
	for _, m := range db.Mappings {
		if m.Type == mappingType {
			result = append(result, m)
		}
	}
	return result
}

// GetRuntimeMappings returns all runtime package mappings
func (db *MappingDatabase) GetRuntimeMappings() []Mapping {
	return db.GetMappingsByType(MappingTypeRuntime)
}

// GetStdlibMappings returns all stdlib mappings
func (db *MappingDatabase) GetStdlibMappings() []Mapping {
	return db.GetMappingsByType(MappingTypeStdlib)
}

// GetAPIMapping retrieves a specific API mapping
func (m *Mapping) GetAPIMapping(npmAPI string) (string, bool) {
	if m.APIMappings == nil {
		return "", false
	}
	goAPI, ok := m.APIMappings[npmAPI]
	return goAPI, ok
}

// IsSupported checks if the mapping is fully or partially supported
func (m *Mapping) IsSupported() bool {
	return m.Status == StatusSupported || m.Status == StatusPartial
}

// IsBuiltin checks if this maps to a Node.js built-in module
func (m *Mapping) IsBuiltin() bool {
	// Built-in modules typically map to runtime or stdlib
	return m.Type == MappingTypeRuntime || m.Type == MappingTypeStdlib
}

// GetGoImportPath returns the Go import path
func (m *Mapping) GetGoImportPath() string {
	return m.Go
}

// Summary returns a human-readable summary of the mapping
func (db *MappingDatabase) Summary() string {
	total := len(db.Mappings)
	supported := 0
	partial := 0
	unsupported := 0

	byType := make(map[MappingType]int)

	for _, m := range db.Mappings {
		switch m.Status {
		case StatusSupported:
			supported++
		case StatusPartial:
			partial++
		case StatusUnsupported:
			unsupported++
		}
		byType[m.Type]++
	}

	return fmt.Sprintf(
		"Mapping Database v%s\n"+
			"Total Packages: %d\n"+
			"  Supported: %d (%.1f%%)\n"+
			"  Partial: %d (%.1f%%)\n"+
			"  Unsupported: %d (%.1f%%)\n"+
			"By Type:\n"+
			"  Runtime: %d\n"+
			"  Stdlib: %d\n"+
			"  Equivalent: %d\n"+
			"  Framework: %d\n"+
			"  Unsupported: %d\n",
		db.Version,
		total,
		supported, float64(supported)/float64(total)*100,
		partial, float64(partial)/float64(total)*100,
		unsupported, float64(unsupported)/float64(total)*100,
		byType[MappingTypeRuntime],
		byType[MappingTypeStdlib],
		byType[MappingTypeEquivalent],
		byType[MappingTypeFramework],
		byType[MappingTypeUnsupported],
	)
}
