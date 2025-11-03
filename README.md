# TS2Go

[![CI - Main](https://github.com/el-j/ts2go/actions/workflows/ci-main.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/ci-main.yml)
[![CI - Develop](https://github.com/el-j/ts2go/actions/workflows/ci-develop.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/ci-develop.yml)
[![Release](https://github.com/el-j/ts2go/actions/workflows/release.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/el-j/ts2go)](https://goreportcard.com/report/github.com/el-j/ts2go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A TypeScript-to-Go transpiler that converts a subset of TypeScript into idiomatic, efficient Go code.

## Quick Start

```bash
# Install dependencies
cd internal/transpiler/parser && npm install && cd ../../..

# Build
make build

# Use CLI
./ts2go convert --in app.ts --out app.go      # Single file conversion
./ts2go transpile <project-dir> --out output  # Full project transpilation
./ts2go analyze <project-dir>                  # Analyze dependencies
./ts2go ui --port 8080 --open                 # Web UI (legacy, for quick testing)
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
- 🎯 **[Next Phase Roadmap](docs/NEXT_PHASE_ROADMAP.md)** - **Detailed plan for upcoming implementation phases**
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

### ✅ Advanced Features (Phases 16-20 Complete)
✅ **Control flow** - if/else, for/while loops, switch statements  
✅ **Error handling** - try/catch/finally, throw statements  
✅ **Async/await** - Goroutine-based async with channels  
✅ **Advanced operators** - typeof, instanceof, in, delete  

### 🆕 Desktop UI (Phase 21 - 52% Complete)
✅ **Tauri + Vue 3 + Monaco Editor** - Professional desktop app  
✅ **Split-pane editor** with TypeScript input and Go output  
✅ **Real-time transpilation** with instant feedback  
✅ **Resizable panes and LogViewer** component  
✅ **39 tests passing** with 100% store coverage  

### In Progress
🔄 Modern JavaScript syntax (arrow functions, template literals, destructuring)  
🔄 Desktop UI completion (Weeks 3-4 remaining)  

❌ Frontend frameworks (React, Vue, Angular) - Out of scope  

## Installation

### From Source
```bash
git clone https://github.com/el-j/ts2go.git
cd ts2go
make build
sudo mv ts2go /usr/local/bin/
```

### From Release (Coming Soon)
Download pre-compiled binaries from the [releases page](https://github.com/el-j/ts2go/releases).

### Using Docker
```bash
docker pull ghcr.io/el-j/ts2go:latest
docker run -v $(pwd):/workspace ghcr.io/el-j/ts2go transpile /workspace/my-project
```

## Testing

```bash
make test
```

## CI/CD & Releases

This project uses GitHub Actions for continuous integration and automated releases:

- **Feature Branches:** Automated builds and tests on every push
- **Develop Branch:** Full test suite with integration tests  
- **Main Branch:** Production-ready builds with quality checks
- **Tagged Releases:** Automatic binary compilation for Linux, macOS, and Windows

See our [workflows](.github/workflows/) for details.

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