package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/el-j/ts2go/internal/transpiler"
)

func TestNewApplication(t *testing.T) {
	app, err := NewApplication()
	if err != nil {
		t.Fatalf("NewApplication failed: %v", err)
	}
	if app == nil {
		t.Fatalf("NewApplication returned nil app")
	}

	if app.compiler == nil || app.fileSystem == nil || app.stateRepo == nil || app.settingsRepo == nil {
		t.Errorf("NewApplication missing dependencies")
	}
}

func TestAdapterAnalyzer_AnalyzeImports(t *testing.T) {
	analyzerAdapter := &adapterAnalyzer{}

	// Create a minimal synthetic AST
	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "ImportDeclaration",
				Children: []transpiler.ASTNode{
					{
						Kind: "StringLiteral",
						Text: "fmt",
					},
					{
						Kind: "ImportClause",
						Children: []transpiler.ASTNode{
							{
								Kind: "Identifier",
								Text: "fmt",
							},
						},
					},
				},
			},
		},
	}

	analysis, err := analyzerAdapter.AnalyzeImports(ast)
	if err != nil {
		t.Fatalf("AnalyzeImports failed: %v", err)
	}

	if analysis == nil {
		t.Fatalf("Expected non-nil analysis")
	}

	// Test invalid type
	_, err = analyzerAdapter.AnalyzeImports("not an ast")
	if err == nil {
		t.Errorf("Expected error when passing string instead of *ASTNode")
	}
}

func TestAdapterCodeGen_ParseTypeScript(t *testing.T) {
	codegenAdapter := &adapterCodeGen{}

	tempDir, err := os.MkdirTemp("", "codegen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tsFile := filepath.Join(tempDir, "test.ts")
	os.WriteFile(tsFile, []byte("import { format } from 'util';\nconst x = 1;"), 0644)

	ast, err := codegenAdapter.ParseTypeScript(tsFile)
	if err != nil {
		t.Skipf("ParseTypeScript failed (Node.js may not be available): %v", err)
	}

	if ast == nil {
		t.Errorf("Expected AST to be non-nil")
	}
}

func TestAdapterCodeGen_GenerateGoCode(t *testing.T) {
	codegenAdapter := &adapterCodeGen{}

	ast := &transpiler.ASTNode{
		Kind: "SourceFile",
		Statements: []transpiler.ASTNode{
			{
				Kind: "VariableStatement",
				Declarations: []transpiler.ASTNode{
					{
						Kind: "VariableDeclaration",
						Name: "x",
						Type: &transpiler.ASTNode{Kind: "NumberKeyword"},
						Initializer: &transpiler.ASTNode{
							Kind: "NumericLiteral",
							Text: "42",
						},
					},
				},
			},
		},
	}

	code, err := codegenAdapter.GenerateGoCode(ast)
	if err != nil {
		t.Fatalf("GenerateGoCode failed: %v", err)
	}

	if code == "" {
		t.Errorf("GenerateGoCode returned empty string")
	}

	// Test invalid type
	_, err = codegenAdapter.GenerateGoCode("not an ast")
	if err == nil {
		t.Errorf("Expected error with invalid AST type")
	}
}

func TestApplication_StateCommand(t *testing.T) {
	app, err := NewApplication()
	if err != nil {
		t.Fatalf("NewApplication failed: %v", err)
	}

	tempDir := t.TempDir()

	// Insufficient args
	if err := app.StateCommand([]string{}); err == nil {
		t.Errorf("expected error for empty args")
	}

	// Unknown action
	if err := app.StateCommand([]string{"invalid_action", tempDir}); err == nil {
		t.Errorf("expected error for unknown action")
	}

	// Get state for non-existent project (returns empty JSON, no error)
	if err := app.StateCommand([]string{"get", tempDir}); err != nil {
		t.Errorf("StateCommand get failed: %v", err)
	}

	// Set state missing json
	if err := app.StateCommand([]string{"set", tempDir}); err == nil {
		t.Errorf("expected error for set without JSON")
	}

	// Set state invalid json
	if err := app.StateCommand([]string{"set", tempDir, "{bad-json"}); err == nil {
		t.Errorf("expected error for set with invalid JSON")
	}

	// Set state valid json
	if err := app.StateCommand([]string{"set", tempDir, `{"projectPath":"` + tempDir + `"}`}); err != nil {
		t.Errorf("StateCommand set failed: %v", err)
	}

	// Get state after set
	if err := app.StateCommand([]string{"get", tempDir}); err != nil {
		t.Errorf("StateCommand get failed: %v", err)
	}

	// Delete state
	if err := app.StateCommand([]string{"delete", tempDir}); err != nil {
		t.Errorf("StateCommand delete failed: %v", err)
	}
}

func TestApplication_SettingsCommand(t *testing.T) {
	app, err := NewApplication()
	if err != nil {
		t.Fatalf("NewApplication failed: %v", err)
	}

	// Insufficient args
	if err := app.SettingsCommand([]string{}); err == nil {
		t.Errorf("expected error for empty args")
	}

	// Unknown action
	if err := app.SettingsCommand([]string{"invalid"}); err == nil {
		t.Errorf("expected error for unknown action")
	}

	// Get settings
	if err := app.SettingsCommand([]string{"get"}); err != nil {
		t.Errorf("SettingsCommand get failed: %v", err)
	}

	// Set settings missing json
	if err := app.SettingsCommand([]string{"set"}); err == nil {
		t.Errorf("expected error for set without JSON")
	}

	// Set settings invalid json
	if err := app.SettingsCommand([]string{"set", "{bad-json"}); err == nil {
		t.Errorf("expected error for set with invalid JSON")
	}

	// Set settings valid json
	if err := app.SettingsCommand([]string{"set", `{"theme":"dark","autoSave":true,"autoSaveDelay":500}`}); err != nil {
		t.Errorf("SettingsCommand set failed: %v", err)
	}
}

func TestApplication_ValidationErrors(t *testing.T) {
	app, err := NewApplication()
	if err != nil {
		t.Fatalf("NewApplication failed: %v", err)
	}

	// TranspileCommand empty
	if err := app.TranspileCommand([]string{}); err == nil {
		t.Errorf("expected error for empty transpile args")
	}

	// AnalyzeCommand empty
	if err := app.AnalyzeCommand([]string{}); err == nil {
		t.Errorf("expected error for empty analyze args")
	}

	// BuildCommand empty
	if err := app.BuildCommand([]string{}); err == nil {
		t.Errorf("expected error for empty build args")
	}

	// TestCommand empty
	if err := app.TestCommand([]string{}); err == nil {
		t.Errorf("expected error for empty test args")
	}
}
