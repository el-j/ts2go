# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **NEXT_STEPS.md** - Focused roadmap for immediate priorities (next 2-4 weeks)
- Comprehensive GitHub Actions workflows for CI/CD
  - Feature branch workflow with cross-platform builds
  - Develop branch workflow with integration tests
  - Main branch workflow with quality checks
  - Release workflow with multi-platform binary compilation
  - Desktop UI build workflow with Tauri support
- CLI entry point (`cmd/ts2go/main.go`)
- Docker support with multi-stage builds
- Next Phase Roadmap documentation
- CI badges in README
- Installation instructions for multiple methods
- **ROADMAP_TO_1.0.0.md** - Comprehensive 12-week roadmap to stable release
- **KNOWN_ISSUES.md** - Documentation of current limitations and workarounds

### Changed
- **Documentation cleanup** - Removed 38 outdated documentation files
  - Removed old phase completion documents (PHASE9-21)
  - Removed old build fix documents
  - Removed superseded roadmap and status documents
  - Consolidated to essential, current documentation only
- Updated README.md with corrected documentation links
- Updated PROJECT_STATUS.md with current dates and next steps
- Fixed `go.work` file to include `cmd/ts2go` module
- Updated README with better documentation structure

### Removed
- 38 outdated markdown documentation files from `/docs` folder
- Duplicate SPEC.md from `/docs` (kept root version)
- Old roadmap and status documents superseded by current versions

### Fixed
- Build system now properly compiles the main CLI binary
- **Go syntax errors in transpiled example files**
  - Fixed struct literals missing type names in `processor.go`
  - Fixed object literals to use `map[string]interface{}` syntax in `server.go`
  - Fixed indentation in `codegen_helpers.go`
  - All files now pass `gofmt -s -l .` check
- **Security: Desktop UI dependency vulnerabilities**
  - Updated happy-dom to v20.0.10+ (fixed critical VM escape vulnerability)
  - Resolved 7 vulnerabilities (1 critical, 6 moderate)

## [0.1.0] - 2025-11-03

### Added
- Initial release
- TypeScript to Go transpilation for types and interfaces
- Class transpilation with full OOP support
- Dependency management and npm package mapping
- Desktop UI (Phase 21 - 52% complete)
- Runtime libraries for Node.js APIs (fs, path, console, etc.)
- CLI with transpile, analyze, and ui commands
- Comprehensive documentation

### Supported Features
- Interface → Struct conversion
- Type aliases, enums, unions, tuples
- Classes with inheritance, static methods, getters/setters
- Function declarations with typed parameters
- Basic expressions and variable declarations
- npm package mapping to Go equivalents

### Known Limitations
- No control flow statements (if/else, loops, switch)
- No modern JavaScript syntax (arrow functions, template literals)
- No async/await support
- No try/catch error handling
- Frontend frameworks not supported (out of scope)

### Coverage
- Approximately 15-20% of typical backend TypeScript codebases

[Unreleased]: https://github.com/el-j/ts2go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/el-j/ts2go/releases/tag/v0.1.0
