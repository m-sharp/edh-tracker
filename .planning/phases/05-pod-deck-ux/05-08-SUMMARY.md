---
phase: 05-pod-deck-ux
plan: 08
subsystem: frontend
tags: [datagrid, pagination, filter, ux]
dependency_graph:
  requires: []
  provides: [pod-decks-client-sort, player-decks-retired-filter, cancel-copy]
  affects: [app/src/routes/pod/index.tsx, app/src/routes/pod/DecksTab.tsx, app/src/routes/player/DecksTab.tsx, app/src/routes/pod/PlayersTab.tsx]
tech_stack:
  added: []
  patterns: [DataGrid initialState.filter, DataGrid client-side pagination]
key_files:
  created: []
  modified:
    - app/src/routes/pod/index.tsx
    - app/src/routes/pod/DecksTab.tsx
    - app/src/routes/player/DecksTab.tsx
    - app/src/routes/pod/PlayersTab.tsx
decisions:
  - PlayerDecksTab empty-state guard checks data.length not visibleRows — a player with only retired decks now sees the grid with default filter active, not the empty state
  - PodDecksTab podId prop retained in interface signature but unused after refactor — harmless and non-breaking
key_decisions:
  - PlayerDecksTab visibleRows removed; DataGrid initialState.filter used instead — retired decks hidden by default but accessible via filter removal
metrics:
  duration: 4min
  completed: "2026-03-27T02:34:34Z"
  tasks: 2
  files: 4
---

# Phase 05 Plan 08: UAT Gap Fixes — Pod Decks Sort, Retired Filter, Cancel Copy Summary

Client-side sort for Pod Decks tab via GetAllDecksForPod, DataGrid initialState.filter hides retired decks by default in PlayerDecksTab, and "Never mind" button copy changed to "Cancel" in PlayersTab dialogs.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Switch Pod Decks tab to client-side pagination | c5e2b93 | app/src/routes/pod/index.tsx, app/src/routes/pod/DecksTab.tsx |
| 2 | Retired filter via DataGrid and Cancel copy fix | afc5aaf | app/src/routes/player/DecksTab.tsx, app/src/routes/pod/PlayersTab.tsx |

## What Was Built

**Task 1 — Pod Decks client-side sort:**
- Replaced `GetDecksForPod` (paginated) with `GetAllDecksForPod` in `pod/index.tsx` loader
- Changed `PodLoaderData.decks` type from `PaginatedResponse<Deck>` to `Deck[]`
- Simplified `PodDecksTab`: removed all server-pagination state (rows, rowCount, loading, paginationModel, handlePaginationChange)
- DataGrid now receives decks array directly; `initialState.sortModel` for `record desc` is now effective

**Task 2 — Retired filter and copy:**
- Removed `visibleRows` hard-filter from `PlayerDecksTab`
- Added "Is Retired" boolean column to the DataGrid columns array
- Added `initialState.filter` with `{ field: "retired", operator: "is", value: "false" }` to hide retired decks by default
- Updated empty-state guard to check `data.length === 0` instead of `visibleRows.length === 0`
- Changed PlayersTab dialog cancel button text from "Never mind" to "Cancel"

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.

## Self-Check: PASSED

- app/src/routes/pod/index.tsx — modified, exists
- app/src/routes/pod/DecksTab.tsx — modified, exists
- app/src/routes/player/DecksTab.tsx — modified, exists
- app/src/routes/pod/PlayersTab.tsx — modified, exists
- Commits c5e2b93 and afc5aaf verified in git log
- tsc --noEmit exits 0 with no errors
