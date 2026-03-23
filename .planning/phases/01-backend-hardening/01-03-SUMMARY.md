---
phase: 01-backend-hardening
plan: 03
subsystem: auth
tags: [go, jwt, authorization, deck, ownership]

# Dependency graph
requires:
  - phase: 01-backend-hardening/01-01
    provides: "lib/errs/errors.go with ErrForbidden sentinel, errors.Is pattern established"
provides:
  - "DeckCreate uses JWT callerPlayerID, body player_id ignored (SEC-02)"
  - "Router-layer ownership checks on deck update and delete via assertCallerOwnsDeck helper"
  - "errors.Is(err, errs.ErrForbidden) replaces strings.HasPrefix forbidden checks in deck router"
  - "callerPlayerID removed from deck Update/SoftDelete/Retire business function signatures"
affects: [frontend, game-model-change]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Router-layer ownership: assertCallerOwnsDeck method on DeckRouter fetches deck via GetByID and checks PlayerID before calling business layer"
    - "JWT-sourced caller identity: CreateDeck ignores body player_id, uses trackerHttp.CallerPlayerID exclusively"
    - "errors.Is over string prefix matching for forbidden error discrimination at router layer"

key-files:
  created: []
  modified:
    - lib/routers/deck.go
    - lib/business/deck/functions.go
    - lib/business/deck/functions_test.go
    - lib/business/deck/types.go
    - lib/routers/deck_test.go

key-decisions:
  - "assertCallerOwnsDeck placed on DeckRouter (not business layer) — ownership is an HTTP-layer concern per D-03/D-05"
  - "Retire function signature also updated (callerPlayerID removed) for consistency even though no Retire endpoint exists yet"
  - "TestDeckUpdate_NotFound and TestDeckUpdate_Forbidden tests removed from business layer — ownership checks no longer there; equivalent coverage moved to router tests"

patterns-established:
  - "Ownership check pattern: router helper method fetches entity via business layer, checks ownership, writes 403 if mismatch — business layer stays free of caller context"

requirements-completed: [SEC-02]

# Metrics
duration: 15min
completed: 2026-03-22
---

# Phase 01 Plan 03: Deck Auth Migration Summary

**Deck authorization moved to router layer: DeckCreate uses JWT player ID, Update/Delete ownership checked via assertCallerOwnsDeck helper, errors.Is replaces string prefix matching**

## Performance

- **Duration:** 15 min
- **Started:** 2026-03-22T00:00:00Z
- **Completed:** 2026-03-22T00:15:00Z
- **Tasks:** 1
- **Files modified:** 5

## Accomplishments

- SEC-02 closed: `POST /api/deck` now uses JWT callerPlayerID, body `player_id` field removed entirely — no player spoofing possible
- D-05 implemented: ownership checks for update/delete moved from business functions to `assertCallerOwnsDeck` on DeckRouter
- D-06 implemented: `DeckCreate` calls `CallerPlayerID` before reading body, consistent with other mutating handlers
- Business layer `Update`/`SoftDelete`/`Retire` signatures simplified (no callerPlayerID), `assertCallerOwnsDeck` function removed entirely from business layer
- `strings.HasPrefix` forbidden checks replaced with typed `errors.Is(err, errs.ErrForbidden)` in deck router

## Task Commits

1. **Task 1: Migrate deck auth to router layer (SEC-02, D-05, D-06)** - `0a90ed6` (feat)

**Plan metadata:** (docs commit follows)

## Files Created/Modified

- `lib/routers/deck.go` - Removed PlayerID from newDeckRequest; DeckCreate calls CallerPlayerID first; assertCallerOwnsDeck method added; UpdateDeck/DeleteDeck call assertCallerOwnsDeck; errors.Is replaces strings.HasPrefix; strings import removed
- `lib/business/deck/functions.go` - assertCallerOwnsDeck removed; Update/SoftDelete/Retire closures drop callerPlayerID parameter; errs import removed
- `lib/business/deck/types.go` - UpdateFunc, SoftDeleteFunc, RetireFunc signatures updated to remove callerPlayerID
- `lib/business/deck/functions_test.go` - TestDeckUpdate tests updated (callerPlayerID arg removed, GetByIdFn mock removed); ownership-specific tests removed from business layer (TestDeckUpdate_NotFound, TestDeckUpdate_Forbidden, TestDeckSoftDelete_NotFound, TestDeckSoftDelete_Forbidden)
- `lib/routers/deck_test.go` - TestDeckRouter_Add_Success: body no longer has PlayerID, uses withAuth; TestDeckRouter_Add_NoAuth added; TestDeckRouter_Add_MissingPlayerID replaced with NoAuth test; Update tests include GetByID mock for ownership; TestDeckRouter_Update_Forbidden uses GetByID returning wrong owner

## Decisions Made

- `assertCallerOwnsDeck` placed on `DeckRouter` (not injected via business layer) — uses `d.decks.GetByID` which returns the full entity including `PlayerID`. This avoids adding a separate repo injection to the router.
- `Retire` function signature updated for consistency even though no Retire HTTP endpoint exists in the router yet.
- Business-layer ownership tests (NotFound/Forbidden for Update/SoftDelete) removed as those code paths no longer exist. Router tests provide equivalent coverage via GetByID mock returning wrong owner.

## Deviations from Plan

None — plan executed exactly as written. The plan anticipated all necessary changes including test updates.

## Issues Encountered

None. Pre-existing test failures in `TestGameRouter_Add_*` (from parallel plan 01-02 work) are unrelated to this plan's scope.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- SEC-02 complete; deck ownership is now enforced at router layer consistently with pod ownership pattern
- DeckCreate no longer accepts attacker-controlled player_id in request body
- Ready for plan 01-04 (game-related hardening) or frontend work referencing deck endpoints

---
*Phase: 01-backend-hardening*
*Completed: 2026-03-22*
