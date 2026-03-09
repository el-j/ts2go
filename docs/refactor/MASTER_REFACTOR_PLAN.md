# TS2Go Hexagonal Architecture Refactoring - Master Plan

**Date Started:** November 16, 2024  
**Target Completion:** v3.0.0  
**Architecture:** Hexagonal (Ports & Adapters)

---

## 🎯 Executive Summary

### Objective
Transform TS2Go from a monolithic architecture into a clean, hexagonal (ports and adapters) architecture that:
- Separates core business logic from infrastructure concerns
- Enables multiple delivery mechanisms (CLI, Desktop UI, Web API)
- Improves testability with dependency injection
- Allows teams to work independently on UI, core, and adapters
- Eliminates code duplication between Rust (Tauri) and Go backends

### Current State
- ✅ Desktop UI with Vue 3 + Tauri
- ✅ CLI tool with basic transpilation
- ❌ Logic duplicated between Rust (main.rs) and Go
- ❌ State management split between localStorage and backend
- ❌ Tight coupling between core logic and infrastructure
- ❌ Difficult to test without real file system/Go compiler

### Target State
```
┌─────────────────────────────────────────────────────────┐
│                    DRIVING ADAPTERS                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐ │
│  │   CLI    │  │ Tauri UI │  │  Web API │  │ gRPC   │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └───┬────┘ │
└───────┼─────────────┼─────────────┼─────────────┼──────┘
        │             │             │             │
        └─────────────┴─────────────┴─────────────┘
                          │
        ┌─────────────────┴─────────────────┐
        │         CORE (HEXAGON)             │
        │  ┌──────────────────────────────┐ │
        │  │  Domain Models & Entities    │ │
        │  ├──────────────────────────────┤ │
        │  │  Business Logic Services     │ │
        │  │  - TranspilationService      │ │
        │  │  - AnalysisService           │ │
        │  │  - RuntimeService            │ │
        │  ├──────────────────────────────┤ │
        │  │  Port Interfaces (Contracts) │ │
        │  └──────────────────────────────┘ │
        └───────────────┬───────────────────┘
                        │
        ┌───────────────┴───────────────┐
        │      DRIVEN ADAPTERS           │
        │  ┌──────────┐  ┌────────────┐ │
        │  │   File   │  │    Go      │ │
        │  │  System  │  │  Compiler  │ │
        │  └──────────┘  └────────────┘ │
        │  ┌──────────┐  ┌────────────┐ │
        │  │   JSON   │  │   SQLite   │ │
        │  │  Store   │  │    DB      │ │
        │  └──────────┘  └────────────┘ │
        └───────────────────────────────┘
```

---

## 📋 Comprehensive Task Breakdown

### Phase 0: Foundation & Architecture Setup
**Status:** 🔄 In Progress  
**Duration:** 1-2 days  
**Deliverables:** Directory structure, port interfaces, domain models

#### Tasks:
- [ ] 0.1: Create new `pkg/` directory structure
  - [ ] 0.1.1: Create `pkg/core/domain/`
  - [ ] 0.1.2: Create `pkg/core/ports/`
  - [ ] 0.1.3: Create `pkg/core/services/`
  - [ ] 0.1.4: Create `pkg/adapters/driven/filesystem/`
  - [ ] 0.1.5: Create `pkg/adapters/driven/gocompiler/`
  - [ ] 0.1.6: Create `pkg/adapters/driven/persistence/`
  - [ ] 0.1.7: Create `pkg/adapters/driving/cli/`
  - [ ] 0.1.8: Create `pkg/adapters/driving/web/`

- [ ] 0.2: Define domain models (`pkg/core/domain/`)
  - [ ] 0.2.1: Create `project.go` - Project entity
  - [ ] 0.2.2: Create `file.go` - File entity
  - [ ] 0.2.3: Create `transpilation.go` - TranspilationResult, AnalysisReport
  - [ ] 0.2.4: Create `runtime.go` - GoVersion, BuildResult, RunResult, TestResult
  - [ ] 0.2.5: Create `state.go` - ProjectState, Settings
  - [ ] 0.2.6: Create `errors.go` - Domain-specific errors

