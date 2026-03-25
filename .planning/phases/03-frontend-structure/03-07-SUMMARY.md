---
phase: 03-frontend-structure
plan: 07
subsystem: ui
tags: [react, cra, react-router, static-assets, spa]

# Dependency graph
requires:
  - phase: 03-frontend-structure
    provides: SPA handler in app/main.go that correctly serves index.html for all non-static routes
provides:
  - CRA build configuration with absolute asset paths rooted at /
affects: [docker-build, app-deployment, spa-refresh]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "CRA homepage field set to '/' — always emit absolute /static/... asset paths regardless of current route"

key-files:
  created: []
  modified:
    - app/package.json

key-decisions:
  - "homepage set to '/' not '.' in CRA config — absolute asset paths prevent sub-route refresh blank screen"

patterns-established:
  - "CRA homepage '/': asset paths in built index.html are absolute (/static/js/...) not relative (./static/js/...), so browser resolves them correctly from any sub-route"

requirements-completed: [FEND-04]

# Metrics
duration: 1min
completed: 2026-03-24
---

# Phase 03 Plan 07: Fix CRA Asset Paths for Sub-Route Refresh Summary

**Changed CRA homepage from "." to "/" so built asset paths are absolute, fixing blank white screen on page refresh at any sub-route (UAT gap closure)**

## Performance

- **Duration:** 1 min
- **Started:** 2026-03-24T15:24:12Z
- **Completed:** 2026-03-24T15:24:37Z
- **Tasks:** 1
- **Files modified:** 1

## Accomplishments

- Fixed root cause of "Uncaught SyntaxError: Unexpected token '<'" on refresh at sub-routes like /pod/123
- CRA now emits absolute asset paths (/static/js/main.abc123.js) instead of relative (./static/js/...)
- Page refresh at any SPA route now loads correctly — browser requests /static/js/... from server root

## Task Commits

Each task was committed atomically:

1. **Task 1: Change homepage from "." to "/" in app/package.json** - `ee80d8a` (fix)

**Plan metadata:** (to be added after docs commit)

## Files Created/Modified

- `app/package.json` - Changed `"homepage": "."` to `"homepage": "/"` — single field change, no other fields touched

## Decisions Made

- `"homepage": "/"` causes CRA to emit absolute asset paths; `"."` causes relative paths that break at sub-routes — the fix is minimal and correct per CRA docs

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Sub-route refresh now works correctly in the built Docker image
- UAT test #2 blocker resolved — app is usable after page refresh in production
- No additional changes required; spaHandler in app/main.go was already correct

---
*Phase: 03-frontend-structure*
*Completed: 2026-03-24*
