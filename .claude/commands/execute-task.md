You are a precision implementation agent for the **figma-vue-bridge** monorepo. Your role is to implement the task identified below — completely, correctly, and verifiably — following all project rules without deviation.

## Task to execute

$ARGUMENTS

---

## Phase 1 — Load and validate the task

1. Read `.claude/orchestrator.json`.
2. Resolve the task ID:
   - If `$ARGUMENTS` is a task ID (e.g., `TASK-007`), use it directly.
   - If `$ARGUMENTS` is a file path, extract the ID from the filename.
3. Look up the task entry in `orchestrator.json` under `tasks["TASK-NNN"]`.
   - If the entry does not exist, check whether `.claude/tasks/TASK-NNN.md` exists. If neither exists, stop: "Task TASK-NNN not found."
   - If the task file exists but is NOT in `orchestrator.json`, read the task file, register it in `orchestrator.json` (status: `todo`, all nullable fields null), and continue.
4. Read the full task file at `.claude/tasks/TASK-NNN.md`.
5. Check status:
   - `done` → stop: "TASK-NNN is already marked done. Nothing to do."
   - `blocked` → stop: "TASK-NNN is blocked. Read .claude/tasks/TASK-NNN.md ## Blocked section for details."
6. Check `dependencies`. For each dependency ID, verify its `status` in `orchestrator.json` is `done`. If any dependency is not done, stop and list which tasks must be completed first.
7. Set status to `in-progress`:
   - In the task file front-matter: `status: in-progress`.
   - In `orchestrator.json`: `tasks["TASK-NNN"].status = "in-progress"`, `startedAt = <current ISO 8601 timestamp>`, `updatedAt = <current ISO 8601 timestamp>`.
   - Write the updated `orchestrator.json`.

---

## Phase 2 — Baseline verification

Run `npm run build 2>&1 | tail -10` and confirm it exits 0 before making any changes. If the baseline is broken, report the errors and stop — do not implement on top of a broken build.

---

## Phase 3 — Deep context read

Read **every file** listed in the task's `## Scope` section before writing a single line of code. Additionally:

- For any file you intend to **modify**: read it completely so your edits are consistent with existing code style, imports, and patterns.
- For any file you intend to **create**: read adjacent files in the same directory to match naming, import style, and export conventions.
- For shared types: confirm the exact import path from `@figma-vue-bridge/shared` — never guess.
- For CLI services: read `packages/cli/src/core/` to understand existing patterns.
- For API routes: read `packages/api/src/routes/` to understand the route + Zod validation pattern.
- For web-ui: read adjacent `.vue` components and Pinia stores.

Do not begin implementation until all relevant files have been read.

---

## Phase 4 — Implementation

Follow the `## Implementation` steps in the task file in the exact order given. While implementing, enforce these project rules absolutely — they are non-negotiable:

### TypeScript rules

- **NEVER** use `any` — use `unknown` with explicit narrowing, or Zod validation.
- **ALWAYS** provide explicit return types on exported functions.
- **ALWAYS** use the correct path aliases:
  - `packages/cli/src/`: `@core/`, `@config/`, `@generators/`, `@transformers/`, `@mappers/`, `@parsers/`, `@cli/`
  - `packages/web-ui/src/`: `@/`
  - Cross-package: `@figma-vue-bridge/shared`, `@figma-vue-bridge/cli`

### Shared types / schemas rules

- **ALWAYS** define shared types in `packages/shared/src/types/` or `packages/shared/src/schemas/`.
- **ALWAYS** export new shared items from `packages/shared/src/index.ts`.
- **NEVER** define shared types inline in `packages/api` or `packages/cli`.

### API rules

- **ALWAYS** use inline Zod validation in route handlers (not generic middleware).
- **ALWAYS** return `{ success: true, data: T }` or `{ success: false, error: { code, message } }`.
- **ALWAYS** use `createApp()` from `packages/api/src/index.ts` in tests (not a live server).
- Error codes: `VALIDATION_ERROR` → 400, `NOT_FOUND` → 404, `CLI_ERROR` | `INTERNAL_ERROR` → 500.

### CLI rules

- **ALWAYS** use pure functional transformers (no side effects).
- **ALWAYS** use `fs-extra` for file operations (`outputFile`, `pathExists`, `ensureDir`).
- **ALWAYS** use atomic writes: write to `.tmp` then rename.
- **ALWAYS** use the `Result<T,E>` pattern for error handling in core logic.
- **NEVER** throw from public service methods — return `Result` types.

### Vue / Web-UI rules

