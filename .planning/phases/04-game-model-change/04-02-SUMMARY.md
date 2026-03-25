---
phase: 04-game-model-change
plan: 02
subsystem: frontend
tags: [game, game-result, ui, components]
requirements: [GAME-01, GAME-02]

dependency_graph:
  requires: []
  provides:
    - "TooltipIconButton with optional color, disabled, size, sx props"
    - "AddResultModal without player picker (deck-only)"
    - "PostGameResult call sends no player_id"
  affects:
    - "app/src/routes/new/index.tsx (Plan 03 will fix remaining TS errors)"

tech_stack:
  added: []
  patterns:
    - "MUI Tooltip + disabled span wrapper pattern for disabled button tooltip support"

key_files:
  created: []
  modified:
    - app/src/components/TooltipIcon.tsx
    - app/src/routes/game/index.tsx

decisions:
  - "span wrapper around disabled IconButton ensures tooltip still renders on hover (MUI Tooltip requirement)"
  - "http.ts PostGameResult required no changes — NewGameResultWithGame player_id was already removed by Plan 01"

metrics:
  duration: "5 minutes"
  completed_date: "2026-03-24T18:24:56Z"
  tasks_completed: 2
  files_modified: 2
---

# Phase 04 Plan 02: GameView Player Field Removal Summary

Player-free AddResultModal with deck-only result entry, plus TooltipIconButton extended with color, disabled, size, and sx props for Plan 03 card remove button.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Extend TooltipIconButton with color, disabled, size, sx props | 2ab1f4b | app/src/components/TooltipIcon.tsx |
| 2 | Remove Player Autocomplete from AddResultModal + clean PostGameResult | da7b0e8 | app/src/routes/game/index.tsx |

## Changes Made

### Task 1: TooltipIconButton Extended Props

`app/src/components/TooltipIcon.tsx`:
- Added `SxProps` and `Theme` imports from `@mui/material`
- Extended `TooltipIconButtonProps` with: `color?`, `size?`, `disabled?`, `sx?` (all optional)
- Updated `TooltipIconButton` to destructure and pass through new props to `IconButton`
- Wrapped `IconButton` in `<span>` — MUI Tooltip requires a DOM element that can receive hover events; disabled IconButtons don't forward events, so the span wrapper ensures tooltip still shows on hover
- Default for `size` is `"medium"` (preserves existing behavior); all other new props default to undefined

### Task 2: AddResultModal Player Field Removed

`app/src/routes/game/index.tsx`:
- Removed `players: PlayerWithRole[]` from `AddResultModalProps` interface
- Removed `playerId` state and `setPlayerId` from `AddResultModal` function
- Removed Player Autocomplete block (`options={players}`) from modal JSX
- Updated `handleAdd` guard: `if (!deckId) return` (was `if (!playerId || !deckId) return`)
- Updated `PostGameResult` call: sends only `{ game_id, deck_id, place, kill_count }` (no `player_id`)
- Updated submit button: `disabled={!deckId}` (was `disabled={!playerId || !deckId}`)
- Removed `players={players}` from `<AddResultModal>` call site in `GameResultsGrid`
- `GameResultsGridProps` still carries `players: PlayerWithRole[]` — used by DataGrid for player name display and manager role check

`app/src/http.ts`: No changes needed — `PostGameResult` already takes `NewGameResultWithGame` which had `player_id` removed in Plan 01.

## Verification

- `cd app && ./node_modules/.bin/tsc --noEmit` produces errors ONLY in `routes/new/index.tsx` (expected — fixed in Plan 03)
- No errors in `game/index.tsx` or `components/TooltipIcon.tsx`
- All existing TooltipIconButton callers (`root.tsx`, `game/index.tsx`) compile without changes

## Deviations from Plan

None — plan executed exactly as written.

The plan mentioned `app/src/http.ts` in `files_modified` but no changes were needed (PostGameResult already takes the updated `NewGameResultWithGame` type from Plan 01). This is expected per the plan's own task 2 action which states "No changes needed to PostGameResult function signature."

## Known Stubs

None — all changes wire directly to existing interfaces with no placeholder data.

## Self-Check: PASSED

- [x] app/src/components/TooltipIcon.tsx — modified, committed at 2ab1f4b
- [x] app/src/routes/game/index.tsx — modified, committed at da7b0e8
- [x] TypeScript errors only in routes/new/index.tsx (expected)
- [x] AddResultModalProps has no players field
- [x] AddResultModal has no playerId state
- [x] handleAdd guard checks only deckId
- [x] PostGameResult call sends no player_id
- [x] TooltipIconButton has color, disabled, size, sx props
