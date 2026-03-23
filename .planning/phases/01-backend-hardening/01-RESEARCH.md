# Phase 1: Backend Hardening - Research

**Researched:** 2026-03-22
**Domain:** Go backend — authorization, transactions, validation, N+1 query elimination
**Confidence:** HIGH

## Summary

This phase is a surgical hardening pass across the existing Go backend. All nine requirements (AUTH-02, SEC-01 through SEC-05, PERF-01, PERF-02, INFRA-02) are well-bounded, self-contained changes to specific files already identified in the CONTEXT.md. No new packages, libraries, or architectural patterns are introduced; all changes extend existing patterns in the codebase.

The dominant pattern throughout is: move auth enforcement up to the router layer, add a sentinel error type for forbidden errors, wrap game creation in a GORM transaction, add field-length guards to existing `Validate()` methods, replace an N+1 stats loop with a batch query, and guard against a weak JWT secret at startup.

The most intricate change is PERF-01 (batch deck stats query), which requires adding a new SQL query, a new method to the `GameResultRepository` interface, a new `GetStatsForDecks` method on `gameResult.Repository`, and rewriting `buildEntitiesWithStats` in the deck business layer. The transaction change (SEC-03) is the second-most complex because `game.Create` must receive `*lib.DBClient` as a constructor argument so it can open a `db.Transaction()` callback.

**Primary recommendation:** Implement in dependency order — SEC-03 (transaction) and PERF-01 (batch stats) first (they modify wiring in `business.go`), then the auth/error changes (SEC-01, SEC-02, INFRA-02, D-07/D-08/D-09), then the validation guards (SEC-04, SEC-05, AUTH-02). Each change is independently testable.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**D-01:** Pass `*lib.DBClient` to the `game.Create` constructor. Inside the `db.Transaction()` callback, create tx-scoped repo copies via `NewRepository(&lib.DBClient{Db: tx})` for both `gameRepo` and `gameResultRepo`. GORM auto-rollbacks on a non-nil error return from the callback.

**D-02:** No context changes needed. No modifications to existing repo methods. Only `game.Create` receives the `*lib.DBClient` constructor arg — no other business functions need it in this phase.

**D-03:** Router layer owns all authorization checks. This is the standard going forward.

**D-04:** Add pod membership check to `GameCreate` at the router layer: extract `callerPlayerID` via `trackerHttp.CallerPlayerID`, call `getPodRole(ctx, req.PodID, callerPlayerID)`, return 403 if the caller is not a member (role = "" or nil). Consistent with the existing `requirePodManager` pattern in the same file.

**D-05:** Phase 1 migrates the existing inconsistency: `assertCallerOwnsDeck` currently lives in the deck business layer (`lib/business/deck/functions.go`). Move ownership enforcement up to `DeckRouter` — router calls `CallerPlayerID` and checks deck ownership before calling business layer `Update`/`SoftDelete`. The `assertCallerOwnsDeck` helper is removed from the business layer.

**D-06:** For SEC-02 (deck create player_id): `DeckCreate` handler ignores `player_id` in the request body entirely. It calls `CallerPlayerID(w, r)` and passes that to `decks.Create`. The `createDeckRequest` struct's `PlayerID` field is removed (or ignored).

**D-07:** Introduce `var ErrForbidden = errors.New("forbidden")` in the `lib/business` package (or a shared `lib/business/errors.go` file).

**D-08:** Update ALL existing `fmt.Errorf("forbidden: ...")` calls throughout the business layer to wrap this sentinel: `fmt.Errorf("forbidden: ...: %w", business.ErrForbidden)`.

**D-09:** Routers that currently check for auth failures use `errors.Is(err, business.ErrForbidden)` to return 403 vs 500. Apply this consistently across all handlers — not just `PromotePlayer`/`KickPlayer`.

**SEC-04:** Add `used_count < max_used_count` check in `JoinByInvite` business function (`lib/business/pod/functions.go`). The `pod_invite` table already has `used_count` and a max — confirm the max field name in the migration.

**SEC-05:** Add max length validation to `Validate()` methods: `player.name` ≤ 256 chars (VARCHAR(256) in migration 2), `pod.name` ≤ 255 chars (migration 8), `deck.name` ≤ 255 chars (migration 14), `game.description` ≤ 256 chars (migration 4). Return 400 with a descriptive message.

