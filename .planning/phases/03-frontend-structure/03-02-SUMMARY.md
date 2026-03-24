---
phase: 03-frontend-structure
plan: 02
subsystem: frontend
tags: [react, mui, routing, loading-states, ui-polish]
dependency_graph:
  requires: ["03-01"]
  provides: ["HomeView in routes/home", "FEND-04 fix", "FEND-05 fix", "Login/Join polish"]
  affects: ["app/src/index.tsx", "app/src/routes/"]
tech_stack:
  added: []
  patterns: ["loading state with useState(true)", "centered loading spinner in Box", "MUI Button as Link"]
key_files:
  created:
    - app/src/routes/home/index.tsx
  modified:
    - app/src/index.tsx
    - app/src/routes/RequireAuth.tsx
    - app/src/routes/login.tsx
    - app/src/routes/join.tsx
decisions:
  - "HomeView loading state initialized to true so CircularProgress renders until fetch resolves (eliminates flash)"
  - "RequireAuth spinner wrapped in Box with justifyContent center — unpositioned spinner was the FEND-04 blank screen root cause"
  - "Button component={Link} pattern used for Go home to get MUI styling on React Router navigation"
metrics:
  duration: 4min
  completed_date: "2026-03-24"
  tasks_completed: 2
  files_changed: 5
---

# Phase 03 Plan 02: HomeView extraction, FEND-04/05 loading fixes, Login/Join polish — Summary

## One-Liner

HomeView extracted to routes/home/ with loading state fix (FEND-05), RequireAuth blank-screen resolved by centering spinner (FEND-04), Login/Join polished with playing cards icon, vertical centering, and MUI buttons.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Extract HomeView + fix FEND-04/FEND-05 | 8973223 | app/src/routes/home/index.tsx (created), app/src/index.tsx, app/src/routes/RequireAuth.tsx |
| 2 | Login/Join UI-SPEC fixes (L-02, L-04, J-01, J-02, J-03) | f15acb7 | app/src/routes/login.tsx, app/src/routes/join.tsx |

## What Was Built

### Task 1: HomeView extraction + loading state fixes

**app/src/routes/home/index.tsx** (new):
- `loading` state initialized to `true` — CircularProgress renders immediately on mount
- After fetch resolves: if pods exist, navigate to first pod; if empty, set `loading(false)` and show message
- On fetch error: `setLoading(false)` so user doesn't get stuck on spinner
- Empty state copy updated: "No pods yet. Create your first pod or ask a manager for an invite link."
- `navigate` added to `useEffect` dependency array (React exhaustive-deps compliance)

**app/src/index.tsx**:
- Inline `HomeView` function removed (was lines 29-48)
- `import HomeView from "./routes/home"` added
- `GetPodsForPlayer`, `useEffect`, `useNavigate`, `Typography` imports cleaned out (no longer used)
- File is now purely routing config + providers

**app/src/routes/RequireAuth.tsx**:
- `CircularProgress` now wrapped in `<Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", pt: 4 }}>`
- The FEND-04 "blank white screen on refresh" was caused by an unpositioned, uncentered 20px spinner that appeared invisible at top-left

### Task 2: Login/Join UI polish

**app/src/routes/login.tsx**:
- Outer Box changed from `mt: 8` to `minHeight: "100vh", justifyContent: "center"` flexbox (L-02 vertical centering)
- `SvgIconPlayingCards fontSize={48}` rendered above the "EDH Tracker" heading (L-04)

**app/src/routes/join.tsx**:
- `SvgIconPlayingCards fontSize={40}` added above messages in both no-code and error states (J-01)
- `<Link to="/">Go home</Link>` replaced with `<Button component={Link} to="/" variant="outlined" size="medium">Go home</Button>` in both states (J-02)
- Error state split into: `<Typography variant="h6">Something went wrong</Typography>` heading + `<Typography variant="body1" color="error">{error}</Typography>` detail (J-03)

## Decisions Made

1. **HomeView loading state initialized to `true`** — initializing to `false` would re-introduce the FEND-05 flash. The `true` default guarantees spinner renders before any fetch starts.

2. **RequireAuth spinner in centered Box** — the FEND-04 blank screen root cause was the spinner rendering with no visible size/position context during auth check, making it effectively invisible.

3. **`Button component={Link}` for Go home** — uses MUI's `component` prop pattern to get proper button styling on a React Router Link, avoiding plain anchor tag navigation.

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all states wire to real data or real navigation.

## Self-Check: PASSED

- `app/src/routes/home/index.tsx` exists: VERIFIED
- `app/src/index.tsx` has no inline `HomeView`: VERIFIED (grep returns 0)
- `app/src/routes/RequireAuth.tsx` has `justifyContent: "center"`: VERIFIED
- Task 1 commit 8973223: VERIFIED
- Task 2 commit f15acb7: VERIFIED
- TypeScript compiles cleanly (`tsc --noEmit` exits 0): VERIFIED
