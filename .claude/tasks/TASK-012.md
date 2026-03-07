---
id: TASK-012
title: 'Fix TypeScript compiler errors in packages/ui-shared'
status: todo
priority: high
dependencies: ['TASK-008']
estimated_effort: M
---

## Goal

Resolve all TypeScript compilation errors in `packages/ui-shared` so `npm run build` and `npm run type-check` exit 0, enabling production builds.

## Context

`make build-desktop` is currently skipped with `⚠️  Desktop build skipped (TypeScript errors)`. Running `cd packages/ui-shared && npm run build` fails. Common TS errors in Vue 3 + Tauri apps include:

- Missing `@tauri-apps/api` type declarations
- Implicit `any` on event handlers
- Unused import/variable warnings promoted to errors
- Vue SFC `defineProps` typing issues

Run `cd packages/ui-shared && npx vue-tsc --noEmit 2>&1` to see the full error list before implementing.

## Scope

### Files to modify

- Various Vue components in `packages/ui-shared/src/` — add explicit types
- `packages/ui-shared/tsconfig.json` — ensure correct lib targets and strictness

### Files to create

- `packages/ui-shared/src/types/tauri.d.ts` — ambient declarations for `invoke` if `@tauri-apps/api` types are incomplete

## Implementation

1. Run `npx vue-tsc --noEmit` to capture all errors
2. Fix each error class in order:
   - Missing type on store actions → add explicit `void` return types
   - Implicit `any` → add type annotations
   - Unused imports → remove
   - Component prop type errors → add `withDefaults(defineProps<>(), {})` patterns
3. Do NOT downgrade `strict: false` — fix the actual types
4. Run `npm run build` — must exit 0

## Acceptance Criteria

- [ ] `cd packages/ui-shared && npx vue-tsc --noEmit` exits 0
- [ ] `npm run build` exits 0 (produces `dist-web/`)
- [ ] `npm run build:tauri` exits 0 (produces `dist-tauri/`)
- [ ] No `any` types introduced to fix TS errors