**PERF-01:** Add `GetStatsForDecks(ctx context.Context, deckIDs []int) (map[int]*Aggregate, error)` to `GameResultRepository` interface and `gameResult.Repository`. Replace the per-deck loop in `buildEntitiesWithStats` with a single batch call. SQL uses `WHERE deck_id IN (?)`.

**PERF-02:** `GET /api/decks` unfiltered path (`default` case in `GetAll` and `getAllPaginated`) returns `400 Bad Request` with message: `"pod_id or player_id query param is required"`. Remove the `GetAll` business function and repo method if unused after this change.

**AUTH-02:** In `main.go` (or `lib/config.go`), after reading `JWT_SECRET`, check `len(cfg.JWTSecret) < 32` and call `log.Fatal` with a descriptive error. This runs before any server starts.

### Claude's Discretion

- Whether to put `ErrForbidden` in `lib/business/errors.go` or inline in `lib/business/business.go` — either works.
- Whether `PERF-02` removes the `GetAll` repo method entirely or just blocks the route. Remove it if no other caller exists.
- `assertCallerOwnsDeck` helper: can be kept as a private helper in the deck router file if reuse within that file is useful, or inlined — Claude's call.

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| AUTH-02 | Server rejects startup if JWT secret shorter than 32 bytes | `main.go` already has post-config hook location; `cfg.Get(lib.JWTSecret)` returns the value |
| SEC-01 | `POST /api/game` verifies caller is pod member before creating | `GameRouter` already has `getPodRole` injected; `requirePodManager` is the exact pattern to follow at a lower threshold (any role vs. manager) |
| SEC-02 | `POST /api/deck` uses caller's JWT player ID, ignores body `player_id` | `DeckCreate` already calls `ValidateCreate(req.PlayerID, ...)` — change to use `CallerPlayerID` and drop `PlayerID` from struct |
| SEC-03 | Game creation wrapped in a single DB transaction | `game.Create` already calls `gameRepo.Add` then `gameResultRepo.BulkAdd` sequentially — wrap both in `db.Transaction()` using tx-scoped repo copies |
| SEC-04 | Pod invite join validates max use count | `pod_invite` table has `used_count` but NO `max_used_count` column — a hardcoded constant or new migration needed |
| SEC-05 | String fields validate max length, return 400 not 500 | Existing `Validate()` methods only check non-empty; add len checks matching VARCHAR sizes from migrations |
| PERF-01 | Deck stats fetched in single batch query | `buildEntitiesWithStats` calls `GetStatsForDeck` once per deck in a loop; batch query approach fully designed |
| PERF-02 | `GET /api/decks` requires filter or returns 400 | Two unfiltered code paths exist (`GetAll` in `GetAll` method and `getAllPaginated` default case) — both must return 400 |
| INFRA-02 | `PromotePlayer`/`KickPlayer` return 403 for auth failures, 500 for DB errors | Currently both always use `http.StatusForbidden`; `errors.Is(err, business.ErrForbidden)` discriminates the two cases |
</phase_requirements>

## Standard Stack

### Core — No New Dependencies

This phase uses only what is already vendored. No new packages needed.

| Library | Version | Purpose |
|---------|---------|---------|
| `gorm.io/gorm` | v1.31.1 | Transactions via `db.Transaction(func(tx *gorm.DB) error {...})` |
| `errors` (stdlib) | — | `errors.New("forbidden")`, `errors.Is(err, ErrForbidden)` |
| `fmt` (stdlib) | — | `fmt.Errorf("...: %w", ErrForbidden)` wrapping |

**Installation:** None required. All dependencies already vendored.

## Architecture Patterns

### GORM Transaction Pattern (SEC-03)

The existing `user.CreatePlayerAndUser` in `lib/repositories/user/repo.go` already uses this pattern. The GORM skill confirms it.

```go
// Source: lib/repositories/user/repo.go (existing precedent) + .claude/skills/gorm/SKILL.md
err := client.GormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    gameRepo := gameRepository.NewRepository(&lib.DBClient{GormDb: tx})
    gameResultRepo := gameResultRepository.NewRepository(&lib.DBClient{GormDb: tx})

    gameID, err := gameRepo.Add(ctx, description, podID, formatID)
    if err != nil {
        return fmt.Errorf("failed to create game: %w", err)
    }
    // ... build results slice ...
    if err := gameResultRepo.BulkAdd(ctx, results); err != nil {
        return fmt.Errorf("failed to create game results: %w", err)
    }
    return nil // commits; non-nil return triggers rollback
})
```

