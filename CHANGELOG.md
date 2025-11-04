# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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

### Changed
- Fixed `go.work` file to include `cmd/ts2go` module
- Updated README with better documentation structure

### Fixed
- Build system now properly compiles the main CLI binary

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
