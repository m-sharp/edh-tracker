---
phase: 01-backend-hardening
plan: "02"
subsystem: backend
tags: [security, authorization, transactions, game]
dependency_graph:
  requires: [01-01]
  provides: [pod-membership-check-on-game-create, atomic-game-creation]
  affects: [lib/routers/game.go, lib/business/game/functions.go, lib/business/business.go, main.go]
tech_stack:
  added: []
  patterns: [GORM-transaction, functional-DI, router-level-auth-check]
key_files:
  created: []
  modified:
    - lib/routers/game.go
    - lib/business/game/functions.go
    - lib/business/game/functions_test.go
    - lib/routers/game_test.go
    - lib/business/business.go
    - main.go
decisions:
  - "Integration tests used for Create success path — transaction wrapper bypasses interface mocks; real DB via testHelpers.NewTestDB (rollback-wrapped)"
  - "nil client is safe for error-path unit tests — function returns before reaching transaction code"
metrics:
  duration: "14 min"
  completed_date: "2026-03-23"
  tasks_completed: 2
  files_modified: 6
requirements_addressed: [SEC-01, SEC-03]
---

# Phase 01 Plan 02: Game Authorization and Atomic Creation Summary

Pod membership check in GameCreate (SEC-01) and transaction-wrapped game + result creation (SEC-03).

## What Was Built

### Task 1: Pod Membership Check in GameCreate (SEC-01)

Modified `GameCreate` handler in `lib/routers/game.go` to enforce pod membership before creating a game:

- After request body parsing and result validation, extracts `callerPlayerID` via `trackerHttp.CallerPlayerID(w, r)`
- Calls `g.getPodRole(ctx, req.PodID, callerPlayerID)` (already injected into `GameRouter`)
- Returns `403 Forbidden` if role is empty (caller is not a pod member)

The check uses the same `getPodRole` function already used by `requirePodManager` — just with a weaker condition (any role, not just manager).

### Task 2: Transaction-Wrapped Game Creation (SEC-03)

Modified `game.Create` constructor in `lib/business/game/functions.go` to wrap game + result inserts in a single GORM transaction:

- Added `client *lib.DBClient` parameter to `game.Create` constructor
- Inside the closure, uses `client.GormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error { ... })`
- Creates tx-scoped repo instances inside the callback: `gameRepository.NewRepository(&lib.DBClient{GormDb: tx})` and `gameResultRepository.NewRepository(&lib.DBClient{GormDb: tx})`
- If `BulkAdd` fails, GORM rolls back the `Add` automatically — no orphaned game rows

Updated `NewBusiness` in `lib/business/business.go`:
- Added `client *lib.DBClient` parameter
- Passes `client` to `game.Create`

Updated `main.go` to pass `client` to `business.NewBusiness`.

### Test Updates

**`lib/business/game/functions_test.go`:**
- `TestCreate_OtherFormat_SkipsDeckFormatCheck` and `TestCreate_MatchingFormat_Success` converted to integration tests using `repoTestHelpers.NewTestDB` (wraps in rollback transaction). Required because the transaction wrapper constructs real repo instances from `*gorm.DB`, bypassing interface-level mocks.
- Error-path tests (`FormatMismatch`, `FormatNotFound`, `InvalidInput`) remain unit tests with `nil` client — they return before reaching the transaction.

**`lib/routers/game_test.go`:**
- `TestGameRouter_Add_Success` and `TestGameRouter_Add_CreateError` updated to use `newFullGameRouter` with `withAuth` and a `getPodRole` mock returning "member".
- Added `TestGameRouter_Add_NonMember_Forbidden` test verifying 403 when `getPodRole` returns empty string.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Router tests failing after pod membership check was added**
- **Found during:** Task 1 verification
- **Issue:** `TestGameRouter_Add_Success` and `TestGameRouter_Add_CreateError` used `newTestGameRouter` (no `getPodRole`, no auth context) — new `CallerPlayerID` check returned 401
- **Fix:** Updated both tests to use `newFullGameRouter` with `withAuth(req, 42)` and a mock `getPodRole` returning "member". Added `TestGameRouter_Add_NonMember_Forbidden` for the new 403 path.
- **Files modified:** `lib/routers/game_test.go`
- **Commit:** bc35ace

**2. [Rule 1 - Bug] Business layer tests for Create failed after transaction refactor**
- **Found during:** Task 2 verification
- **Issue:** `Create` tests passed mock `gameRepo.AddFn` and `gameResultRepo.BulkAddFn`, but the transaction wrapper creates tx-scoped concrete repos, bypassing the interface mocks
- **Fix:** Converted success-path tests to integration tests using `repoTestHelpers.NewTestDB` (with rollback). Error-path tests retain `nil` client since they return before the transaction.
- **Files modified:** `lib/business/game/functions_test.go`
- **Commit:** bc35ace

## Known Stubs

None — plan goals fully achieved.

## Self-Check: PASSED
