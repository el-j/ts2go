# TS2Go Architecture

## Overview

TS2Go is a TypeScript-to-Go transpiler designed as a mono-repo with clear separation of concerns.

## Project Structure

```
ts2go/
├── cmd/ts2go/           # CLI entry point
├── internal/transpiler/ # Core transpilation logic
│   ├── parser/         # Node.js TypeScript parser
│   ├── ast.go          # AST node definitions
│   ├── codegen.go      # Go code generation
│   └── transpiler.go   # Main transpiler orchestration
├── runtime/            # Go runtime library for Node.js APIs
│   ├── console/        # console.log, console.error, etc.
│   ├── fs/             # File system operations
│   └── path/           # Path manipulation
├── tests/              # Integration tests
│   └── fixtures/       # Test TypeScript files
└── docs/               # Documentation
```

## Components

### 1. CLI (`cmd/ts2go/`)

Simple command-line interface using Go's `flag` package:
- Parses `--in` and `--out` arguments
- Calls the transpiler
- Handles errors and exit codes

### 2. Parser (`internal/transpiler/parser/`)

Node.js script that:
- Uses TypeScript's official compiler API
- Reads TypeScript source files
- Generates a simplified AST as JSON
- Outputs to stdout for consumption by Go

**Why Node.js?** TypeScript's compiler API is the authoritative way to parse TypeScript. Rather than rewrite a parser, we leverage the existing tool.

### 3. Transpiler (`internal/transpiler/`)

Core Go package with three files:

**ast.go**: Defines the AST node structure
- `ASTNode` struct maps to TypeScript's SyntaxKind
- Constants for common node types

**codegen.go**: Code generation logic
- `CodeGenerator` walks the AST
- Generates Go code string by string
- Handles indentation and formatting
- Maps TypeScript constructs to Go equivalents

**transpiler.go**: Orchestration
- `Transpile()` is the main entry point
- Calls Node.js parser
- Invokes code generator
- Writes output file

### 4. Runtime (`runtime/`)

Go implementations of Node.js standard library APIs:

- **fs/**: File system operations (`readFile`, etc.)
- **console/**: Console output (`console.log`, `console.error`)
- **path/**: Path manipulation (`join`, `dirname`, etc.)

When transpiled code imports Node.js modules, they're rewritten to use these runtime packages.

## Data Flow

```
TypeScript File
      ↓
   [parser.js]  ← TypeScript Compiler API
      ↓
   JSON AST
      ↓
   [transpiler.go]
      ↓
   [codegen.go]
      ↓
   Go Source Code
```

## Key Design Decisions

### Type Mapping

| TypeScript | Go |
|------------|-----|
| `string` | `string` |
| `number` | `float64` |
| `boolean` | `bool` |
| `Array<T>` | `[]T` |
| `interface` | `struct` |
| `any` | ❌ Disallowed |

### Name Conventions

- TypeScript: `camelCase` for functions and variables
- Go: `PascalCase` for exported functions, `camelCase` for local variables
- Transpiler automatically converts interface/function names to PascalCase

### Error Handling

- TypeScript: `try/catch`
- Go: Multiple return values `(result, error)`
- Transpiler maps exceptions to error returns

### Imports

- Node.js APIs: Rewritten to runtime package imports
- Local files: Transpiled and imported as Go packages
- npm modules: ❌ Not supported (must be rewritten)

## Testing Strategy

### Unit Tests
Test individual transpiler functions:
- Type mapping
- Name conversion
- AST node handling

### Integration Tests
Full end-to-end tests:
1. Create `.ts` fixture file
2. Run transpiler
3. Compile generated Go code
4. Run and verify output

## Future Enhancements

See [implementationPlan.md](../implementationPlan.md) for:
- Phase 5: Classes, modules, async/await
- Phase 6: Comprehensive testing and release