**Constructor signature change:** `game.Create` receives `*lib.DBClient` as a new argument alongside the existing `repos.GameRepository` and `repos.GameResultRepository`. The `CreateFunc` type alias is UNCHANGED (callers only see the closure). Only `business.go` is updated to pass `client` to `game.Create`.

**CRITICAL observation:** `gameResult.Repository.NewRepository` accepts `*lib.DBClient` and uses `client.GormDb`. The DBClient struct has field `GormDb *gorm.DB`. Creating a tx-scoped instance requires `&lib.DBClient{GormDb: tx}` — the `log` field will be nil, which is acceptable per project convention (nil log field is fine for repositories).

### ErrForbidden Sentinel Pattern (D-07, D-08, D-09)

```go
// New file: lib/business/errors.go
package business

import "errors"

// ErrForbidden is returned by business functions when the caller lacks permission.
var ErrForbidden = errors.New("forbidden")
```

Existing forbidden strings to migrate (confirmed by reading source):
- `lib/business/pod/functions.go`: `PromoteToManager` — `fmt.Errorf("forbidden: caller is not a manager of pod %d", podID)`
- `lib/business/pod/functions.go`: `RemovePlayer` — `fmt.Errorf("forbidden: caller is not a manager of pod %d", podID)`
- `lib/business/pod/functions.go`: `Leave` — `fmt.Errorf("forbidden: cannot leave pod as the only manager...")`
- `lib/business/deck/functions.go`: `assertCallerOwnsDeck` — `fmt.Errorf("forbidden: deck %d does not belong to caller", deckID)` (this helper is being migrated up to the router, so this may be removed rather than migrated)

Router callers currently using `strings.HasPrefix(err.Error(), "forbidden:")`:
- `lib/routers/deck.go`: `UpdateDeck` and `DeleteDeck`
- These switch to `errors.Is(err, business.ErrForbidden)`

Router callers currently using hardcoded `http.StatusForbidden` for ALL errors:
- `lib/routers/pod.go`: `PromotePlayer`, `KickPlayer`, `LeavePod`
- These change to `errors.Is` check — 403 for ErrForbidden, 500 for others

### Batch Stats Query Pattern (PERF-01)

The existing `getStatsForDeck` SQL constant in `lib/repositories/gameResult/repo.go` queries one deck at a time. The batch version queries all decks in one round trip.

```go
// Source: lib/repositories/gameResult/repo.go (adapted from getStatsForDeck)
const getStatsForDecks = `
    SELECT deck.id AS deck_id,
           game_result.game_id, game_result.place, game_result.kill_count,
           (SELECT COUNT(*) FROM game_result gr2
             WHERE gr2.game_id = game_result.game_id
               AND gr2.deleted_at IS NULL) AS player_count
      FROM game_result
      INNER JOIN deck ON game_result.deck_id = deck.id
     WHERE deck.id IN ?
       AND game_result.deleted_at IS NULL
`
```

The scan struct needs a `DeckID` field to group results. Returns `map[int]*Aggregate` keyed by deck ID. Missing deck IDs (no games played) map to a zero-value `Aggregate`.

**New interface method:**
```go
// lib/repositories/interfaces.go — GameResultRepository
GetStatsForDecks(ctx context.Context, deckIDs []int) (map[int]*Aggregate, error)
```

**Updated `buildEntitiesWithStats`:**
```go
func buildEntitiesWithStats(ctx context.Context, decks []deckRepository.Model, gameResultRepo repos.GameResultRepository) ([]EntityWithStats, error) {
    deckIDs := make([]int, len(decks))
    for i, d := range decks { deckIDs[i] = d.ID }

    statsMap, err := gameResultRepo.GetStatsForDecks(ctx, deckIDs)
    if err != nil {
        return nil, fmt.Errorf("failed to get batch stats: %w", err)
    }

    result := make([]EntityWithStats, 0, len(decks))
    for _, d := range decks {
        entity := ToEntity(d, d.Player.Name, d.Format.Name, commanderInfoFromModel(d))
        result = append(result, ToEntityWithStats(entity, statsMap[d.ID]))
    }
    return result, nil
}
```

