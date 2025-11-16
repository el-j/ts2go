# Internal Package - DEPRECATED

**⚠️ DEPRECATION NOTICE**

This `internal/` directory contains the **legacy pre-hexagonal architecture** implementation. It is maintained for backward compatibility but **should not be used for new development**.

## Status

- **Deprecated:** Yes
- **Reason:** Replaced by hexagonal architecture in `pkg/core/` and `pkg/adapters/`
- **Still Used By:**
  - `ts2go convert` command (single-file transpilation)
  - `examples/multipackage-demo/` (legacy example)
  - Integration tests in `tests/`

## Migration Path

**For New Code:** Use the hexagonal architecture:

```
pkg/
├── core/
│   ├── domain/      # Business entities
│   ├── ports/       # Interfaces
│   └── services/    # Business logic
└── adapters/
    ├── driven/      # Infrastructure (filesystem, compiler, database)
    └── driving/     # UI/API (CLI, web)
```

**For Existing Code:** The legacy `convert` command still works but will be replaced in future versions.

## Why Keep This?

1. **Backward Compatibility** - Existing users depend on `convert` command
2. **Examples** - Demo projects still reference these packages
3. **Tests** - Integration tests use this code
4. **Reference** - Shows what was replaced during refactor

## Future Plans

- **v0.3.x:** Mark `convert` command as deprecated
- **v1.0.0:** Remove `internal/` entirely, migrate all functionality to hexagonal architecture
- **Migration Guide:** Will be provided in `docs/MIGRATION_GUIDE.md`

## Package Overview

### internal/transpiler/
Original transpilation engine (AST parsing, code generation)

### internal/analyzer/
TypeScript code analysis (imports, dependencies)

### internal/mapper/
NPM to Go package mapping

### internal/optimizer/
Code optimization passes

### internal/orchestrator/
Multi-package coordination

### internal/project/
Project structure management

### internal/module/
Module dependency resolution

## See Also

- [Hexagonal Architecture Documentation](../docs/HEXAGONAL_ARCHITECTURE_COMPLETE.md)
- [Phase 7: Cleanup](../docs/PHASE7_CLEANUP_COMPLETE.md)
- [Current State](../CURRENT_STATE.md)
