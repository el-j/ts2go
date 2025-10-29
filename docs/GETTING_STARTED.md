# Getting Started with TS2Go

## Installation

### Prerequisites
- Go 1.21 or higher
- Node.js and npm (for the TypeScript parser)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/your-org/ts2go.git
cd ts2go

# Install npm dependencies for the parser
cd internal/transpiler/parser
npm install
cd ../../..

# Build the CLI tool
make build

# The binary will be at ./ts2go
```

## Basic Usage

### Transpile a TypeScript File

```bash
./ts2go --in example.ts --out example.go
```

### Example

**Input (example.ts):**
```typescript
interface Person {
  name: string;
  age: number;
}

function greet(person: Person): string {
  return "Hello, " + person.name;
}

console.log(greet({ name: "Alice", age: 30 }));
```

**Output (example.go):**
```go
package main

import "fmt"

type Person struct {
	Name string `json:"name"`
	Age float64 `json:"age"`
}

func Greet(person Person) string {
	return "Hello, " + person.Name
}

func main() {
	fmt.Println(Greet(Person{Name: "Alice", Age: 30}))
}
```

## Supported Features

See [SPEC.md](../SPEC.md) for a complete list of supported TypeScript features.

### Key Mappings

- `interface` → `struct`
- `type` alias → `type` alias
- `string` → `string`
- `number` → `float64`
- `boolean` → `bool`
- `Array<T>` → `[]T`
- `console.log()` → `fmt.Println()`

## Limitations

- No `any` or `unknown` types
- No union types (yet)
- No `class` support (use `interface` instead)
- No `async/await` or `Promise`
- No npm dependencies transpilation
- Error handling uses Go's return pattern instead of `try/catch`

## Architecture

TS2Go consists of three main components:

1. **Parser** (`internal/transpiler/parser/`): Node.js script that uses TypeScript's compiler API to generate an AST
2. **Transpiler** (`internal/transpiler/`): Go code that walks the AST and generates Go code
3. **Runtime** (`runtime/`): Go implementations of Node.js APIs (fs, console, path)

## Running Tests

```bash
make test
```

## Contributing

See the [implementation plan](../implementationPlan.md) for the roadmap and architecture details.
