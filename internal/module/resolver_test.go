package module

import (
	"strings"
	"testing"
)

func TestImportResolver_ResolveImports(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module with imports
	module := &Module{
		Path: "/project/src/services/api.ts",
		Imports: []Import{
			{
				Name:         "User",
				ImportedName: "User",
				Type:         ImportNamed,
				Source:       "../models/user",
			},
			{
				Name:         "createUser",
				ImportedName: "createUser",
				Type:         ImportNamed,
				Source:       "../models/user",
			},
			{
				Name:         "config",
				ImportedName: "default",
				Type:         ImportDefault,
				Source:       "../config",
			},
		},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/services/api.ts")

	resolved, err := resolver.ResolveImports()
	if err != nil {
		t.Fatalf("ResolveImports failed: %v", err)
	}

	// Check we have 2 imports (models and config)
	if len(resolved) != 2 {
		t.Errorf("Expected 2 resolved imports, got %d", len(resolved))
		for i, r := range resolved {
			t.Logf("  [%d] %s", i, r.PackagePath)
		}
	}

	// Find models import
	var modelsImport *ResolvedImport
	for i := range resolved {
		if strings.Contains(resolved[i].PackagePath, "models") {
			modelsImport = &resolved[i]
			break
		}
	}

	if modelsImport == nil {
		t.Fatal("Expected models import not found")
	}

	if modelsImport.PackagePath != "github.com/user/project/models" {
		t.Errorf("Expected package path github.com/user/project/models, got %s", modelsImport.PackagePath)
	}

	// Check symbols
	if len(modelsImport.Symbols) != 2 {
		t.Errorf("Expected 2 symbols from models, got %d", len(modelsImport.Symbols))
	}
}

func TestImportResolver_ResolveImports_SamePackage(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module with import from same package
	module := &Module{
		Path: "/project/src/models/user.ts",
		Imports: []Import{
			{
				Name:         "Post",
				ImportedName: "Post",
				Type:         ImportNamed,
				Source:       "./post", // Same package
			},
		},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/models/user.ts")

	resolved, err := resolver.ResolveImports()
	if err != nil {
		t.Fatalf("ResolveImports failed: %v", err)
	}

	// Should have 0 imports (same package)
	if len(resolved) != 0 {
		t.Errorf("Expected 0 resolved imports (same package), got %d", len(resolved))
	}
}

func TestImportResolver_ResolveImports_NpmPackages(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module with npm imports
	module := &Module{
		Path: "/project/src/index.ts",
		Imports: []Import{
			{
				Name:         "readFile",
				ImportedName: "readFile",
				Type:         ImportNamed,
				Source:       "fs",
			},
			{
				Name:         "axios",
				ImportedName: "default",
				Type:         ImportDefault,
				Source:       "axios",
			},
		},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/index.ts")

	resolved, err := resolver.ResolveImports()
	if err != nil {
		t.Fatalf("ResolveImports failed: %v", err)
	}

	// Should have 2 imports (fs and axios)
	if len(resolved) != 2 {
		t.Errorf("Expected 2 resolved imports, got %d", len(resolved))
	}

	// Check for runtime imports
	hasRuntime := false
	for _, imp := range resolved {
		if imp.IsRuntime {
			hasRuntime = true
			break
		}
	}

	if !hasRuntime {
		t.Error("Expected at least one runtime import")
	}
}

func TestImportResolver_GenerateImportBlock(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module with mixed imports
	module := &Module{
		Path: "/project/src/services/api.ts",
		Imports: []Import{
			{
				Name:         "User",
				ImportedName: "User",
				Type:         ImportNamed,
				Source:       "../models/user",
			},
			{
				Name:         "readFile",
				ImportedName: "readFile",
				Type:         ImportNamed,
				Source:       "fs",
			},
		},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/services/api.ts")

	importBlock, err := resolver.GenerateImportBlock()
	if err != nil {
		t.Fatalf("GenerateImportBlock failed: %v", err)
	}

	if importBlock == "" {
		t.Fatal("Expected import block, got empty string")
	}

	// Check that it starts with "import ("
	if !strings.HasPrefix(importBlock, "import (") {
		t.Error("Import block should start with 'import ('")
	}

	// Check that it ends with ")"
	if !strings.HasSuffix(importBlock, ")") {
		t.Error("Import block should end with ')'")
	}

	// Check that it contains both imports
	if !strings.Contains(importBlock, "github.com/el-j/ts2go-runtime/fs") {
		t.Error("Import block should contain fs runtime import")
	}

	if !strings.Contains(importBlock, "github.com/user/project/models") {
		t.Error("Import block should contain models import")
	}
}

func TestImportResolver_GenerateImportBlock_Empty(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module with no imports
	module := &Module{
		Path:    "/project/src/models/user.ts",
		Imports: []Import{},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/models/user.ts")

	importBlock, err := resolver.GenerateImportBlock()
	if err != nil {
		t.Fatalf("GenerateImportBlock failed: %v", err)
	}

	if importBlock != "" {
		t.Errorf("Expected empty import block, got: %s", importBlock)
	}
}

func TestImportResolver_GetImportedSymbols(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module
	module := &Module{
		Path: "/project/src/services/api.ts",
		Imports: []Import{
			{
				Name:         "User",
				ImportedName: "User",
				Type:         ImportNamed,
				Source:       "../models/user",
			},
			{
				Name:         "createUser",
				ImportedName: "createUser",
				Type:         ImportNamed,
				Source:       "../models/user",
			},
		},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/services/api.ts")

	symbols := resolver.GetImportedSymbols("github.com/user/project/models")

	if len(symbols) != 2 {
		t.Errorf("Expected 2 symbols, got %d", len(symbols))
	}

	// Check symbol names
	symbolMap := make(map[string]bool)
	for _, s := range symbols {
		symbolMap[s] = true
	}

	if !symbolMap["User"] || !symbolMap["createUser"] {
		t.Error("Expected User and createUser symbols")
	}
}

func TestImportResolver_ShouldImportPackage(t *testing.T) {
	// Setup
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")

	// Create a module
	module := &Module{
		Path: "/project/src/services/api.ts",
		Imports: []Import{
			{
				Name:         "User",
				ImportedName: "User",
				Type:         ImportNamed,
				Source:       "../models/user",
			},
		},
	}

	registry.AddModule(module)

	resolver := NewImportResolver(registry, mapper)
	resolver.SetCurrentFile("/project/src/services/api.ts")

	// Should import models package
	if !resolver.ShouldImportPackage("github.com/user/project/models") {
		t.Error("Should import models package")
	}

	// Should not import same package
	if resolver.ShouldImportPackage("github.com/user/project/services") {
		t.Error("Should not import same package")
	}

	// Should not import unused package
	if resolver.ShouldImportPackage("github.com/user/project/utils") {
		t.Error("Should not import unused package")
	}
}

func TestIsStdLibPackage(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		expected bool
	}{
		{"fmt package", "fmt", true},
		{"os package", "os", true},
		{"path package", "path", true},
		{"path/filepath package", "path/filepath", true},
		{"net/http package", "net/http", true},
		{"encoding/json package", "encoding/json", true},
		{"third party package", "github.com/user/project", false},
		{"runtime package", "github.com/el-j/ts2go-runtime/fs", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isStdLibPackage(tt.pkg)
			if result != tt.expected {
				t.Errorf("isStdLibPackage(%s) = %v, want %v", tt.pkg, result, tt.expected)
			}
		})
	}
}

func TestResolveNpmPackage(t *testing.T) {
	registry := NewExportRegistry()
	mapper := NewPackageMapper("/project", "github.com/user/project", "src")
	resolver := NewImportResolver(registry, mapper)

	tests := []struct {
		name        string
		npmPackage  string
		expectedNil bool
	}{
		{"fs package", "fs", false},
		{"path package", "path", false},
		{"axios package", "axios", false},
		{"lodash package", "lodash", false},
		{"unknown package", "some-random-package", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.resolveNpmPackage(tt.npmPackage)
			if tt.expectedNil && result != "" {
				t.Errorf("Expected empty result for %s, got %s", tt.npmPackage, result)
			}
			if !tt.expectedNil && result == "" {
				t.Errorf("Expected non-empty result for %s", tt.npmPackage)
			}
		})
	}
}