- [ ] 0.3: Define all port interfaces (`pkg/core/ports/`)
  - [ ] 0.3.1: Create `transpilation.go` - TranspilationService interface
  - [ ] 0.3.2: Create `runtime.go` - GoRuntimeService interface
  - [ ] 0.3.3: Create `state.go` - StateService interface
  - [ ] 0.3.4: Create `filesystem.go` - FileSystem interface
  - [ ] 0.3.5: Create `compiler.go` - GoCompiler interface
  - [ ] 0.3.6: Create `repository.go` - StateRepository, SettingsRepository interfaces

- [ ] 0.4: Documentation
  - [ ] 0.4.1: Create `docs/refactor/PHASE_0_FOUNDATION.md`
  - [ ] 0.4.2: Document architecture decisions
  - [ ] 0.4.3: Create interface usage examples

---

### Phase 1: Core Business Logic Migration
**Status:** ⏳ Not Started  
**Duration:** 3-5 days  
**Deliverables:** Pure business logic in `pkg/core/services/`, independent of infrastructure

#### Tasks:
- [ ] 1.1: Migrate Transpiler Logic
  - [ ] 1.1.1: Copy `internal/transpiler/` to `pkg/core/services/transpiler/`
  - [ ] 1.1.2: Create `transpilation_service.go` with interface dependencies
  - [ ] 1.1.3: Replace all `os.*` calls with `FileSystem` port calls
  - [ ] 1.1.4: Replace all `exec.Command("go")` with `GoCompiler` port calls
  - [ ] 1.1.5: Remove all infrastructure imports (os, exec, etc.)
  - [ ] 1.1.6: Add unit tests with mock adapters

- [ ] 1.2: Migrate Analyzer Logic
  - [ ] 1.2.1: Copy `internal/analyzer/` to `pkg/core/services/analyzer/`
  - [ ] 1.2.2: Create `analysis_service.go`
  - [ ] 1.2.3: Refactor to use FileSystem port
  - [ ] 1.2.4: Add comprehensive tests

- [ ] 1.3: Migrate Mapper Logic
  - [ ] 1.3.1: Copy `internal/mapper/` to `pkg/core/services/mapper/`
  - [ ] 1.3.2: Ensure pure mapping logic (no I/O)
  - [ ] 1.3.3: Add mapping tests

- [ ] 1.4: Create Runtime Service
  - [ ] 1.4.1: Create `runtime_service.go`
  - [ ] 1.4.2: Implement DetectGoInstallation using GoCompiler port
  - [ ] 1.4.3: Implement BuildProject using GoCompiler port
  - [ ] 1.4.4: Implement RunProject using GoCompiler port
  - [ ] 1.4.5: Implement TestProject using GoCompiler port
  - [ ] 1.4.6: Add tests with mock compiler

- [ ] 1.5: Create State Service
  - [ ] 1.5.1: Create `state_service.go`
  - [ ] 1.5.2: Implement GetProjectState using StateRepository port
  - [ ] 1.5.3: Implement SaveProjectState using StateRepository port
  - [ ] 1.5.4: Implement GetSettings using SettingsRepository port
  - [ ] 1.5.5: Implement SaveSettings using SettingsRepository port
  - [ ] 1.5.6: Add persistence tests

- [ ] 1.6: Documentation
  - [ ] 1.6.1: Create `docs/refactor/PHASE_1_CORE_MIGRATION.md`
  - [ ] 1.6.2: Document service APIs
  - [ ] 1.6.3: Add architecture diagrams

---

### Phase 2: Infrastructure Adapters Implementation
**Status:** ⏳ Not Started  
**Duration:** 2-3 days  
**Deliverables:** Concrete adapter implementations for all driven ports

#### Tasks:
- [ ] 2.1: FileSystem Adapter
  - [ ] 2.1.1: Create `pkg/adapters/driven/filesystem/disk_fs.go`
  - [ ] 2.1.2: Implement ReadFile using os.ReadFile
  - [ ] 2.1.3: Implement WriteFile using os.WriteFile
  - [ ] 2.1.4: Implement ScanDirectory using filepath.Walk
  - [ ] 2.1.5: Implement DirectoryExists using os.Stat
  - [ ] 2.1.6: Add integration tests

