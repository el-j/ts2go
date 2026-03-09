# Project Structure Restructure Plan

**Date:** November 16, 2025  
**Status:** Planning  
**Goal:** Create a clean, intuitive project structure that separates concerns

## 🎯 Vision

A first-time visitor should immediately understand:
1. **Go CLI/Core** - The main transpilation engine
2. **Desktop UI** - The Tauri + Vue application  
3. **Build System** - Simple Makefile orchestrating everything

## 📊 Current Problems

### Root Directory Chaos
```
❌ cmd/, pkg/, internal/ mixed with desktop-ui/, docs/, examples/
❌ go.mod, go.work at root but also nested
❌ Hard to tell what's Go vs TypeScript vs Rust
❌ No clear entry point for new developers
❌ Makefile hidden among many files
```

### Build Confusion
```
❌ Multiple go.mod files (root, cmd/ts2go, pkg/, runtime/, tests/)
❌ Desktop UI buried in root
❌ Tauri (Rust) hidden inside desktop-ui/src-tauri/
❌ Runtime library scattered
```

## ✨ Proposed Structure

```
ts2go/
├── Makefile                    # 🎯 Main orchestration
├── README.md                   # 📖 Clear entry point
├── LICENSE
├── VERSION
│
├── go/                         # 🔵 Go CLI & Core Engine
│   ├── go.mod                  # Single Go module
│   ├── go.sum
│   ├── cmd/
│   │   ├── ts2go/             # CLI binary
│   │   └── ts2go-web/         # Web server
│   ├── core/                   # Hexagonal architecture
│   │   ├── domain/            # Business entities
│   │   ├── ports/             # Interface contracts
│   │   └── services/          # Business logic
│   ├── adapters/              # Infrastructure
│   │   ├── driven/            # Filesystem, compiler, persistence
│   │   └── driving/           # CLI, web API
│   ├── runtime/               # Go runtime library
│   │   ├── array/
│   │   ├── buffer/
│   │   ├── console/
│   │   └── ...
│   └── internal/              # Legacy (deprecated)
│
├── desktop/                    # 🟢 Desktop Application
│   ├── README.md              # Desktop-specific docs
│   ├── Makefile               # Desktop build commands
│   ├── package.json           # Vue/TypeScript deps
│   ├── ui/                    # Vue frontend
│   │   ├── src/
│   │   ├── public/
│   │   ├── index.html
│   │   └── vite.config.ts
│   └── tauri/                 # Rust backend
│       ├── Cargo.toml
│       ├── src/
│       └── target/
│
├── docs/                       # 📚 Documentation
│   ├── getting-started.md
│   ├── architecture.md
│   ├── api-reference.md
│   └── archive/               # Historical docs
│
├── examples/                   # 💡 Usage Examples
│   ├── simple-function/
│   ├── class-conversion/
│   └── multipackage-demo/
│
├── scripts/                    # 🔧 Build & CI Scripts
│   ├── bump-version.sh
│   ├── release.sh
│   └── install-deps.sh
│
└── .github/                    # 🤖 CI/CD
    └── workflows/
```

## 📋 Detailed Changes

### 1. Go CLI Consolidation

**Create:** `go/` directory for all Go code

