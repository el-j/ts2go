# TS2Go Mastermind Roadmap & Tracking Matrix

> **Current Status**: Milestone 1 in progress (CLI Core Transpilation & Pipeline Fixes)  
> **Target Definition of Done**: 100% test coverage, 0 stub code, 0 placeholders, 100% in-code godoc, 100% e2e verified, 0 type issues, 0 lint warnings.

---

## 🗺️ Milestone Overview

```
[Milestone 1: CLI Core Transpilation] ──▶ [Milestone 2: Code Quality & Boss Files]
                 │                                        │
                 ▼                                        ▼
[Milestone 3: 100% Test Coverage & CI] ◀──────────────────┘
                 │
                 ▼
[Milestone 4: Web API & GUI Modernization]
```

| Milestone | Title | Priority | GitHub Milestone | Status | Issues |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **M1** | CLI Core Transpilation & Pipeline Fixes | **P0 (CLI First)** | [Milestone 1](https://github.com/el-j/ts2go/milestone/1) | 🟢 Completed | #14, #15, #16, #17 |
| **M2** | Code Quality, Refactoring & Boss-Files | **P1 (Quality)** | [Milestone 2](https://github.com/el-j/ts2go/milestone/2) | 🟢 Completed | #18, #19, #20, #21 |
| **M3** | 100% Test Coverage, E2E & Mutation | **P1 (Testing)** | [Milestone 3](https://github.com/el-j/ts2go/milestone/3) | 🟢 Completed | #22, #23, #24, #25 |
| **M4** | Web API & GUI Modernization | **P2 (GUI After)** | [Milestone 4](https://github.com/el-j/ts2go/milestone/4) | 🟡 Next Up | #26, #27, #28 |

---

## 📋 Detailed Task Breakdown & Issue Status

### Milestone 1: CLI Core Transpilation & Pipeline Fixes (Priority: P0)
*Goal: Fix all core transpilation bugs, get Go test suite green, eliminate placeholder mappers, wire AST nodes.*

- [x] **[#14: Fix Core E2E Transpilation Test Failures (Object Literals & Case Sensitivity)](https://github.com/el-j/ts2go/issues/14)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `go/internal/transpiler/codegen_expressions.go`, `go/internal/transpiler/codegen_statements.go`, `go/internal/tests/e2e_test.go`
  - **DoD**:
    - [x] Struct instantiation emitted instead of `map[string]interface{}` when return type is a struct.
    - [x] Function call casing matches PascalCased function declarations (`CreateUser` vs `createUser`).
    - [x] Top-level `const` emitted as `const PI = ...` instead of `:=` inside `func main()`.
    - [x] `go test ./internal/tests` passes 100%.

- [x] **[#15: Wire Real NPM-to-Go Mapper into Hexagonal CLI & Fix Integration Tests](https://github.com/el-j/ts2go/issues/15)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `go/adapters/driving/cli/application.go`, `go/mappings/npm-to-go.yaml`, `go/internal/mapper/integration_test.go`
  - **DoD**:
    - [x] Replace `placeholderMapper` in CLI with real `internal/mapper` instance.
    - [x] Fix builtins typing in `npm-to-go.yaml` and `integration_test.go` (`os`, `http`, `url`).
    - [x] Replace `github.com/yourusername/...` and `github.com/ts2go/...` with `github.com/el-j/ts2go/...`.
    - [x] `go test ./internal/mapper` passes 100%.

- [x] **[#16: Connect Unlinked AST Node Generators & Fix Expression Fallthroughs](https://github.com/el-j/ts2go/issues/16)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `go/internal/transpiler/codegen_expressions.go`, `go/internal/transpiler/codegen_statements.go`
  - **DoD**:
    - [x] Wire `generateElementAccess` (`arr[i]`, `obj[k]`) into `generateExpression`.
    - [x] Wire `generateSpreadElement` (`...args`) into `generateExpression`.
    - [x] Handle `ParenthesizedExpression`, `NonNullExpression`, `AsExpression`, `UndefinedKeyword`.
    - [x] Eliminate silent drops; return typed `TranspilationError`.

- [x] **[#17: Multi-File & Project Orchestrator Integration in CLI](https://github.com/el-j/ts2go/issues/17)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `go/core/services/transpilation_service.go`, `go/internal/project/module.go`, `go/internal/orchestrator/`
  - **DoD**:
    - [x] Connect dependency graph scanner into `TranspileProject`.
    - [x] Infer package names from folder hierarchy (no `package main` in library packages).
    - [x] Generate real `go.mod` for transpiled output.

---

### Milestone 2: Code Quality, Refactoring & Boss-Files Breakdown (Priority: P1)
*Goal: Break down monolithic files (>400 lines), eliminate dead trees, 0 linter warnings, 100% godoc.*

- [x] **[#18: Break Down Transpiler Boss Files into Modular Subgenerators](https://github.com/el-j/ts2go/issues/18)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `codegen_expressions.go` (split into literals, calls, operators, expressions), `codegen_statements.go` (split into declarations, loops, control, statements), `codegen_classes.go` (split into constructors, methods, classes)
  - **DoD**:
    - [x] Decompose all files to <300 lines (target was <350 lines).
    - [x] 100% transpiler and E2E tests pass.

- [x] **[#19: Decommission Stale Duplicate Directories & Dead Code](https://github.com/el-j/ts2go/issues/19)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `desktop/ui/` (deleted), `desktop/tauri/` (deleted), `go/cli/` (deleted), port `UICommand` into hexagonal adapter.
  - **DoD**:
    - [x] Single canonical frontend (`packages/ui-shared`) and backend (`desktop/tauri-backend`).
    - [x] Dead legacy CLI directory removed; CLI fully uses hexagonal architecture.

- [x] **[#20: Resolve All 75 golangci-lint Issues & Deprecated Packages](https://github.com/el-j/ts2go/issues/20)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: Entire `go/` codebase, `.golangci.yml`
  - **DoD**:
    - [x] Upgraded `.golangci.yml` to v2 format.
    - [x] Replaced deprecated `io/ioutil` with `io` and `os`.
    - [x] Handled or explicitly acknowledged all error returns (`errcheck`).
    - [x] Resolved all staticcheck warnings (`QF1012`, redundant nil checks).
    - [x] `golangci-lint run ./go/...` passes with 0 errors and 0 warnings.

- [x] **[#21: 100% In-Code Documentation & Godoc Standards](https://github.com/el-j/ts2go/issues/21)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: All Go packages in `go/`, `cmd/ts2go/main.go`, `Makefile`
  - **DoD**:
    - [x] Added `doc.go` to all packages across domain, ports, services, adapters, internal, and runtime.
    - [x] Dynamic version reading from `VERSION` or build-time injection via `-ldflags`.
    - [x] CLI `printUsage()` updated with `state` and `settings` documentation.

---

### Milestone 3: 100% Test Coverage, E2E & Mutation Testing (Priority: P1)
*Goal: Fulfill rigorous testing DoD: 100% test coverage, comprehensive E2E matrix, mutation testing, strict CI.*

- [x] **[#22: Complete Unit Tests for 0% and Low Coverage Packages](https://github.com/el-j/ts2go/issues/22)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `runtime/console`, `runtime/fs`, `runtime/path`, `cmd/ts2go`, `adapters/driving/cli`, `adapters/driving/web`
  - **DoD**:
    - [x] Achieved 100% coverage on `runtime/console` and `runtime/fs`.
    - [x] Achieved >92% coverage on `runtime/path`.
    - [x] Comprehensive tests for CLI app and web server handlers.

- [x] **[#23: Comprehensive End-to-End (E2E) CLI Test Matrix](https://github.com/el-j/ts2go/issues/23)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `go/internal/tests/e2e_matrix_test.go`
  - **DoD**:
    - [x] Automated matrix comparing Go compilation output against Node.js runtime.
    - [x] Verified arithmetic, loops, switch statements, classes/constructors, arrow functions, interfaces, and nullish coalescing.

- [x] **[#24: Set Up Mutation Testing Suite (gremlins / go-mutesting)](https://github.com/el-j/ts2go/issues/24)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: Root `Makefile`, `scripts/mutation-test.sh`
  - **DoD**:
    - [x] Integrated Gremlins mutation engine with `make mutation-test`.
    - [x] 100% test efficacy and 100% mutator coverage verified on runtime packages.

- [x] **[#25: Fix CI/CD Pipeline (ci.yml) to Enforce Strict Quality Gate](https://github.com/el-j/ts2go/issues/25)**
  - **Status**: 🟢 Resolved & Verified
  - **Target Files**: `.github/workflows/ci.yml`
  - **DoD**:
    - [x] Enforced Go 1.24+ using `go.work`.
    - [x] Zero tolerance for linting failures (`golangci-lint` without `|| true`).
    - [x] Strict test gate executing all Go packages with race detector enabled (`go test -v -race ./...`).

---

### Milestone 4: Web API & GUI Modernization (Priority: P2)
*Goal: Modularize UI boss files, ensure error-free builds, connect Web API backend.*

- [ ] **[#26: Fix Frontend Build, TS Errors & Vitest Runner](https://github.com/el-j/ts2go/issues/26)**
  - **Status**: 🟢 Completed in Initial Quality Pass (TS5103 fixed, `vitest run` configured, prettier added).

- [ ] **[#27: Break Down UI Boss Files (ProjectView.vue & SettingsView.vue)](https://github.com/el-j/ts2go/issues/27)**
  - **Target Files**: `packages/ui-shared/src/views/ProjectView.vue` (1,954 lines)
  - **DoD**: Break into subcomponents <400 lines each.

- [ ] **[#28: Implement Real Web API Backend (Replace Placeholder Stubs in server.go)](https://github.com/el-j/ts2go/issues/28)**
  - **Target Files**: `go/adapters/driving/web/server.go`
  - **DoD**: Replace `placeholderCodeGen` and `placeholderMapper` with real implementations; test all endpoints.
