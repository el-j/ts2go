# ts2go API Reference

## Command Line Interface

### ts2go convert

Convert a single TypeScript file to Go.

**Syntax:**
```bash
ts2go convert <input.ts> <output.go>
```

**Arguments:**
- `input.ts` - Path to TypeScript source file
- `output.go` - Path to output Go file

**Example:**
```bash
ts2go convert src/utils.ts output/utils.go
```

**Exit Codes:**
- `0` - Success
- `1` - Transpilation error

---

### ts2go transpile

Transpile an entire TypeScript project to Go.

**Syntax:**
```bash
ts2go transpile <project-dir> [options]
```

**Arguments:**
- `project-dir` - Path to TypeScript project directory

**Options:**
- `--out, -o <dir>` - Output directory (default: `output`)
- `--verbose, -v` - Show detailed progress information
- `--quiet, -q` - Suppress all progress output (errors only)
- `--watch, -w` - Watch for file changes and auto-transpile

**Examples:**
```bash
# Basic transpilation
ts2go transpile ./my-project

# Custom output directory
ts2go transpile ./my-project --out ./go-output

# Verbose mode
ts2go transpile ./my-project --verbose

# Watch mode
ts2go transpile ./my-project --watch

# Quiet mode (CI/CD)
ts2go transpile ./my-project --quiet
```

**Output:**
- Transpiled Go files in output directory
- `go.mod` with required dependencies
- Build order report
- Dependency summary

**Exit Codes:**
- `0` - Success
- `1` - Transpilation errors

---

### ts2go analyze

Analyze a TypeScript project without transpiling.

**Syntax:**
```bash
ts2go analyze <project-dir>
```

**Arguments:**
- `project-dir` - Path to TypeScript project directory

**Output:**
Shows:
- Number of files found
- Import analysis
- Dependency classification
- Package.json dependencies
- Estimated Go dependencies

**Example:**
```bash
ts2go analyze ./my-project
```

---

### ts2go help

Display help information.

**Syntax:**
```bash
ts2go help
ts2go --help
ts2go -h
```

---

## Go API

### Transpiler Package

Import:
```go
import "github.com/yourusername/ts2go/internal/transpiler"
```

#### Transpile

Transpile a single TypeScript file to Go with default options.

```go
func Transpile(input, output string) error
```

**Parameters:**
- `input` - Path to TypeScript source file
- `output` - Path to output Go file

**Returns:**
- `error` - Transpilation error, or nil on success

**Example:**
```go
err := transpiler.Transpile("app.ts", "app.go")
if err != nil {
    log.Fatal(err)
}
```

---

#### TranspileWithOptions

Transpile with custom options.

```go
func TranspileWithOptions(input, output string, options TranspileOptions) error
```

**Parameters:**
- `input` - Path to TypeScript source file
- `output` - Path to output Go file
- `options` - TranspileOptions configuration

**TranspileOptions:**
```go
type TranspileOptions struct {
    Optimize bool // Enable code optimization (default: true)
}
```

**Example:**
```go
opts := transpiler.TranspileOptions{
    Optimize: false, // Disable optimization
}
err := transpiler.TranspileWithOptions("app.ts", "app.go", opts)
```

---

#### ParseTypeScript

Parse TypeScript file into AST.

```go
func ParseTypeScript(inputFile string) (*ASTNode, error)
```

**Parameters:**
- `inputFile` - Path to TypeScript file

**Returns:**
- `*ASTNode` - Parsed AST
- `error` - Parse error, or nil

**Example:**
```go
ast, err := transpiler.ParseTypeScript("app.ts")
if err != nil {
    log.Fatal(err)
}
```

---

#### OptimizeCode

Optimize generated Go code.

```go
func OptimizeCode(code string) (string, error)
```

**Parameters:**
- `code` - Go source code string

**Returns:**
- `string` - Optimized Go code
- `error` - Optimization error, or nil

**Example:**
```go
optimized, err := transpiler.OptimizeCode(goCode)
if err != nil {
    log.Fatal(err)
}
```

---

### Error Types

#### TranspilationError

Detailed transpilation error with context.

```go
type TranspilationError struct {
    File       string   // Source file path
    Line       int      // Line number (1-based)
    Column     int      // Column number (1-based)
    Message    string   // Error message
    Code       string   // Error code
    Suggestion string   // Suggested fix
    Context    []string // Code context lines
}
```

