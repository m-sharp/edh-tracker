---
phase: 05-pod-deck-ux
plan: "07"
subsystem: backend-api, frontend-http
tags: [commander, api, http, bug-fix]
dependency_graph:
  requires: []
  provides: [commander-create-returns-id]
  affects: [deck/new, deck/SettingsTab]
tech_stack:
  added: []
  patterns: [json-body-201-response, typed-http-function]
key_files:
  created: []
  modified:
    - lib/routers/commander.go
    - app/src/http.ts
    - app/src/routes/deck/new/index.tsx
    - app/src/routes/deck/SettingsTab.tsx
decisions:
  - CommanderCreate mirrors DeckCreate pattern exactly (capture ID, set Content-Type, WriteHeader 201, json.NewEncoder Encode)
  - PostCommander mirrors PostDeck/PostPod pattern (throw on non-ok, return parsed JSON)
  - Caller sites simplified to direct destructuring since PostCommander now throws internally
metrics:
  duration: 2min
  completed: "2026-03-26"
  tasks: 2
  files: 4
---

# Phase 05 Plan 07: Commander Creation Round-Trip Fix Summary

Fixed the commander creation round-trip: `POST /api/commander` now returns `{"id": N}` with 201, and `PostCommander` in http.ts returns `Promise<{id: number}>` — ending the SyntaxError crash when creating new commanders in the deck form.

## Tasks Completed

| Task | Description | Commit | Files |
|------|-------------|--------|-------|
| 1 | CommanderCreate captures ID and writes JSON body | 6adaf21 | lib/routers/commander.go |
| 2 | PostCommander returns Promise<{id: number}> and update callers | 9a7bfe1 | app/src/http.ts, deck/new/index.tsx, SettingsTab.tsx |

## What Was Built

**Backend (`lib/routers/commander.go`):** `CommanderCreate` now captures the int returned by `c.commanders.Create`, sets `Content-Type: application/json`, writes HTTP 201, and encodes `{"id": id}` using `json.NewEncoder(w).Encode`. This mirrors the existing `DeckCreate` pattern exactly.

**Frontend (`app/src/http.ts`):** `PostCommander` replaced with a typed async function returning `Promise<{ id: number }>` — same pattern as `PostDeck` and `PostPod`. Internal `res.ok` check throws on failure instead of returning a raw `Response`.

**Caller updates:**
- `app/src/routes/deck/new/index.tsx`: `handleCommanderSelect` simplified — `res.ok` guard and `res.json()` call removed; now `const { id } = await PostCommander(newName)` directly.
- `app/src/routes/deck/SettingsTab.tsx`: both primary commander and partner commander `onChange` handlers updated with the same simplification.

## Deviations from Plan

None — plan executed exactly as written.

## Verification

- `go vet ./lib/routers/...`: pre-existing error in `lib/business/business.go:91` (unrelated to this plan, in another parallel agent's domain). Commander router changes are syntactically correct and mirror the established pattern.
- `cd app && ./node_modules/.bin/tsc --noEmit`: passes with no errors.

## Self-Check: PASSED

- `lib/routers/commander.go` modified and committed (6adaf21)
- `app/src/http.ts` modified and committed (9a7bfe1)
- `app/src/routes/deck/new/index.tsx` modified and committed (9a7bfe1)
- `app/src/routes/deck/SettingsTab.tsx` modified and committed (9a7bfe1)
