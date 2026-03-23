# Codebase Concerns

**Analysis Date:** 2026-03-22

## Tech Debt

**Unprotected Commander and Deck CREATE endpoints:**
- Issue: `POST /api/commander` and `POST /api/deck` are state-changing routes with no `RequireAuth: true` and no `NoAuth: true` override. By the convention in `api.go`, all POST/PATCH/DELETE routes are automatically auth-protected — but this relies on correct wiring at registration time. Deck create does call `CallerPlayerID` inside `DeckCreate` for PATCH/DELETE but not in the actual create handler, meaning an unauthenticated POST can reach the handler until `RequireAuth` middleware fires. For commander create, there is no ownership concept at all; any authenticated user can add a commander.
- Files: `lib/routers/commander.go`, `lib/routers/deck.go`
- Impact: Commander creation is effectively unrestricted beyond being authenticated. Deck creation accepts a `player_id` from the request body — a caller can create a deck on behalf of any player without ownership verification.
- Fix approach: Add explicit `player_id`-from-context ownership check in `DeckCreate`, matching the pattern used in `UpdateDeck` and `DeleteDeck`. Consider restricting commander creation to managers or to any authenticated user but adding rate limiting.

**`GET /api/decks` without filter is acknowledged as broken:**
- Issue: `d.decks.GetAll(ctx)` is called when no `pod_id` or `player_id` is provided. Two separate TODO comments in `lib/routers/deck.go` (lines 81 and 113) flag this as slow and say it "should probably not exist."
- Files: `lib/routers/deck.go` (lines 81, 113), `lib/business/deck/functions.go`
- Impact: The `buildEntitiesWithStats` helper fires one `GetStatsForDeck` query per deck (N+1). This scales linearly with the total number of decks across all users. The paginated fallback path also calls the unfiltered `GetAll` then paginates in memory — entirely defeating pagination.
- Fix approach: Require at least one filter (`pod_id` or `player_id`) on all deck list endpoints, matching the pattern already enforced on `GET /api/games`. Remove or 404 the unfiltered path.

**N+1 stats queries on deck list endpoints:**
- Issue: `buildEntitiesWithStats` in `lib/business/deck/functions.go` (line 186) issues one `GetStatsForDeck` call per deck. All filtered list functions (`GetAllForPlayer`, `GetAllByPod`, both paginated variants) call this helper.
- Files: `lib/business/deck/functions.go`, `lib/repositories/gameResult/repo.go`
- Impact: A pod with 20 members each with 5 decks triggers 100 sequential stat queries per deck-list request. This was noted in the existing TODO in `lib/repositories/gameResult/repo.go` (line 14).
- Fix approach: Add a batch stats query to `GameResultRepository` (e.g. `GetStatsForDecks(ctx, []int) map[int]*Aggregate`) and replace the per-deck loop with a single query.

**`enrichGameModels` silently drops failed games:**
- Issue: In `lib/business/game/functions.go` (line 231), if result enrichment fails for a game, that game is silently dropped from the response with only a warning log. The caller receives a shorter-than-expected list with no error signal.
- Files: `lib/business/game/functions.go` (lines 230–234)
- Impact: A transient DB error in result enrichment causes data loss from the API response without surfacing a 500.
- Fix approach: Return the error rather than silently dropping, or add a sentinel value that the frontend can display.

**`SoftDelete` for pod ignores `callerPlayerID` parameter:**
- Issue: `pod.SoftDelete` in `lib/business/pod/functions.go` (line 173) accepts a `callerPlayerID` argument but never uses it — it just calls `podRepo.SoftDelete(ctx, podID)`.
- Files: `lib/business/pod/functions.go` (lines 172–176)
- Impact: The ownership check is performed by `requireManager` in the router layer before calling this function, so security is not compromised in practice. However, the unused parameter is misleading and the business layer does not enforce ownership independently.
- Fix approach: Either remove the `callerPlayerID` parameter from the function signature, or add a re-check inside the function for defence-in-depth.

