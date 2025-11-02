package transpiler

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yourusername/ts2go/pkg/optimizer"
)

// TranspileOptions configures the transpilation process
type TranspileOptions struct {
	Optimize bool // Enable code optimization (dead code elimination, etc.)
}

// Transpile takes a TS file and outputs Go code
func Transpile(input, output string) error {
	return TranspileWithOptions(input, output, TranspileOptions{Optimize: true})
}

// TranspileWithOptions takes a TS file and outputs Go code with custom options
func TranspileWithOptions(input, output string, options TranspileOptions) error {
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

	// Step 3: Optimize if requested
	if options.Optimize {
		goCode, err = OptimizeCode(goCode)
		if err != nil {
			// Don't fail on optimization errors - just use unoptimized code
			fmt.Fprintf(os.Stderr, "Warning: optimization failed: %v\n", err)
		}
	}

	// Step 4: Write output file
	if err := os.WriteFile(output, []byte(goCode), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// ParseTypeScript uses the Node.js parser to convert TS to AST (public API)
func ParseTypeScript(inputFile string) (*ASTNode, error) {
	return parseTypeScript(inputFile)
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
		return nil, ParseError(inputFile, "parser.js not found. Run 'npm install' in internal/transpiler/parser/")
	}

	// Execute the Node.js parser
	cmd := exec.Command("node", parserPath, inputFile)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg := string(exitErr.Stderr)
			return nil, ParseError(inputFile, errMsg)
		}
		return nil, ParseError(inputFile, err.Error())
	}

	// Parse JSON output
	var ast ASTNode
	if err := json.Unmarshal(output, &ast); err != nil {
		return nil, ParseError(inputFile, fmt.Sprintf("failed to parse AST JSON: %v", err))
	}

	return &ast, nil
}

// OptimizeCode performs optimizations on generated Go code
func OptimizeCode(code string) (string, error) {
	opt := optimizer.NewOptimizer()
	return opt.Optimize(code)
}