- [ ] 2.2: GoCompiler Adapter
  - [ ] 2.2.1: Create `pkg/adapters/driven/gocompiler/os_go.go`
  - [ ] 2.2.2: Port `detect_go_installation` from Tauri main.rs
  - [ ] 2.2.3: Port `run_go_code` from Tauri main.rs
  - [ ] 2.2.4: Port `build_go_file` from Tauri main.rs
  - [ ] 2.2.5: Port `test_go_project` from Tauri main.rs
  - [ ] 2.2.6: Implement Format using `go fmt`
  - [ ] 2.2.7: Add integration tests with real Go compiler

- [ ] 2.3: Persistence Adapters
  - [ ] 2.3.1: Create `pkg/adapters/driven/persistence/json_state.go`
  - [ ] 2.3.2: Implement StateRepository with JSON file storage
  - [ ] 2.3.3: Create `pkg/adapters/driven/persistence/json_settings.go`
  - [ ] 2.3.4: Implement SettingsRepository with JSON file storage
  - [ ] 2.3.5: Add file locking for concurrent access
  - [ ] 2.3.6: Add migration logic for existing localStorage data
  - [ ] 2.3.7: Add integration tests

- [ ] 2.4: Mock Adapters for Testing
  - [ ] 2.4.1: Create `pkg/adapters/driven/filesystem/mock_fs.go`
  - [ ] 2.4.2: Create `pkg/adapters/driven/gocompiler/mock_compiler.go`
  - [ ] 2.4.3: Create `pkg/adapters/driven/persistence/mock_repo.go`
  - [ ] 2.4.4: Add in-memory implementations for fast testing

- [ ] 2.5: Documentation
  - [ ] 2.5.1: Create `docs/refactor/PHASE_2_ADAPTERS.md`
  - [ ] 2.5.2: Document adapter implementations
  - [ ] 2.5.3: Add adapter usage examples

---

### Phase 3: CLI Refactoring & Dependency Injection
**Status:** ⏳ Not Started  
**Duration:** 2-3 days  
**Deliverables:** Refactored CLI using hexagonal core, clean DI setup

#### Tasks:
- [ ] 3.1: Create CLI Adapter
  - [ ] 3.1.1: Create `pkg/adapters/driving/cli/app.go`
  - [ ] 3.1.2: Define CLI application struct with service dependencies
  - [ ] 3.1.3: Create command handlers as methods
  - [ ] 3.1.4: Add output formatting helpers

- [ ] 3.2: Refactor main.go (Composition Root)
  - [ ] 3.2.1: Update `cmd/ts2go/main.go`
  - [ ] 3.2.2: Instantiate all driven adapters
  - [ ] 3.2.3: Inject adapters into core services
  - [ ] 3.2.4: Inject services into CLI app
  - [ ] 3.2.5: Wire up Cobra commands to app methods

- [ ] 3.3: Update CLI Commands
  - [ ] 3.3.1: Refactor `transpile` command to call TranspilationService
  - [ ] 3.3.2: Refactor `analyze` command to call AnalysisService
  - [ ] 3.3.3: Refactor `run` command to call RuntimeService
  - [ ] 3.3.4: Refactor `build` command to call RuntimeService
  - [ ] 3.3.5: Refactor `test` command to call RuntimeService
  - [ ] 3.3.6: Add `detect-go` command for UI support
  - [ ] 3.3.7: Add `state list` command for UI support
  - [ ] 3.3.8: Add `state get` command for UI support
  - [ ] 3.3.9: Add `settings get` command for UI support
  - [ ] 3.3.10: Add `settings set` command for UI support

- [ ] 3.4: Testing
  - [ ] 3.4.1: Add CLI integration tests
  - [ ] 3.4.2: Test with mock adapters
  - [ ] 3.4.3: Test with real adapters (e2e)

- [ ] 3.5: Documentation
  - [ ] 3.5.1: Create `docs/refactor/PHASE_3_CLI_REFACTOR.md`
  - [ ] 3.5.2: Update CLI usage documentation
  - [ ] 3.5.3: Add DI pattern examples

---

### Phase 4: Tauri Backend Simplification
**Status:** ⏳ Not Started  
**Duration:** 2-3 days  
**Deliverables:** Simplified Tauri backend that delegates to CLI, no duplicate logic

