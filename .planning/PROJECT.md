# EDH Tracker — Launch Preparation

## What This Is

EDH Tracker is a Magic: The Gathering Commander (EDH) game tracking app for small playgroups. It tracks players, decks, commanders, and game results with a points system based on kills and finish position. The app is being prepared for soft launch with a small friend group, with the intent to eventually open it more broadly.

## Core Value

A pod can sit down, record a game in under a minute, and immediately see accurate standings — on their phones.

## Requirements

### Validated

- ✓ Google OAuth login and JWT session management — existing
- ✓ Pod creation, invite links, and manager/member roles — existing
- ✓ Player profiles with computed stats (wins, kills, points) — existing
- ✓ Deck management (create, update, retire) with commander assignment — existing
- ✓ Game recording with per-deck kills and finish position — existing
- ✓ Points system: kills + place-based bonuses — existing
- ✓ Seeded player linking (connect existing player records to Google accounts) — existing
- ✓ Paginated game and deck list endpoints — existing

### Active

**Frontend design overhaul (highest priority):**
- [x] Define and apply an overarching visual design language for the app (mobile-first, MUI conventions solidified) — Validated in Phase 2: Design Language
- [ ] Mobile-friendly layout and interaction patterns across all views — app must work at the table on a phone
- [x] Refactor large route files into per-view subdirectories (`pod/`, `player/`, `deck/`, `game/` under `app/src/routes/`) — Validated in Phase 3: Frontend Structure
- [x] Shared tab component with query-string-persisted active tab, shared loading/error handling — Validated in Phase 3: Frontend Structure
- [x] Shared tooltip icon and tooltip icon button components — Validated in Phase 3: Frontend Structure
- [x] Fix: page refresh causes blank white screen (React Router client-side routing issue) — Validated in Phase 3: Frontend Structure
- [x] Fix: "No pods yet" flash before data loads on HomeView — Validated in Phase 3: Frontend Structure
- [x] Move pod creation out of player settings and into the pod page — Validated in Phase 5: pod-deck-ux
- [x] New user onboarding: clear path and UX for a user with no pods yet — Validated in Phase 5: pod-deck-ux
- [ ] Rebuild CLAUDE.md context section on frontend patterns (after refactor settles)

**Game model change:**
- [x] Remove player requirement from game entry — games track decks only (player is implicit via deck ownership) — Validated in Phase 4: Game Model Change
- [x] Deck picker in game form displays owner name (e.g., "Rakdos, Lord of Riots (Mike)") — Validated in Phase 4: Game Model Change
- [x] Remove/hide player field from game creation and result forms — Validated in Phase 4: Game Model Change

**Functional gaps:**
- [x] Player can create new decks via the UI — Validated in Phase 5: pod-deck-ux
- [x] Pod → Decks tab sorted by record by default — Validated in Phase 5: pod-deck-ux
- [x] Pod → Players view shows player records and points within the pod — Validated in Phase 5: pod-deck-ux
- [x] New game form complete redesign (currently "looks terrible") — Validated in Phase 4: Game Model Change
- [x] Record renderer supports any number of places (currently hardcoded to 4) — Validated in Phase 5: pod-deck-ux
- [x] Tooltip on deck commander update: "This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead." — Validated in Phase 5: Pod & Deck UX
- [x] Investigate and define retired deck behavior across all views (pod, player, game) — Validated in Phase 5: Pod & Deck UX

**Auth and session:**
- [ ] 401 interceptor in `http.ts`: call `logout()` from auth context on any 401 response, redirect to `/login`

