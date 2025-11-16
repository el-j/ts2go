# Root Directory Cleanup Plan

**Date:** November 16, 2025  
**Status:** In Progress  
**Goal:** Clean up cluttered root directory after hexagonal refactor

## 🎯 Objectives

1. **Archive outdated documentation** - Move phase completion docs to archive
2. **Consolidate roadmaps** - Keep only current/active roadmap
3. **Remove obsolete binaries** - ts2go and ts2go-web executables (generated files)
4. **Evaluate internal/** - Determine if still needed after hexagonal refactor
5. **Clean documentation structure** - Organize docs/ directory better

## 📋 Analysis

### Root Directory Status (BEFORE)

```
Total files in root: ~30 files
- 10+ markdown files (roadmaps, phase completions, status files)
- 2 binaries (ts2go, ts2go-web)
- Multiple roadmap files
- Phase completion documents
```

### What to Keep

**Essential Files:**
✅ `README.md` - Main documentation
✅ `LICENSE` - License
✅ `CHANGELOG.md` - Version history
✅ `CONTRIBUTING.md` - Contribution guide
✅ `Makefile` - Build system
✅ `Dockerfile` - Container support
✅ `VERSION` - Version tracking
✅ `go.mod`, `go.sum`, `go.work` - Go modules
✅ `.gitignore`, `.dockerignore` - VCS config

**Current Status:**
✅ `CURRENT_STATE.md` - Keep (main status)
✅ `PROJECT_STATUS.md` - Keep (backup status view)
✅ `SPEC.md` - Keep (TypeScript feature spec)
✅ `KNOWN_ISSUES.md` - Keep (important for users)

### What to Archive

**Phase Completion Documents (Move to docs/archive/):**
- `PHASE0_GO_FMT_IMPLEMENTATION.md`
- `PHASE1_AND_2_IMPLEMENTATION_COMPLETE.md`
- `ALL_TASKS_COMPLETE.md`
- `TESTING_GUIDE_PHASE1_AND_2.md`

**Old Roadmaps (Move to docs/archive/):**
- `ROADMAP_TO_1.0.0.md`
- `ROADMAP_TO_3.0.0-hexagonal.md`
- `ROADMAP_TO_3.0.0-moreinfo.md`
- `ROADMAP_TO_ALPHA_2.0.1.md`

### What to Remove

**Generated Binaries (Should be in .gitignore):**
- `ts2go` (executable, should be built not committed)
- `ts2go-web` (executable, should be built not committed)

**Note:** Check if these are in `.gitignore`

### Internal/ Directory

**Status:** KEEP for now
**Reason:** 
- Still used by `convert` command in main.go
- Used by examples/multipackage-demo
- Used by integration tests
- Phase 7 docs say "deprecated but kept for compatibility"

**Action:** Add clear deprecation notice

## 🔧 Cleanup Actions

### 1. Create Archive Directory ✅

```bash
mkdir -p docs/archive
```

### 2. Move Phase Documents ✅

```bash
# Phase completion docs
git mv PHASE0_GO_FMT_IMPLEMENTATION.md docs/archive/
git mv PHASE1_AND_2_IMPLEMENTATION_COMPLETE.md docs/archive/
git mv ALL_TASKS_COMPLETE.md docs/archive/
git mv TESTING_GUIDE_PHASE1_AND_2.md docs/archive/
```

### 3. Move Old Roadmaps ✅

```bash
# Keep only most current roadmap in root
git mv ROADMAP_TO_1.0.0.md docs/archive/
git mv ROADMAP_TO_3.0.0-hexagonal.md docs/archive/
git mv ROADMAP_TO_3.0.0-moreinfo.md docs/archive/
git mv ROADMAP_TO_ALPHA_2.0.1.md docs/archive/
```

### 4. Update .gitignore ✅

Added built binaries:
```
/ts2go
/ts2go-web
```

### 5. Add Deprecation Notice to internal/ ✅

Created `internal/README.md` with deprecation warning

### 6. Remove Build Artifacts ✅

```bash
rm -f ts2go ts2go-web
```

### 7. Verification ✅

- ✅ `make build` works
- ✅ Tests pass (pkg/core/domain, pkg/core/services)
- ✅ Archive created with README
- ✅ Root cleaned: 28 → 20 files (30% reduction)

## 📊 Expected Results (AFTER)

### Root Directory
```
✅ README.md              - Main entry point
✅ LICENSE                - Legal
✅ CHANGELOG.md           - History
✅ CONTRIBUTING.md        - How to contribute
✅ CURRENT_STATE.md       - Current status (canonical)
✅ PROJECT_STATUS.md      - Alternative status view
✅ SPEC.md                - Feature specification
✅ KNOWN_ISSUES.md        - User-facing issues
✅ Makefile               - Build commands
✅ Dockerfile             - Container build
✅ VERSION                - Version number
✅ go.mod, go.sum, go.work - Go modules

Total: ~15 files (down from ~30)
```

### docs/archive/
```
📦 PHASE0_GO_FMT_IMPLEMENTATION.md
📦 PHASE1_AND_2_IMPLEMENTATION_COMPLETE.md
📦 ALL_TASKS_COMPLETE.md
📦 TESTING_GUIDE_PHASE1_AND_2.md
📦 ROADMAP_TO_1.0.0.md
📦 ROADMAP_TO_3.0.0-hexagonal.md
📦 ROADMAP_TO_3.0.0-moreinfo.md
📦 ROADMAP_TO_ALPHA_2.0.1.md
```

## ✅ Validation

After cleanup:
- ✅ `make build` still works
- ✅ `./ts2go convert` still works (uses internal/)
- ✅ Desktop UI not affected
- ✅ Tests still pass (pkg/core/*)
- ✅ Documentation archived with index
- ✅ Root directory: 28 → 20 files (30% reduction)

## 📊 Results

### Before Cleanup
```
28 files in root directory
- 10+ markdown files (roadmaps, phase docs, status)
- 2 binary executables (ts2go, ts2go-web)
- Cluttered, hard to navigate
```

### After Cleanup
```
20 files in root directory
- 8 markdown files (essential only)
- 0 binaries (moved to .gitignore)
- Clean, organized, easy to find what you need
```

### Files Remaining in Root
```
✅ README.md              - Main documentation
✅ LICENSE                - Legal
✅ CHANGELOG.md           - Version history
✅ CONTRIBUTING.md        - Contribution guide
✅ CURRENT_STATE.md       - Current status
✅ PROJECT_STATUS.md      - Project overview
✅ SPEC.md                - Feature specification
✅ KNOWN_ISSUES.md        - Known limitations
✅ CLEANUP_PLAN.md        - This document
✅ Makefile               - Build commands
✅ Dockerfile             - Container build
✅ VERSION                - Version number
✅ .gitignore             - Git config
✅ .dockerignore          - Docker config
✅ go.mod, go.sum, go.work - Go modules
```

### Archived Files
See `docs/archive/README.md` for full list:
- 4 phase completion documents
- 4 historical roadmaps
- All preserved for historical reference
