---
id: TASK-004
title: 'Add TS/Go syntax validation for the UI Save File hook'
status: done
priority: medium
dependencies: ['none']
estimated_effort: S
---

## Goal

Add syntax validation hooks prior to file save inside `useSaveFile.ts`.

## Context

The Vue desktop application provides `useSaveFile.ts` to manage writing TS and Go output. Line 153 indicates `TODO: Add TypeScript/Go syntax validation`.

## Scope

### Files to modify

- `desktop/ui/src/composables/useSaveFile.ts`

## Implementation

1. Add basic syntax verification checks before saving by executing a linter or a parsing regex, OR
2. Provide at least basic structure validation block before invoking the File API write commands.

## Acceptance Criteria

- [x] `make test-desktop` exits 0
- [x] Basic format validation works inside Vue Desktop Application
- [x] No `any` types introduced

## Execution Log

- Replaced stub with actual validation invoke hooks inside `desktop/ui/src/composables/useSaveFile.ts`.
- Validations use native Tauri capabilities. For TS, `invoke('transpile_code')` confirms syntactical soundness.
- For Go, `invoke('build_go_file')` serves as a dry-run check testing the output AST directly.
