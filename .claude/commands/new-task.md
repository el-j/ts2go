You are a task-definition writer for the **figma-vue-bridge** monorepo. Given the description below, produce one precise, self-contained task file and register it in `.claude/orchestrator.json`.

## Task description

$ARGUMENTS

---

## Step 1 — Load orchestrator state and determine the next task ID

Read `.claude/orchestrator.json`.

- If the file does not exist, initialise it with the empty structure:
  ```json
  {
    "version": "1",
    "updatedAt": "<ISO 8601>",
    "activePlanId": null,
    "plans": {},
    "tasks": {}
  }
  ```
- Find the highest existing `TASK-NNN` key in the `tasks` map. The new task ID is that number plus one, zero-padded to three digits (e.g., if the highest is `TASK-007`, create `TASK-008`).
- If `tasks` is empty, start at `TASK-001`.
- If there is an `activePlanId`, note it — the new task will be linked to that plan. If there is no active plan, set `planId` to `null` in the registry entry.

---

## Step 2 — Gather codebase context

Before writing the task, read the files most relevant to the description. Use CLAUDE.md as your map. At minimum:

- If the task involves a new shared type or schema → read `packages/shared/src/types/` and `packages/shared/src/schemas/` to check if it already exists.
- If the task involves the CLI → read `packages/cli/src/core/` for the relevant service(s).
- If the task involves the API → read `packages/api/src/routes/` and `packages/api/src/services/`.
- If the task involves the Web UI → read `packages/web-ui/src/views/` and `packages/web-ui/src/stores/`.
- If the task involves the Figma Plugin → read `packages/figma-plugin/src/code.ts` and `packages/figma-plugin/src/extractors.ts`.
- If the task involves the VS Code extension → read `packages/vscode-extension/src/extension.ts`.

Do **not** skip this. The task file must reference real file paths and real symbol names.

---

## Step 3 — Write the task file

Create `.claude/tasks/TASK-NNN.md` using this **exact** template:

```
---
id: TASK-NNN
title: "<concise imperative title, max 60 chars>"
status: todo
priority: <critical|high|medium|low>
dependencies: [<TASK-NNN, ...> or "none"]
estimated_effort: <XS 15min | S 30min | M 1h | L 2h | XL 4h+>
---

## Goal

One sentence: what this task achieves and why it is needed.

## Context

Key facts the executing agent must know (type names, service names, existing patterns to follow, pitfalls to avoid). Be specific. Reference exact file paths and the correct path aliases.

## Scope

### Files to modify
- `path/to/file.ts` — what change and why

### Files to create
- `path/to/new-file.ts` — what it contains

### Files to delete
- `path/to/old-file.ts` — reason

## Implementation

Step-by-step instructions in order. Be explicit about:
- Which package/path alias to use for imports
- Which existing types/schemas to import from `@figma-vue-bridge/shared`
- Which existing utilities/services to reuse
- What NOT to do

## Acceptance Criteria

- [ ] `npm run build` exits 0 with no new TypeScript errors
- [ ] `npm run test --workspace=@figma-vue-bridge/<pkg>` exits 0
- [ ] <specific functional criterion 1>
- [ ] <specific functional criterion 2>
- [ ] No `any` types introduced
```

---

## Step 4 — Register in orchestrator.json

After saving the task file, update `.claude/orchestrator.json`:

1. Add an entry to `tasks`:

```json
"TASK-NNN": {
  "id": "TASK-NNN",
  "planId": "<activePlanId or null>",
  "title": "<title from task file>",
  "status": "todo",
  "priority": "<priority>",
  "dependencies": ["TASK-NNN"],
  "estimatedEffort": "<XS|S|M|L|XL>",
  "createdAt": "<ISO 8601 timestamp>",
  "startedAt": null,
  "completedAt": null,
  "filesChanged": [],
  "buildStatus": null,
  "testStatus": null,
  "notes": ""
}
```

2. If there is an `activePlanId`, append `TASK-NNN` to that plan's `taskIds` array.
3. Set `"updatedAt"` to the current ISO 8601 timestamp.
4. Write the full updated JSON back to `.claude/orchestrator.json`.

---

## Step 5 — Output confirmation

Print:

```
Created:      .claude/tasks/TASK-NNN.md
Registered:   .claude/orchestrator.json → tasks.TASK-NNN
Title:        <title>
Priority:     <priority>
Effort:       <effort>
Dependencies: <list or none>
Plan:         <PLAN-NNN or unlinked>

Run `/execute-task TASK-NNN` to implement this task.
```

---

## Constraints

- One task file only — do not create multiple tasks.
- If the description is too vague to write a scoped task, state what information is missing and stop. Do not invent scope.
- NEVER add `any` types to implementation instructions.
- NEVER point shared types to `packages/api` or `packages/cli` — always `packages/shared/src/`.
- NEVER use `npx vitest` — always `npm run test`.
- NEVER use `console.log` in implementation instructions — always the `logger` from `@figma-vue-bridge/shared`.
- ALWAYS check existing `tasks` in `orchestrator.json` before creating a duplicate — task descriptions that are already covered must be rejected with an explanation.
- ALWAYS follow project coding standards from CLAUDE.md.