**Edge case:** When `len(deckIDs) == 0`, skip the query and return `[]EntityWithStats{}` immediately — GORM's `IN ?` with an empty slice produces invalid SQL in some drivers.

### Max Length Validation Pattern (SEC-05)

The existing `Validate()` methods only check for empty strings. The pattern is to add an additional `len()` check.

**Verified column sizes from migrations:**
- `player.name`: `VARCHAR(256)` (migration 2) — NOTE: CONTEXT.md says 50 chars, but migration 2 says `VARCHAR(256)`. Use 256.
- `pod.name`: `VARCHAR(255)` (migration 8)
- `deck.name`: `VARCHAR(255)` (migration 14)
- `game.description`: `VARCHAR(256)` (migration 4) — this field is in `createGameRequest`, not an entity `Validate()` — validate at the router level

```go
// Example — player entity (lib/business/player/entity.go)
func (e Entity) Validate() error {
    if e.Name == "" {
        return fmt.Errorf("name is required")
    }
    if len(e.Name) > 256 {
        return fmt.Errorf("name must be 256 characters or fewer")
    }
    return nil
}
```

For `game.description`, the validation must happen in `GameCreate` router handler before calling `games.Create`, since there is no `game.Entity.Validate()` method — only `gameResult.InputEntity.Validate()`.

### SEC-04 — Invite Max Use Count

**Critical finding:** The `pod_invite` table (migration 18) has `used_count INT NOT NULL DEFAULT 0` but has NO `max_used_count` column. The `podInvite.Model` struct confirms this — it has only `UsedCount`, no max field.

The CONTEXT.md says "a max — confirm the max field name in the migration." There is no such field. Two options:
1. Use a hardcoded constant (e.g., `const maxInviteUses = 10`)
2. Add a new migration adding a `max_used_count` column to `pod_invite`

Since there is no existing max field in the DB schema or model, this requires either a hardcoded constant or a new migration. A hardcoded constant avoids a schema migration.

### AUTH-02 — JWT Secret Length Guard

The config is fully loaded before any server starts, and `main.go` calls `lib.NewConfig(requireCfgs...)` then proceeds to `lib.NewDBClient`. The guard can go in `main.go` right after `NewConfig`:

```go
// main.go — after cfg creation
jwtSecret, _ := cfg.Get(lib.JWTSecret)
if len(jwtSecret) < 32 {
    log.Fatalf("JWT_SECRET must be at least 32 bytes; got %d", len(jwtSecret))
}
```

### PERF-02 — Remove Unfiltered Deck Path

Two code paths must return 400 (confirmed by reading `deck.go`):

1. `GetAll` handler — the `default` case in the `switch` block currently calls `d.decks.GetAll(ctx)`
2. `getAllPaginated` — the `default` case currently also calls `d.decks.GetAll(ctx)`

Both change to:
```go
default:
    http.Error(w, "pod_id or player_id query param is required", http.StatusBadRequest)
    return
```

After this change, `deck.GetAll` business function and `DeckRepository.GetAll` interface method should be removed if no other callers exist. Verified: no callers of `d.decks.GetAll` exist outside `deck.go`. The business `GetAll` constructor in `business.go` can be removed. The `DeckRepository.GetAll` interface method and repo implementation can be removed if no other code references it — check with `grep` at plan time.

### SEC-02 — DeckCreate Player ID

The `newDeckRequest` struct has a `PlayerID int` field and `DeckCreate` passes it to `decks.Create`. The fix:
1. Either remove `PlayerID` from `newDeckRequest` or keep it (ignored) — removing is cleaner
2. Extract `callerPlayerID` via `trackerHttp.CallerPlayerID(w, r)`
3. Pass `callerPlayerID` to `decks.Create`
4. Update `deck.ValidateCreate` call — it currently receives `req.PlayerID`; change to `callerPlayerID`

### D-05 — assertCallerOwnsDeck Migration

Currently `deck.Update` and `deck.SoftDelete` both call `assertCallerOwnsDeck` (a private function in `lib/business/deck/functions.go`) before performing their operation. The migration:

