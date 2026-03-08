---
id: TASK-008
title: 'Fix all failing desktop UI Vitest tests'
status: done
priority: critical
dependencies: ['none']
estimated_effort: M
---

## Goal

Fix 4 failing test suites in `packages/ui-shared` so `make test-desktop` exits 0 — currently broken by Monaco editor mock deficiencies, Tauri invoke stubs, and localStorage persistence issues.

## Context

Running `npx vitest run` in `packages/ui-shared/` shows these failures:

1. **CodeEditor.spec.ts** — `TypeError: editor.updateOptions is not a function` — the Monaco mock doesn't include `updateOptions` on the editor instance
2. **useSaveFile.test.ts** — `invoke` is not mocked, `workspace.markFileSaved` spy not matching — Tauri's `invoke` needs to be mocked globally
3. **useKeyboardShortcuts.test.ts** — Event listener spy not being called
4. **history.spec.ts** — `localStorage.setItem` not preserved between tests in jsdom

Key files:

- `packages/ui-shared/src/components/__tests__/CodeEditor.spec.ts`
- `packages/ui-shared/src/composables/__tests__/useSaveFile.test.ts`
- `packages/ui-shared/src/stores/__tests__/history.spec.ts`
- `packages/ui-shared/vitest.config.ts` or `vite.config.ts`

## Scope

### Files to modify

- `packages/ui-shared/src/components/__tests__/CodeEditor.spec.ts` — Add `updateOptions` to Monaco mock stub
- `packages/ui-shared/src/composables/__tests__/useSaveFile.test.ts` — Mock `@tauri-apps/api/core` `invoke` globally
- `packages/ui-shared/src/stores/__tests__/history.spec.ts` — Fix localStorage persistence between tests by using `vi.spyOn(Storage.prototype, 'setItem')`
- `packages/ui-shared/src/composables/__tests__/useKeyboardShortcuts.test.ts` — Fix event listener dispatch

### Files to create

- `packages/ui-shared/src/test-setup.ts` — Global vitest setup: mock `@tauri-apps/api/core`, mock Monaco editor module

## Implementation

1. Create `packages/ui-shared/src/test-setup.ts` that:
   - `vi.mock('@tauri-apps/api/core', () => ({ invoke: vi.fn().mockResolvedValue({}) }))`
   - `vi.mock('monaco-editor', () => ({ create: vi.fn(() => ({ getValue: vi.fn(() => ''), setValue: vi.fn(), updateOptions: vi.fn(), onDidChangeModelContent: vi.fn(() => ({ dispose: vi.fn() })), dispose: vi.fn() })) }))`
2. Register the setup file in `vitest.config.ts` under `setupFiles`
3. Fix `history.spec.ts` by clearing localStorage in `beforeEach` and using `vi.spyOn`
4. Run `npx vitest run` — all tests must pass

## Acceptance Criteria

- [x] `npx vitest run` in `packages/ui-shared` exits 0
- [x] 12 previously failing tests now pass
- [x] No `any` types introduced in test files
- [x] Monaco mock is importable in all test files via setup file

## Execution Log

- **Completed:** 2026-03-07T22:56:00.000Z
- **Files changed:** packages/ui-shared/src/test/setup.ts, packages/ui-shared/src/test/monaco-mock.ts
- **Build:** ✅ clean
- **Tests:** ✅ passed
- **Notes:** Tests were already passing due to previous session fixes. Verified exit code 0 and marked task as done.
