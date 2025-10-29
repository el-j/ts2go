# TS2Go

A TypeScript-to-Go transpiler that converts a subset of TypeScript into idiomatic, efficient Go code.

## Quick Start

```bash
# Install dependencies
cd internal/transpiler/parser && npm install && cd ../../..

# Build
make build

# Use
./ts2go --in input.ts --out output.go
```

## Example

**TypeScript:**
```typescript
interface Person {
  name: string;
  age: number;
}

function greet(person: Person): string {
  return "Hello, " + person.name;
}
```

**Generated Go:**
```go
type Person struct {
	Name string `json:"name"`
	Age float64 `json:"age"`
}

func Greet(person Person) string {
	return "Hello, " + person.Name
}
```

## Documentation

- [Getting Started](docs/GETTING_STARTED.md) - Installation and basic usage
- [Architecture](docs/ARCHITECTURE.md) - Project structure and design
- [Specification](SPEC.md) - Supported TypeScript features
- [Implementation Plan](implementationPlan.md) - Roadmap and development phases

## Features

✅ Interface to struct transpilation  
✅ Type aliases  
✅ Function declarations with typed parameters  
✅ Basic types (string, number, boolean, arrays)  
✅ Console.log mapping  
✅ Go runtime library for Node.js APIs  

❌ Classes (coming soon)  
❌ Async/await  
❌ Union types  
❌ npm dependencies  

## Testing

```bash
make test
```

## Project Structure

```
ts2go/
├── cmd/ts2go/              # CLI tool
├── internal/transpiler/    # Core transpilation engine
│   └── parser/            # TypeScript AST parser (Node.js)
├── runtime/               # Go implementations of Node.js APIs
├── tests/                 # Integration tests
└── docs/                  # Documentation
```

## License

MIT