1. The router already calls `trackerHttp.CallerPlayerID` before `UpdateDeck` and `DeleteDeck`
2. Add `requireDeckOwner(w, r, deckID, callerPlayerID)` helper to `DeckRouter` (analogous to `requirePodManager` in `GameRouter`)
3. The helper calls `decks.GetByID(ctx, deckID)` and checks `entity.PlayerID == callerPlayerID`
4. Remove `callerPlayerID int` parameter from `deck.Update` and `deck.SoftDelete` signatures (or keep it as unused — removing is cleaner)
5. Remove `assertCallerOwnsDeck` from `functions.go`

**Warning:** Removing `callerPlayerID` from `deck.Update` and `deck.SoftDelete` signatures requires updating `types.go` (`UpdateFunc`, `SoftDeleteFunc`), `business.go` (constructor call), and all test mocks.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead |
|---------|-------------|-------------|
| DB transaction with automatic rollback | Manual begin/rollback/commit | `db.Transaction(func(tx *gorm.DB) error {...})` — GORM handles rollback on error, commit on nil |
| Error type discrimination | String prefix matching | `errors.Is(err, business.ErrForbidden)` — sentinel error wrapping with `%w` |
| IN clause batch query | Looped single queries | GORM `Where("deck_id IN ?", deckIDs)` — one round trip |

## Runtime State Inventory

Step 2.5: SKIPPED — this is a greenfield hardening phase with no rename, rebrand, or migration of stored string values. No runtime state inventory required.

## Environment Availability

Step 2.6: SKIPPED — this phase makes only code changes within the Go backend. All required tools (Go, Docker) are already confirmed present from existing project setup. No new external dependencies introduced.

## Common Pitfalls

### Pitfall 1: DBClient.GormDb nil in tx-scoped repository
**What goes wrong:** `&lib.DBClient{GormDb: tx}` leaves the `log` field nil. If any repository method tries to use the logger, it panics.
**Why it happens:** `lib.DBClient` has an unexported `log *zap.Logger` field.
**How to avoid:** Repository methods in `gameResult` and `game` do not use the logger — confirmed by reading `repo.go`. Safe to pass nil log.
**Warning signs:** If a future repo method calls `r.db.log`, it will panic. Add a note to the PR.

### Pitfall 2: GORM IN clause with empty slice
**What goes wrong:** `Where("deck_id IN ?", []int{})` generates invalid SQL (`deck_id IN ()`) in MySQL.
**Why it happens:** GORM does not guard against empty IN lists.
**How to avoid:** In `GetStatsForDecks`, guard at the top: if `len(deckIDs) == 0`, return empty map immediately.
**Warning signs:** SQL error on pod/player with no decks.

### Pitfall 3: player.name VARCHAR(256) not VARCHAR(50)
**What goes wrong:** The CONTEXT.md note says "≤ 50 chars (matches `VARCHAR(50)` in migration 7)". Migration 7 creates the `user` table — not the `player` table. Migration 2 creates the `player` table with `VARCHAR(256)`.
**Why it happens:** CONTEXT.md cited the wrong migration number for player name.
**How to avoid:** Use 256 as the max for `player.name`, matching migration 2's actual `VARCHAR(256)`.

### Pitfall 4: SEC-04 missing max_used_count column
**What goes wrong:** CONTEXT.md says "The `pod_invite` table already has `used_count` and a max". But migration 18 shows only `used_count` — no max column exists.
**Why it happens:** CONTEXT.md assumed a max column exists; the migration does not have one.
**How to avoid:** Use a hardcoded constant for max invite uses, or add a migration. Do NOT assume the field exists in the Model struct.

### Pitfall 5: Removing callerPlayerID from deck.Update/SoftDelete signatures
**What goes wrong:** Downstream callers (including test mocks in `testHelpers`) will fail to compile.
**Why it happens:** The `UpdateFunc` and `SoftDeleteFunc` type aliases reference the current signature.
**How to avoid:** Update `types.go`, all mock implementations in `testHelpers`, and `business.go` together in the same change.

### Pitfall 6: PromotePlayer/KickPlayer always return 403 today
**What goes wrong:** `PromotePlayer` and `KickPlayer` in `pod.go` call `trackerHttp.WriteError(p.log, w, http.StatusForbidden, err, ...)` for ALL errors. This sends 403 even when the DB fails.
**Why it happens:** No error discrimination exists today.
**How to avoid:** After adding `ErrForbidden` wrapping in `pod.PromoteToManager` and `pod.RemovePlayer`, update the routers to use `errors.Is(err, business.ErrForbidden)` — 403 if true, 500 otherwise.

