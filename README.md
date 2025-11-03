# TS2Go

A TypeScript-to-Go transpiler that converts a subset of TypeScript into idiomatic, efficient Go code.

## Quick Start

```bash
# Install dependencies
cd internal/transpiler/parser && npm install && cd ../../..

# Build
make build

# Use CLI
./ts2go transpile <project-dir> --out output

# Or use the Desktop UI (NEW!)
./ts2go ui --port 8080 --open
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

### Getting Started
- 🚀 [Getting Started](docs/GETTING_STARTED.md) - Installation and basic usage
- 🖥️ **[Phase 21: Desktop UI](docs/PHASE21_DESKTOP_UI.md)** - **NEW! Web UI documentation**
- 📚 [Migration Guide](docs/MIGRATION_GUIDE.md) - Complete guide to migrating projects
- 💡 [Examples](docs/EXAMPLES.md) - Before/after transpilation examples

### Reference
- 📋 [Specification](SPEC.md) - Currently supported TypeScript features
- 📦 [Package Mappings](docs/PACKAGE_MAPPINGS.md) - npm to Go quick reference (100+ packages)
- 📊 [Status](docs/STATUS.md) - Current implementation status

### Advanced
- 🗺️ **[Roadmap](docs/ROADMAP.md)** - **Complete plan for all TypeScript features + dependencies**
- 🔧 **[Dependency Guide](docs/DEPENDENCY_GUIDE.md)** - **Implementation guide for dependency resolution**
- 🏗️ [Architecture](docs/ARCHITECTURE.md) - Project structure and design
- 📝 [Implementation Plan](implementationPlan.md) - Original development phases

## Features

### Core Transpilation
✅ Interface to struct transpilation  
✅ Type aliases  
✅ Function declarations with typed parameters  
✅ Basic types (string, number, boolean, arrays)  
✅ Classes with inheritance, static methods, getters/setters  
✅ Union types and enums  
✅ Console.log mapping  
✅ Go runtime library for Node.js APIs  

### 🆕 Desktop UI (Phase 21)
✅ **Web-based UI** for easy transpilation  
✅ **Split-pane editor** with TypeScript input and Go output  
✅ **Real-time transpilation** with instant feedback  
✅ **Built-in examples** showcasing key features  
✅ **Responsive design** works on all devices  

### In Progress
🔄 Control flow statements (if/else, loops, switch)  
🔄 Modern JavaScript syntax (arrow functions, template literals)  
🔄 Async/await support  

❌ Frontend frameworks (React, Vue, Angular) - Out of scope  

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