---
id: TASK-009
title: 'Add SaaS backend unit tests (0% coverage)'
status: todo
priority: critical
dependencies: ['none']
estimated_effort: XL
---

## Goal

Write a comprehensive unit test suite for the SaaS backend (`saas/backend/`) which currently has **zero test files** across all 19 packages, covering authentication, transpilation handlers, storage handlers, and background jobs.

## Context

`find saas/backend -name "*_test.go"` returns empty. The backend uses PostgreSQL and Redis — tests should use `sqlmock` and `redismock` to avoid external dependencies.

Key packages to test:

- `saas/backend/auth/` — JWT token creation/validation
- `saas/backend/transpilation/handlers.go` — TASK-001 fixed ownership checks
- `saas/backend/storage/handlers.go` — TASK-001 fixed project access checks
- `saas/backend/background/jobs.go` — TASK-002 fixed Redis cleanup
- `saas/backend/middleware/` — auth, rate limiting, metrics
- `saas/backend/projects/repository.go` — `CheckAccess` logic

Dependencies to add to `saas/backend/go.mod`:

- `github.com/DATA-DOG/go-sqlmock v1.x`
- `github.com/go-redis/redismock/v9`

## Scope

### Files to create

- `saas/backend/auth/auth_test.go`
- `saas/backend/transpilation/handlers_test.go`
- `saas/backend/storage/handlers_test.go`
- `saas/backend/background/jobs_test.go`
- `saas/backend/middleware/auth_test.go`
- `saas/backend/projects/repository_test.go`

## Implementation

1. Add test dependencies: `go get github.com/DATA-DOG/go-sqlmock@v1` and `go get github.com/go-redis/redismock/v9`
2. **auth_test.go**: Test `GenerateToken`, `ValidateToken`, and token expiry
3. **handlers_test.go** (transpilation): Use `httptest.NewRecorder()`. Mock DB with sqlmock. Test that `GetJobStatus` returns 403 when user != job owner (the TASK-001 fix). Test `CancelJob` similarly.
4. **handlers_test.go** (storage): Test `DownloadFile` and `GetPresignedURL` return 403 when user doesn't own project
5. **jobs_test.go**: Test `CleanupOldJobsJob` calls both DB DELETE with RETURNING and Redis DEL (the TASK-002 fix) using redismock
6. **auth_test.go** (middleware): Test that `RequireAuth` returns 401 on missing/invalid token
7. **repository_test.go**: Test `CheckAccess` returns true for owner, true for team member, false for unauthorized user

## Acceptance Criteria

- [ ] `go test ./... ` in `saas/backend` exits 0
- [ ] Transpilation handler ownership check test passes
- [ ] Storage handler project access check test passes
- [ ] Background jobs Redis cleanup test passes
- [ ] Auth middleware 401 test passes
- [ ] Coverage ≥ 60% for tested files
