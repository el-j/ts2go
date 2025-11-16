# Project Restructure - Migration Guide

**Date:** November 16, 2025  
**Status:** ✅ Complete

## What Changed?

The project has been restructured for clarity and maintainability. The new structure separates Go, Desktop, and documentation into clear top-level directories.

## New Structure

### Before (Cluttered)
```
ts2go/
├── pkg/core/              # Hexagonal architecture
├── pkg/adapters/          # Adapters
├── pkg/cli/               # CLI (was this an adapter?)
├── cmd/                   # Binaries
├── runtime/               # Runtime library
├── internal/              # Legacy code
├── desktop-ui/            # Desktop app
│   └── src-tauri/         # Tauri (hidden inside!)
├── docs/                  # Documentation
├── examples/              # Examples
├── CURRENT_STATE.md       # Many root docs
├── PROJECT_STATUS.md
├── SPEC.md
└── ... 20+ other files
```

### After (Clean)
```
ts2go/
├── go/                    # 🔵 All Go code
│   ├── core/              # Hexagonal architecture
│   ├── adapters/          # Infrastructure & UI adapters
│   ├── cmd/               # Binaries
│   ├── runtime/           # Runtime library
│   ├── internal/          # Legacy (deprecated)
│   ├── examples/          # Go examples
│   └── go.mod             # Single module
│
├── desktop/               # 🟢 Desktop application
│   ├── ui/                # Vue.js frontend
│   └── tauri/             # Rust backend
│
├── docs/                  # 📚 All documentation
│   ├── README.md          # Docs index
│   ├── guides/            # User guides
│   ├── development/       # Developer docs
│   └── archive/           # Historical docs
│
├── bin/                   # 🔨 Compiled binaries
├── scripts/               # 🔧 Build scripts
├── Makefile               # 🎯 Main orchestration
├── README.md              # Project overview
├── CHANGELOG.md           # Release notes
├── CONTRIBUTING.md        # How to contribute
├── LICENSE                # Legal
└── VERSION                # Version number
```

## Migration Steps for Contributors

### 1. Update Your Local Repository

```bash
# Fetch latest changes
git pull origin main

# Your local work is safe - Git tracked all moves
# Verify new structure
ls -F
# Should see: go/ desktop/ docs/ bin/ ...
```

### 2. Update Import Paths (If Contributing Go Code)

**Old imports:**
```go
import "github.com/el-j/ts2go/pkg/core/domain"
import "github.com/el-j/ts2go/internal/transpiler"
```

**New imports:**
```go
import "github.com/el-j/ts2go/core/domain"
import "github.com/el-j/ts2go/internal/transpiler"
```

Note: Module path is still `github.com/el-j/ts2go`, just files moved to `go/` directory.

### 3. Update Build Commands

**Old commands:**
```bash
go build -o ts2go ./cmd/ts2go
cd desktop-ui && npm run tauri:dev
```

**New commands:**
```bash
make build-go              # Build CLI
make dev-desktop           # Run desktop in dev mode
make test                  # Run all tests
make help                  # See all commands
```

### 4. Update File Paths

| Old Location | New Location |
|--------------|--------------|
| `pkg/core/` | `go/core/` |
| `pkg/adapters/` | `go/adapters/` |
| `pkg/cli/` | `go/cli/` (or `go/adapters/driving/cli/`) |
| `cmd/` | `go/cmd/` |
| `runtime/` | `go/runtime/` |
| `internal/` | `go/internal/` |
| `desktop-ui/` | `desktop/ui/` |
| `desktop-ui/src-tauri/` | `desktop/tauri/` |
| `examples/` | `go/examples/` |
| `CURRENT_STATE.md` | `docs/CURRENT_STATE.md` |
| `PROJECT_STATUS.md` | `docs/PROJECT_STATUS.md` |
| `SPEC.md` | `docs/SPEC.md` |

### 5. Desktop Development

**Old:**
```bash
cd desktop-ui
npm install
npm run tauri:dev
```

**New:**
```bash
cd desktop
make dev          # or: make build
```

Or from root:
```bash
make dev-desktop
```

## Breaking Changes

### For External Users

**If you imported ts2go as a Go module:**
- Module path unchanged: `github.com/el-j/ts2go`
- Import paths unchanged (we didn't change the module name)
- No action needed! ✅

### For Contributors

**If you have open PRs:**
1. Rebase your branch on latest main
2. Update any file paths in your changes
3. Run `make test` to verify

**If you have local branches:**
```bash
git fetch origin
git rebase origin/main
# Resolve any conflicts (mostly file moves)
```

## Benefits

### For New Contributors
- ✅ Immediately understand project structure
- ✅ Know where Go code lives (`go/`)
- ✅ Know where Desktop app lives (`desktop/`)
- ✅ Clear separation of concerns

### For Existing Contributors
- ✅ Easier navigation
- ✅ Better module boundaries
- ✅ Simpler build commands (`make help`)
- ✅ Obvious where new code belongs

### For Maintainers
- ✅ Professional structure (industry standard)
- ✅ Scalable (easy to add new components)
- ✅ Better documentation organization
- ✅ Single Go module (simpler dependency management)

## FAQ

### Q: Why did you change the structure?

**A:** The old structure was cluttered and confusing. Mixing Go, TypeScript, Rust, and docs at the same level made it hard for newcomers to understand the project. The new structure follows industry best practices.

### Q: Do I need to update my code?

**A:** If you're using ts2go as a CLI tool or importing it as a module, **no changes needed**. If you're contributing code, you may need to update file paths in your PRs.

### Q: Where did `pkg/` go?

**A:** Moved to `go/` directory. We consolidated all Go code under one directory for clarity.

### Q: Where did `desktop-ui/` go?

**A:** Renamed to `desktop/` and restructured:
- UI code: `desktop/ui/`
- Tauri backend: `desktop/tauri/`

### Q: How do I build now?

**A:** Use the Makefile:
```bash
make build-go        # Build CLI
make build-desktop   # Build desktop app
make test            # Run tests
make help            # See all commands
```

### Q: Are there any API changes?

**A:** No! The public API is unchanged. Only internal structure changed.

### Q: What about `internal/`?

**A:** Still there in `go/internal/`, marked as deprecated. Will be removed in v1.0.0 after full migration to hexagonal architecture.

## Support

- **Issues:** Open a GitHub issue
- **Questions:** Check [docs/](docs/) or ask in discussions
- **Contributing:** See [CONTRIBUTING.md](../CONTRIBUTING.md)

---

**Migration completed:** November 16, 2025  
**Backup branch:** `backup-before-restructure-20251116`