- **ALWAYS** use `<script setup lang="ts">` with Composition API.
- **ALWAYS** use Pinia stores with **setup function syntax** (`defineStore("name", () => { ... })`).
- **ALWAYS** use PrimeVue 4 components and Tailwind 4 utility classes.
- **NEVER** use Options API.

### Testing rules

- **ALWAYS** use `npm run test --workspace=@figma-vue-bridge/<pkg>` for targeted tests.
- **ALWAYS** use `vi.mock()` for external dependencies in unit tests.
- **ALWAYS** use `createApp()` from `packages/api/src/index.ts` for API integration tests.
- For web-ui component tests: use `@vue/test-utils` + `happy-dom`.

### Logging rules

- **NEVER** use `console.log`, `console.error`, or `console.warn` in production code — import and use `logger` from `@figma-vue-bridge/shared`.

### General rules

- **NEVER** use `sleep` or arbitrary delays.
- **NEVER** create unnecessary documentation markdown files.
- **ALWAYS** use descriptive, full variable and function names — no abbreviations.

---

## Phase 5 — Verification

Run these commands in order. Wait for each to complete before running the next.

### 5a. Build

```
npm run build 2>&1 | grep -E "(error|Error|success|✓)" | head -40
```

The build **must** exit 0 with no TypeScript errors. If errors exist, fix them before proceeding.

### 5b. Tests

Run the most targeted test command that covers the changed code:

```
npm run test --workspace=@figma-vue-bridge/<affected-package>
```

If no specific package is obvious, run the full suite:

```
npm run test
```

All tests **must** pass. If tests fail due to your changes, fix them before proceeding. If pre-existing tests fail (unrelated to your changes), note them explicitly but continue.

---

## Phase 6 — Mark task complete

Only after both build and tests are green:

1. Update the task file front-matter: `status: in-progress` → `status: done`.
2. Append a `## Execution Log` section to the bottom of the task file:

```markdown
## Execution Log

- **Completed:** <ISO date>
- **Files changed:** <comma-separated list of files actually modified/created/deleted>
- **Build:** ✅ clean
- **Tests:** ✅ passed (or note any pre-existing failures)
- **Notes:** <anything relevant for future reference, or "none">
```

3. Update `orchestrator.json`:
   - `tasks["TASK-NNN"].status = "done"`
   - `tasks["TASK-NNN"].completedAt = <current ISO 8601 timestamp>`
   - `tasks["TASK-NNN"].filesChanged = [<array of all files actually changed>]`
   - `tasks["TASK-NNN"].buildStatus = "clean"`
   - `tasks["TASK-NNN"].testStatus = "passed"` (or `"passed-with-preexisting-failures"` if applicable)
   - `tasks["TASK-NNN"].notes = <any relevant notes, or empty string>`
   - `updatedAt = <current ISO 8601 timestamp>`
4. Check if all tasks in the plan's `taskIds` array are now `done`. If so:
   - Set `plans["PLAN-NNN"].status = "done"`
   - Set `activePlanId = null` (unless another plan is still `in-progress`)
5. Write the full updated `orchestrator.json`.

---

## Phase 7 — Output a completion report

Print to stdout:

```
✅ TASK-NNN complete — <title>

Files changed:
  modified: path/to/file.ts
  created:  path/to/new-file.ts

Build: ✅ clean
Tests: ✅ <N tests passed>

Plan PLAN-NNN progress: <N done> / <total> tasks
  <TASK-NNN> ✅ done
  <TASK-NNN> ⏳ todo   ← next
  <TASK-NNN> ⏳ todo

State saved to .claude/orchestrator.json

Next task: TASK-NNN (<title>) — run `/execute-task TASK-NNN`
         or: 🎉 Plan PLAN-NNN is complete — all tasks done.
```

---

## Failure protocol

If at any point you cannot complete the task (blocked by missing information, a dependency, or an unfixable error):

1. Set `status: blocked` in the task file front-matter.
2. Append a `## Blocked` section to the task file explaining exactly what is missing and what must happen to unblock it.
3. Update `orchestrator.json`:
   - `tasks["TASK-NNN"].status = "blocked"`
   - `tasks["TASK-NNN"].notes = "<concise reason>"`
   - `updatedAt = <current ISO 8601 timestamp>`
   - Write the updated file.
4. Print a clear error report:
   ```
   ❌ TASK-NNN blocked — <title>
   Reason: <concise reason>
   To unblock: <what must happen>
   State saved to .claude/orchestrator.json
   ```
5. **Do not** leave partial or broken code in the repository. Revert any incomplete changes.
