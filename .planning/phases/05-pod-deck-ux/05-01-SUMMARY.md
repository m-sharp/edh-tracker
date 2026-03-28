---
phase: 05-pod-deck-ux
plan: 01
subsystem: api
tags: [go, gorm, mysql, batch-query, stats]

# Dependency graph
requires:
  - phase: 01-backend-hardening
    provides: functional DI pattern, interface-driven repos, GetStatsForDecks batch pattern
provides:
  - GetStatsForPlayersInPod batch query returning pod-scoped stats per player
  - POST /api/pod returns 201 with JSON body {"id": N}
  - POST /api/deck returns 201 with JSON body {"id": N}
affects: [05-02, 05-03, 05-04, frontend-pod-creation, frontend-deck-creation]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Batch pod-scoped stats: single SQL query with game.pod_id filter replaces N+1 per-player GetStatsForPlayer calls"
    - "Create endpoints return resource ID: capture return value, set Content-Type header before WriteHeader, encode JSON body"

key-files:
  created: []
  modified:
    - lib/repositories/gameResult/stats.go
    - lib/repositories/gameResult/repo.go
    - lib/repositories/interfaces.go
    - lib/business/player/functions.go
    - lib/business/testHelpers/mocks.go
    - lib/routers/pod.go
    - lib/routers/deck.go
    - lib/business/player/functions_test.go

key-decisions:
  - "gameStatWithPlayer uses struct embedding of gameStat (linter-applied pattern, consistent with gameStatWithDeck)"
  - "GetStatsForPlayersInPod called before per-player loop — single batch replaces N per-member queries"
  - "w.Header().Set before w.WriteHeader — Go http requires headers set before status code"

patterns-established:
  - "Batch stats pattern: collect IDs from member list, single IN-query with pod filter, map result by ID"

requirements-completed: [POD-04]

# Metrics
duration: 15min
completed: 2026-03-25
---

# Phase 5 Plan 01: Backend API Enhancements Summary

**Pod-scoped player stats batch query plus pod/deck create endpoints now returning resource IDs in JSON response body**

## Performance

- **Duration:** ~15 min
- **Started:** 2026-03-25T00:00:00Z
- **Completed:** 2026-03-25T00:15:00Z
- **Tasks:** 2
- **Files modified:** 8

## Accomplishments

- Added `GetStatsForPlayersInPod` repo method with a single SQL query filtering by `game.pod_id` — eliminates N+1 per-player stats calls in `GetAllByPod`
- `POST /api/pod` now returns `201 Created` with `{"id": N}` instead of empty body
- `POST /api/deck` now returns `201 Created` with `{"id": N}` instead of empty body

## Task Commits

1. **Task 1: Pod-scoped player stats — repo method, interface, business layer** - `2089af5` (feat)
2. **Task 2: Pod and deck create endpoints return resource IDs** - `5c61b9b` (feat)

## Files Created/Modified

- `lib/repositories/gameResult/stats.go` - Added `gameStatWithPlayer` scan struct
- `lib/repositories/gameResult/repo.go` - Added `getStatsForPlayersInPod` SQL constant and `GetStatsForPlayersInPod` method
- `lib/repositories/interfaces.go` - Added `GetStatsForPlayersInPod` to `GameResultRepository` interface
- `lib/business/player/functions.go` - Rewrote `GetAllByPod` to use batch pod-scoped stats
- `lib/business/testHelpers/mocks.go` - Added `GetStatsForPlayersInPodFn` field and method to `MockGameResultRepo`
- `lib/routers/pod.go` - `PodCreate` captures and returns pod ID in JSON body
- `lib/routers/deck.go` - `DeckCreate` captures and returns deck ID in JSON body
- `lib/business/player/functions_test.go` - Updated `TestGetAllByPod_*` tests to use new batch mock

## Decisions Made

- `gameStatWithPlayer` uses struct embedding of `gameStat` — the linter applied this automatically, consistent with `gameStatWithDeck`. Kept as-is.
- `GetStatsForPlayersInPod` placed before the per-player loop in `GetAllByPod` so the batch fires once regardless of member count
- `w.Header().Set("Content-Type", "application/json")` placed before `w.WriteHeader(http.StatusCreated)` — Go's `http.ResponseWriter` requires headers before status code

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Added GetStatsForPlayersInPod to MockGameResultRepo**
- **Found during:** Task 1 (after adding interface method)
- **Issue:** `MockGameResultRepo` in `lib/business/testHelpers/mocks.go` did not implement the updated `GameResultRepository` interface — compile-time check failed
- **Fix:** Added `GetStatsForPlayersInPodFn` field and `GetStatsForPlayersInPod` method to mock following the existing fn-field pattern
- **Files modified:** `lib/business/testHelpers/mocks.go`
- **Verification:** `go vet ./lib/...` exits 0
- **Committed in:** `2089af5` (Task 1 commit)

**2. [Rule 1 - Bug] Updated GetAllByPod test fixtures for changed function signature**
- **Found during:** Task 2 verification (`go test ./lib/business/player/...`)
- **Issue:** `TestGetAllByPod_Success` used `GetStatsForPlayerFn` which panics — function now calls `GetStatsForPlayersInPod`. `TestGetAllByPod_PlayerNotFound_Skipped` passed `nil` for gameResultRepo but the batch call runs before per-player lookup
- **Fix:** Replaced mock setup in both tests to use `GetStatsForPlayersInPodFn`; wired gameResultRepo in `PlayerNotFound` test
- **Files modified:** `lib/business/player/functions_test.go`
- **Verification:** `go test ./lib/business/player/...` passes
- **Committed in:** `5c61b9b` (Task 2 commit)

---

**Total deviations:** 2 auto-fixed (1 blocking interface satisfaction, 1 test update for changed behavior)
**Impact on plan:** Both fixes required for correct compilation and test coverage. No scope creep.

## Issues Encountered

None beyond the deviations documented above.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Backend is ready: pod/deck creation returns IDs for frontend navigation
- Pod Players tab will receive pod-scoped stats (kills/record filtered to that pod only) instead of global player stats
- Frontend plans (05-02 through 05-05) can proceed

---
*Phase: 05-pod-deck-ux*
*Completed: 2026-03-25*