**Methods:**
```go
func (e *TranspilationError) Error() string
func (e *TranspilationError) WithCode(code string) *TranspilationError
func (e *TranspilationError) WithSuggestion(suggestion string) *TranspilationError
func (e *TranspilationError) WithContext(lines []string) *TranspilationError
```

**Error Constructors:**
```go
func UnsupportedFeatureError(file string, line, column int, feature string) *TranspilationError
func TypeConversionError(file string, line, column int, fromType, toType string) *TranspilationError
func ParseError(file string, message string) *TranspilationError
func ImportResolutionError(file string, line, column int, importPath string) *TranspilationError
```

**Example:**
```go
err := transpiler.UnsupportedFeatureError(
    "app.ts",
    10,
    5,
    "decorators",
)
fmt.Println(err.Error())
// Output:
// app.ts:10:5: Unsupported feature: decorators
//   Code: UNSUPPORTED_FEATURE
// 💡 Suggestion: Check the documentation for supported features
```

---

### Optimizer Package

Import:
```go
import "github.com/yourusername/ts2go/pkg/optimizer"
```

#### Optimizer

```go
type Optimizer struct {
    // private fields
}
```

**Constructor:**
```go
func NewOptimizer() *Optimizer
```

**Methods:**

##### Optimize

Perform all optimizations on Go code.

```go
func (o *Optimizer) Optimize(code string) (string, error)
```

**Optimizations performed:**
- Remove unused imports
- Remove unused functions
- Remove unused variables
- Clean up empty blocks

**Example:**
```go
opt := optimizer.NewOptimizer()
optimized, err := opt.Optimize(goCode)
if err != nil {
    log.Fatal(err)
}
```

##### OptimizeImports

Optimize imports only.

```go
func (o *Optimizer) OptimizeImports(code string) (string, error)
```

**Example:**
```go
opt := optimizer.NewOptimizer()
cleaned, err := opt.OptimizeImports(goCode)
```

---

### CLI Package

Import:
```go
import "github.com/yourusername/ts2go/pkg/cli"
```

#### Progress Reporting

##### ProgressLevel

```go
type ProgressLevel int

const (
    ProgressQuiet   ProgressLevel = iota // Only errors
    ProgressNormal                       // Standard output
    ProgressVerbose                      // Detailed output
)
```

##### ProgressReporter

```go
type ProgressReporter struct {
    // private fields
}
```

**Constructor:**
```go
func NewProgressReporter(level ProgressLevel) *ProgressReporter
```

**Methods:**
```go
func (p *ProgressReporter) Start(operation string, totalSteps int)
func (p *ProgressReporter) Step(description string)
func (p *ProgressReporter) StepVerbose(description string, details string)
func (p *ProgressReporter) Success(description string)
func (p *ProgressReporter) Warning(description string)
func (p *ProgressReporter) Error(description string)
func (p *ProgressReporter) Info(format string, args ...interface{})
func (p *ProgressReporter) Verbose(format string, args ...interface{})
func (p *ProgressReporter) Complete(message string)
func (p *ProgressReporter) SetLevel(level ProgressLevel)
```

**Example:**
```go
progress := cli.NewProgressReporter(cli.ProgressVerbose)
progress.Start("Processing Files", 100)

for i := 0; i < 100; i++ {
    progress.Step(fmt.Sprintf("Processing file %d", i))
    // do work
}

progress.Complete("All files processed")
```

##### ProgressBar

Visual progress bar for terminal.

```go
type ProgressBar struct {
    // private fields
}
```

**Constructor:**
```go
func NewProgressBar(total int, width int) *ProgressBar
```

**Methods:**
```go
func (pb *ProgressBar) Update(current int)
func (pb *ProgressBar) Increment()
func (pb *ProgressBar) Render()
func (pb *ProgressBar) Clear()
```

**Example:**
```go
bar := cli.NewProgressBar(100, 40)
for i := 0; i <= 100; i++ {
    bar.Update(i)
    time.Sleep(50 * time.Millisecond)
}
```

##### Spinner

Animated spinner for long operations.

```go
type Spinner struct {
    // private fields
}
```

**Constructor:**
```go
func NewSpinner() *Spinner
```

