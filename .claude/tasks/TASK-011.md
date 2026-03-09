---
id: TASK-011
title: 'Create E2E transpilation test pipeline'
status: todo
priority: high
dependencies: ['TASK-005', 'TASK-006', 'TASK-007']
estimated_effort: L
---

## Goal

Build a proper E2E test file (`go/tests/e2e_pipeline_test.go`) that transpiles 10+ TypeScript fixture files end-to-end, compiles the generated Go output with `go build`, and asserts output is valid runnable Go code.

## Context

`go/internal/tests/e2e_test.go` exists but only tests 2 fixtures and doesn't build-check most of them. The fixtures in `go/tests/fixtures/*.ts` are the true test assets. The goal is a pipeline that:

1. Takes each `.ts` fixture file
2. Transpiles via `transpiler.Transpile()`
3. Runs `go build` on the output
4. Optionally runs the output if a `.expected` file exists

The E2E suite tests the FULL stack vs. the unit tests which only test individual functions.

## Scope

### Files to create

- `go/tests/e2e_pipeline_test.go` — table-driven tests for all fixtures

### Files to modify

- `go/tests/fixtures/simple.ts` et al — ensure each fixture has a valid `.ts` source

## Implementation

1. Create `go/tests/e2e_pipeline_test.go` with a `TestE2EPipeline` table test
2. Discover all `.ts` files in `go/tests/fixtures/` via `os.ReadDir`
3. Skip files in `debugOnly` list (e.g., `debug-forof.ts`, `debug-while.ts`)
4. For each `.ts` file:
   - Call `transpiler.Transpile(tsPath, goOutPath)`
   - Parse generated Go with `go/parser.ParseFile()` to verify AST is valid
   - Try `exec.Command("go", "build", "-o", "/dev/null", goOutPath)` for main-package files
5. Run with short timeout per test (30s max): `t.Skip` if Node.js unavailable
6. Track how many fixture files PASS build — should be ≥ 80%

## Acceptance Criteria

- [ ] `go test ./tests/... -run TestE2EPipeline -v` exits 0
- [ ] ≥ 10 fixture files tested end-to-end
- [ ] ≥ 80% pass the `go build` check
- [ ] Test reports which fixtures fail and why (using `t.Logf`)
