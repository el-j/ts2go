# Agentic Orchestration & Mastermind System

Welcome to the **ts2go** Agentic Orchestration Hub. This directory is the canonical local source of truth for planning, tracking, task execution, and architectural decisions across **any AI model, agent, or human developer** (Google Antigravity/Gemini, Anthropic Claude, GitHub Copilot, OpenAI, Cursor, etc.).

---

## 🎯 Core Operating Principles

Regardless of which AI model or tool is driving development, the following standards are **strictly non-negotiable**:

1. **CLI First, GUI After**:
   - The Go CLI and core transpilation engine are the bedrock of the project.
   - No GUI features or visual improvements take priority over CLI stability, core transpilation correctness, and green Go tests.
2. **Definition of Done (DoD)**:
   - **100% Test Coverage**: All packages must have unit tests; no 0% coverage packages allowed.
   - **Zero Code Gaps or Stub Code**: No `placeholderMapper` returning empty strings or dummy implementations.
   - **Zero Masked Failures**: No `|| true` or `continue-on-error: true` in CI or scripts.
   - **Zero Warnings / Zero Type Issues**: `make fmt-check`, `make type-check`, and `make lint` must exit code 0.
   - **100% In-Code Documentation**: All exported types, interfaces, functions, and packages must have godoc comments.
   - **E2E Validated**: Every transpiled feature must build with `go build` and execute with `go run`, matching Node.js output.
3. **Dual Synchronization (GitHub + Local Docs)**:
   - **GitHub Tracker**: Canonical external source of truth: [el-j/ts2go Issues & Milestones](https://github.com/el-j/ts2go/issues).
   - **Local Mastermind**: Canonical local source of truth: [ROADMAP.md](ROADMAP.md) and [HISTORY.md](HISTORY.md).
   - Every completed task must update both GitHub (`gh issue close <id>`) and local documentation.

---

## 📂 Mastermind Directory Map

| Document | Purpose |
| :--- | :--- |
| **[ROADMAP.md](ROADMAP.md)** | **The Living Mastermind**: All 4 Milestones, 15 Issues, dependencies, progress bars, and status. |
| **[AGENT_PLAYBOOK.md](AGENT_PLAYBOOK.md)** | **Standard Operating Procedure**: Step-by-step instructions for any AI model picking up work. |
| **[HISTORY.md](HISTORY.md)** | **Orchestration Log**: Detailed chronological record of all phases, past AI work, audit conclusions, and milestones. |
| **[DECISIONS.md](DECISIONS.md)** | **Architecture Decision Records (ADRs)**: Key technical constraints and architectural decisions. |

---

## 🔄 How to Work with This System (Any AI Agent)

When starting or continuing a session:
1. **Locate Active Work**: Open `docs/orchestration/ROADMAP.md` and check the currently active Milestone and Issue.
2. **Review GitHub Issue**: Run `gh issue view <issue_number>` to get the exact problem description and acceptance criteria.
3. **Implement Following DoD**: Follow the instructions in `docs/orchestration/AGENT_PLAYBOOK.md`.
4. **Enforce Local Quality Gate**:
   ```bash
   make fmt         # Format Go & Frontend code
   make fmt-check   # Verify formatting
   make type-check  # Verify TypeScript
   make test-go     # Run Go tests
   ```
5. **Mark Done**: Close the issue with `gh issue close <issue_number> -c "..."` and update `docs/orchestration/ROADMAP.md` and `docs/orchestration/HISTORY.md`.
