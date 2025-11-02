package analyzer

import (
	"testing"

	"github.com/yourusername/ts2go/internal/transpiler"
)

func TestClassifyImportSource(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected ImportType
	}{
		{"npm package", "axios", ImportTypePackage},
		{"scoped package", "@babel/core", ImportTypePackage},
		{"builtin fs", "fs", ImportTypeBuiltin},
		{"builtin path", "path", ImportTypeBuiltin},
		{"builtin with node prefix", "node:fs", ImportTypeBuiltin},
		{"local relative", "./utils/helper", ImportTypeLocal},
		{"local parent", "../models/user", ImportTypeLocal},
		{"local absolute", "/absolute/path", ImportTypeLocal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyImportSource(tt.source)
			if result != tt.expected {
				t.Errorf("classifyImportSource('%s') = %v, expected %v", tt.source, result, tt.expected)
			}
		})
	}
}

func TestAnalyzeImportsDefault(t *testing.T) {
	// Create AST for: import axios from 'axios'
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "Identifier",
								Text: "axios",
							},
						},
					},
					{
						Kind: "StringLiteral",
						Text: "axios",
					},
				},
			},
		},
	}

	analysis, err := AnalyzeImports(ast)
	if err != nil {
		t.Fatalf("AnalyzeImports failed: %v", err)
	}

	if len(analysis.Imports) != 1 {
		t.Fatalf("Expected 1 import, got %d", len(analysis.Imports))
	}

	imp := analysis.Imports[0]
	if imp.Source != "axios" {
		t.Errorf("Expected source 'axios', got '%s'", imp.Source)
	}

	if imp.Type != ImportTypePackage {
		t.Errorf("Expected ImportTypePackage, got %v", imp.Type)
	}

	if !imp.IsDefault {
		t.Error("Expected IsDefault to be true")
	}

	if len(imp.Symbols) != 1 || imp.Symbols[0] != "axios" {
		t.Errorf("Expected symbols ['axios'], got %v", imp.Symbols)
	}

	if len(analysis.NpmPackages) != 1 || analysis.NpmPackages[0] != "axios" {
		t.Errorf("Expected npm packages ['axios'], got %v", analysis.NpmPackages)
	}
}

func TestAnalyzeImportsNamed(t *testing.T) {
	// Create AST for: import { join, dirname } from 'path'
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "NamedImports",
								Children: []transpiler.ASTNode{
									{
										Kind: "ImportSpecifier",
										Name: "join",
									},
									{
										Kind: "ImportSpecifier",
										Name: "dirname",
									},
								},
							},
						},
					},
					{
						Kind: "StringLiteral",
						Text: "path",
					},
				},
			},
		},
	}

	analysis, err := AnalyzeImports(ast)
	if err != nil {
		t.Fatalf("AnalyzeImports failed: %v", err)
	}

	if len(analysis.Imports) != 1 {
		t.Fatalf("Expected 1 import, got %d", len(analysis.Imports))
	}

	imp := analysis.Imports[0]
	if imp.Source != "path" {
		t.Errorf("Expected source 'path', got '%s'", imp.Source)
	}

	if imp.Type != ImportTypeBuiltin {
		t.Errorf("Expected ImportTypeBuiltin, got %v", imp.Type)
	}

	if len(imp.Symbols) != 2 {
		t.Fatalf("Expected 2 symbols, got %d", len(imp.Symbols))
	}

	if imp.Symbols[0] != "join" || imp.Symbols[1] != "dirname" {
		t.Errorf("Expected symbols ['join', 'dirname'], got %v", imp.Symbols)
	}

	if len(analysis.BuiltinModules) != 1 || analysis.BuiltinModules[0] != "path" {
		t.Errorf("Expected builtin modules ['path'], got %v", analysis.BuiltinModules)
	}
}

func TestAnalyzeImportsNamespace(t *testing.T) {
	// Create AST for: import * as fs from 'fs'
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "NamespaceImport",
								Children: []transpiler.ASTNode{
									{
										Kind: "Identifier",
										Text: "fs",
									},
								},
							},
						},
					},
					{
						Kind: "StringLiteral",
						Text: "fs",
					},
				},
			},
		},
	}

	analysis, err := AnalyzeImports(ast)
	if err != nil {
		t.Fatalf("AnalyzeImports failed: %v", err)
	}

	if len(analysis.Imports) != 1 {
		t.Fatalf("Expected 1 import, got %d", len(analysis.Imports))
	}

	imp := analysis.Imports[0]
	if !imp.IsNamespace {
		t.Error("Expected IsNamespace to be true")
	}

	if imp.Alias != "fs" {
		t.Errorf("Expected alias 'fs', got '%s'", imp.Alias)
	}

	if len(imp.Symbols) != 1 || imp.Symbols[0] != "*" {
		t.Errorf("Expected symbols ['*'], got %v", imp.Symbols)
	}
}

