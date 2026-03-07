---
id: TASK-005
title: 'Fix broken tests/fixtures generated Go files'
status: todo
priority: critical
dependencies: ['none']
estimated_effort: S
---

## Goal

Remove or fix the broken generated `.go` files in `go/tests/fixtures/` that cause `go build ./...` and coverage to fail, blocking the entire CI pipeline.

## Context

`go/tests/fixtures/integration-simple.go` has `if /* unsupported expression */ {` at line 18 — invalid Go syntax. Several other fixtures also have compile-time errors including `ternary-test.go`, `template-test.go`, `optional-chaining.go`. These are **generated outputs** (not source), so the correct fix is to either regenerate them from their `.ts` counterparts using the current transpiler or move them to a dedicated `testdata/` directory with a build tag that excludes them from `go build ./...`.

Exact broken files found by `go build ./...`:

- `tests/fixtures/integration-simple.go` (line 18: missing condition)
- `tests/fixtures/ternary-test.go` (type assertion errors)
- `tests/fixtures/template-test.go` (`fmt` undefined, unused vars)

## Scope

### Files to modify

- `go/tests/fixtures/*.go` — Add `//go:build ignore` tag to all fixture generated files OR regenerate them

### Files to delete

- `go/internal/transpiler/debug_test.go` — Temporary debug test leftover from auditing session, must be removed before shipping

## Implementation

1. Add `//go:build ignore` build tag as the first line of every `.go` file under `go/tests/fixtures/` that is a generated output (not a manually-written expected file)
2. Delete `go/internal/transpiler/debug_test.go`
3. Run `go build ./...` — must exit 0
4. Run `make test-go` — must exit 0

## Acceptance Criteria

- [ ] `go build ./...` exits 0
- [ ] `make test-go` exits 0
- [ ] `debug_test.go` is deleted
- [ ] All fixture `.go` files have `//go:build ignore` tag