**Move:**
- `pkg/core/` → `go/core/`
- `pkg/adapters/` → `go/adapters/`
- `pkg/cli/` → `go/adapters/driving/cli/` (it's a driving adapter!)
- `pkg/optimizer/` → `go/core/services/optimizer/`
- `cmd/` → `go/cmd/`
- `runtime/` → `go/runtime/`
- `internal/` → `go/internal/` (keep deprecated marker)

**Single go.mod:**
- Consolidate all Go modules into `go/go.mod`
- Remove nested go.mod files
- Update import paths: `github.com/el-j/ts2go/go/core/...`

### 2. Desktop Application Separation

**Create:** `desktop/` directory

**Move:**
- `desktop-ui/` → `desktop/`
- Restructure inside desktop:
  - `desktop-ui/src/` → `desktop/ui/src/`
  - `desktop-ui/src-tauri/` → `desktop/tauri/`
  - `desktop-ui/package.json` → `desktop/package.json`
  - `desktop-ui/docs/` → `desktop/docs/` or `docs/desktop/`

**Clean separation:**
- `desktop/ui/` - Pure Vue/TypeScript frontend
- `desktop/tauri/` - Pure Rust Tauri backend
- `desktop/Makefile` - Desktop-specific commands
- `desktop/README.md` - Desktop-specific documentation

### 3. Root Directory Cleanup

**Keep in root (10-12 files):**
```
✅ Makefile          # Main orchestration
✅ README.md         # Project overview
✅ LICENSE           # Legal
✅ VERSION           # Version number
✅ CHANGELOG.md      # Release notes
✅ CONTRIBUTING.md   # How to contribute
✅ .gitignore        # Git config
✅ .dockerignore     # Docker config
```

**Move out of root:**
```
📁 CURRENT_STATE.md → docs/current-state.md
📁 PROJECT_STATUS.md → docs/project-status.md
📁 SPEC.md → docs/specification.md
📁 KNOWN_ISSUES.md → docs/known-issues.md
📁 CLEANUP_PLAN.md → docs/archive/cleanup-plan.md
```

### 4. Documentation Reorganization

**Structure:**
```
docs/
├── README.md                   # Docs index
├── getting-started.md          # Quick start
├── architecture.md             # System design
├── specification.md            # TypeScript features
├── api-reference.md            # API docs
├── known-issues.md             # Limitations
├── current-state.md            # Implementation status
├── project-status.md           # High-level overview
│
├── guides/                     # User guides
│   ├── migration-guide.md
│   ├── package-mappings.md
│   └── desktop-user-guide.md
│
├── development/                # Developer docs
│   ├── contributing.md
│   ├── testing-guide.md
│   └── release-process.md
│
└── archive/                    # Historical
    ├── phase0-complete.md
    └── ...
```

### 5. Makefile Simplification

**New root Makefile:**
```makefile
# ts2go Root Makefile
# Orchestrates Go CLI and Desktop application builds

.PHONY: all build test clean install

# Default: build everything
all: build-go build-desktop

# Go CLI
build-go:
	@echo "Building Go CLI..."
	cd go && go build -o ../bin/ts2go ./cmd/ts2go

test-go:
	@echo "Testing Go code..."
	cd go && go test ./...

# Desktop Application
build-desktop:
	@echo "Building Desktop app..."
	cd desktop && $(MAKE) build

run-desktop:
	cd desktop && $(MAKE) dev

# Combined
test: test-go test-desktop
	@echo "All tests passed!"

clean:
	rm -rf bin/
	cd go && go clean
	cd desktop && $(MAKE) clean

install:
	./scripts/install-deps.sh

# Development
dev-go:
	cd go && go run ./cmd/ts2go

dev-desktop:
	cd desktop && $(MAKE) dev

# Help
help:
	@echo "ts2go Build System"
	@echo ""
	@echo "Targets:"
	@echo "  all           - Build everything"
	@echo "  build-go      - Build Go CLI"
	@echo "  build-desktop - Build Desktop app"
	@echo "  test          - Run all tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install dependencies"
	@echo "  dev-go        - Run Go CLI in dev mode"
	@echo "  dev-desktop   - Run Desktop app in dev mode"
```

## 🎯 Migration Tasks

### Phase 1: Preparation
- [ ] Create new directory structure (go/, desktop/)
- [ ] Backup current state
- [ ] Update .gitignore for new paths

### Phase 2: Go Consolidation
- [ ] Create go/go.mod with correct module path
- [ ] Move pkg/core/ → go/core/
- [ ] Move pkg/adapters/ → go/adapters/
- [ ] Move pkg/cli/ → go/adapters/driving/cli/
- [ ] Move cmd/ → go/cmd/
- [ ] Move runtime/ → go/runtime/
- [ ] Move internal/ → go/internal/
- [ ] Update all import paths in Go files
- [ ] Remove old pkg/ directory
- [ ] Verify go build works

### Phase 3: Desktop Separation
- [ ] Create desktop/ directory structure
- [ ] Move desktop-ui/ content → desktop/
- [ ] Separate ui/ and tauri/ subdirectories
- [ ] Create desktop/Makefile
- [ ] Update desktop package.json paths
- [ ] Update Tauri config paths
- [ ] Verify desktop build works

### Phase 4: Documentation
- [ ] Move root .md files → docs/
- [ ] Create docs/README.md index
- [ ] Organize into guides/ and development/
- [ ] Update all documentation links
- [ ] Create comprehensive docs/getting-started.md

### Phase 5: Build System
- [ ] Create new root Makefile
- [ ] Create desktop/Makefile
- [ ] Create scripts/install-deps.sh
- [ ] Update CI/CD workflows for new paths
- [ ] Test all make targets

### Phase 6: Final Cleanup
- [ ] Remove old directories (pkg/, desktop-ui/)
- [ ] Update README.md with new structure
- [ ] Update CONTRIBUTING.md
- [ ] Run full test suite
- [ ] Create migration guide for contributors

## ✅ Success Criteria

### For New Contributors
✅ Understand project at a glance  
✅ Know where Go code lives (`go/`)  
✅ Know where Desktop app lives (`desktop/`)  
✅ Clear build commands (`make build-go`, `make build-desktop`)  
✅ Obvious entry point (README.md → Makefile)  

### For Existing Users
✅ All commands still work  
✅ Tests still pass  
✅ Builds still succeed  
✅ Migration guide available  
✅ No breaking changes to published APIs  

### For Maintainers
✅ Clearer module boundaries  
✅ Easier to navigate codebase  
✅ Better separation of concerns  
✅ Simpler build orchestration  
✅ Obvious where new code belongs  

## 📈 Benefits

1. **Clarity** - Immediate understanding of project structure
2. **Separation** - Go, TypeScript, Rust clearly separated
3. **Scalability** - Easy to add new components
4. **Onboarding** - New contributors productive faster
5. **Maintenance** - Easier to update individual parts
6. **Professional** - Matches industry standards

## 🚀 Implementation Order

1. **Phase 1: Preparation** (15 min)
2. **Phase 2: Go Consolidation** (45 min) - Most import path updates
3. **Phase 3: Desktop Separation** (30 min)
4. **Phase 4: Documentation** (20 min)
5. **Phase 5: Build System** (25 min)
6. **Phase 6: Final Cleanup** (15 min)

**Total estimated time:** ~2.5 hours

## ⚠️ Risks & Mitigations

**Risk:** Breaking existing tools/scripts  
**Mitigation:** Create symlinks for transition period

**Risk:** Import path changes break external users  
**Mitigation:** This is pre-1.0, acceptable to change

**Risk:** Desktop build breaks  
**Mitigation:** Test after each move, have rollback plan

**Risk:** CI/CD breaks  
**Mitigation:** Update workflows incrementally

## 🎬 Let's Do This!

Ready to execute? This will transform the project from cluttered to clean, from confusing to crystal clear.
