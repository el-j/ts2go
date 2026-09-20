You are the plan archivist for the `virtual-agency` monorepo. Your role is to seal a completed plan by moving it out of the active registry into a permanent archive file — keeping `orchestrator.json` slim permanently.

> **When to run this**: After `/execute-task` reports "🎉 Plan PLAN-NNN complete" — all tasks in the plan are done.

## Plan to archive

$ARGUMENTS

---

## Step 1 — Validate the plan is ready to archive

1. Read `.claude/orchestrator.json`.
2. Resolve the plan ID from `$ARGUMENTS` (e.g., `PLAN-039`).
3. Look up the plan in `orchestrator.json` `plans["PLAN-NNN"]`.
   - If not found: stop — "PLAN-NNN not found in active registry. It may already be archived in .claude/plans/PLAN-NNN.json."
4. Collect all task IDs from `plans["PLAN-NNN"].taskIds` (or `plans["PLAN-NNN"].tasks`).
5. For each task ID, check `orchestrator.json` `tasks["TASK-NNN"].status`:
   - All tasks must be `done`. If any task is `todo`, `in-progress`, or `blocked`, stop:
     ```
     ❌ Cannot archive PLAN-NNN — tasks still pending:
       TASK-NNN: <title> (status: <status>)
     Complete or defer those tasks first.
     ```
   - Tasks marked `deferred` or `superseded` are acceptable — they will be included in the archive with their current status.

---

## Step 2 — Build the archive document

Construct the archive data:

```json
{
  "id": "PLAN-NNN",
  "goal": "<goal from plans entry>",
  "status": "done",
  "createdAt": "<from plan entry>",
  "completedAt": "<completedAt from plan, or current ISO 8601 if missing>",
  "archivedAt": "<current ISO 8601 timestamp>",
  "note": "<any note from plan entry, or empty string>",
  "criticalPath": ["TASK-NNN", ...],
  "tasks": {
    "TASK-NNN": { <full task object from orchestrator.json tasks map> },
    ...
  }
}
```

Include **all** task objects for this plan (from the `tasks` map in `orchestrator.json`). Include the full object — not just IDs.

---

## Step 3 — Write the archive file

Write the archive JSON to `.claude/plans/PLAN-NNN.json` (create the directory if it does not exist).

Do **not** overwrite an existing archive file. If `.claude/plans/PLAN-NNN.json` already exists, stop:

```
❌ .claude/plans/PLAN-NNN.json already exists. Archive is already sealed.
```

---

## Step 4 — Update orchestrator-index.md

Read `.claude/orchestrator-index.md`. If it does not exist, create it with the header:

```markdown
# Orchestrator Plan Index

_This file is for human reference only. Commands do not load it._

| Plan | Goal | Status | Tasks | Completed |
| ---- | ---- | ------ | ----- | --------- |
```

Append a new row for the archived plan:

```
| PLAN-NNN | <goal truncated to 70 chars> | done | <task count> | <completedAt date YYYY-MM-DD> |
```

Write the updated file.

---

## Step 5 — Trim orchestrator.json

Remove the archived plan and its tasks from `orchestrator.json`:

1. Delete `plans["PLAN-NNN"]` from the `plans` map.
2. Delete every `tasks["TASK-NNN"]` that belonged to this plan.
3. If `activePlanId === "PLAN-NNN"`, set `activePlanId = null`.
4. Set `updatedAt = <current ISO 8601 timestamp>`.
5. Add or update `notes` to mention the newly archived plan:
   `"Plans ... PLAN-NNN archived in .claude/plans/. See orchestrator-index.md for history."`
6. Write the updated `orchestrator.json`.

---

## Step 6 — Output confirmation

Print:

```
✅ PLAN-NNN archived

Archive:    .claude/plans/PLAN-NNN.json  (<task count> tasks)
Index:      .claude/orchestrator-index.md  (updated)
Active registry trimmed: orchestrator.json now <new size estimate>

orchestrator.json now contains:
  Plans:  <count remaining>
  Tasks:  <count remaining>

History: <total archived plan count> plans in .claude/plans/
```

---

## Constraints

- NEVER archive a plan with non-done tasks — validate first (Step 1).
- NEVER overwrite an existing `.claude/plans/PLAN-NNN.json`.
- NEVER modify the task `.md` files in `.claude/tasks/` — leave them as-is.
- The archive is a permanent, immutable record. Do not add interpretation or summaries.
