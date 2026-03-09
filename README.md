# TS2Go

[![CI - Main](https://github.com/el-j/ts2go/actions/workflows/ci-main.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/ci-main.yml)
[![CI - Develop](https://github.com/el-j/ts2go/actions/workflows/ci-develop.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/ci-develop.yml)
[![Release](https://github.com/el-j/ts2go/actions/workflows/release.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/el-j/ts2go)](https://goreportcard.com/report/github.com/el-j/ts2go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A TypeScript-to-Go transpiler that converts a subset of TypeScript into idiomatic, efficient Go code.

## Quick Start

```bash
# Build the CLI
make build-go

# Run CLI
./bin/ts2go transpile <project-dir> --out output

# Or use Makefile targets
make help           # See all available commands
make build          # Build everything
make test           # Run all tests
make dev-desktop    # Run desktop app in dev mode
```

## Project Structure

```
ts2go/
├── go/              # 🔵 Go CLI & Core Engine
│   ├── core/        # Hexagonal architecture (domain, ports, services)
│   ├── adapters/    # Infrastructure adapters
│   ├── cmd/         # CLI binaries
│   ├── runtime/     # Go runtime library
│   └── go.mod       # Single Go module
│
├── desktop/         # 🟢 Desktop Application
│   ├── ui/          # Vue.js + TypeScript frontend
│   └── tauri/       # Rust + Tauri backend
│
├── docs/            # 📚 Documentation
├── bin/             # 🔨 Compiled binaries
└── Makefile         # 🎯 Main build orchestration
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

**Main Documentation:** See [docs/README.md](docs/README.md) for complete documentation index.

### Quick Links
- 📊 [Current State](docs/CURRENT_STATE.md) - Implementation status
- 🚀 [Getting Started](docs/GETTING_STARTED_v2.md) - Installation and usage
- 📋 [Specification](docs/SPEC.md) - Supported TypeScript features
- 🖥️ [Desktop Guide](desktop/ui/USER_GUIDE.md) - Desktop app user guide
- 📦 [Package Mappings](docs/PACKAGE_MAPPINGS.md) - NPM to Go mappings
- 🏗️ [Architecture](docs/ARCHITECTURE.md) - System design
- ⚠️ [Known Issues](docs/KNOWN_ISSUES.md) - Limitations and workarounds

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