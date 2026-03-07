---
id: TASK-001
title: 'Fix backend authentication and ownership TODOs'
status: done
priority: high
dependencies: ['none']
estimated_effort: S
---

## Goal

Ensure API handlers correctly verify user ownership and project access before mutating data.

## Context

Inside `saas/backend/transpilation/handlers.go` and `saas/backend/storage/handlers.go`, endpoints currently omit critical access checks (marked by TODO comments). They need to verify the user owns the transpilation job and the target project respectfully.

## Scope

### Files to modify

- `saas/backend/transpilation/handlers.go` — Add job ownership checks.
- `saas/backend/storage/handlers.go` — Add project access checks.

## Implementation

1. Retrieve the authenticated user ID from the request context or middleware.
2. Ensure queries and permission checks restrict access only to resources the user owns or is granted access to.
3. Return `403 Forbidden` if unauthorized access is attempted.

## Acceptance Criteria

- [x] `npm run build` exits 0 with no new errors
- [x] `make test` exits 0
- [x] Unauthorized users cannot mutate jobs or projects they don't own
- [x] No `any` types introduced

## Execution Log

- **Completed:** 2026-03-07T13:55:00.000Z
- **Files changed:** saas/backend/transpilation/handlers.go, saas/backend/storage/handlers.go, saas/backend/api/main.go
- **Build:** ✅ clean
- **Tests:** ✅ passed
- **Notes:** Added Job UserID checks and injected projects.Repository into storage.NewHandler to fulfill verify checks.
