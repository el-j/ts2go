---
id: TASK-007
title: 'Add unit tests for CLI driving adapter'
status: done
priority: high
dependencies: ['TASK-005']
estimated_effort: M
---

## Goal

Write unit tests for `go/adapters/driving/cli/` which currently has no test files, ensuring the application wiring, command handlers, and service integrations are verified.

## Context

`go/adapters/driving/cli/application.go` is the main DI container wiring all services together. It was recently updated (TASK-003) to use real `adapterAnalyzer`, `adapterCodeGen`, and `adapterMapper` implementations. The CLI commands live in `go/cli/` as cobra commands connected through the application object.

Key entry points to test:

- `NewApplication()` — constructs all services correctly
- `adapterAnalyzer.AnalyzeImports()` — calls through to internal/analyzer
- `adapterCodeGen.ParseTypeScript()` / `GenerateGoCode()` — invokes transpiler
- CLI `convert` command end-to-end (subprocess invocation)

## Scope

### Files to create

- `go/adapters/driving/cli/application_test.go` — unit tests for NewApplication and adapters
- `go/adapters/driving/cli/command_test.go` — cobra command invocation tests

## Implementation

1. Test that `NewApplication()` returns non-nil service (smoke test)
2. Test `adapterAnalyzer.AnalyzeImports()` with a small inline TypeScript snippet written to temp file
3. Test `adapterCodeGen.ParseTypeScript()` with valid TS written to temp file (requires Node.js parser — add `t.Skip` if unavailable)
4. Test `adapterCodeGen.GenerateGoCode()` with a synthetic `ASTNode` tree
5. For `command_test.go`, spawn `go run ./cmd/ts2go convert --in fixtures/simple.ts --out /tmp/test.go` as subprocess and assert exit 0

## Acceptance Criteria

- [x] `go test ./adapters/driving/cli/... -v` exits 0
- [x] `NewApplication()` smoke test passes
- [x] At least 5 unit tests covering distinct code paths
- [x] No `any` types introduced

## Execution Log

- **Completed:** 2026-03-07T22:50:00.000Z
- **Files changed:** go/adapters/driving/cli/application_test.go, go/adapters/driving/cli/command_test.go, go/internal/transpiler/transpiler.go
- **Build:** ✅ clean
- **Tests:** ✅ passed
- **Notes:** Added tests for application.go dependencies wiring, AST processing adapters, and end-to-end convert command execution. Had to update transpiler.go to load parser.js from `../../../` during `go test`.
