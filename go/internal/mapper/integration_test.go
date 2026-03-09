package mapper

import (
	"path/filepath"
	"testing"
)

func TestLoadActualMappings(t *testing.T) {
	// Load the actual npm-to-go.yaml file
	mappingFile := filepath.Join("..", "..", "mappings", "npm-to-go.yaml")

	db, err := LoadMappings(mappingFile)
	if err != nil {
		t.Fatalf("Failed to load actual mappings: %v", err)
	}

	if db.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", db.Version)
	}

	// Should have many mappings (49 in current version)
	if len(db.Mappings) < 40 {
		t.Errorf("Expected at least 40 mappings, got %d", len(db.Mappings))
	}

	t.Logf("Loaded %d package mappings", len(db.Mappings))
	t.Logf("\n%s", db.Summary())
}

func TestActualNodeBuiltins(t *testing.T) {
	mappingFile := filepath.Join("..", "..", "mappings", "npm-to-go.yaml")
	db, err := LoadMappings(mappingFile)
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Test Node.js built-ins
	builtins := []struct {
		npm        string
		expectType MappingType
	}{
		{"fs", MappingTypeRuntime},
		{"path", MappingTypeRuntime},
		{"os", MappingTypeStdlib},
		{"http", MappingTypeStdlib},
		{"crypto", MappingTypeStdlib},
		{"url", MappingTypeStdlib},
	}

	for _, builtin := range builtins {
		t.Run(builtin.npm, func(t *testing.T) {
			mapping, err := db.GetMapping(builtin.npm)
			if err != nil {
				t.Errorf("Expected to find mapping for %s: %v", builtin.npm, err)
				return
			}

			if mapping.Type != builtin.expectType {
				t.Errorf("Expected %s to be type %s, got %s", builtin.npm, builtin.expectType, mapping.Type)
			}

			if !mapping.IsBuiltin() {
				t.Errorf("Expected %s to be builtin", builtin.npm)
			}

			t.Logf("%s -> %s (%s, %s)", mapping.Npm, mapping.Go, mapping.Type, mapping.Status)
		})
	}
}

func TestActualPopularPackages(t *testing.T) {
	mappingFile := filepath.Join("..", "..", "mappings", "npm-to-go.yaml")
	db, err := LoadMappings(mappingFile)
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Test popular packages
	packages := []struct {
		npm          string
		expectType   MappingType
		expectStatus MappingStatus
	}{
		{"axios", MappingTypeEquivalent, StatusSupported},
		{"express", MappingTypeFramework, StatusPartial},
		{"lodash", MappingTypeEquivalent, StatusSupported},
		{"uuid", MappingTypeEquivalent, StatusSupported},
		{"bcrypt", MappingTypeStdlib, StatusSupported},
		{"jsonwebtoken", MappingTypeEquivalent, StatusSupported},
	}

	for _, pkg := range packages {
		t.Run(pkg.npm, func(t *testing.T) {
			mapping, err := db.GetMapping(pkg.npm)
			if err != nil {
				t.Errorf("Expected to find mapping for %s: %v", pkg.npm, err)
				return
			}

			if mapping.Type != pkg.expectType {
				t.Errorf("Expected %s to be type %s, got %s", pkg.npm, pkg.expectType, mapping.Type)
			}

			if mapping.Status != pkg.expectStatus {
				t.Errorf("Expected %s to have status %s, got %s", pkg.npm, pkg.expectStatus, mapping.Status)
			}

			t.Logf("%s -> %s (%s, %s)", mapping.Npm, mapping.Go, mapping.Type, mapping.Status)
			if len(mapping.APIMappings) > 0 {
				t.Logf("  API mappings: %d", len(mapping.APIMappings))
			}
		})
	}
}

func TestActualUnsupportedPackages(t *testing.T) {
	mappingFile := filepath.Join("..", "..", "mappings", "npm-to-go.yaml")
	db, err := LoadMappings(mappingFile)
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Test unsupported packages
	unsupported := []string{"react", "vue", "angular", "webpack", "babel"}

	for _, pkg := range unsupported {
		t.Run(pkg, func(t *testing.T) {
			mapping, err := db.GetMapping(pkg)
			if err != nil {
				t.Errorf("Expected to find mapping for %s: %v", pkg, err)
				return
			}

			if mapping.Status != StatusUnsupported {
				t.Errorf("Expected %s to be unsupported, got %s", pkg, mapping.Status)
			}

			if mapping.Reason == "" {
				t.Errorf("Expected reason for unsupported package %s", pkg)
			}

			t.Logf("%s: %s (suggestion: %s)", pkg, mapping.Reason, mapping.Suggestion)
		})
	}
}

func TestActualAPIMappings(t *testing.T) {
	mappingFile := filepath.Join("..", "..", "mappings", "npm-to-go.yaml")
	db, err := LoadMappings(mappingFile)
	if err != nil {
		t.Fatalf("Failed to load mappings: %v", err)
	}

	// Test packages with API mappings
	tests := []struct {
		npm    string
		npmAPI string
	}{
		{"axios", "axios.get"},
		{"lodash", "_.map"},
		{"express", "app.get"},
		{"os", "os.platform"},
	}

	for _, test := range tests {
		t.Run(test.npm+"."+test.npmAPI, func(t *testing.T) {
			mapping, err := db.GetMapping(test.npm)
			if err != nil {
				t.Fatalf("Failed to get mapping for %s: %v", test.npm, err)
			}

			goAPI, ok := mapping.GetAPIMapping(test.npmAPI)
			if !ok {
				t.Errorf("Expected to find API mapping for %s", test.npmAPI)
				return
			}

			t.Logf("%s -> %s", test.npmAPI, goAPI)
		})
	}
}
