---
id: TASK-016
title: 'Enable golangci-lint and fix all lint warnings'
status: todo
priority: medium
dependencies: ['TASK-005']
estimated_effort: S
---

## Goal

Configure `golangci-lint` across the monorepo (both `go/` and `saas/backend/`) and fix all reported warnings, including the "len() for nil slices" issues and unused imports.

## Context

`go vet ./...` in `go/` currently passes (0 issues). However, several linting warnings were noted during the audit:

- `codegen_expressions.go`, `codegen_statements.go`, `codegen_classes.go`, `codegen_types.go` have "should omit nil check; len() for nil slices is defined as zero" warnings
- Some unused imports in test files

`golangci-lint` provides a superset including: `nilnil`, `exhaustive`, `godot`, `errcheck`

## Scope

### Files to create

- `go/.golangci.yml` — golangci-lint configuration
- `saas/backend/.golangci.yml` — golangci-lint configuration

### Files to modify

- `go/internal/transpiler/codegen_expressions.go` — fix nil checks before len()
- `go/internal/transpiler/codegen_statements.go` — fix nil checks before len()
- `go/internal/transpiler/codegen_classes.go` — fix nil checks before len()
- `go/internal/transpiler/codegen_types.go` — fix nil checks before len()
- `Makefile` — add `lint` target

## Implementation

1. Install `golangci-lint` via homebrew or `go install`
2. Create `.golangci.yml` with key linters: `govet`, `errcheck`, `staticcheck`, `unused`, `gosimple`, `ineffassign`
3. Fix nil-before-len patterns: replace `if x != nil && len(x) > 0` with `if len(x) > 0`
4. Add `make lint` target: `golangci-lint run ./...`
5. CI TASK-014 should add `golangci-lint` as a CI step

## Acceptance Criteria

- [ ] `golangci-lint run ./...` exits 0 in both `go/` and `saas/backend/`
- [ ] No nil-before-len patterns remain
- [ ] `Makefile` has `lint`, `lint-go`, `lint-backend` targets
