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

### Project Status & Planning
- 📊 **[Current State](CURRENT_STATE.md)** - **Actual implementation status (from source code)**
- 🎯 **[Roadmap to Alpha 2.0.1](ROADMAP_TO_ALPHA_2.0.1.md)** - **Plan to add build & test features**
- ⚠️ [Known Issues](KNOWN_ISSUES.md) - Current limitations and workarounds

### Getting Started
- 🚀 [Getting Started](docs/GETTING_STARTED_v2.md) - Installation and basic usage
- 🖥️ [Desktop UI Guide](desktop-ui/USER_GUIDE.md) - Desktop application user guide
- 📚 [Migration Guide](docs/MIGRATION_GUIDE.md) - Complete guide to migrating projects
- 💡 [Examples](docs/EXAMPLES.md) - Before/after transpilation examples

### Reference
- 📋 [Specification](SPEC.md) - Currently supported TypeScript features
- 📦 [Package Mappings](docs/PACKAGE_MAPPINGS.md) - npm to Go quick reference (100+ packages)
- 🏗️ [Architecture](docs/ARCHITECTURE.md) - Project structure and design
- 🔧 [Dependency Guide](docs/DEPENDENCY_GUIDE.md) - Implementation guide for dependency resolution

## Features

### ✅ Desktop UI (85-90% Complete)
✅ **Full project workspace** with file tree, multi-file tabs, Monaco editor  
✅ **Real-time transpilation** with progress tracking  
✅ **Go code execution** - Run single files or complete projects  
✅ **Settings & persistence** - Customizable theme, editor, and project settings  
✅ **Recent projects** - Quick access to recently opened projects  
✅ **Example gallery** - 6 built-in examples to learn from  
✅ **History tracking** - View past transpilations  
✅ **Keyboard shortcuts** - Efficient workflow with hotkeys  

### ✅ Core Transpilation (80% Complete)
✅ Interface to struct transpilation  
✅ Type aliases, enums, unions, tuples  
✅ Function declarations with typed parameters  
✅ Basic & advanced types (string, number, boolean, arrays, generics)  
✅ Classes with inheritance, static methods, getters/setters  
✅ Console.log and Node.js API mappings  
✅ **Control flow** - if/else, for/while loops, switch statements  
✅ **Error handling** - try/catch/finally, throw statements  
✅ **Async/await** - Goroutine-based async with channels  
✅ **Modern JS** - Arrow functions, template literals, destructuring, spread operators  

### ✅ CLI Tools (70% Complete)
✅ `convert` - Single file transpilation  
✅ `transpile` - Full project transpilation  
✅ `analyze` - Dependency analysis  
✅ `watch` - Auto-transpile on file changes  
✅ `ui` - Legacy web UI server  

### 🔄 Alpha 2.0.1 Features (In Progress → Complete!)
✅ **Go build integration** - Compile transpiled code into binaries (CLI + GUI)
✅ **Go test integration** - Run tests on transpiled code (CLI + GUI)
✅ **Build artifacts** - Track and manage compiled binaries
✅ **Test results UI** - Display test pass/fail in GUI with detailed output

**New Commands:**
- `ts2go build --source <dir> --output <binary>` - Compile Go code
- `ts2go test --source <dir> --verbose --coverage` - Run Go tests

**GUI Features:**
- 🔨 Build button in ProjectView
- 🧪 Test button in ProjectView
- Build/test results display in log panel
- Artifact tracking with localStorage persistence

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