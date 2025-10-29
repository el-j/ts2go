package transpiler

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Transpile takes a TS file and outputs Go code
func Transpile(input, output string) error {
	// Step 1: Parse TypeScript using Node.js parser
	ast, err := parseTypeScript(input)
	if err != nil {
		return fmt.Errorf("failed to parse TypeScript: %w", err)
	}

	// Step 2: Generate Go code from AST
	generator := NewCodeGenerator()
	goCode, err := generator.Generate(ast)
	if err != nil {
		return fmt.Errorf("failed to generate Go code: %w", err)
	}

	// Step 3: Write output file
	if err := os.WriteFile(output, []byte(goCode), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// parseTypeScript uses the Node.js parser to convert TS to AST
func parseTypeScript(inputFile string) (*ASTNode, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Try to find parser.js in several locations
	parserPaths := []string{
		"./internal/transpiler/parser/parser.js",
		"../internal/transpiler/parser/parser.js",
		"../../internal/transpiler/parser/parser.js",
		filepath.Join(cwd, "internal", "transpiler", "parser", "parser.js"),
		filepath.Join(filepath.Dir(cwd), "internal", "transpiler", "parser", "parser.js"),
	}

	// Also try relative to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		parserPaths = append(parserPaths,
			filepath.Join(exeDir, "internal", "transpiler", "parser", "parser.js"),
			filepath.Join(exeDir, "..", "internal", "transpiler", "parser", "parser.js"),
		)
	}

	parserPath := ""
	for _, path := range parserPaths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			parserPath = absPath
			break
		}
	}

	if parserPath == "" {
		return nil, fmt.Errorf("parser.js not found. Run 'npm install' in internal/transpiler/parser/")
	}

	// Execute the Node.js parser
	cmd := exec.Command("node", parserPath, inputFile)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("parser failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute parser: %w", err)
	}

	// Parse JSON output
	var ast ASTNode
	if err := json.Unmarshal(output, &ast); err != nil {
		return nil, fmt.Errorf("failed to parse AST JSON: %w", err)
	}

	return &ast, nil
}
