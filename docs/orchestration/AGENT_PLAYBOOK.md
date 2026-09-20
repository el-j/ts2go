# Agent Playbook: Standard Operating Procedure

This document defines the exact operational procedure that any AI model or contributor must follow when working on the **ts2go** repository.

---

## 🔁 The 5-Step Execution Cycle

Every unit of work (issue/task) must proceed through this cycle:

```
[1. Select Issue] ──▶ [2. Research & Plan] ──▶ [3. Implement TDD] ──▶ [4. Verify & Guard] ──▶ [5. Sync & Close]
```

### Step 1: Select & Claim Issue
1. Consult [ROADMAP.md](ROADMAP.md) for the active milestone.
2. Select the next unblocked issue with the highest priority (P0 before P1, P1 before P2).
3. Read the full issue context from GitHub:
   ```bash
   gh issue view <issue_number>
   ```

### Step 2: Research & Plan
1. Locate all related source files and existing tests.
2. Verify why the test is failing or where the feature gap exists.
3. Check for repository-specific patterns (e.g., hexagonal architecture, AST traversal, parser bridge).

### Step 3: Implement with Test-Driven Development (TDD)
1. Write or expand the unit/integration test first to reproduce the gap or bug.
2. Implement the minimal clean code change to satisfy the test.
3. Preserve all existing docstrings and comments; write godoc comments for any newly added symbols.
4. **Zero Stub Rule**: Never write placeholder methods returning empty values, dummy data, or unhandled errors.

### Step 4: Verify & Run Quality Gates
Before any work can be considered complete, run all verification commands locally:
```bash
make fmt          # Format all code
make fmt-check    # Verify zero unformatted files
make type-check   # Verify zero TypeScript issues
make test-go      # Run Go tests
./.git/hooks/pre-commit # Verify pre-commit hook passes
```

### Step 5: Dual Synchronization & Closing
1. Update [ROADMAP.md](ROADMAP.md) and mark the DoD check items as completed `[x]`.
2. Append a brief entry in [HISTORY.md](HISTORY.md) detailing:
   - What was changed
   - Which tests were added/fixed
   - Verification command results
3. Close the GitHub issue with a summary comment:
   ```bash
   gh issue close <issue_number> --comment "Resolved in commit <hash>. All DoD criteria verified."
   ```

---

## 🚫 Strictly Forbidden Practices
- **No Test Masking**: Never use `|| true`, `continue-on-error: true`, or commented-out test assertions.
- **No Untracked Changes**: Never make major changes without tying them to an issue or milestone.
- **No Broken Builds**: Never commit code that breaks `make fmt-check`, `make build-go`, or `make type-check`.
- **No Scope Creep**: Keep pull requests and commits atomic, focused strictly on the issue at hand.
