package mapper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMappings(t *testing.T) {
	// Create temporary mapping file
	tmpDir := t.TempDir()
	mappingFile := filepath.Join(tmpDir, "test-mappings.yaml")

	content := `version: "1.0"
mappings:
  - npm: "fs"
    go: "github.com/el-j/ts2go/runtime/fs"
    type: "runtime"
    status: "supported"
    complexity: "simple"
    description: "File system operations"
    
  - npm: "axios"
    go: "github.com/go-resty/resty/v2"
    type: "equivalent"
    status: "supported"
    complexity: "medium"
    description: "HTTP client"
    api_mappings:
      "axios.get": "resty.R().Get"
      "axios.post": "resty.R().Post"
      
  - npm: "react"
    type: "unsupported"
    status: "unsupported"
    reason: "Frontend framework"
    suggestion: "Use Go templates"
`

	if err := os.WriteFile(mappingFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test mapping file: %v", err)
	}

	db, err := LoadMappings(mappingFile)
	if err != nil {
		t.Fatalf("LoadMappings failed: %v", err)
	}

	if db.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", db.Version)
	}

	if len(db.Mappings) != 3 {
		t.Errorf("Expected 3 mappings, got %d", len(db.Mappings))
	}

	if len(db.index) != 3 {
		t.Errorf("Expected index size 3, got %d", len(db.index))
	}
}

func TestGetMapping(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{
				Npm:        "fs",
				Go:         "github.com/el-j/ts2go/runtime/fs",
				Type:       MappingTypeRuntime,
				Status:     StatusSupported,
				Complexity: ComplexitySimple,
			},
			{
				Npm:        "axios",
				Go:         "github.com/go-resty/resty/v2",
				Type:       MappingTypeEquivalent,
				Status:     StatusSupported,
				Complexity: ComplexityMedium,
			},
		},
	}

	// Build index
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	tests := []struct {
		name        string
		npmPackage  string
		expectError bool
		expectedGo  string
	}{
		{
			name:        "existing runtime package",
			npmPackage:  "fs",
			expectError: false,
			expectedGo:  "github.com/el-j/ts2go/runtime/fs",
		},
		{
			name:        "existing equivalent package",
			npmPackage:  "axios",
			expectError: false,
			expectedGo:  "github.com/go-resty/resty/v2",
		},
		{
			name:        "non-existing package",
			npmPackage:  "unknown",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping, err := db.GetMapping(tt.npmPackage)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if mapping.Go != tt.expectedGo {
					t.Errorf("Expected Go package %s, got %s", tt.expectedGo, mapping.Go)
				}
			}
		})
	}
}

func TestHasMapping(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Go: "runtime/fs"},
			{Npm: "axios", Go: "github.com/go-resty/resty/v2"},
		},
	}

	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	if !db.HasMapping("fs") {
		t.Error("Expected fs to have mapping")
	}

	if !db.HasMapping("axios") {
		t.Error("Expected axios to have mapping")
	}

	if db.HasMapping("unknown") {
		t.Error("Expected unknown to not have mapping")
	}
}

func TestGetSupportedMappings(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Status: StatusSupported},
			{Npm: "http", Status: StatusPartial},
			{Npm: "react", Status: StatusUnsupported},
		},
	}

	supported := db.GetSupportedMappings()

	if len(supported) != 2 {
		t.Errorf("Expected 2 supported mappings, got %d", len(supported))
	}

	for _, m := range supported {
		if m.Status == StatusUnsupported {
			t.Errorf("Unsupported mapping %s found in supported list", m.Npm)
		}
	}
}

func TestGetMappingsByType(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Type: MappingTypeRuntime},
			{Npm: "path", Type: MappingTypeRuntime},
			{Npm: "os", Type: MappingTypeStdlib},
			{Npm: "axios", Type: MappingTypeEquivalent},
			{Npm: "express", Type: MappingTypeFramework},
			{Npm: "react", Type: MappingTypeUnsupported},
		},
	}

	tests := []struct {
		name          string
		mappingType   MappingType
		expectedCount int
	}{
		{"runtime mappings", MappingTypeRuntime, 2},
		{"stdlib mappings", MappingTypeStdlib, 1},
		{"equivalent mappings", MappingTypeEquivalent, 1},
		{"framework mappings", MappingTypeFramework, 1},
		{"unsupported mappings", MappingTypeUnsupported, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mappings := db.GetMappingsByType(tt.mappingType)
			if len(mappings) != tt.expectedCount {
				t.Errorf("Expected %d mappings, got %d", tt.expectedCount, len(mappings))
			}
		})
	}
}