#### Tasks:
- [ ] 4.1: Gut Duplicate Rust Logic
  - [ ] 4.1.1: Remove `detect_go_installation` from main.rs
  - [ ] 4.1.2: Remove `run_go_code` from main.rs
  - [ ] 4.1.3: Remove `build_go_file` from main.rs
  - [ ] 4.1.4: Remove `test_go_project` from main.rs
  - [ ] 4.1.5: Remove `check_go_fmt` from main.rs
  - [ ] 4.1.6: Remove all `exec.Command("go")` calls

- [ ] 4.2: Create CLI Delegation Layer
  - [ ] 4.2.1: Create `desktop-ui/src-tauri/src/cli_bridge.rs`
  - [ ] 4.2.2: Add helper to find ts2go-cli binary
  - [ ] 4.2.3: Add helper to execute CLI commands
  - [ ] 4.2.4: Add JSON parsing for CLI output
  - [ ] 4.2.5: Add error handling and logging

- [ ] 4.3: Refactor Tauri Commands
  - [ ] 4.3.1: Refactor `transpile_project` to call `ts2go-cli transpile`
  - [ ] 4.3.2: Refactor `detect_go_installation` to call `ts2go-cli detect-go`
  - [ ] 4.3.3: Refactor `run_go_project` to call `ts2go-cli run`
  - [ ] 4.3.4: Refactor `build_go_file` to call `ts2go-cli build`
  - [ ] 4.3.5: Refactor `test_go_project` to call `ts2go-cli test`
  - [ ] 4.3.6: Add `get_project_state` calling `ts2go-cli state get`
  - [ ] 4.3.7: Add `get_all_states` calling `ts2go-cli state list`
  - [ ] 4.3.8: Add `get_settings` calling `ts2go-cli settings get`
  - [ ] 4.3.9: Add `save_settings` calling `ts2go-cli settings set`

- [ ] 4.4: Bundle CLI with Desktop App
  - [ ] 4.4.1: Update Tauri build config to include CLI binary
  - [ ] 4.4.2: Add CLI to app resources
  - [ ] 4.4.3: Update path resolution for bundled CLI

- [ ] 4.5: Testing
  - [ ] 4.5.1: Add Tauri command tests
  - [ ] 4.5.2: Test CLI integration
  - [ ] 4.5.3: Test on all platforms (macOS, Windows, Linux)

- [ ] 4.6: Documentation
  - [ ] 4.6.1: Create `docs/refactor/PHASE_4_TAURI_SIMPLIFICATION.md`
  - [ ] 4.6.2: Document Tauri-CLI bridge
  - [ ] 4.6.3: Update build documentation

---

### Phase 5: Frontend State Migration
**Status:** ⏳ Not Started  
**Duration:** 1-2 days  
**Deliverables:** Vue stores using backend persistence, no localStorage

#### Tasks:
- [ ] 5.1: Refactor Workspace Store
  - [ ] 5.1.1: Update `desktop-ui/src/stores/workspace.ts`
  - [ ] 5.1.2: Remove localStorage logic
  - [ ] 5.1.3: Replace with Tauri command calls
  - [ ] 5.1.4: Add state synchronization logic

- [ ] 5.2: Refactor Settings Store
  - [ ] 5.2.1: Update `desktop-ui/src/stores/settings.ts`
  - [ ] 5.2.2: Remove localStorage logic
  - [ ] 5.2.3: Replace with Tauri command calls for get_settings/save_settings
  - [ ] 5.2.4: Add settings migration from localStorage

- [ ] 5.3: Data Migration Tool
  - [ ] 5.3.1: Create migration script to export localStorage data
  - [ ] 5.3.2: Create CLI command `ts2go-cli migrate import`
  - [ ] 5.3.3: Add migration UI in desktop app
  - [ ] 5.3.4: Add rollback capability

- [ ] 5.4: Testing
  - [ ] 5.4.1: Test state persistence across app restarts
  - [ ] 5.4.2: Test migration from localStorage
  - [ ] 5.4.3: Test concurrent state updates

