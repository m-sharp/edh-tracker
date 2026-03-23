---
phase: 01-backend-hardening
plan: "04"
subsystem: backend/repositories,backend/business,backend/routers
tags: [performance, batch-query, n+1, deck-stats, api-validation]
dependency_graph:
  requires: []
  provides: [GetStatsForDecks-batch-method, unfiltered-deck-400-guard]
  affects: [lib/repositories/gameResult, lib/business/deck, lib/routers/deck]
tech_stack:
  added: []
  patterns: [batch-IN-query, early-return-400]
key_files:
  created: []
  modified:
    - lib/repositories/gameResult/repo.go
    - lib/repositories/gameResult/stats.go
    - lib/repositories/interfaces.go
    - lib/business/deck/functions.go
    - lib/business/deck/types.go
    - lib/business/business.go
    - lib/business/testHelpers/mocks.go
    - lib/routers/deck.go
    - lib/routers/deck_test.go
    - lib/business/deck/functions_test.go
decisions:
  - "Removed deck.GetAll business function entirely — no other callers existed once router default path was removed"
  - "Added GetStatsForDecks to GameResultRepository interface and mock — forces compile-time satisfaction on concrete type"
metrics:
  duration_minutes: 8
  completed_date: "2026-03-22"
  tasks_completed: 2
  tasks_total: 2
  files_modified: 10
---

# Phase 01 Plan 04: Batch Deck Stats + Block Unfiltered Deck Endpoint Summary

Replaced N+1 deck stats query with a single `WHERE deck.id IN ?` batch query, and blocked the unfiltered `GET /api/decks` endpoint with a 400 response.

## What Was Built

**Task 1 — Batch deck stats query (PERF-01)**

Added `GetStatsForDecks(ctx, []int) (map[int]*Aggregate, error)` to the `gameResult` repository. This issues a single SQL query with `WHERE deck.id IN ?` instead of one query per deck. `buildEntitiesWithStats` in the deck business layer now collects all deck IDs, calls the batch method once, and maps results back to each deck. An empty deck slice returns early without any DB call.

Supporting changes:
- Added `getStatsForDecks` SQL constant with `deck.id AS deck_id` to allow grouping in Go
- Added `gameStatWithDeck` struct in `stats.go` for scanning batch rows
- Added `GetStatsForDecks` to `GameResultRepository` interface
- Updated `MockGameResultRepo` and 3 deck business tests to use the new batch signature

**Task 2 — Block unfiltered GET /api/decks (PERF-02)**

Both `GetAll` and `getAllPaginated` handlers in `DeckRouter` now return `400 Bad Request` with `"pod_id or player_id query param is required"` when neither filter is provided. The dead-code `deck.GetAll` business function and its `GetAllFunc` type were removed. The `business.go` wiring and router tests were updated accordingly.

## Decisions Made

- Removed `deck.GetAll` entirely rather than keeping it — grep confirmed no other callers existed outside the two removed default-path branches
- Added `GetStatsForDecks` to the repository interface (not just the concrete type) to maintain the functional DI pattern and keep the mock-based tests working

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Updated MockGameResultRepo and deck business tests for new batch signature**
- **Found during:** Task 1 verification (`go vet` failure)
- **Issue:** `MockGameResultRepo` in `testHelpers/mocks.go` did not implement the new `GetStatsForDecks` method; 3 deck business tests used `GetStatsForDeckFn` (the old per-deck mock) which would panic at runtime since `buildEntitiesWithStats` now calls the batch method
- **Fix:** Added `GetStatsForDecksFn` field and `GetStatsForDecks` method to `MockGameResultRepo`; updated the 3 test mocks to return a populated `map[int]*Aggregate`
- **Files modified:** `lib/business/testHelpers/mocks.go`, `lib/business/deck/functions_test.go`
- **Commit:** b3c1f11

**2. [Rule 1 - Bug] Updated deck router tests for new 400 behavior**
- **Found during:** Task 2
- **Issue:** `TestDeckRouter_GetAll_Success` and `TestDeckRouter_GetAll_Error` tested the old unfiltered path (expected 200/500); after the route change those test cases were semantically wrong
- **Fix:** Replaced both tests with `TestDeckRouter_GetAll_NoFilter_Returns400` and `TestDeckRouter_GetAll_Paginated_NoFilter_Returns400` that assert 400 with correct error message
- **Files modified:** `lib/routers/deck_test.go`
- **Commit:** b843e33

## Test Results

All tests pass: `go vet ./lib/... && go vet . && go test ./lib/...`

## Known Stubs

None.

## Self-Check: PASSED

- `lib/repositories/gameResult/repo.go` — contains `GetStatsForDecks` method: FOUND
- `lib/repositories/interfaces.go` — contains `GetStatsForDecks` in interface: FOUND
- `lib/business/deck/functions.go` — calls `GetStatsForDecks` (not per-deck loop): FOUND
- `lib/routers/deck.go` — contains `"pod_id or player_id query param is required"` in both paths: FOUND
- Commits b3c1f11 and b843e33: FOUND in git log
