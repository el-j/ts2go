---
id: TASK-014
title: 'Set up CI/CD GitHub Actions pipeline'
status: todo
priority: high
dependencies: ['TASK-005', 'TASK-008', 'TASK-012']
estimated_effort: M
---

## Goal

Create a production-grade GitHub Actions CI pipeline that runs all tests on every push and PR, builds all binaries, and blocks merging if any step fails.

## Context

There is no `.github/workflows/` directory. The project needs:

1. Go tests (`make test-go`)
2. Desktop/UI tests (`make test-desktop`)
3. Go build (`make build-go`)
4. TypeScript type-check (`npm run type-check`)
5. Go lint (`go vet ./...`)

Matrix: ubuntu-latest, goes-1.22

## Scope

### Files to create

- `.github/workflows/ci.yml` — main CI workflow

## Implementation

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Install Node.js (for TS parser)
        uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install parser deps
        run: cd go/internal/transpiler/parser && npm ci
      - name: Go vet
        run: cd go && go vet ./...
      - name: Go tests
        run: make test-go
      - name: Go build
        run: make build-go

  test-desktop:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install deps
        run: cd packages/ui-shared && npm ci
      - name: TypeScript check
        run: cd packages/ui-shared && npm run type-check
      - name: Vitest
        run: cd packages/ui-shared && npm run test
```

The `saas/backend` tests should be a separate job requiring Postgres + Redis in services.

## Acceptance Criteria

- [ ] `.github/workflows/ci.yml` created and syntactically valid
- [ ] CI passes on `main` branch after all fixes from TASK-005 through TASK-013 are merged
- [ ] PR checks block merge on test failure
- [ ] Go and Node.js are both properly set up in the workflow
