---
id: TASK-015
title: 'Harden SaaS backend security and rate limiting'
status: todo
priority: high
dependencies: ['TASK-009']
estimated_effort: M
---

## Goal

Ensure all SaaS backend security mechanisms are fully operational: `RequireAPIKey`, `OptionalAuth` middleware, rate limiting via `LimitByEndpoint`, and proper input validation on all transpilation requests.

## Context

From `deadcode` output, these functions exist but are never called from the main API router:

- `middleware/auth.go:62` — `Middleware.RequireAPIKey`
- `middleware/auth.go:99` — `Middleware.OptionalAuth`
- `middleware/ratelimit.go:70` — `RateLimiter.LimitByEndpoint`
- `middleware/ratelimit.go:151` — `StrictRateLimitConfig`
- `middleware/ratelimit.go:159` — `GenerousRateLimitConfig`

Check `saas/backend/api/main.go` route registrations. If these middleware functions are defined but not wired to routes, they need to be applied.

Additionally, the `/transpile` endpoint accepts arbitrary code — needs size limits (max 1MB) and content-type validation.

## Scope

### Files to modify

- `saas/backend/api/main.go` — Apply `RateLimiter.LimitByEndpoint` to transpile and upload routes
- `saas/backend/transpilation/handlers.go` — Add request size limit middleware
- `saas/backend/middleware/ratelimit.go` — Verify `StrictRateLimitConfig` and `GenerousRateLimitConfig` are used in `main.go`

## Implementation

1. In `saas/backend/api/main.go`, locate the transpile route group and apply:
   - `rateLimiter.LimitByEndpoint()` (strict config — 10 req/min per user)
   - `middleware.RequireAuth()` (already likely applied but verify)
2. Add a Gin middleware to the transpile route that rejects requests with `Content-Length > 1MB`
3. Verify `GetUserID`, `GetEmail`, `GetUsername` helper functions are used by handlers (they're currently unreachable — either use them or remove them)
4. Test with `TASK-009` — run `handlers_test.go` to verify 429 rate limit response

## Acceptance Criteria

- [ ] `go vet ./...` exits 0 in `saas/backend`
- [ ] Rate limiting middleware is wired to at least the `/transpile` endpoint
- [ ] Request size limit of 1MB enforced on file upload endpoints
- [ ] Either `GetUserID`/`GetEmail`/`GetUsername` are used in handlers OR removed
- [ ] `go test ./...` still passes