func TestAnalyzeImportsLocal(t *testing.T) {
	// Create AST for: import { User } from './models/user'
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "NamedImports",
								Children: []transpiler.ASTNode{
									{
										Kind: "ImportSpecifier",
										Name: "User",
									},
								},
							},
						},
					},
					{
						Kind: "StringLiteral",
						Text: "./models/user",
					},
				},
			},
		},
	}

	analysis, err := AnalyzeImports(ast)
	if err != nil {
		t.Fatalf("AnalyzeImports failed: %v", err)
	}

	if len(analysis.Imports) != 1 {
		t.Fatalf("Expected 1 import, got %d", len(analysis.Imports))
	}

	imp := analysis.Imports[0]
	if imp.Type != ImportTypeLocal {
		t.Errorf("Expected ImportTypeLocal, got %v", imp.Type)
	}

	if len(analysis.LocalFiles) != 1 || analysis.LocalFiles[0] != "./models/user" {
		t.Errorf("Expected local files ['./models/user'], got %v", analysis.LocalFiles)
	}
}

func TestAnalyzeImportsMultiple(t *testing.T) {
	// Create AST with multiple imports
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{Kind: "Identifier", Text: "axios"},
						},
					},
					{Kind: "StringLiteral", Text: "axios"},
				},
			},
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "NamedImports",
								Children: []transpiler.ASTNode{
									{Kind: "ImportSpecifier", Name: "readFile"},
								},
							},
						},
					},
					{Kind: "StringLiteral", Text: "fs"},
				},
			},
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "NamedImports",
								Children: []transpiler.ASTNode{
									{Kind: "ImportSpecifier", Name: "User"},
								},
							},
						},
					},
					{Kind: "StringLiteral", Text: "./models/user"},
				},
			},
		},
	}

	analysis, err := AnalyzeImports(ast)
	if err != nil {
		t.Fatalf("AnalyzeImports failed: %v", err)
	}

	if len(analysis.Imports) != 3 {
		t.Fatalf("Expected 3 imports, got %d", len(analysis.Imports))
	}

	if len(analysis.NpmPackages) != 1 {
		t.Errorf("Expected 1 npm package, got %d", len(analysis.NpmPackages))
	}

	if len(analysis.BuiltinModules) != 1 {
		t.Errorf("Expected 1 builtin module, got %d", len(analysis.BuiltinModules))
	}

	if len(analysis.LocalFiles) != 1 {
		t.Errorf("Expected 1 local file, got %d", len(analysis.LocalFiles))
	}
}

func TestGetImportBySource(t *testing.T) {
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{Kind: "Identifier", Text: "axios"},
						},
					},
					{Kind: "StringLiteral", Text: "axios"},
				},
			},
		},
	}

	analysis, _ := AnalyzeImports(ast)

	imp, exists := analysis.GetImportBySource("axios")
	if !exists {
		t.Error("Expected to find import for 'axios'")
	}

	if imp.Source != "axios" {
		t.Errorf("Expected source 'axios', got '%s'", imp.Source)
	}

	_, exists = analysis.GetImportBySource("nonexistent")
	if exists {
		t.Error("Should not find nonexistent import")
	}
}

func TestHasPackage(t *testing.T) {
	analysis := &ImportAnalysis{
		NpmPackages: []string{"axios", "lodash"},
	}

	if !analysis.HasPackage("axios") {
		t.Error("Expected to have 'axios'")
	}

	if analysis.HasPackage("nonexistent") {
		t.Error("Should not have 'nonexistent'")
	}
}

func TestHasBuiltin(t *testing.T) {
	analysis := &ImportAnalysis{
		BuiltinModules: []string{"fs", "path"},
	}

	if !analysis.HasBuiltin("fs") {
		t.Error("Expected to have 'fs'")
	}

	if analysis.HasBuiltin("axios") {
		t.Error("Should not have 'axios' as builtin")
	}
}

func TestGetDependencies(t *testing.T) {
	analysis := &ImportAnalysis{
		NpmPackages:    []string{"axios", "lodash"},
		BuiltinModules: []string{"fs", "path"},
	}

	deps := analysis.GetDependencies()
	if len(deps) != 4 {
		t.Errorf("Expected 4 dependencies, got %d", len(deps))
	}
}

func TestImportTypesString(t *testing.T) {
	tests := []struct {
		importType ImportType
		expected   string
	}{
		{ImportTypePackage, "package"},
		{ImportTypeLocal, "local"},
		{ImportTypeBuiltin, "builtin"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.importType.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