**`GetAll` endpoints expose all data without authentication:**
- Issue: `GET /api/players`, `GET /api/player`, `GET /api/formats`, `GET /api/commanders` are all unauthenticated GET routes. They do not set `RequireAuth: true`. Any unauthenticated request can enumerate all players, all formats, and all commanders.
- Files: `lib/routers/player.go`, `lib/routers/format.go`, `lib/routers/commander.go`
- Impact: Low severity for this app's use case (closed group), but notable if the app is ever exposed to the public internet. Player names and pod membership are enumerable without logging in.
- Fix approach: Decide on an access model. If the app is private-group only, add `RequireAuth: true` to all GET routes. If public browsing is intentional, document it explicitly.

**`POST /api/game` does not verify caller is a pod member:**
- Issue: `GameCreate` in `lib/routers/game.go` (line 220) checks authentication but does not verify the caller is a member or manager of the `pod_id` in the request. Any authenticated user can create a game in any pod they know the ID of.
- Files: `lib/routers/game.go` (line 220), `lib/business/game/functions.go`
- Impact: Users can submit game records to pods they are not members of.
- Fix approach: Add a pod-membership check (or `requireManager` check) at the start of `GameCreate`, consistent with how `UpdateGame`, `DeleteGame`, and game result mutations are handled.

**Unvalidated `player_id` in `AddResult` / game creation results:**
- Issue: When creating a game, game result inputs contain a `DeckID` and `Place`/`KillCount` but there is no server-side check that the deck belongs to a player who is a member of the pod the game is being recorded for.
- Files: `lib/business/game/functions.go` (line 102–153), `lib/business/gameResult/entity.go`
- Impact: It is possible to record game results using decks from unrelated players or pods.
- Fix approach: Add a check that each deck's `player_id` is a member of the game's `pod_id`.

## Known Bugs

**Home view loading blip:**
- Symptoms: The "No pods yet" message briefly flashes before redirect when a user has pods.
- Files: `app/src/index.tsx` (lines 42–44)
- Trigger: User is authenticated and has pods; the `useEffect` fires after first render.
- Workaround: None currently; acknowledged in TODO comment.

**`PromotePlayer` and `KickPlayer` use 403 for all errors:**
- Symptoms: Both handlers write `http.StatusForbidden` for _any_ error from `pods.PromoteToManager` / `pods.RemovePlayer`, including internal/DB errors.
- Files: `lib/routers/pod.go` (lines 315, 348)
- Trigger: A database failure during role check or removal returns 403 instead of 500.
- Workaround: None; distinguish forbidden errors from business logic vs. infrastructure errors.

## Security Considerations

**CORS origin is a single exact-match string:**
- Risk: `CORSMiddleware` in `lib/trackerHttp/http.go` sets `Access-Control-Allow-Origin` to the value of `FRONTEND_URL`. If that URL ever changes (e.g., staging vs production), the config must be updated. There is no wildcard, which is intentional and correct. However, the CORS middleware is applied universally but the `FRONTEND_URL` is the same string used for OAuth redirect; they are the same concept but conflated.
- Files: `lib/trackerHttp/http.go`, `api.go`, `lib/config.go`
- Current mitigation: Single origin only; `secure` cookie flag driven by `DEV` env var.
- Recommendations: Document that `FRONTEND_URL` is dual-purpose (CORS + OAuth redirect). Consider splitting into two separate env vars if origins diverge.

**JWT secret length not validated:**
- Risk: `JWT_SECRET` is read from env with no minimum-length or entropy check.
- Files: `lib/config.go`, `lib/trackerHttp/auth.go`
- Current mitigation: HS256 is used, which is HMAC-based; a short secret reduces security but doesn't break the auth flow.
- Recommendations: Add a startup check that rejects secrets shorter than 32 bytes.

**Silent JWT reissue failure:**
- Risk: `reissueSession` in `lib/trackerHttp/auth.go` (line 107) silently swallows the error when JWT re-signing fails. The session proceeds on the old token until expiry.
- Files: `lib/trackerHttp/auth.go` (lines 104–111)
- Current mitigation: The existing token remains valid.
- Recommendations: Log the error at warn level; currently it fails completely silently.