### Pitfall 7: Deck auth check location during D-05 migration
**What goes wrong:** If `assertCallerOwnsDeck` is removed from `deck.Update` before the router check is added, the ownership check is lost entirely.
**Why it happens:** Two-step refactor with a window of vulnerability.
**How to avoid:** Add the router-level check first, verify it works, then remove the business-layer helper.

## Code Examples

### Transaction: GORM with tx-scoped repo copies

```go
// Source: lib/business/game/functions.go (modified per D-01)
func Create(
    log *zap.Logger,
    client *lib.DBClient,           // NEW arg
    gameRepo repos.GameRepository,
    gameResultRepo repos.GameResultRepository,
    deckRepo repos.DeckRepository,
    getFormat format.GetByIDFunc,
) CreateFunc {
    return func(ctx context.Context, description string, podID, formatID int, inputs []gameResult.InputEntity) error {
        // ... format + deck validation unchanged ...

        return client.GormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
            txClient := &lib.DBClient{GormDb: tx}
            txGameRepo := gameRepository.NewRepository(txClient)
            txGameResultRepo := gameResultRepository.NewRepository(txClient)

            gameID, err := txGameRepo.Add(ctx, description, podID, formatID)
            if err != nil {
                return fmt.Errorf("failed to create game: %w", err)
            }
            results := make([]gameResultRepository.Model, 0, len(inputs))
            for _, input := range inputs {
                results = append(results, gameResultRepository.Model{
                    GameID: gameID, DeckID: input.DeckID,
                    Place: input.Place, KillCount: input.Kills,
                })
            }
            if err := txGameResultRepo.BulkAdd(ctx, results); err != nil {
                return fmt.Errorf("failed to create game results: %w", err)
            }
            return nil
        })
    }
}
```

### ErrForbidden sentinel

```go
// Source: new file lib/business/errors.go
package business

import "errors"

var ErrForbidden = errors.New("forbidden")

// Usage in pod/functions.go:
// return fmt.Errorf("forbidden: caller is not a manager of pod %d: %w", podID, ErrForbidden)

// Usage in routers (replaces strings.HasPrefix check):
// if errors.Is(err, business.ErrForbidden) {
//     http.Error(w, err.Error(), http.StatusForbidden)
//     return
// }
// trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "...", "...")
```

### Batch stats query

```go
// Source: lib/repositories/gameResult/repo.go (new method)
const getStatsForDecks = `
    SELECT deck.id AS deck_id,
           game_result.game_id, game_result.place, game_result.kill_count,
           (SELECT COUNT(*) FROM game_result gr2
             WHERE gr2.game_id = game_result.game_id
               AND gr2.deleted_at IS NULL) AS player_count
      FROM game_result
      INNER JOIN deck ON game_result.deck_id = deck.id
     WHERE deck.id IN ?
       AND game_result.deleted_at IS NULL`

type deckGameStat struct {
    DeckID      int
    GameID      int
    Place       int
    KillCount   int
    PlayerCount int
}

func (r *Repository) GetStatsForDecks(ctx context.Context, deckIDs []int) (map[int]*Aggregate, error) {
    if len(deckIDs) == 0 {
        return map[int]*Aggregate{}, nil
    }
    var rows []deckGameStat
    if err := r.db.WithContext(ctx).Raw(getStatsForDecks, deckIDs).Scan(&rows).Error; err != nil {
        return nil, fmt.Errorf("failed to get batch stats for decks: %w", err)
    }
    // group by deck_id, then use existing gameStats.toAggregate() pattern
    grouped := map[int]gameStats{}
    for _, row := range rows {
        grouped[row.DeckID] = append(grouped[row.DeckID], gameStat{
            GameID: row.GameID, Place: row.Place,
            KillCount: row.KillCount, PlayerCount: row.PlayerCount,
        })
    }
    result := map[int]*Aggregate{}
    for _, id := range deckIDs {
        agg := grouped[id].toAggregate()
        result[id] = &agg
    }
    return result, nil
}
```

### Pod membership check in GameCreate (SEC-01)

