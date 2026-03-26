---
phase: 05-pod-deck-ux
plan: "02"
subsystem: frontend
tags: [ui, ux, decks, settings, new-game]
dependency_graph:
  requires: []
  provides: [record-min-4-places, pod-decks-default-sort, playing-cards-link, retired-deck-filter, create-pod-removed, discard-button]
  affects: [app/src/components/stats.tsx, app/src/routes/pod/DecksTab.tsx, app/src/routes/root.tsx, app/src/routes/player/DecksTab.tsx, app/src/routes/player/SettingsTab.tsx, app/src/routes/new/index.tsx]
tech_stack:
  added: []
  patterns: [component=Link on MUI Button for navigation]
key_files:
  created: []
  modified:
    - app/src/components/stats.tsx
    - app/src/routes/pod/DecksTab.tsx
    - app/src/routes/root.tsx
    - app/src/routes/player/DecksTab.tsx
    - app/src/routes/player/SettingsTab.tsx
    - app/src/routes/new/index.tsx
decisions:
  - "Removed unused Divider, useNavigate, and PostPod imports from SettingsTab after Create Pod section was removed — Rule 2 auto-fix to prevent TS errors"
  - "Discard button uses component=Link pattern (established project convention) rather than useNavigate"
  - "Discard button placed in a flex row with Submit, flex: 1 vs flex: 2 ratio to keep submit visually dominant"
metrics:
  duration: "7 minutes"
  completed_date: "2026-03-26T02:41:50Z"
  tasks_completed: 2
  files_modified: 6
---

# Phase 05 Plan 02: Frontend Quick Fixes Summary

Six independent frontend improvements shipped in Wave 1 alongside backend work. Record always shows 4 place columns, playing cards icon navigates home, Pod Decks tab defaults to record-descending sort, Player Decks tab hides retired decks, Create Pod is removed from player settings, and new game form has a Discard button.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Record min 4 places + Pod Decks default sort + playing cards icon link | a3339a4 | stats.tsx, pod/DecksTab.tsx, root.tsx |
| 2 | Retired deck filter + remove Create Pod + new game cancel button | 4d2d14b | player/DecksTab.tsx, player/SettingsTab.tsx, new/index.tsx |

## What Was Built

**Task 1 — Three AppBar and display fixes:**
- `stats.tsx`: Changed `Math.max(...keys, 1)` to `Math.max(...keys, 4)` — Record now always renders at least 4 place columns (e.g., "0 / 3 / 0 / 0" for a 1-win deck in a 4-player group)
- `pod/DecksTab.tsx`: Added `initialState={{ sorting: { sortModel: [{ field: "record", sort: "desc" }] } }}` to DataGrid — Pod Decks tab opens sorted by record descending
- `root.tsx`: Wrapped `SvgIconPlayingCards` in `<Link to="/">` inside the existing Box — icon navigates to home on click

**Task 2 — Three form and settings fixes:**
- `player/DecksTab.tsx`: Added `visibleRows` filtered by `!d.retired`; passes `visibleRows` to DataGrid instead of raw `data`; removed the `Is Retired` boolean column (hidden decks make the column meaningless)
- `player/SettingsTab.tsx`: Removed Create New Pod section entirely (state, handler, TextField, Button, error text, Divider); removed `PostPod` from imports; removed `useNavigate` and `navigate` (was only used in the removed handler)
- `new/index.tsx`: Added `Link` to react-router-dom imports; added Discard button using `component={Link} to={/pod/${podId}}` in a flex row alongside the Submit button

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing] Removed stale imports from SettingsTab after Create Pod deletion**
- **Found during:** Task 2
- **Issue:** Removing the Create Pod section left `useNavigate`, `navigate`, `PostPod`, and `Divider` as unused — TypeScript would error on unused imports
- **Fix:** Removed `useNavigate` from react-router-dom imports, removed the `navigate` const, removed `PostPod` from http imports, removed `Divider` from MUI imports
- **Files modified:** app/src/routes/player/SettingsTab.tsx
- **Commit:** 4d2d14b (bundled with task)

## Requirements Coverage

- POD-03: Pod Decks tab default sort — complete (initialState on DataGrid)
- DECK-02: Commander tooltip — already present (no work needed, verified in plan)
- DECK-03: Retired deck behavior — complete (filter in Player Decks tab)

## Known Stubs

None — all changes wire to real data. Player Decks filter operates on actual `retired` field from the API.

## Self-Check: PASSED

All 6 modified files confirmed present on disk. Both task commits (a3339a4, 4d2d14b) confirmed in git log.
