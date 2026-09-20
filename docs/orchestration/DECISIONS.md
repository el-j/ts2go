# Architecture Decision Records (ADRs)

Key architectural decisions governing development in **ts2go**.

---

## ADR-001: CLI First, GUI After
- **Status**: Accepted (2026-09-20)
- **Context**: The project spans a Go CLI transpiler, desktop Tauri app, and SaaS backend. Previous efforts prioritized UI state and feature docs while core CLI transpilation tests were failing.
- **Decision**: All development must prioritize the Go CLI transpiler and pipeline correctness (Milestone 1). The GUI and Web API are secondary consumers of the core library and will be addressed in Milestone 4.
- **Consequences**: No pull requests or changes for GUI aesthetics will be merged if CLI tests or coverage goals are not met.

---

## ADR-002: Hexagonal Architecture as Standard
- **Status**: Accepted
- **Context**: Decoupling the AST parsing, mapping, optimization, and code generation logic from presentation drivers (CLI, Desktop IPC, Web REST).
- **Decision**:
  - `core/domain`: Pure data models and entities (`ProjectState`, `TranspilationResult`, `File`).
  - `core/ports`: Interfaces defining primary (driving) and secondary (driven) boundaries.
  - `core/services`: Pure business logic implementing ports.
  - `adapters/driven`: Infrastructure implementations (e.g. `filesystem`, `gocompiler`, `persistence`).
  - `adapters/driving`: User entry points (e.g. `cli`, `web`).
- **Consequences**: Business logic must never import concrete adapters. All dependencies must be injected via interfaces.

---

## ADR-003: Model-Agnostic Agentic Orchestration
- **Status**: Accepted (2026-09-20)
- **Context**: Multiple AI assistants (Claude, Copilot, Gemini/Antigravity, etc.) are used across sessions. Provider-specific planning formats created divergence and unmaintained artifacts.
- **Decision**:
  - GitHub Issues and Milestones serve as the canonical external tracking authority.
  - `docs/orchestration/` serves as the canonical local tracking authority.
  - All AI assistants must adhere to [AGENT_PLAYBOOK.md](AGENT_PLAYBOOK.md).
- **Consequences**: Workflow consistency is maintained regardless of which AI model or IDE is active.

---

## ADR-004: Strict Quality Guardrails and Pre-Commit Hooks
- **Status**: Accepted (2026-09-20)
- **Context**: Code formatting, linting errors, and type warnings accumulated without gatekeeping.
- **Decision**:
  - A pre-commit hook is installed at `.git/hooks/pre-commit`.
  - Commits with unformatted Go code (`gofmt -s`), unformatted frontend code (`prettier`), `go vet` errors, or `vue-tsc` type errors are blocked automatically.
- **Consequences**: No malformed or type-broken code can enter the git history.
