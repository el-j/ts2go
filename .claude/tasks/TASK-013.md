---
id: TASK-013
title: 'Add web adapter unit tests and API route tests'
status: todo
priority: high
dependencies: ['TASK-005']
estimated_effort: M
---

## Goal

Write unit tests for `go/adapters/driving/web/` (the HTTP/web server adapter) which currently has no test files.

## Context

`go/adapters/driving/web/` contains the web server adapter for serving the transpiler as an HTTP API. Currently no tests exist. Use `net/http/httptest` for HTTP handler testing.

Look inside the web adapter directory for handler functions that need coverage. Identify all routes and test:

- Healthy GET endpoints (200 + correct JSON)
- Error cases (400 for bad input, 404 for missing resources)
- Auth middleware integration

## Scope

### Files to create

- `go/adapters/driving/web/handlers_test.go` — HTTP handler tests using `httptest`
- `go/adapters/driving/web/routes_test.go` — route registration smoke tests

## Implementation

1. Read `go/adapters/driving/web/*.go` to identify exact handler functions and router setup
2. For each handler, create test using `httptest.NewServer()` or `httptest.NewRecorder()`
3. Test valid request → 200 + correct JSON structure
4. Test invalid/missing parameters → 400/422
5. Test unauthenticated routes if auth middleware is configured
6. For the transpile endpoint, provide a minimal TS snippet and assert response contains Go code

## Acceptance Criteria

- [ ] `go test ./adapters/driving/web/... -v` exits 0
- [ ] Coverage ≥ 60% for web adapter handlers
- [ ] At least one test per HTTP handler function
- [ ] No `any` types introduced
