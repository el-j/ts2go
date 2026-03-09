package mapper

import (
	"fmt"
	"testing"

	"github.com/el-j/ts2go/internal/analyzer"
)

func TestRewriteImport(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	tests := []struct {
		name           string
		imp            *analyzer.Import
		expectGoImport string
		expectAlias    string
		expectError    bool
	}{
		{
			name: "runtime package",
			imp: &analyzer.Import{
				Source:  "fs",
				Type:    analyzer.ImportTypeBuiltin,
				Symbols: []string{"default"},
			},
			expectGoImport: "github.com/el-j/ts2go/runtime/fs",
			expectAlias:    "fs",
			expectError:    false,
		},
		{
			name: "equivalent package",
			imp: &analyzer.Import{
				Source:  "axios",
				Type:    analyzer.ImportTypePackage,
				Symbols: []string{"default"},
			},
			expectGoImport: "github.com/go-resty/resty/v2",
			expectAlias:    "resty",
			expectError:    false,
		},
		{
			name: "local file import",
			imp: &analyzer.Import{
				Source:  "./models/user",
				Type:    analyzer.ImportTypeLocal,
				Symbols: []string{"User"},
			},
			expectGoImport: "myproject/models",
			expectError:    false,
		},
		{
			name: "unsupported package",
			imp: &analyzer.Import{
				Source:  "react",
				Type:    analyzer.ImportTypePackage,
				Symbols: []string{"default"},
			},
			expectError: true,
		},
		{
			name: "unknown package",
			imp: &analyzer.Import{
				Source:  "unknown-pkg",
				Type:    analyzer.ImportTypePackage,
				Symbols: []string{"default"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rewriter.RewriteImport(tt.imp)

			if tt.expectError {
				if result.Error == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if result.Error != nil {
				t.Errorf("Unexpected error: %v", result.Error)
				return
			}

			if result.GoImport != tt.expectGoImport {
				t.Errorf("Expected Go import %s, got %s", tt.expectGoImport, result.GoImport)
			}

			if tt.expectAlias != "" && result.PackageAlias != tt.expectAlias {
				t.Errorf("Expected alias %s, got %s", tt.expectAlias, result.PackageAlias)
			}

			t.Logf("Rewrite: %s -> %s (alias: %s)", result.OriginalSource, result.GoImport, result.PackageAlias)
		})
	}
}

func TestRewriteLocalImport(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	tests := []struct {
		name           string
		source         string
		expectGoImport string
		expectLocal    bool
	}{
		{
			name:           "current directory",
			source:         "./models/user",
			expectGoImport: "myproject/models",
			expectLocal:    true,
		},
		{
			name:           "parent directory",
			source:         "../utils/helper",
			expectGoImport: "myproject/utils",
			expectLocal:    true,
		},
		{
			name:           "nested directory",
			source:         "./api/handlers/user",
			expectGoImport: "myproject/api/handlers",
			expectLocal:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imp := &analyzer.Import{
				Source:  tt.source,
				Type:    analyzer.ImportTypeLocal,
				Symbols: []string{"default"},
			}

			result := rewriter.RewriteImport(imp)

			if !result.IsLocal {
				t.Error("Expected IsLocal to be true")
			}

			if !result.NeedsTranspile {
				t.Error("Expected NeedsTranspile to be true for local imports")
			}

			if result.GoImport != tt.expectGoImport {
				t.Errorf("Expected Go import %s, got %s", tt.expectGoImport, result.GoImport)
			}

			t.Logf("Local import: %s -> %s", tt.source, result.GoImport)
		})
	}
}

func TestMapSymbols(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	tests := []struct {
		name            string
		imp             *analyzer.Import
		mappingSource   string
		expectedSymbols map[string]string
	}{
		{
			name: "default import",
			imp: &analyzer.Import{
				Source:    "axios",
				Symbols:   []string{"default"},
				IsDefault: true,
			},
			mappingSource: "axios",
			expectedSymbols: map[string]string{
				"default": "resty",
			},
		},
		{
			name: "namespace import",
			imp: &analyzer.Import{
				Source:      "fs",
				Symbols:     []string{"*"},
				IsNamespace: true,
			},
			mappingSource: "fs",
			expectedSymbols: map[string]string{
				"*": "fs",
			},
		},
		{
			name: "named imports",
			imp: &analyzer.Import{
				Source:  "path",
				Symbols: []string{"join", "dirname"},
			},
			mappingSource: "path",
			expectedSymbols: map[string]string{
				"join":    "Join",
				"dirname": "Dirname",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping, _ := db.GetMapping(tt.mappingSource)
			symbolMappings := rewriter.mapSymbols(tt.imp, mapping)

			for tsSymbol, expectedGoSymbol := range tt.expectedSymbols {
				goSymbol, ok := symbolMappings[tsSymbol]
				if !ok {
					t.Errorf("Expected symbol mapping for %s", tsSymbol)
					continue
				}

				if goSymbol != expectedGoSymbol {
					t.Errorf("Symbol %s: expected %s, got %s", tsSymbol, expectedGoSymbol, goSymbol)
				}

				t.Logf("Symbol mapping: %s -> %s", tsSymbol, goSymbol)
			}
		})
	}
}

func TestGenerateGoImports(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	results := []*RewriteResult{
		{
			OriginalSource: "fs",
			GoImport:       "github.com/el-j/ts2go/runtime/fs",
		},
		{
			OriginalSource: "axios",
			GoImport:       "github.com/go-resty/resty/v2",
		},
		{
			OriginalSource: "./models/user",
			GoImport:       "myproject/models",
		},
	}

	goImports := rewriter.GenerateGoImports(results)

	if goImports == "" {
		t.Fatal("Expected non-empty import block")
	}

	// Should contain import statement
	if !contains(goImports, "import (") {
		t.Error("Expected import block to start with 'import ('")
	}

	// Should contain all packages
	expectedPkgs := []string{
		"github.com/el-j/ts2go/runtime/fs",
		"github.com/go-resty/resty/v2",
		"myproject/models",
	}

	for _, pkg := range expectedPkgs {
		if !contains(goImports, pkg) {
			t.Errorf("Expected import block to contain %s", pkg)
		}
	}

	t.Logf("Generated imports:\n%s", goImports)
}

func TestGenerateGoImportsWithErrors(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	results := []*RewriteResult{
		{
			OriginalSource: "fs",
			GoImport:       "github.com/el-j/ts2go/runtime/fs",
		},
		{
			OriginalSource: "react",
			Error:          fmt.Errorf("package not supported: react"),
		},
	}

	goImports := rewriter.GenerateGoImports(results)

	if !contains(goImports, "fs") {
		t.Error("Expected successful import to be included")
	}

	if !contains(goImports, "Error") {
		t.Error("Expected error comment to be included")
	}

	t.Logf("Generated imports with errors:\n%s", goImports)
}

func TestGetRequiredPackages(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	results := []*RewriteResult{
		{OriginalSource: "fs", GoImport: "runtime/fs"},
		{OriginalSource: "axios", GoImport: "github.com/go-resty/resty/v2"},
		{OriginalSource: "fs", GoImport: "runtime/fs"}, // Duplicate
		{OriginalSource: "react", Error: fmt.Errorf("unsupported")},
	}

	packages := rewriter.GetRequiredPackages(results)

	if len(packages) != 2 {
		t.Errorf("Expected 2 unique packages, got %d", len(packages))
	}

	// Check that duplicates are removed
	seen := make(map[string]bool)
	for _, pkg := range packages {
		if seen[pkg] {
			t.Errorf("Duplicate package found: %s", pkg)
		}
		seen[pkg] = true
	}

	t.Logf("Required packages: %v", packages)
}

func TestRewriteImports(t *testing.T) {
	db := createTestMappingDB()
	rewriter := NewImportRewriter(db, "myproject")

	imports := []*analyzer.Import{
		{Source: "fs", Type: analyzer.ImportTypeBuiltin, Symbols: []string{"default"}},
		{Source: "axios", Type: analyzer.ImportTypePackage, Symbols: []string{"default"}},
		{Source: "./models/user", Type: analyzer.ImportTypeLocal, Symbols: []string{"User"}},
	}

	results := rewriter.RewriteImports(imports)

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	successCount := 0
	for _, result := range results {
		if result.Error == nil {
			successCount++
		}
	}

	if successCount != 3 {
		t.Errorf("Expected 3 successful rewrites, got %d", successCount)
	}
}

func TestExtractFunctionName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"resty.R().Get", "Get"},
		{"lo.Map", "Map"},
		{"path.Join", "Join"},
		{"SimpleFunction", "SimpleFunction"},
		{"", ""},
	}

	for _, tt := range tests {
		result := extractFunctionName(tt.input)
		if result != tt.expected {
			t.Errorf("extractFunctionName(%s): expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"a", "A"},
		{"", ""},
		{"Hello", "Hello"},
	}

	for _, tt := range tests {
		result := capitalizeFirst(tt.input)
		if result != tt.expected {
			t.Errorf("capitalizeFirst(%s): expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

// Helper function to create test mapping database
func createTestMappingDB() *MappingDatabase {
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
				APIMappings: map[string]string{
					"axios.get": "resty.R().Get",
				},
			},
			{
				Npm:        "path",
				Go:         "github.com/el-j/ts2go/runtime/path",
				Type:       MappingTypeRuntime,
				Status:     StatusSupported,
				Complexity: ComplexitySimple,
			},
			{
				Npm:        "react",
				Type:       MappingTypeUnsupported,
				Status:     StatusUnsupported,
				Reason:     "Frontend framework",
				Suggestion: "Use Go templates",
			},
		},
	}

	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	return db
}
