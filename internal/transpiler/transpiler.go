package transpiler

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/el-j/ts2go/pkg/optimizer"
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

// findNodeBinary locates the Node.js binary, preferring bundled version
func findNodeBinary() string {
	// Try to find bundled node first (for standalone app)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)

		// For macOS app bundle: .app/Contents/MacOS/ts2go-cli -> .app/Contents/Resources/bin/node
		if runtime.GOOS == "darwin" {
			if contentsDir := filepath.Dir(exeDir); filepath.Base(contentsDir) == "Contents" {
				bundledNode := filepath.Join(contentsDir, "Resources", "bin", "node")
				if _, err := os.Stat(bundledNode); err == nil {
					return bundledNode
				}
			}
		}

		// For Windows/Linux: near executable
		bundledNode := filepath.Join(exeDir, "node")
		if runtime.GOOS == "windows" {
			bundledNode += ".exe"
		}
		if _, err := os.Stat(bundledNode); err == nil {
			return bundledNode
		}

		// Also try in bin directory next to executable
		bundledNode = filepath.Join(exeDir, "..", "bin", "node")
		if runtime.GOOS == "windows" {
			bundledNode += ".exe"
		}
		if _, err := os.Stat(bundledNode); err == nil {
			return bundledNode
		}
	}

	// Fallback to system node
	// Try common installation paths
	commonPaths := []string{
		"/usr/local/bin/node",
		"/usr/bin/node",
		"/opt/homebrew/bin/node",
	}

	// Also check user's home directory for nvm
	if home, err := os.UserHomeDir(); err == nil {
		commonPaths = append(commonPaths,
			filepath.Join(home, ".nvm", "versions", "node", "v20.11.1", "bin", "node"),
			filepath.Join(home, ".nvm", "versions", "node", "v22.0.0", "bin", "node"),
			filepath.Join(home, ".nvm", "versions", "node", "v18.0.0", "bin", "node"),
		)
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Ultimate fallback: just use "node" and hope it's in PATH
	return "node"
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

	// Also try relative to executable (for bundled app)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)

		// For development and normal builds
		parserPaths = append(parserPaths,
			filepath.Join(exeDir, "internal", "transpiler", "parser", "parser.js"),
			filepath.Join(exeDir, "..", "internal", "transpiler", "parser", "parser.js"),
		)

		// For macOS app bundle: .app/Contents/MacOS/ts2go-cli -> .app/Contents/Resources/bin/parser/
		if runtime.GOOS == "darwin" {
			// Check if we're in a .app bundle structure
			if contentsDir := filepath.Dir(exeDir); filepath.Base(contentsDir) == "Contents" {
				appResourcesParser := filepath.Join(contentsDir, "Resources", "bin", "parser", "parser.js")
				parserPaths = append([]string{appResourcesParser}, parserPaths...)
			}
		}

		// For Windows/Linux bundled location: near executable
		parserPaths = append(parserPaths,
			filepath.Join(exeDir, "parser", "parser.js"),
			filepath.Join(exeDir, "..", "parser", "parser.js"),
			filepath.Join(exeDir, "..", "bin", "parser", "parser.js"),
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

	// Find Node.js binary (bundled or system)
	nodeBinary := findNodeBinary()

	// Execute the Node.js parser
	cmd := exec.Command(nodeBinary, parserPath, inputFile)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg := string(exitErr.Stderr)
			return nil, ParseError(inputFile, fmt.Sprintf("exec: %q: %s", nodeBinary, errMsg))
		}
		return nil, ParseError(inputFile, fmt.Sprintf("exec: %q: %s", nodeBinary, err.Error()))
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