- [ ] 5.5: Documentation
  - [ ] 5.5.1: Create `docs/refactor/PHASE_5_FRONTEND_MIGRATION.md`
  - [ ] 5.5.2: Document state management patterns
  - [ ] 5.5.3: Add migration guide

---

### Phase 6: Web API Adapter (Proof of Concept)
**Status:** ⏳ Not Started  
**Duration:** 2-3 days  
**Deliverables:** Working REST API demonstrating hexagonal architecture benefits

#### Tasks:
- [ ] 6.1: Create Web API Adapter
  - [ ] 6.1.1: Create `pkg/adapters/driving/web/`
  - [ ] 6.1.2: Create `server.go` with gin or net/http
  - [ ] 6.1.3: Create `handlers.go` for HTTP endpoints
  - [ ] 6.1.4: Create `middleware.go` for auth, logging
  - [ ] 6.1.5: Add CORS configuration

- [ ] 6.2: Implement REST Endpoints
  - [ ] 6.2.1: POST `/api/v1/transpile` - Transpile project
  - [ ] 6.2.2: GET `/api/v1/analyze/:projectId` - Analyze project
  - [ ] 6.2.3: POST `/api/v1/run` - Run Go project
  - [ ] 6.2.4: POST `/api/v1/build` - Build Go project
  - [ ] 6.2.5: POST `/api/v1/test` - Test Go project
  - [ ] 6.2.6: GET `/api/v1/go/version` - Detect Go installation
  - [ ] 6.2.7: GET `/api/v1/state` - Get all project states
  - [ ] 6.2.8: GET `/api/v1/settings` - Get settings
  - [ ] 6.2.9: PUT `/api/v1/settings` - Update settings

- [ ] 6.3: Create Web Main
  - [ ] 6.3.1: Create `cmd/ts2go-web/main.go`
  - [ ] 6.3.2: Copy composition root from CLI main.go
  - [ ] 6.3.3: Instantiate web adapter with core services
  - [ ] 6.3.4: Add graceful shutdown
  - [ ] 6.3.5: Add configuration (port, host, etc.)

- [ ] 6.4: API Documentation
  - [ ] 6.4.1: Add OpenAPI/Swagger spec
  - [ ] 6.4.2: Generate API documentation
  - [ ] 6.4.3: Add Postman collection

- [ ] 6.5: Testing
  - [ ] 6.5.1: Add API integration tests
  - [ ] 6.5.2: Add load tests
  - [ ] 6.5.3: Test authentication/authorization

- [ ] 6.6: Documentation
  - [ ] 6.6.1: Create `docs/refactor/PHASE_6_WEB_API.md`
  - [ ] 6.6.2: Add API usage examples
  - [ ] 6.6.3: Document deployment guide

---

### Phase 7: Cleanup & Legacy Removal
**Status:** ⏳ Not Started  
**Duration:** 1-2 days  
**Deliverables:** Clean codebase, legacy code removed

#### Tasks:
- [ ] 7.1: Remove Old Code
  - [ ] 7.1.1: Delete `internal/transpiler/` (moved to pkg/core)
  - [ ] 7.1.2: Delete `internal/analyzer/` (moved to pkg/core)
  - [ ] 7.1.3: Delete `internal/mapper/` (moved to pkg/core)
  - [ ] 7.1.4: Delete `internal/orchestrator/` (if unused)
  - [ ] 7.1.5: Delete old CLI commands (if any)

- [ ] 7.2: Update Imports
  - [ ] 7.2.1: Run `gofmt -w .` on entire codebase
  - [ ] 7.2.2: Run `goimports -w .` to fix imports
  - [ ] 7.2.3: Update go.mod dependencies
  - [ ] 7.2.4: Run `go mod tidy`

- [ ] 7.3: Code Quality
  - [ ] 7.3.1: Run `golangci-lint` on entire codebase
  - [ ] 7.3.2: Fix all linter warnings
  - [ ] 7.3.3: Add missing godoc comments
  - [ ] 7.3.4: Improve error messages

- [ ] 7.4: Final Testing
  - [ ] 7.4.1: Run full test suite
  - [ ] 7.4.2: Run integration tests
  - [ ] 7.4.3: Test CLI on all platforms
  - [ ] 7.4.4: Test Desktop UI on all platforms
  - [ ] 7.4.5: Test Web API

