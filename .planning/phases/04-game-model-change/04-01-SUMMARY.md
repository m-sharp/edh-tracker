---
phase: 04-game-model-change
plan: "01"
subsystem: backend-api, frontend-types, frontend-components
tags: [game-model, api-contract, typescript, react]
dependency_graph:
  requires: []
  provides: [AddResult-no-playerID, NewGameResult-no-player_id, dynamic-Record]
  affects: [lib/business/game, lib/routers/game, app/src/types, app/src/components/stats]
tech_stack:
  added: []
  patterns: [functional-DI-type-update, dynamic-array-rendering]
key_files:
  created: []
  modified:
    - lib/business/game/types.go
    - lib/business/game/functions.go
    - lib/business/game/functions_test.go
    - lib/routers/game.go
    - lib/routers/game_test.go
    - app/src/types.ts
    - app/src/components/stats.tsx
decisions:
  - "playerID removed from AddResult chain — player implicit via deck ownership (aligns with GAME-01 contract)"
  - "Record/RecordComparator use Math.max over keys to handle variable pod sizes dynamically"
metrics:
  duration: "~5 min"
  completed: "2026-03-24"
  tasks: 2
  files: 7
---

# Phase 04 Plan 01: Remove playerID from AddResult chain + dynamic Record Summary

**One-liner:** Removed dead `playerID` parameter from Go AddResult chain and `player_id` from TypeScript game result types; replaced hardcoded 4-place Record display with a Math.max-driven dynamic version.

## Tasks Completed

| # | Name | Commit | Files |
|---|------|--------|-------|
| 1 | Remove playerID from backend AddResult chain | f3b9369 | types.go, functions.go, functions_test.go, game.go, game_test.go |
| 2 | Remove player_id from frontend types + dynamic Record/RecordComparator | 8d2c47e | types.ts, stats.tsx |

## What Was Built

**Task 1 — Backend AddResult cleanup:**
- `AddResultFunc` type alias in `lib/business/game/types.go` now has signature `func(ctx, gameID, deckID, place, killCount int) (int, error)` — `playerID` removed
- Closure in `lib/business/game/functions.go` updated to match (body was already building the model without PlayerID)
- `addGameResultRequest` struct in `lib/routers/game.go` no longer has `PlayerID` field
- Call site updated: `g.games.AddResult(ctx, req.GameID, req.DeckID, req.Place, req.KillCount)`
- Tests in `game_test.go` and `functions_test.go` updated to match new 4-arg signature

**Task 2 — Frontend type cleanup + dynamic Record:**
- `NewGameResult`: removed `player_id: number`
- `NewGameResultWithGame`: removed `player_id: number`
- `NewGameData`: removed `players: Array<Player>`
- `Record` component: replaced with `Math.max(...Object.keys(record).map(Number), 1)` to find max place, then `Array.from` to build parts array and `.join(" / ")` for display
- `RecordComparator`: replaced with dynamic `for (let place = 1; place <= maxPlace; place++)` loop
- Removed `getter` helper function and TODO comment

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed functions_test.go TestAddResult_Success calling old 5-arg signature**
- **Found during:** Task 1 (go vet revealed the error)
- **Issue:** `lib/business/game/functions_test.go` line 265 called `fn(context.Background(), 1, 10, 42, 2, 1)` — still passing 5 numeric args after AddResultFunc was reduced to 4
- **Fix:** Changed call to `fn(context.Background(), 1, 10, 2, 1)` (dropped the playerID `42`)
- **Files modified:** lib/business/game/functions_test.go
- **Commit:** f3b9369

## Verification

- `go vet ./lib/...` exits 0
- `go test ./lib/routers/ -run TestGameRouter_AddGameResult` passes
- `tsc --noEmit` produces errors ONLY in `routes/new/index.tsx` (expected — fixed in Plan 02); no errors in stats.tsx or types.ts

## Known Stubs

None — all changes are clean type contract removals and dynamic rendering logic.

## Self-Check: PASSED

- lib/business/game/types.go: FOUND
- lib/business/game/functions.go: FOUND
- lib/routers/game.go: FOUND
- app/src/types.ts: FOUND
- app/src/components/stats.tsx: FOUND
- Commit f3b9369: FOUND
- Commit 8d2c47e: FOUND
