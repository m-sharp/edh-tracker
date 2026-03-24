---
phase: 03-frontend-structure
plan: 03
subsystem: ui
tags: [react, typescript, mui, react-router, mobile]

# Dependency graph
requires:
  - phase: 03-frontend-structure plan 01
    provides: TabbedLayout component and app/src/components/ directory

provides:
  - PodView restructured into pod/ subdirectory with 5 per-tab files
  - TabbedLayout integration with queryKey "podTab" for query-string tab persistence
  - P-01 fix: DataGrid heights responsive (xs: 400, sm: 600) in DecksTab and GamesTab
  - P-05 fix: Promote/Remove buttons have 44px minimum touch target in PlayersTab
  - P-06 fix: Settings form rows wrap on narrow viewports in SettingsTab
  - P-07 fix: AppBar "EDH Tracker" title hidden on mobile (xs: none, sm: flex)

affects:
  - Phase 03 plans 04-06 that import from pod/ subdirectory

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Pod subdirectory: PodView split into pod/index.tsx + 4 tab files, data loaded in index.tsx and passed as props"
    - "Responsive DataGrid height: sx={{ height: { xs: 400, sm: 600 } }} pattern for mobile-friendly grids"
    - "Touch target fix: sx={{ minHeight: 44 }} on small Button components for 44px minimum tap area"
    - "Settings row wrap: display:flex + flexWrap:wrap + minWidth on TextField for narrow viewport wrap"
    - "AppBar responsive title: display: { xs: 'none', sm: 'flex' } hides text on mobile, icon remains"

key-files:
  created:
    - app/src/routes/pod/index.tsx
    - app/src/routes/pod/DecksTab.tsx
    - app/src/routes/pod/PlayersTab.tsx
    - app/src/routes/pod/GamesTab.tsx
    - app/src/routes/pod/SettingsTab.tsx
  modified:
    - app/src/routes/root.tsx
  deleted:
    - app/src/routes/pod.tsx

key-decisions:
  - "PodView passes all data as props to tab components — data loading stays in index.tsx, tabs are pure display"
  - "AppBar title uses xs:none/sm:flex to prevent crowding with PodSelector + Avatar + Logout on 375px viewports"

patterns-established:
  - "Pod tab files accept props from parent index.tsx (no independent data fetching in tab files)"
  - "Responsive grid height pattern: xs:400/sm:600 applied consistently across DecksTab and GamesTab"

requirements-completed: [FEND-01, DSNG-04, FEND-02]

# Metrics
duration: 15min
completed: 2026-03-24
---

# Phase 03 Plan 03: Pod View Restructure Summary

**pod.tsx (355 lines) split into pod/index.tsx + 4 tab files using TabbedLayout with query-string persistence; P-01/P-05/P-06/P-07 mobile fixes applied**

## Performance

- **Duration:** 15 min
- **Started:** 2026-03-24T03:12:00Z
- **Completed:** 2026-03-24T03:27:00Z
- **Tasks:** 2
- **Files modified:** 7 (5 created, 1 modified, 1 deleted)

## Accomplishments
- Restructured monolithic pod.tsx (355 lines) into pod/ subdirectory with 5 files
- Integrated TabbedLayout with queryKey="podTab" for query-string-persisted tab state
- Applied P-01: DataGrid heights responsive (xs: 400, sm: 600) in DecksTab and GamesTab
- Applied P-05: 44px minimum touch targets on Promote/Remove buttons in PlayersTab
- Applied P-06: Settings form rows wrap via flexWrap in SettingsTab
- Applied P-07: AppBar "EDH Tracker" text hidden on mobile (xs: "none") to prevent crowding

## Task Commits

Each task was committed atomically:

1. **Task 1: Split pod.tsx into pod/ subdirectory with TabbedLayout migration** - `9f1a397` (feat, includes dependency cherry-picks)
2. **Task 2: Fix AppBar title visibility on mobile (P-07)** - `54aa8fe` (feat)

**Plan metadata:** (docs commit follows)

## Files Created/Modified
- `app/src/routes/pod/index.tsx` - PodView + podLoader using TabbedLayout with queryKey="podTab"; Settings tab hidden for non-managers
- `app/src/routes/pod/DecksTab.tsx` - Pod Decks tab with server-side paginated DataGrid; responsive height (xs: 400, sm: 600)
- `app/src/routes/pod/PlayersTab.tsx` - Pod Players tab with List; Promote/Remove buttons have sx={{ minHeight: 44 }}
- `app/src/routes/pod/GamesTab.tsx` - Pod Games tab with server-side paginated DataGrid; responsive height (xs: 400, sm: 600)
- `app/src/routes/pod/SettingsTab.tsx` - Pod Settings tab with name edit, invite link, delete pod; form rows use flexWrap + minWidth: 160
- `app/src/routes/root.tsx` - AppBar Typography display changed from "flex" to { xs: "none", sm: "flex" } for P-07 mobile fix
- `app/src/routes/pod.tsx` - DELETED (replaced by pod/ subdirectory)

## Decisions Made
- Data loading stays in index.tsx; tab components receive data as props. This keeps each tab file focused on rendering.
- AppBar title uses xs:"none"/sm:"flex" rather than abbreviation — complete removal on mobile is cleaner than a short title.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

The worktree (`agent-a2568049`) did not have `node_modules` since they're not tracked in git. Created a symlink from worktree `app/node_modules` to the main project's `app/node_modules` to enable TypeScript verification. TypeScript compiled cleanly after symlink.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- pod/ subdirectory established; future plans targeting pod views can import from pod/ files directly
- TabbedLayout pattern demonstrated with 4 tabs and conditional hidden tab — ready for use in player/ and deck/ restructure (plans 03-04 and 03-05)
- All 4 mobile usability fixes for DSNG-04 applied: P-01, P-05, P-06, P-07

## Self-Check: PASSED

All 5 pod/ files confirmed present. pod.tsx confirmed deleted. root.tsx contains xs:"none". tsc --noEmit exits 0. Both task commits (9f1a397, 54aa8fe) verified in git log.

---
*Phase: 03-frontend-structure*
*Completed: 2026-03-24*
