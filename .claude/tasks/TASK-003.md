---
id: TASK-003
title: 'Implement CLI services and JS object property exclusion'
status: todo
priority: medium
dependencies: ['none']
estimated_effort: M
---

## Goal

Resolve Phase 1 continuation stubs in Go CLI and implement basic JS property exclusions during transpiling.

## Context

We have Phase 1 continuation logic inside `go/adapters/driving/cli/application.go` and property exclusion stub inside `go/internal/transpiler/codegen_statements.go`.

## Scope

### Files to modify

- `go/adapters/driving/cli/application.go`
- `go/internal/transpiler/codegen_statements.go`

## Implementation

1. Fulfill uninitialized services or remove the stub if services are already adequately hydrated inside `application.go`.
2. Introduce simple filtering logic inside `codegen_statements.go` to support property exclusions.

## Acceptance Criteria

- [ ] `make build-go` exits 0
- [ ] `make test-go` exits 0
- [ ] TS2Go exclusions work correctly and appropriately
- [ ] No `any` types introduced
