---
phase: 01-backend-hardening
plan: 01
subsystem: api
tags: [go, error-handling, http-status, sentinel-errors]

# Dependency graph
requires: []
provides:
  - "lib/errs/errors.go with ErrForbidden sentinel importable by all packages"
  - "pod business layer wraps forbidden errors with errs.ErrForbidden"
  - "deck business layer wraps forbidden errors with errs.ErrForbidden"
  - "pod router discriminates 403 (ErrForbidden) vs 500 (DB/other) using errors.Is"
affects:
  - 01-backend-hardening
  - any future plan adding authorization checks in business or router layers

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Sentinel error pattern: define shared ErrForbidden in lib/errs, wrap with fmt.Errorf(..., errs.ErrForbidden), discriminate via errors.Is at router layer"

key-files:
  created:
    - lib/errs/errors.go
  modified:
    - lib/business/pod/functions.go
    - lib/business/deck/functions.go
    - lib/routers/pod.go
    - lib/routers/pod_test.go

key-decisions:
  - "ErrForbidden placed in lib/errs package (not lib/business) to avoid circular import since sub-packages pod/deck are under lib/business"
  - "Router uses errors.Is to discriminate 403 vs 500; plain errors (DB failures) now correctly return 500 instead of 403"

patterns-established:
  - "Sentinel error pattern: lib/errs/ErrForbidden is the canonical forbidden sentinel; wrap via fmt.Errorf; check via errors.Is"

requirements-completed: [INFRA-02]

# Metrics
duration: 9min
completed: 2026-03-22
---

# Phase 01 Plan 01: ErrForbidden Sentinel and 403/500 Discrimination Summary

**ErrForbidden sentinel in lib/errs with errors.Is-based 403 vs 500 discrimination in pod router handlers, fixing PromotePlayer/KickPlayer/LeavePod returning 403 for DB errors**

## Performance

- **Duration:** 9 min
- **Started:** 2026-03-22T02:27:04Z
- **Completed:** 2026-03-22T02:35:44Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments
- Created `lib/errs/errors.go` with the shared `ErrForbidden` sentinel, importable by all packages without circular dependency
- Updated pod and deck business functions (4 forbidden error sites) to wrap `errs.ErrForbidden` so errors.Is works transitively
- Updated PromotePlayer, KickPlayer, LeavePod handlers to return 403 only for `ErrForbidden`-wrapped errors and 500 for DB/other errors
- Updated pod router tests: renamed error tests for clarity, added DB-error 500 tests for all three handlers

## Task Commits

Each task was committed atomically:

1. **Task 1: Create ErrForbidden sentinel and update business-layer forbidden errors** - `05f10e1` (feat)
2. **Task 2: Update pod router to discriminate 403 vs 500 using errors.Is** - `b843e33` (fix, included in parallel agent commit)

## Files Created/Modified
- `lib/errs/errors.go` - New shared package with ErrForbidden sentinel
- `lib/business/pod/functions.go` - PromoteToManager, RemovePlayer, Leave now wrap errs.ErrForbidden
- `lib/business/deck/functions.go` - assertCallerOwnsDeck now wraps errs.ErrForbidden
- `lib/routers/pod.go` - PromotePlayer, KickPlayer, LeavePod use errors.Is for 403 vs 500 discrimination
- `lib/routers/pod_test.go` - Updated tests to wrap ErrForbidden, added 500 DB error test cases

## Decisions Made
- **ErrForbidden in lib/errs not lib/business:** Sub-packages (lib/business/pod, lib/business/deck) cannot import their parent package lib/business in Go — placing the sentinel in a separate lib/errs package breaks the circular import while remaining importable by both business sub-packages and routers
- **errors.Is not string prefix check:** Using errors.Is ensures the sentinel identity is checked (not string matching), making the discrimination correct even through multiple error wrapping layers

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Parallel agent committed Task 2 changes**
- **Found during:** Task 2 commit
- **Issue:** Parallel agent 01-04 committed pod.go and pod_test.go changes (which included my Task 2 edits) to the branch in commit b843e33 before I could commit them
- **Fix:** Verified the committed changes matched the plan exactly; no additional action needed
- **Files modified:** lib/routers/pod.go, lib/routers/pod_test.go
- **Verification:** go test ./lib/routers/... passes; errors.Is checks confirmed via grep
- **Committed in:** b843e33

---

**Total deviations:** 1 (parallel agent commit overlap — no scope creep, all changes correct)
**Impact on plan:** No impact on correctness. All plan acceptance criteria met.

## Issues Encountered
- Git index.lock contention from parallel agents required waiting briefly before staging
- ErrForbidden method (`GetStatsForDecks`) missing from MockGameResultRepo — pre-existing compile break fixed by linter before my vet run; mock was already partially updated (field present, method absent) and the linter completed the implementation

## Next Phase Readiness
- `lib/errs/ErrForbidden` sentinel is in place and tested — subsequent plans (deck router authorization, game creation auth) can import and use it immediately
- Pattern established: business functions return `fmt.Errorf("...: %w", errs.ErrForbidden)`, routers check `errors.Is(err, errs.ErrForbidden)`

## Known Stubs
None — all changes wire to real sentinel values and real error paths.

---
*Phase: 01-backend-hardening*
*Completed: 2026-03-22*
