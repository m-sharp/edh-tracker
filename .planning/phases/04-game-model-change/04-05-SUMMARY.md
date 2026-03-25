---
phase: 04-game-model-change
plan: 05
subsystem: frontend
tags: [ui, layout, new-game, pod-header]
dependency_graph:
  requires: []
  provides: [inline-remove-button, pod-header-new-game-button]
  affects: [app/src/routes/new/index.tsx, app/src/routes/pod/index.tsx, app/src/routes/pod/GamesTab.tsx]
tech_stack:
  added: []
  patterns: [MUI Button component=Link, flex-row card layout]
key_files:
  created: []
  modified:
    - app/src/routes/new/index.tsx
    - app/src/routes/pod/index.tsx
    - app/src/routes/pod/GamesTab.tsx
decisions:
  - "New Game button uses component={Link} to={...} on MUI Button — MUI styling with React Router navigation, no useNavigate needed in pod/index.tsx"
  - "Card fields column uses flex: 1 to fill available width; remove button uses fixed size pinned top-right via alignItems: flex-start"
metrics:
  duration: 4min
  completed_date: "2026-03-24"
  tasks_completed: 2
  files_modified: 3
---

# Phase 04 Plan 05: New Game Card Layout + Pod Header Button Summary

**One-liner:** Inline the remove button with card fields and surface the New Game entry point in the pod header above all tabs.

## Tasks Completed

| Task | Description | Commit |
|------|-------------|--------|
| 1 | Restructure New Game card to inline remove button with fields | 00a90dc |
| 2 | Move New Game button to pod header, remove from GamesTab | 9ae288b |

## What Was Built

### Task 1 — New Game card layout restructure

The card interior in `app/src/routes/new/index.tsx` was restructured from a two-row layout (remove button in its own Box above fields) to a single flex row:

- Outer `Box` with `display: flex`, `alignItems: flex-start`, `gap: 1`
- Left column (`flex: 1`, `flexDirection: column`, `gap: 1.5`) holds the Autocomplete (Deck picker) and the Place/Kills row
- Right side: `TooltipIconButton` (remove) pinned to the top of the row next to the Deck field
- All existing props on Autocomplete, TextField, and TooltipIconButton preserved unchanged

The `mb: 1.5` margin that was on the Autocomplete's `sx` prop was removed — the parent column's `gap: 1.5` now handles the spacing between Autocomplete and Place/Kills row.

### Task 2 — New Game button to pod header

**GamesTab.tsx:**
- Removed the `Box` containing the New Game `Button` (lines 56-60 in original)
- Removed `useNavigate` import and the `navigate` variable (no longer needed)
- Removed `Button` from MUI imports (no longer used)

**pod/index.tsx:**
- Added `Link` to `react-router-dom` import
- Added `Button` to `@mui/material` import
- Replaced bare `<Typography variant="h4" sx={{ mb: 2 }}>` with a flex row:
  - `Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between", mb: 2 }}`
  - Pod name Typography on the left
  - `Button variant="contained" component={Link} to={...}` on the right
- New Game button now appears above the tab strip, visible on all tabs (Decks, Players, Games, Settings)

## Deviations from Plan

None - plan executed exactly as written.

## Known Stubs

None.

## Self-Check: PASSED

Files exist:
- app/src/routes/new/index.tsx — FOUND
- app/src/routes/pod/index.tsx — FOUND
- app/src/routes/pod/GamesTab.tsx — FOUND

Commits:
- 00a90dc — FOUND (feat(04-05): restructure New Game card)
- 9ae288b — FOUND (feat(04-05): move New Game button to pod header)