```go
// Source: lib/routers/game.go (modified per D-04)
func (g *GameRouter) GameCreate(w http.ResponseWriter, r *http.Request) {
    // ... existing body parsing + validation ...

    callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
    if !ok {
        return
    }

    role, err := g.getPodRole(r.Context(), req.PodID, callerPlayerID)
    if err != nil {
        trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to check pod membership", "internal error")
        return
    }
    if role == "" {
        http.Error(w, "Forbidden: must be a pod member to create games", http.StatusForbidden)
        return
    }

    // ... existing games.Create call ...
}
```

## State of the Art

| Area | Current State | After Phase 1 |
|------|--------------|---------------|
| Game creation auth | No pod membership check — any authenticated user can POST /api/game | 403 for non-members |
| Deck create player_id | Body `player_id` used verbatim — caller can create decks for others | Body `player_id` ignored; JWT player ID used |
| Game creation atomicity | Two sequential DB calls — game row can persist if result insert fails | Single transaction; automatic rollback on failure |
| Invite max use | Only expiry checked | Both expiry and use count checked |
| String length | Unvalidated — MySQL truncation or error at 500 | 400 with descriptive message before DB call |
| Deck stats loading | N+1 query (one per deck) | Single batch query for all decks |
| Unfiltered deck list | Slow full-table scan accessible | 400 — filter required |
| JWT secret | Any length accepted | Startup fatal if < 32 bytes |
| Forbidden vs. infra errors | `strings.HasPrefix` check or always 403 | `errors.Is(err, ErrForbidden)` — 403 vs. 500 |

## Open Questions

1. **SEC-04: max invite uses — constant or migration?**
   - What we know: No `max_used_count` column exists in `pod_invite` table or model. The CONTEXT.md assumption was incorrect.
   - What's unclear: Whether the product wants per-invite max counts (migration needed) or a global max (constant suffices).
   - Recommendation: Use a hardcoded package-level constant `const maxInviteUses = 10` in `lib/business/pod/functions.go`. This satisfies the security requirement without a schema migration and can be elevated to a DB column in a future phase.