**Invite codes have no max-use limit:**
- Risk: `pod_invite` tracks `used_count` but `JoinByInvite` in `lib/business/pod/functions.go` (line 108) only checks expiry — it never checks `used_count` against a maximum.
- Files: `lib/business/pod/functions.go`, `lib/repositories/podInvite/repo.go`
- Current mitigation: Codes expire after 7 days (`expiresAt`).
- Recommendations: Add a max use count (e.g., 10) and check it before allowing a join.

**Redirect cookie path not validated:**
- Risk: In `lib/routers/auth.go` (lines 124–136), the `redirect` query parameter is stored directly as a cookie value and later joined to `FRONTEND_URL`. An attacker could craft a redirect to a path within the frontend URL. Cross-origin redirect is prevented by `url.JoinPath` using `FRONTEND_URL` as the base, but there is no allowlist of valid paths.
- Files: `lib/routers/auth.go`
- Current mitigation: `url.JoinPath` ensures the redirect stays within `FRONTEND_URL`.
- Recommendations: Low risk given the current implementation; document the behaviour.

## Performance Bottlenecks

**N+1 queries on all deck list calls:**
- Problem: Every deck in a list result triggers an independent `GetStatsForDeck` DB query.
- Files: `lib/business/deck/functions.go` (line 186), `lib/repositories/gameResult/repo.go` (line 118)
- Cause: Stats are computed per-deck in a loop with no batching.
- Improvement path: Replace the loop with a single `WHERE deck_id IN (?)` query and build an in-memory map.

**`GET /api/games` without soft-delete filter on count query:**
- Problem: Paginated game count queries in `lib/repositories/game/repo.go` count _all_ records for a pod. GORM's soft-delete on the model adds `deleted_at IS NULL` to Find queries but may not apply to explicit `.Count()` calls unless the model is specified.
- Files: `lib/repositories/game/repo.go`
- Cause: Need to confirm `.Count()` respects soft-delete scope.
- Improvement path: Verify the count query in `GetAllByPodPaginated` returns only undeleted records; add explicit `Where("deleted_at IS NULL")` if needed.

**Unfiltered `GetAll` deck endpoint loads every deck:**
- Problem: `GET /api/decks` with no filter calls `deckRepo.GetAll` which fetches every non-deleted deck from the database with full GORM preloads, then fires N stats queries.
- Files: `lib/routers/deck.go` (lines 81, 113), `lib/business/deck/functions.go` (line 46)
- Cause: No required filter params; acknowledged in TODO comments.
- Improvement path: Require `pod_id` or `player_id`; remove unfiltered path.

## Fragile Areas

**Migration count-based tracking:**
- Files: `lib/migrations/migrate.go`
- Why fragile: Migrations are tracked by row count in the `migration` table, not by migration number. The system compares the current count to the sorted index of all migrations. Inserting or deleting a migration record manually (or reordering migration keys in `getAllMigrations`) can cause incorrect migrations to be skipped or re-run.
- Safe modification: Always add migrations at the highest sequential number. Never renumber or remove existing migration files. The skipped-19 situation (documented: migration 19 was intentionally written and registered) means Migration 19 is now a real file and is registered — this is no longer a gap.
- Test coverage: No tests for the migration runner logic itself.

**Game creation is not transactional:**
- Files: `lib/business/game/functions.go` (lines 132–151)
- Why fragile: `gameRepo.Add` inserts the game row, then `gameResultRepo.BulkAdd` inserts results. If `BulkAdd` fails, the orphaned game row remains in the database with no results.
- Safe modification: Wrap both operations in a DB transaction.
- Test coverage: Business-layer test covers success path; partial-failure path is not tested.

**`GetAllByPod` deck function fetches player IDs then does a second query:**
- Files: `lib/business/deck/functions.go` (lines 139–159)
- Why fragile: `GetPlayerIDs` is called first; if the pod is large, this is a separate round trip before the deck query. No transaction or snapshot isolation; membership could change between the two queries.
- Safe modification: Acceptable for the app's scale; note the two-query pattern when modifying pod membership logic.

