---
phase: 03-frontend-structure
plan: 06
subsystem: ui
tags: [react, typescript, mui, react-router, game-view]

# Dependency graph
requires:
  - phase: 03-frontend-structure
    plan: 01
    provides: TooltipIconButton component used in GameView
provides:
  - app/src/routes/game/index.tsx (GameView + gameLoader)
  - app/src/routes/new/index.tsx (NewGameView + newGameLoader + createGame)
affects:
  - app/src/index.tsx (import resolution unchanged — auto-resolves to index.tsx)

# Tech tracking
tech_stack:
  added: []
  patterns:
    - TooltipIconButton wraps all icon buttons in GameView for consistent accessible tooltips
    - DataGrid autoHeight replaces fixed-height wrapper Box
    - Typography replaces raw HTML elements (h1, em) throughout GameView

# Key files
key_files:
  created:
    - app/src/routes/game/index.tsx
    - app/src/routes/new/index.tsx
  modified: []
  deleted:
    - app/src/routes/game.tsx
    - app/src/routes/new.tsx

# Decisions
decisions:
  - "NewGameView moved as straight file move only — no UI changes applied per plan (Phase 4 redesigns it entirely)"
  - "Save button text updated to Save Description per copywriting contract"
  - "No-description Typography uses plain text instead of raw em element to stay consistent with MUI"

# Metrics
metrics:
  duration: 7min
  completed_date: "2026-03-24"
  tasks: 2
  files: 4
---

# Phase 03 Plan 06: GameView and NewGameView Restructure Summary

GameView moved to game/index.tsx with all 7 G-series UI-SPEC fixes; NewGameView moved to new/index.tsx with import path updates only.

## Tasks Completed

| # | Task | Commit | Files |
|---|------|--------|-------|
| 1 | Move GameView to game/index.tsx with UI-SPEC fixes G-01 through G-07 | 1d1d170 | app/src/routes/game/index.tsx (created), app/src/routes/game.tsx (deleted) |
| 2 | Move NewGameView to new/index.tsx | 1630f21 | app/src/routes/new/index.tsx (created), app/src/routes/new.tsx (deleted) |

## What Was Built

### Task 1: GameView restructure + UI-SPEC fixes

- **G-01**: Raw `<h1>` replaced with `<Typography variant="h4">`
- **G-02**: Raw `<em>` date replaced with `<Typography variant="body2" color="text.secondary">`
- **G-03**: Edit description `<IconButton>` wrapped with `<TooltipIconButton title="Edit description" />`
- **G-04**: Edit/remove result `<IconButton>` elements wrapped with `<TooltipIconButton title="Edit result" />` and `<TooltipIconButton title="Remove result" />`
- **G-05**: Outer Box `alignItems` changed from `"center"` to `"flex-start"`
- **G-06**: Fixed `height: 355` wrapper Box removed; `autoHeight` prop added to DataGrid
- **G-07**: Delete Game button `mt: 2` changed to `mt: 3`
- Import paths updated to `../../auth`, `../../http`, `../../types`, `../../components/TooltipIcon`
- `IconButton` import removed (replaced by TooltipIconButton throughout)

### Task 2: NewGameView straight move

- `new.tsx` moved to `new/index.tsx` with `../http` and `../types` updated to `../../http` and `../../types`
- All exports preserved: `default` (View), `newGameLoader`, `createGame`
- No content changes — Phase 4 redesigns this view entirely

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing critical detail] Save button text updated to "Save Description"**
- **Found during:** Task 1
- **Issue:** Plan noted "verify the description save button says 'Save Description'" — original had plain "Save"
- **Fix:** Updated button text from "Save" to "Save Description"
- **Files modified:** app/src/routes/game/index.tsx
- **Commit:** 1d1d170

**2. [Rule 1 - Minor cleanup] No-description em element removed in GameDescription**
- **Found during:** Task 1
- **Issue:** `<Typography color="text.secondary"><em>No description</em></Typography>` still contained a raw em inside Typography
- **Fix:** Changed to `<Typography color="text.secondary">No description</Typography>` (plain text, no em)
- **Files modified:** app/src/routes/game/index.tsx
- **Commit:** 1d1d170

## Known Stubs

None. All UI elements are wired to real data.

## Self-Check: PASSED

- app/src/routes/game/index.tsx: FOUND
- app/src/routes/new/index.tsx: FOUND
- app/src/routes/game.tsx: CONFIRMED DELETED
- app/src/routes/new.tsx: CONFIRMED DELETED
- Commit 1d1d170: FOUND
- Commit 1630f21: FOUND
- TypeScript: tsc --noEmit exits 0