**Methods:**
```go
func (s *Spinner) Start(message string)
func (s *Spinner) Tick()
func (s *Spinner) Stop()
```

**Example:**
```go
spinner := cli.NewSpinner()
spinner.Start("Loading...")

ticker := time.NewTicker(100 * time.Millisecond)
go func() {
    for range ticker.C {
        spinner.Tick()
    }
}()

// do long operation
time.Sleep(3 * time.Second)

ticker.Stop()
spinner.Stop()
```

---

#### Watch Mode

##### WatchProject

Watch a project directory for changes and auto-transpile.

```go
func WatchProject(projectDir, outputDir string, progressLevel ProgressLevel) error
```

**Parameters:**
- `projectDir` - TypeScript project directory
- `outputDir` - Go output directory
- `progressLevel` - Progress reporting level

**Features:**
- Watches all .ts files (except .d.ts)
- 500ms debouncing for rapid changes
- Initial full transpilation on start
- Ignores node_modules and hidden directories

**Example:**
```go
err := cli.WatchProject(
    "./ts-project",
    "./go-output",
    cli.ProgressVerbose,
)
if err != nil {
    log.Fatal(err)
}
```

---

### Analyzer Package

Import:
```go
import "github.com/yourusername/ts2go/internal/analyzer"
```

#### AnalyzeImports

Analyze imports from TypeScript AST.

```go
func AnalyzeImports(astNode *transpiler.ASTNode) (*ImportAnalysis, error)
```

**Returns:**
```go
type ImportAnalysis struct {
    Imports      []Import
    PackageFiles []string
    LocalFiles   []string
}

type Import struct {
    Path       string
    ImportType ImportType
    Symbols    []string
    IsDefault  bool
    IsNamespace bool
}
```

**Example:**
```go
ast, _ := transpiler.ParseTypeScript("app.ts")
analysis, err := analyzer.AnalyzeImports(ast)
if err != nil {
    log.Fatal(err)
}

for _, imp := range analysis.Imports {
    fmt.Printf("%s: %s\n", imp.ImportType, imp.Path)
}
```

#### AnalyzeFile

Convenience function to analyze a file directly.

```go
func AnalyzeFile(filePath string) (*ImportAnalysis, error)
```

**Example:**
```go
analysis, err := analyzer.AnalyzeFile("app.ts")
```

---

## Environment Variables

### PARSER_PATH

Override the default parser.js location.

```bash
export PARSER_PATH=/custom/path/to/parser.js
ts2go convert app.ts app.go
```

---

## Configuration Files

### go.mod Generation

ts2go automatically generates `go.mod` with:
- Module path (based on project name)
- Go version (1.22.5+)
- Required dependencies (based on npm packages)

**Example generated go.mod:**
```go
module github.com/user/myproject

go 1.22.5

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/google/uuid v1.4.0
)
```

---

## Exit Codes

All ts2go commands use standard exit codes:

- `0` - Success
- `1` - General error
- `2` - Command line usage error

---

## Logging

### Verbosity Levels

1. **Quiet** (`--quiet`): Only errors
2. **Normal** (default): Standard progress
3. **Verbose** (`--verbose`): Detailed information

### Output Format

**Normal mode:**
```
[1/5] (20%) Scanning project...
  ✓ Found 25 TypeScript files
[2/5] (40%) Analyzing imports...
```

**Verbose mode:**
```
[1/5] (20%) Scanning project...
    Found: src/index.ts
    Found: src/utils.ts
  ✓ Found 25 TypeScript files
```

**Quiet mode:**
```
(no output unless errors occur)
```

---

## Performance

### Benchmarks

Average transpilation times (M1 Mac):
- Single file: ~150ms
- Small project (10 files): ~1.5s
- Medium project (100 files): ~12s
- Large project (1000 files): ~2min

With optimization enabled (default):
- Adds ~10-15% overhead
- Results in 10-30% smaller output files

---

## Version Information

Check version:
```bash
ts2go --version  # Coming soon
```

---

## See Also

- [Getting Started Guide](GETTING_STARTED_v2.md)
- [Migration Guide](MIGRATION_GUIDE.md)
- [Package Mappings](PACKAGE_MAPPINGS.md)
- [Architecture](ARCHITECTURE.md)
