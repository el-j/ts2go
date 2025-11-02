# Getting Started with ts2go

## Overview

**ts2go** is a TypeScript-to-Go transpiler that automatically converts TypeScript code into idiomatic Go code. It handles type conversions, class-to-struct transformations, async patterns, and npm package mapping.

**Time to first transpilation**: ~10 minutes

## Prerequisites

- **Go**: 1.22.5 or later
- **Node.js**: 16+ (for TypeScript parsing)
- **npm**: For installing parser dependencies

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/ts2go.git
cd ts2go
```

### 2. Install Dependencies

```bash
# Install Node.js parser dependencies
cd internal/transpiler/parser
npm install
cd ../../..

# Build the transpiler
make build
# OR
go build -o ts2go ./cmd/ts2go
```

### 3. Verify Installation

```bash
./ts2go help
```

You should see the help message with available commands.

## Quick Start

### Your First Transpilation

**1. Create a simple TypeScript file (`hello.ts`):**

```typescript
function greet(name: string): string {
    return `Hello, ${name}!`;
}

console.log(greet("World"));
```

**2. Transpile it:**

```bash
./ts2go convert hello.ts hello.go
```

**3. Run the Go code:**

```bash
go run hello.go
```

Output:
```
Hello, World!
```

## Basic Usage

### Single File Conversion

```bash
ts2go convert <input.ts> <output.go>
```

Example:
```bash
ts2go convert src/utils.ts output/utils.go
```

### Project Transpilation

```bash
ts2go transpile <project-directory> [options]
```

Options:
- `--out, -o <dir>`: Output directory (default: `output`)
- `--verbose, -v`: Show detailed progress
- `--quiet, -q`: Suppress progress output
- `--watch, -w`: Watch for changes and auto-transpile

Example:
```bash
# Basic transpilation
ts2go transpile ./my-ts-project

# Custom output directory
ts2go transpile ./my-ts-project --out ./go-output

# Watch mode for development
ts2go transpile ./my-ts-project --watch --verbose
```

### Project Analysis

```bash
ts2go analyze <project-directory>
```

Analyzes your TypeScript project and shows:
- File count and structure
- Import dependencies
- npm package usage
- Type complexity

## Type Mapping

### Primitive Types

| TypeScript | Go |
|------------|-----|
| `number` | `float64` |
| `string` | `string` |
| `boolean` | `bool` |
| `any` | `interface{}` |
| `void` | (no return value) |
| `null`, `undefined` | `nil` |

### Complex Types

| TypeScript | Go |
|------------|-----|
| `Array<T>` | `[]T` |
| `T[]` | `[]T` |
| `{ [key: string]: T }` | `map[string]T` |
| `interface` | `struct` |
| `enum` | `const` block or `type` |
| `tuple` | `struct` with numbered fields |

### Example: Interface to Struct

**TypeScript:**
```typescript
interface User {
    id: number;
    name: string;
    email?: string;
}
```

**Go:**
```go
type User struct {
    ID    float64
    Name  string
    Email *string // optional field becomes pointer
}
```

## Feature Support

### ✅ Fully Supported

- Basic types (number, string, boolean)
- Functions and arrow functions
- Interfaces → Structs
- Classes → Structs with methods
- Enums → Constants
- Import/Export statements
- Template literals → fmt.Sprintf
- Optional chaining (`?.`)
- Nullish coalescing (`??`)
- Type assertions
- Union types
- Generics (basic support)
- Tuples → Structs

### ⚠️ Partial Support

- Async/await → Goroutines + channels (requires manual review)
- Decorators (limited support)
- Advanced generics (complex constraints)
- Mixins (converted to composition)

### ❌ Not Supported

- JSX/TSX (React components)
- Namespaces (use packages instead)
- Triple-slash directives
- Ambient declarations

## NPM Package Mapping

ts2go automatically maps common npm packages to their Go equivalents:

| npm Package | Go Package |
|-------------|------------|
| `express` | `github.com/gin-gonic/gin` |
| `axios` | `net/http` |
| `lodash` | `github.com/samber/lo` |
| `uuid` | `github.com/google/uuid` |
| `moment` | `time` package |

See [PACKAGE_MAPPINGS.md](PACKAGE_MAPPINGS.md) for complete list.

## Common Patterns

### 1. Error Handling

**TypeScript:**
```typescript
function divide(a: number, b: number): number {
    if (b === 0) {
        throw new Error("Division by zero");
    }
    return a / b;
}
```

**Go:**
```go
func Divide(a float64, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("Division by zero")
    }
    return a / b, nil
}
```

### 2. Classes

**TypeScript:**
```typescript
class Counter {
    private count: number = 0;
    
