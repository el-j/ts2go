---
id: TASK-006
title: 'Add unit tests for adapters/driven (0% coverage)'
status: todo
priority: critical
dependencies: ['TASK-005']
estimated_effort: L
---

## Goal

Write unit tests for all three `go/adapters/driven/` packages which currently have 0% coverage, covering the filesystem adapter, Go compiler adapter, and persistence (JSON repository) adapter.

## Context

Dead code analysis and coverage reports show all functions in:

- `go/adapters/driven/filesystem/filesystem.go` → 0% coverage
- `go/adapters/driven/gocompiler/compiler.go` → 0% coverage (functions: `Build`, `Run`, `Test`, `Format`, `ModInit`, `ModTidy`, `Get`)
- `go/adapters/driven/persistence/json_repository.go` → 0% coverage (functions: `GetState`, `SaveState`, `GetAllStates`, `DeleteState`, etc.)

The existing tests in `go/core/services/transpilation_service_test.go` mock these adapters but the concrete implementations are never tested directly. Use `os.MkdirTemp` for isolation.

## Scope

### Files to create

- `go/adapters/driven/filesystem/filesystem_test.go`
- `go/adapters/driven/gocompiler/compiler_test.go`
- `go/adapters/driven/persistence/json_repository_test.go`

## Implementation

1. **filesystem_test.go**: Test `WriteFile`, `ReadFile`, `Exists`, `DeleteFile`, `ListFiles`, `CreateDir` using `os.MkdirTemp` for sandbox isolation
2. **compiler_test.go**: Use `testdata/` fixtures with minimal Go programs. Test `Build` with a valid program (expect success) and invalid program (expect error). Skip if `go` binary not in PATH.
3. **json_repository_test.go**: Test full CRUD lifecycle — `SaveState` → `GetState` → `GetAllStates` → `DeleteState`. The JSON file should be created/read/deleted in a temp dir.
4. Target **≥ 80% coverage** per file
5. Run `go test ./adapters/driven/... -v -cover`

## Acceptance Criteria

- [ ] `go test ./adapters/driven/... -cover` exits 0
- [ ] Coverage ≥ 80% for each file
- [ ] No `any` types in tests
- [ ] Tests are isolated (no disk state leaks between runs)
