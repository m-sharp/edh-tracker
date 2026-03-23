# Phase 1: Backend Hardening - Context

**Gathered:** 2026-03-22
**Status:** Ready for planning

<domain>
## Phase Boundary

Close authorization, transaction, validation, and performance gaps in the API. All backend — no frontend work in this phase.

**In scope:**
- SEC-01: Pod membership check on `POST /api/game`
- SEC-02: `POST /api/deck` ignores `player_id` from body, uses caller's JWT player ID
- SEC-03: Game creation wrapped in a single DB transaction
- SEC-04: Pod invite join validates `used_count` against a max (not just expiry)
- SEC-05: String field inputs validated for max length — returns 400 not 500
- PERF-01: Deck stats fetched in a single batch query (no N+1 loop)
- PERF-02: `GET /api/decks` requires `pod_id` or `player_id` filter — unfiltered path removed or returns 400
- AUTH-02: Server startup rejects JWT secrets shorter than 32 bytes
- INFRA-02: `PromotePlayer`/`KickPlayer` return 403 for auth failures and 500 for DB errors

**Out of scope:** All frontend work, all test coverage (Phase 6), production readiness (Phase 7).

</domain>

<decisions>
## Implementation Decisions

### Transactions (SEC-03)

- **D-01:** Pass `*lib.DBClient` to the `game.Create` constructor. Inside the `db.Transaction()` callback, create tx-scoped repo copies via `NewRepository(&lib.DBClient{Db: tx})` for both `gameRepo` and `gameResultRepo`. GORM auto-rollbacks on a non-nil error return from the callback.
- **D-02:** No context changes needed. No modifications to existing repo methods. Only `game.Create` receives the `*lib.DBClient` constructor arg — no other business functions need it in this phase.

### Authorization Layer Standardization (SEC-01, SEC-02)

- **D-03:** **Router layer owns all authorization checks.** This is the standard going forward.
- **D-04:** Add pod membership check to `GameCreate` at the router layer: extract `callerPlayerID` via `trackerHttp.CallerPlayerID`, call `getPodRole(ctx, req.PodID, callerPlayerID)`, return 403 if the caller is not a member (role = "" or nil). Consistent with the existing `requirePodManager` pattern in the same file.
- **D-05:** **Phase 1 migrates the existing inconsistency:** `assertCallerOwnsDeck` currently lives in the deck business layer (`lib/business/deck/functions.go`). Move ownership enforcement up to `DeckRouter` — router calls `CallerPlayerID` and checks deck ownership before calling business layer `Update`/`SoftDelete`. The `assertCallerOwnsDeck` helper is removed from the business layer.
- **D-06:** For SEC-02 (deck create player_id): `DeckCreate` handler ignores `player_id` in the request body entirely. It calls `CallerPlayerID(w, r)` and passes that to `decks.Create`. The `createDeckRequest` struct's `PlayerID` field is removed (or ignored).

### Error Code Discrimination (INFRA-02)

- **D-07:** Introduce `var ErrForbidden = errors.New("forbidden")` in the `lib/business` package (or a shared `lib/business/errors.go` file).
- **D-08:** Update ALL existing `fmt.Errorf("forbidden: ...")` calls throughout the business layer to wrap this sentinel: `fmt.Errorf("forbidden: ...: %w", business.ErrForbidden)`.
- **D-09:** Routers that currently check for auth failures use `errors.Is(err, business.ErrForbidden)` to return 403 vs 500. Apply this consistently across all handlers — not just `PromotePlayer`/`KickPlayer`.

### Remaining Requirements (Implementation Notes for Planner)

These requirements are well-specified in REQUIREMENTS.md and do not have significant design ambiguity. Planner should implement per the requirements directly:

- **SEC-04:** Add `used_count < max_used_count` check in `JoinByInvite` business function (`lib/business/pod/functions.go`). The `pod_invite` table already has `used_count` and a max — confirm the max field name in the migration.
- **SEC-05:** Add max length validation to `Validate()` methods: `player.name` ≤ 50 chars (matches `VARCHAR(50)` in migration 7), `pod.name` ≤ 255 chars (migration 8), `deck.name` ≤ 255 chars (migration 14), `game.description` ≤ 256 chars (migration 4). Return 400 with a descriptive message.
- **PERF-01:** Add `GetStatsForDecks(ctx context.Context, deckIDs []int) (map[int]*Aggregate, error)` to `GameResultRepository` interface and `gameResult.Repository`. Replace the per-deck loop in `buildEntitiesWithStats` with a single batch call. SQL uses `WHERE deck_id IN (?)`.
- **PERF-02:** `GET /api/decks` unfiltered path (`default` case in `GetAll` and `getAllPaginated`) returns `400 Bad Request` with message: `"pod_id or player_id query param is required"`. Remove the `GetAll` business function and repo method if unused after this change.
- **AUTH-02:** In `main.go` (or `lib/config.go`), after reading `JWT_SECRET`, check `len(cfg.JWTSecret) < 32` and call `log.Fatal` with a descriptive error. This runs before any server starts.