func TestGetRuntimeMappings(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Type: MappingTypeRuntime},
			{Npm: "path", Type: MappingTypeRuntime},
			{Npm: "os", Type: MappingTypeStdlib},
		},
	}

	runtime := db.GetRuntimeMappings()

	if len(runtime) != 2 {
		t.Errorf("Expected 2 runtime mappings, got %d", len(runtime))
	}

	for _, m := range runtime {
		if m.Type != MappingTypeRuntime {
			t.Errorf("Non-runtime mapping %s found in runtime list", m.Npm)
		}
	}
}

func TestGetStdlibMappings(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Type: MappingTypeRuntime},
			{Npm: "os", Type: MappingTypeStdlib},
			{Npm: "crypto", Type: MappingTypeStdlib},
		},
	}

	stdlib := db.GetStdlibMappings()

	if len(stdlib) != 2 {
		t.Errorf("Expected 2 stdlib mappings, got %d", len(stdlib))
	}

	for _, m := range stdlib {
		if m.Type != MappingTypeStdlib {
			t.Errorf("Non-stdlib mapping %s found in stdlib list", m.Npm)
		}
	}
}

func TestGetAPIMapping(t *testing.T) {
	mapping := Mapping{
		Npm: "axios",
		APIMappings: map[string]string{
			"axios.get":  "resty.R().Get",
			"axios.post": "resty.R().Post",
		},
	}

	goAPI, ok := mapping.GetAPIMapping("axios.get")
	if !ok {
		t.Error("Expected to find axios.get mapping")
	}
	if goAPI != "resty.R().Get" {
		t.Errorf("Expected resty.R().Get, got %s", goAPI)
	}

	_, ok = mapping.GetAPIMapping("axios.unknown")
	if ok {
		t.Error("Expected not to find axios.unknown mapping")
	}

	// Test mapping with no API mappings
	emptyMapping := Mapping{Npm: "fs"}
	_, ok = emptyMapping.GetAPIMapping("anything")
	if ok {
		t.Error("Expected not to find mapping when APIMappings is nil")
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name     string
		status   MappingStatus
		expected bool
	}{
		{"supported", StatusSupported, true},
		{"partial", StatusPartial, true},
		{"unsupported", StatusUnsupported, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := Mapping{Status: tt.status}
			if mapping.IsSupported() != tt.expected {
				t.Errorf("Expected IsSupported=%v for status %s", tt.expected, tt.status)
			}
		})
	}
}

func TestIsBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		mapType  MappingType
		expected bool
	}{
		{"runtime is builtin", MappingTypeRuntime, true},
		{"stdlib is builtin", MappingTypeStdlib, true},
		{"equivalent is not builtin", MappingTypeEquivalent, false},
		{"framework is not builtin", MappingTypeFramework, false},
		{"unsupported is not builtin", MappingTypeUnsupported, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := Mapping{Type: tt.mapType}
			if mapping.IsBuiltin() != tt.expected {
				t.Errorf("Expected IsBuiltin=%v for type %s", tt.expected, tt.mapType)
			}
		})
	}
}

func TestGetGoImportPath(t *testing.T) {
	mapping := Mapping{
		Npm: "axios",
		Go:  "github.com/go-resty/resty/v2",
	}

	if path := mapping.GetGoImportPath(); path != "github.com/go-resty/resty/v2" {
		t.Errorf("Expected github.com/go-resty/resty/v2, got %s", path)
	}
}

func TestSummary(t *testing.T) {
	db := &MappingDatabase{
		Version: "1.0",
		Mappings: []Mapping{
			{Npm: "fs", Type: MappingTypeRuntime, Status: StatusSupported},
			{Npm: "path", Type: MappingTypeRuntime, Status: StatusSupported},
			{Npm: "os", Type: MappingTypeStdlib, Status: StatusSupported},
			{Npm: "http", Type: MappingTypeStdlib, Status: StatusPartial},
			{Npm: "axios", Type: MappingTypeEquivalent, Status: StatusSupported},
			{Npm: "express", Type: MappingTypeFramework, Status: StatusPartial},
			{Npm: "react", Type: MappingTypeUnsupported, Status: StatusUnsupported},
		},
	}

	summary := db.Summary()

	// Check that summary contains key information
	expectedStrings := []string{
		"v1.0",
		"Total Packages: 7",
		"Supported: 4",
		"Partial: 2",
		"Unsupported: 1",
		"Runtime: 2",
		"Stdlib: 2",
		"Equivalent: 1",
		"Framework: 1",
	}

	for _, expected := range expectedStrings {
		if !contains(summary, expected) {
			t.Errorf("Summary missing expected string: %s\nGot: %s", expected, summary)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
