You are a task-definition writer for the `virtual-agency` monorepo. Given the description below, produce one precise, self-contained task file and register it in `.claude/orchestrator.json`.

> **v2 Architecture**: `orchestrator.json` is a slim active registry. Use `counters.nextTaskId` for the new task ID — do **not** scan task keys to compute the max.

## Task description

$ARGUMENTS

---

## Step 1 — Load orchestrator state and determine the next task ID

Read `.claude/orchestrator.json`.

- If the file does not exist, initialise it with the v2 empty structure:
  ```json
  {
    "version": "2",
    "updatedAt": "<ISO 8601>",
    "counters": { "nextTaskId": 1, "nextPlanId": 1 },
    "activePlanId": null,
    "notes": "",
    "plans": {},
    "tasks": {}
  }
  ```
- Read `counters.nextTaskId` — this is the ID number to use (zero-padded to 3 digits, e.g. `204` → `TASK-204`).
- If there is an `activePlanId`, note it — the new task will be linked to that plan.
  If there is no active plan, set `planId` to `null` in the registry entry.
- Scan the active `tasks` map keys only (not archives) to confirm no identical task already exists.

---

## Step 2 — Gather codebase context

Before writing the task, read the files most relevant to the description. At minimum:

- If the task involves a new schema → read `packages/shared/src/schemas/` to check if it already exists and to understand naming conventions.
- If the task involves a service → read `packages/core/src/services/` to understand existing patterns.
- If the task involves an agent → read `packages/core/src/agents/` to understand `BaseAgent` and existing implementations.
- If the task involves the API → read `packages/core/src/api/`.
- If the task involves the dashboard → read `apps/dashboard/src/`.
- If the task involves the API client → read `packages/api-client/src/`.

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
role: <backend|frontend|schema|devops|ui|ux|qa|verify|planning>
dependencies: [<TASK-NNN, ...> or "none"]
estimated_effort: <XS 15min | S 30min | M 1h | L 2h | XL 4h+>
---

## Goal

One sentence: what this task achieves and why it is needed.

## Context

Key facts the executing agent must know (schema names, service names, existing patterns to follow, pitfalls to avoid). Be specific. Reference exact file paths.

## Scope

### Files to modify
- `path/to/file.ts` — what change and why

### Files to create
- `path/to/new-file.ts` — what it contains

### Files to delete
- `path/to/old-file.ts` — reason

## Implementation

Step-by-step instructions in order. Be explicit about:
- Which Effect-TS patterns to use (Layer, Context.GenericTag, Effect.tryPromise, etc.)
- Which schemas to import from `@virtual-agency/shared`
- Which existing utilities/services to reuse
- What NOT to do

## Acceptance Criteria

- [ ] `npm run build` exits 0 with no new TypeScript errors
- [ ] `npm test` exits 0 (or `npm test -- <path>` for targeted run)
- [ ] <specific functional criterion 1>
- [ ] <specific functional criterion 2>
- [ ] No `any` types introduced
- [ ] No raw Promises introduced
- [ ] No try/catch introduced
```

### Role field → specialist agent

| role       | Agent              | Use for                                      |
| ---------- | ------------------ | -------------------------------------------- |
| `backend`  | Senior Developer   | Services, agents, API routes, DB, queue      |
| `frontend` | Senior Developer   | Vue components, dashboard, composables       |
| `schema`   | Senior Developer   | Effect.Schema definitions in packages/shared |
| `devops`   | Senior Developer   | Docker, CI, migrations, infra                |
| `ui`       | UI Designer        | Visual design, component styling             |
| `ux`       | UX Architect       | User flows, information architecture         |
| `qa`       | Evidence Collector | Writing tests, documenting issues            |
| `verify`   | Reality Checker    | Integration/build/test certification         |
| `planning` | Senior PM          | Breaking down ambiguous goals                |

---

## Step 4 — Register in orchestrator.json

After saving the task file, update `.claude/orchestrator.json`:

1. Add an entry to `tasks`:

```json
"TASK-NNN": {
  "id": "TASK-NNN",
  "planId": "<activePlanId or null>",
  "title": "<title from task file>",
  "role": "<role from task file>",
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

2. Increment `counters.nextTaskId` by 1.
3. If there is an `activePlanId`, append `TASK-NNN` to that plan's `taskIds` array.
4. Set `"updatedAt"` to the current ISO 8601 timestamp.
5. Write the full updated JSON back to `.claude/orchestrator.json`.

---

## Step 5 — Output confirmation

Print:

```
Created:      .claude/tasks/TASK-NNN.md
Registered:   .claude/orchestrator.json → tasks.TASK-NNN (nextTaskId now NNN+1)
Title:        <title>
Role:         <role> → will use <agent persona>
Priority:     <priority>
Effort:       <effort>
Dependencies: <list or none>
Plan:         <PLAN-NNN or unlinked>

Run `/execute-task TASK-NNN` to implement this task.
```

---

## Constraints

- One task file only — do not create multiple tasks.
- If the description is too vague to write a scoped task, state what information is missing and stop.
- NEVER add `any` types, raw Promises, try/catch, or `console.log` to implementation instructions.
- NEVER point schemas to `packages/core` — always `packages/shared/src/schemas/`.
- ALWAYS assign a `role:` field.
- ALWAYS check active `tasks` in `orchestrator.json` before creating a duplicate.