### Claude's Discretion

- Whether to put `ErrForbidden` in `lib/business/errors.go` or inline in `lib/business/business.go` — either works.
- Whether `PERF-02` removes the `GetAll` repo method entirely or just blocks the route. Remove it if no other caller exists.
- `assertCallerOwnsDeck` helper: can be kept as a private helper in the deck router file if reuse within that file is useful, or inlined — Claude's call.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

No external specs or ADRs exist for this phase — requirements are fully captured in the decisions above and the requirements document.

### Requirements
- `.planning/REQUIREMENTS.md` — Full requirement specs for SEC-01 through INFRA-02 (Phase 1 requirements)

### Codebase Reference Points
- `lib/routers/game.go` — `GameCreate`, `requirePodManager`, existing `getPodRole` usage
- `lib/routers/deck.go` — `DeckCreate`, `DeckUpdate`, `DeckSoftDelete`
- `lib/business/deck/functions.go` — `assertCallerOwnsDeck` (to be migrated up to router)
- `lib/business/game/functions.go` — `Create` closure (transaction target)
- `lib/business/pod/functions.go` — `JoinByInvite` (SEC-04), existing `fmt.Errorf("forbidden: ...")` calls
- `lib/repositories/gameResult/repo.go` — `GetStatsForDeck` (PERF-01 base), interface additions needed
- `lib/repositories/interfaces.go` — `GameResultRepository` interface (PERF-01 new method)
- `lib/migrations/` — Column sizes for SEC-05 length validation (migrations 2, 4, 7, 8, 14)
- `main.go` — JWT secret startup guard (AUTH-02)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `trackerHttp.CallerPlayerID(w, r)` — extracts playerID from JWT context, writes 401 and returns false if missing. Used in `UpdateGame`, `DeleteGame`, all deck mutations — add to `GameCreate` and deck router for ownership checks.
- `g.getPodRole` — already injected in `GameRouter` via `biz.Pods.GetRole`. Usable directly in `GameCreate`.
- `requirePodManager` in `GameRouter` — pattern to follow for the new pod membership check (but membership check is weaker than manager check — just `role != ""`).
- `gameRepository.NewRepository(client)` / `gameResultRepository.NewRepository(client)` — accept `*lib.DBClient`, enabling tx-scoped instances inline.

### Established Patterns
- Business function constructors receive repo interfaces as args — `game.Create` adds `*lib.DBClient` as one more arg.
- Sentinel string errors: `fmt.Errorf("forbidden: ...: %w", ErrForbidden)` replaces existing bare `fmt.Errorf("forbidden: ...")` calls.
- `require.NoError` / table-driven tests in `lib/business/testHelpers` and `lib/repositories/testHelpers` — follow for any new tests.

### Integration Points
- `lib/business/business.go` — `NewBusiness(log, repos)` wires all constructors. `game.Create` constructor call here gets the `*lib.DBClient` arg added.
- `lib/repositories/interfaces.go` — `GameResultRepository` interface gets `GetStatsForDecks` added.
- `lib/business/game/types.go` — `CreateFunc` signature unchanged (callerPlayerID check stays at router, not in business).

</code_context>

<specifics>
## Specific Ideas

- Transaction pattern explicitly chosen over context injection: "Pass `*lib.DBClient` to Create constructor, create tx-scoped repo copies inline" — do not use context-based tx propagation for this phase.
- Auth standardization is a Phase 1 deliverable, not deferred: the `assertCallerOwnsDeck` migration from business to router layer is in scope.
- `ErrForbidden` sentinel applied broadly to all existing "forbidden:" errors — not just the `PromotePlayer`/`KickPlayer` fix. Full consistent rollout.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-backend-hardening*
*Context gathered: 2026-03-22*