2. **`deck.Update` / `deck.SoftDelete` signature change scope**
   - What we know: `callerPlayerID` is currently a parameter on both. Removing it requires updating `types.go`, mock helpers, and `business.go`.
   - What's unclear: Whether to remove the parameter or keep it unused.
   - Recommendation: Remove it — the business layer should not be making authorization decisions per D-03/D-05.

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing + testify v1.11.1 |
| Config file | none (standard `go test`) |
| Quick run command | `go test ./lib/business/... ./lib/repositories/...` |
| Full suite command | `go test ./lib/...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| AUTH-02 | `log.Fatal` if JWT secret < 32 bytes | unit | `go test ./lib/... -run TestJWTSecretLength` | ❌ Wave 0 |
| SEC-01 | `GameCreate` returns 403 for non-pod-member | unit (router) | `go test ./lib/routers/ -run TestGameCreate_NonMember` | ❌ Wave 0 |
| SEC-02 | `DeckCreate` uses JWT player ID, ignores body | unit (router) | `go test ./lib/routers/ -run TestDeckCreate_UsesCallerID` | ❌ Wave 0 |
| SEC-03 | Game row absent if result insert fails | unit (business) | `go test ./lib/business/game/ -run TestCreate_RollbackOnResultFailure` | ❌ Wave 0 |
| SEC-04 | `JoinByInvite` returns error if used_count >= max | unit (business) | `go test ./lib/business/pod/ -run TestJoinByInvite_MaxUsed` | ❌ Wave 0 |
| SEC-05 | Long name returns 400 not 500 | unit (entity validate) | `go test ./lib/business/player/ ./lib/business/pod/ ./lib/business/deck/ -run TestValidate_MaxLength` | ❌ Wave 0 |
| PERF-01 | `buildEntitiesWithStats` issues one batch query | unit (business) | `go test ./lib/business/deck/ -run TestGetAllByPod_BatchStats` | ❌ Wave 0 |
| PERF-02 | Unfiltered `/api/decks` returns 400 | unit (router) | `go test ./lib/routers/ -run TestDeckGetAll_NoFilter` | ❌ Wave 0 |
| INFRA-02 | `PromotePlayer` returns 403 for auth, 500 for DB error | unit (router) | `go test ./lib/routers/ -run TestPromotePlayer_ErrorCodes` | ❌ Wave 0 |

### Sampling Rate

- **Per task commit:** `go vet ./lib/...` (fast compile check)
- **Per wave merge:** `go test ./lib/...`
- **Phase gate:** `go test ./lib/...` green before `/gsd:verify-work`

### Wave 0 Gaps

All test files for this phase are new — existing test files cover the pre-change behavior. The new tests must be added as part of implementation, not as a separate Wave 0. Per project conventions:

- `lib/routers/game_test.go` — covers SEC-01, SEC-02, PERF-02
- `lib/routers/deck_test.go` — covers SEC-02, PERF-02, INFRA-02 patterns
- `lib/routers/pod_test.go` — covers INFRA-02 (PromotePlayer, KickPlayer error codes)
- `lib/business/game/functions_test.go` — already exists; add SEC-03 rollback case
- `lib/business/pod/functions_test.go` — already exists; add SEC-04 max use case
- `lib/business/player/entity_test.go` — already exists; add SEC-05 length cases
- `lib/business/pod/entity_test.go` — already exists; add SEC-05 length case
- `lib/business/deck/entity_test.go` — already exists; add SEC-05 length case

## Project Constraints (from CLAUDE.md)

- **No framework changes**: Go + Gorilla Mux + GORM + MySQL backend; React + MUI + React Router v6 frontend. Phase 1 is backend-only.
- **No breaking changes**: Existing game/player/deck data in the DB must remain accessible. PERF-02 changes the API contract for unfiltered `GET /api/decks` — any frontend callers that send no filter will break. Verify no existing frontend code calls `/api/decks` without a filter before removing the path. (The frontend uses `GetAllByPod` and `GetAllForPlayer` which always supply a filter — safe to remove.)
- **Auth: Google OAuth only**: No new auth changes in this phase.
- **Docker deployment**: No changes to Dockerfiles needed.
- **After API changes**: Run `/smoke-test` skill (rebuild Docker image and verify core endpoints respond).
- **`go vet ./lib/...`** for compile checks — never `go build ./...` or `go build ./lib/...`.
- **`go test ./lib/...`** for tests.

## Sources

### Primary (HIGH confidence)

- `lib/routers/game.go` — GameCreate, requirePodManager, getPodRole injection
- `lib/routers/deck.go` — DeckCreate, UpdateDeck, DeleteDeck, assertCallerOwnsDeck usage pattern
- `lib/routers/pod.go` — PromotePlayer, KickPlayer current behavior (always 403)
- `lib/business/deck/functions.go` — assertCallerOwnsDeck, buildEntitiesWithStats N+1 loop
- `lib/business/game/functions.go` — Create, two sequential DB calls without transaction
- `lib/business/pod/functions.go` — JoinByInvite, existing forbidden strings
- `lib/repositories/gameResult/repo.go` — GetStatsForDeck, existing SQL query pattern
- `lib/repositories/interfaces.go` — GameResultRepository interface
- `lib/migrations/2.go` — player.name VARCHAR(256)
- `lib/migrations/4.go` — game.description VARCHAR(256)
- `lib/migrations/8.go` — pod.name VARCHAR(255)
- `lib/migrations/14.go` — deck.name VARCHAR(255)
- `lib/migrations/18.go` — pod_invite table, used_count only (no max_used_count)
- `lib/repositories/podInvite/model.go` — Model struct confirming no MaxUsedCount field
- `lib/business/pod/entity.go` — existing Validate() methods
- `lib/business/player/entity.go` — existing Validate() method
- `lib/business/deck/entity.go` — ValidateCreate function
- `lib/business/game/types.go` — CreateFunc signature
- `lib/business/business.go` — NewBusiness wiring, game.Create constructor call
- `main.go` — startup sequence, JWT secret config location
- `lib/config.go` — Config.Get, JWTSecret key constant
- `.claude/skills/gorm/SKILL.md` — Transaction pattern, DBClient struct field names

### Secondary (MEDIUM confidence)

- GORM transaction docs (via project skill) — transaction auto-rollback semantics confirmed

### Tertiary (LOW confidence)

None.

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — no new dependencies, all patterns verified in existing code
- Architecture: HIGH — all changes traced to specific lines in specific files
- Pitfalls: HIGH — all pitfalls found by direct inspection of source code, not inference

**Research date:** 2026-03-22
**Valid until:** 2026-06-22 (stable Go codebase — no moving targets)