    increment(): void {
        this.count++;
    }
    
    getCount(): number {
        return this.count;
    }
}
```

**Go:**
```go
type Counter struct {
    count float64
}

func NewCounter() *Counter {
    return &Counter{count: 0}
}

func (c *Counter) Increment() {
    c.count++
}

func (c *Counter) GetCount() float64 {
    return c.count
}
```

### 3. Async Operations

**TypeScript:**
```typescript
async function fetchData(url: string): Promise<string> {
    const response = await fetch(url);
    return response.text();
}
```

**Go:**
```go
func FetchData(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    return string(body), err
}
```

## Workflow

### Development Workflow

1. **Initial Transpilation**
   ```bash
   ts2go transpile ./ts-project --out ./go-project
   ```

2. **Review Generated Code**
   - Check for any warnings
   - Review async/await conversions
   - Verify error handling

3. **Setup Go Module**
   ```bash
   cd go-project
   go mod tidy
   ```

4. **Build and Test**
   ```bash
   go build
   go test ./...
   ```

5. **Watch Mode for Development**
   ```bash
   ts2go transpile ./ts-project --watch
   ```

### Production Workflow

1. **Transpile with Optimization**
   ```bash
   ts2go transpile ./ts-project --out ./production
   ```
   (Optimization is enabled by default)

2. **Run Tests**
   ```bash
   cd production
   go test ./...
   ```

3. **Build Release**
   ```bash
   go build -ldflags="-s -w" -o app
   ```

## Troubleshooting

### Parser Not Found

**Error:** `parser.js not found`

**Solution:**
```bash
cd internal/transpiler/parser
npm install
```

### Import Resolution Failed

**Error:** `Import resolution error`

**Solution:**
- Check that npm package has a Go equivalent
- See [PACKAGE_MAPPINGS.md](PACKAGE_MAPPINGS.md)
- Add custom mapping if needed

### Build Errors in Generated Code

**Common issues:**
1. **Undefined variables**: Check variable scoping
2. **Type mismatches**: May need manual type adjustments
3. **Missing imports**: Add required Go imports

**Solution:**
- Use `--verbose` flag to see detailed errors
- Review the specific generated file
- Consult [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)

## Next Steps

- Read the [Migration Guide](MIGRATION_GUIDE.md) for detailed conversion patterns
- Check [EXAMPLES.md](EXAMPLES.md) for more complex examples
- See [PACKAGE_MAPPINGS.md](PACKAGE_MAPPINGS.md) for npm package equivalents
- Review [ARCHITECTURE.md](ARCHITECTURE.md) to understand how it works

## Getting Help

- **Issues**: Open an issue on GitHub
- **Discussions**: Check GitHub Discussions
- **Examples**: See `examples/` directory
- **Documentation**: Full docs in `docs/` directory

## Example Project

See the `examples/rest-api` directory for a complete example of transpiling an Express.js API to Go with Gin.

```bash
cd examples/rest-api
ts2go transpile ./ts-src --out ./go-src
cd go-src
go run main.go
```

---

**Ready to transpile?** Start with `ts2go convert` for a single file, or `ts2go transpile` for a full project!