## Scaling Limits

**No request body size limit:**
- Current capacity: No `http.MaxBytesReader` is applied to any handler.
- Limit: A large request body (e.g., a game with thousands of results) will be read entirely into memory by `io.ReadAll`.
- Scaling path: Add `r.Body = http.MaxBytesReader(w, r.Body, maxBytes)` in `api.go` or in individual handlers.

**No rate limiting on invite code generation or pod join:**
- Current capacity: Unlimited invite generation and join attempts.
- Limit: An authenticated user can generate unlimited invite codes or attempt unlimited join operations.
- Scaling path: Add per-user rate limiting middleware on invite and join endpoints.

## Test Coverage Gaps

**No router tests for Commander or Format routers:**
- What's not tested: `CommanderRouter` handler logic (`GetAllCommanders`, `GetCommanderById`, `CommanderCreate`) and `FormatRouter` have no test files under `lib/routers/`.
- Files: `lib/routers/commander.go`, `lib/routers/format.go`
- Risk: Regressions in commander/format handlers would not be caught by automated tests.
- Priority: Medium

**`commander.Repository.GetAll` acknowledged as lacking tests:**
- What's not tested: The `GetAll` method has an inline TODO comment: "Should have tests."
- Files: `lib/repositories/commander/repo.go` (line 60)
- Risk: Refactors to the commander repository could silently break the GetAll path.
- Priority: Low

**Auth router tests are limited in scope:**
- What's not tested: `Login` handler (CSRF nonce generation, cookie setting, redirect), `Logout` handler, and `Me` handler are not covered by tests in `lib/routers/auth_test.go`. Only the `Callback` path has tests.
- Files: `lib/routers/auth_test.go`
- Risk: Changes to login/logout flow or the Me endpoint would go untested.
- Priority: Medium

**No frontend tests:**
- What's not tested: The entire `app/src/` directory has zero test files. All React components, the HTTP client layer (`app/src/http.ts`), auth context (`app/src/auth.tsx`), and route loaders are untested.
- Files: All files under `app/src/`
- Risk: Frontend regressions require manual verification; no safety net for refactors.
- Priority: High

**Game creation authorization path not tested:**
- What's not tested: No test verifies that `POST /api/game` rejects a caller who is not a member of the pod. The missing pod-membership check (see Security section) means this gap is both a coverage gap and an active bug.
- Files: `lib/routers/game_test.go`
- Risk: Authorization bypass would not be caught in CI.
- Priority: High

## Missing Critical Features

**No input length validation on string fields:**
- Problem: Player names, pod names, commander names, and deck names have no server-side maximum length enforcement. The DB schema may impose column limits, but the API returns a 500 instead of a 400 on overflow.
- Blocks: Clean user-facing error messages for oversized input.

**No pagination on player, commander, or format list endpoints:**
- Problem: `GET /api/players`, `GET /api/commanders`, and `GET /api/formats` are unbounded list endpoints with no pagination support.
- Blocks: These are low-risk now (small data sets) but would need pagination before the app scales beyond a small group.

**Retired deck visibility is undefined:**
- Problem: Deck retirement is implemented (`Retire` function, `retired` flag on the model) but the comment in `app/src/routes/deck.tsx` (line 175) notes this behaviour is unresolved: "Check retirement behavior - discuss where retired decks should and should not show."
- Files: `app/src/routes/deck.tsx`, `lib/business/deck/functions.go`
- Blocks: Consistent UI treatment of retired decks across pod, player, and game views.

**`app/main.go` frontend server has open TODOs:**
- Problem: `app/main.go` (lines 65–68) has TODOs for `sitemap.xml`, `robots.txt`, `favicon.ico`, and a question about whether the server is needed at all.
- Files: `app/main.go`
- Blocks: Minor polish items; the server is used for Docker deployment.

---

*Concerns audit: 2026-03-22*
