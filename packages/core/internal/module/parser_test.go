package module

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParser_ParseFile_NamedExports(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
export interface User {
	id: string;
	name: string;
}

export function createUser(name: string): User {
	return { id: '123', name };
}

export const DEFAULT_ROLE = 'user';
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if module == nil {
		t.Fatal("Expected module, got nil")
	}

	// Check exports
	if len(module.Exports) != 3 {
		t.Errorf("Expected 3 exports, got %d", len(module.Exports))
	}

	// Verify export names
	exportNames := make(map[string]bool)
	for _, exp := range module.Exports {
		exportNames[exp.Name] = true
	}

	expected := []string{"User", "createUser", "DEFAULT_ROLE"}
	for _, name := range expected {
		if !exportNames[name] {
			t.Errorf("Expected export %s not found", name)
		}
	}
}

func TestParser_ParseFile_DefaultExport(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
interface Config {
	apiUrl: string;
}

const config: Config = {
	apiUrl: 'https://api.example.com'
};

export default config;
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check for default export
	hasDefault := false
	for _, exp := range module.Exports {
		if exp.Type == ExportDefault {
			hasDefault = true
			if exp.Name != "default" {
				t.Errorf("Expected default export name 'default', got %s", exp.Name)
			}
		}
	}

	if !hasDefault {
		t.Error("Expected default export not found")
	}
}

func TestParser_ParseFile_NamedImports(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
import { User, createUser } from './models/user';
import { fetchData } from './api';

function main() {
	const user: User = createUser('Alice');
	fetchData('/users');
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check imports
	if len(module.Imports) != 3 {
		t.Errorf("Expected 3 imports, got %d", len(module.Imports))
	}

	// Verify import names and sources
	importMap := make(map[string]string) // name -> source
	for _, imp := range module.Imports {
		importMap[imp.Name] = imp.Source
	}

	expected := map[string]string{
		"User":       "./models/user",
		"createUser": "./models/user",
		"fetchData":  "./api",
	}

	for name, expectedSource := range expected {
		if source, ok := importMap[name]; !ok {
			t.Errorf("Expected import %s not found", name)
		} else if source != expectedSource {
			t.Errorf("Import %s: expected source %s, got %s", name, expectedSource, source)
		}
	}
}

func TestParser_ParseFile_DefaultImport(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
import config from './config';
import axios from 'axios';

console.log(config.apiUrl);
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check for default imports
	defaultImports := 0
	for _, imp := range module.Imports {
		if imp.Type == ImportDefault {
			defaultImports++
		}
	}

	if defaultImports != 2 {
		t.Errorf("Expected 2 default imports, got %d", defaultImports)
	}
}

func TestParser_ParseFile_NamespaceImport(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
import * as utils from './utils';
import * as path from 'path';

utils.formatDate(new Date());
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check for namespace imports
	namespaceImports := 0
	for _, imp := range module.Imports {
		if imp.Type == ImportNamespace {
			namespaceImports++
		}
	}

	if namespaceImports != 2 {
		t.Errorf("Expected 2 namespace imports, got %d", namespaceImports)
	}
}

func TestParser_ParseFile_ReExports(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
export { User, createUser } from './models/user';
export * from './utils';
export * as api from './api';
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check exports
	if len(module.Exports) < 3 {
		t.Errorf("Expected at least 3 exports, got %d", len(module.Exports))
	}

	// Check for different export types
	hasNamed := false
	hasAll := false
	hasAllAs := false

	for _, exp := range module.Exports {
		switch exp.Type {
		case ExportNamed:
			hasNamed = true
		case ExportAll:
			hasAll = true
		case ExportAllAs:
			hasAllAs = true
		}
	}

	if !hasNamed {
		t.Error("Expected named re-export not found")
	}
	if !hasAll {
		t.Error("Expected export * not found")
	}
	if !hasAllAs {
		t.Error("Expected export * as not found")
	}
}

func TestParser_ParseFile_MixedExports(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	content := `
export interface User {
	id: string;
}

export function createUser(): User {
	return { id: '123' };
}

export default class UserService {
	getUser() { }
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()
	module, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check exports
	if len(module.Exports) != 3 {
		t.Errorf("Expected 3 exports, got %d", len(module.Exports))
	}

	// Check for default export
	hasDefault := false
	for _, exp := range module.Exports {
		if exp.Type == ExportDefault {
			hasDefault = true
		}
	}

	if !hasDefault {
		t.Error("Expected default export not found")
	}
}

func TestExportRegistry_AddAndGetModule(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/test/user.ts",
		Exports: []Export{
			{Name: "User", Type: ExportNamed},
		},
	}

	registry.AddModule(module)

	retrieved := registry.GetModule("/test/user.ts")
	if retrieved == nil {
		t.Fatal("Expected module, got nil")
	}

	if retrieved.Path != module.Path {
		t.Errorf("Expected path %s, got %s", module.Path, retrieved.Path)
	}
}

func TestExportRegistry_GetExports(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/test/user.ts",
		Exports: []Export{
			{Name: "User", Type: ExportNamed},
			{Name: "createUser", Type: ExportNamed},
		},
	}

	registry.AddModule(module)

	exports := registry.GetExports("/test/user.ts")
	if len(exports) != 2 {
		t.Errorf("Expected 2 exports, got %d", len(exports))
	}
}

func TestExportRegistry_AllModules(t *testing.T) {
	registry := NewExportRegistry()

	module1 := &Module{Path: "/test/user.ts"}
	module2 := &Module{Path: "/test/api.ts"}

	registry.AddModule(module1)
	registry.AddModule(module2)

	modules := registry.AllModules()
	if len(modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(modules))
	}
}