- [ ] 7.5: Documentation
  - [ ] 7.5.1: Create `docs/refactor/PHASE_7_CLEANUP.md`
  - [ ] 7.5.2: Update main README.md
  - [ ] 7.5.3: Update ARCHITECTURE.md
  - [ ] 7.5.4: Create migration guide for contributors

---

## 📊 Success Metrics

### Code Quality Metrics
- [ ] 100% of core business logic is pure (no I/O dependencies)
- [ ] 90%+ test coverage on core services
- [ ] 80%+ test coverage on adapters
- [ ] Zero golangci-lint warnings
- [ ] All godoc comments present

### Architecture Metrics
- [ ] Zero direct infrastructure dependencies in core
- [ ] All services use dependency injection
- [ ] All adapters implement port interfaces
- [ ] Multiple driving adapters working (CLI, Tauri, Web)

### Performance Metrics
- [ ] CLI performance matches or exceeds old implementation
- [ ] Desktop UI responsiveness unchanged
- [ ] Web API handles 100+ concurrent requests

### Developer Experience Metrics
- [ ] New team members can understand architecture in < 1 day
- [ ] Adding new adapter takes < 1 day
- [ ] Adding new feature to core takes < 2 days
- [ ] Tests run in < 30 seconds

---

## 🚧 Risk Assessment & Mitigation

### High Risks
1. **Breaking existing functionality**
   - Mitigation: Comprehensive testing at each phase
   - Mitigation: Keep old code until new code is verified

2. **Performance regression**
   - Mitigation: Benchmark before and after each phase
   - Mitigation: Profile critical paths

3. **Desktop UI instability during Tauri changes**
   - Mitigation: Test on all platforms after each change
   - Mitigation: Keep fallback to direct Rust implementation

### Medium Risks
1. **Data migration issues (localStorage → backend)**
   - Mitigation: Create robust migration tool
   - Mitigation: Allow rollback to localStorage

2. **Complexity of dependency injection**
   - Mitigation: Document DI patterns clearly
   - Mitigation: Use simple, explicit wiring (no magic)

3. **Learning curve for contributors**
   - Mitigation: Excellent documentation
   - Mitigation: Pair programming sessions
   - Mitigation: Video walkthroughs

---

## 📅 Timeline

| Phase | Duration | Start Date | End Date | Status |
|-------|----------|------------|----------|--------|
| Phase 0 | 1-2 days | Nov 16 | Nov 17 | 🔄 In Progress |
| Phase 1 | 3-5 days | Nov 18 | Nov 22 | ⏳ Pending |
| Phase 2 | 2-3 days | Nov 23 | Nov 25 | ⏳ Pending |
| Phase 3 | 2-3 days | Nov 26 | Nov 28 | ⏳ Pending |
| Phase 4 | 2-3 days | Nov 29 | Dec 1 | ⏳ Pending |
| Phase 5 | 1-2 days | Dec 2 | Dec 3 | ⏳ Pending |
| Phase 6 | 2-3 days | Dec 4 | Dec 6 | ⏳ Pending |
| Phase 7 | 1-2 days | Dec 7 | Dec 8 | ⏳ Pending |
| **Total** | **14-23 days** | **Nov 16** | **Dec 8** | |

---

## 🎓 Learning Resources

### Hexagonal Architecture
- [Alistair Cockburn - Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Netflix - Ready for changes with Hexagonal Architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749)

### Dependency Injection in Go
- [Go Patterns - Dependency Injection](https://www.practical-go-lessons.com/chap-33-dependency-injection)
- [Uber Go Guide - Dependency Injection](https://github.com/uber-go/guide/blob/master/style.md#dependency-injection)

### Testing Strategies
- [Go Testing Best Practices](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testify - Go Testing Framework](https://github.com/stretchr/testify)

---

## 📝 Notes

- Each phase should be completed and verified before moving to the next
- All tests must pass before proceeding
- Documentation should be updated continuously
- Rollback plan exists for each phase
- Team review required before Phase 7 (cleanup)

---

**Last Updated:** November 16, 2024  
**Document Owner:** @el-j  
**Status:** Living Document - Update as refactoring progresses