**Backend correctness:**
- [x] Add pod-membership check on `POST /api/game` (currently any authenticated user can create a game in any pod) — Validated in Phase 1: Backend Hardening
- [x] Add `player_id`-from-context ownership check in `DeckCreate` (currently accepts `player_id` from body) — Validated in Phase 1: Backend Hardening
- [x] Wrap game creation in a DB transaction (orphaned game rows if result insert fails) — Validated in Phase 1: Backend Hardening
- [x] Fix `PromotePlayer`/`KickPlayer` returning 403 for all errors including DB errors — Validated in Phase 1: Backend Hardening
- [x] Add startup check rejecting JWT secrets shorter than 32 bytes — Validated in Phase 1: Backend Hardening
- [x] Validate `used_count` against a max on pod invite join (currently only expiry is checked) — Validated in Phase 1: Backend Hardening
- [x] Add input length validation on string fields (names, titles) — return 400 not 500 — Validated in Phase 1: Backend Hardening

**Performance:**
- [x] Batch deck stats queries — replace N+1 per-deck loop with single `WHERE deck_id IN (?)` query — Validated in Phase 1: Backend Hardening
- [x] Require at least one filter on `GET /api/decks`; remove or 404 the unfiltered path — Validated in Phase 1: Backend Hardening

**Test coverage:**
- [ ] `commander.Repository.GetAll` tests (acknowledged TODO in repo)
- [ ] Router tests for Commander and Format routers
- [ ] Auth router tests for Login, Logout, Me handlers
- [ ] Game creation authorization test (pod membership check)

**Production readiness:**
- [ ] Investigate CORS / nginx setup: determine if current config breaks cookies in deployed environment
- [ ] If needed: add dev proxy in `app/package.json` and nginx reverse proxy in Docker Compose
- [ ] Various small inline TODOs (see `.claude/plans/outstanding-todos.md` for full list)

### Out of Scope

- Public marketing or advertising — soft launch only; growth via word of mouth
- Pagination on player/commander/format list endpoints — acceptable at current playgroup scale
- Rate limiting on invite generation and pod join — post-launch concern
- Switching framework or bundler (CRA → Vite) — not needed for launch
- Third-party API access / external CORS callers — app is single-frontend only

## Context

**Codebase state:** Mature 4-layer Go backend (routers → business → repositories → DB), functional DI pattern, GORM + MySQL. Frontend is React 18 + TypeScript + MUI v5 + React Router v6. Full auth (Google OAuth + JWT cookies) and pod permission system are complete.

**Soft launch target:** Small friend group who will use the app at the table during games. Mobile experience is therefore critical — the primary use case is recording a game on a phone immediately after it ends.

**Frontend design gap:** The frontend is functional but visually bare. No overarching design system has been established. The restyling initiative must define the design language first, then implement it — mobile-first, using MUI's component system properly and consistently.

**Game model change (decided):** Games track decks, not players. Player is implicit via deck ownership. The deck picker in the game form must display the owner name alongside the commander name to avoid ambiguity when multiple players share similar decks.

**CORS situation:** Current backend adds CORS headers using `FRONTEND_URL`. This works but conflicts with credential-bearing cookies (browser rejects `* + credentials: true`). Current config uses a single exact-match origin, which is correct — but whether this actually breaks in the deployed environment needs to be confirmed before launch.

**Known tech debt to address before launch:** N+1 deck stats queries, unfiltered deck endpoint, game creation not transactional, missing authorization checks on game create and deck create.

## Constraints

- **Tech stack**: Go + Gorilla Mux + GORM + MySQL backend; React + MUI + React Router v6 frontend — no framework changes
- **Auth**: Google OAuth only — no email/password auth
- **Deployment**: Docker (separate images for API, React app, MySQL) — deployment shape must remain compatible
- **Compatibility**: No breaking changes to existing game/player/deck data already in the database

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Games track decks only, not players | Player is implicit via deck ownership; requiring player on each game slot added friction without value | — Implemented (Phase 4) |
| Deck picker displays owner name | Multiple players may use similar commanders; owner context prevents mis-selection | — Implemented (Phase 4) |
| Soft launch before full polish | Friend group provides real feedback; better to iterate on real usage than over-engineer before first user | — Pending |
| Frontend design language to be defined before implementation | Retrofitting a design system is harder than building to one; define first, implement per-phase | — Implemented (Phase 2) |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-27 after Phase 05 (pod-deck-ux) completion*
