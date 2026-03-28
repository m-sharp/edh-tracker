# Roadmap: EDH Tracker — Launch Preparation

## Overview

The backend is mature and auth is complete, but several correctness gaps (missing authorization checks, non-transactional game creation, input validation), frontend structural debt, and a missing design language block the soft launch. This roadmap hardens the backend first, establishes the design language, restructures the frontend, ships the game model change, rounds out pod and deck UX, closes auth and test gaps, and finally verifies deployment — delivering a mobile-ready app that a small friend group can use at the table.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Backend Hardening** - Close authorization, transaction, validation, and performance gaps in the API (completed 2026-03-23)
- [x] **Phase 2: Design Language** - Define and apply the visual design system before any UI work begins (completed 2026-03-23)
- [x] **Phase 3: Frontend Structure** - Refactor route files, build shared components, and fix routing bugs (completed 2026-03-24)
- [x] **Phase 4: Game Model Change** - Remove player field from game entry and redesign the game form (completed 2026-03-24)
- [x] **Phase 5: Pod & Deck UX** - Complete pod and deck feature gaps and onboarding flow (completed 2026-03-27)
- [ ] **Phase 6: Auth, Session & Test Coverage** - 401 interceptor, JWT startup guard, and missing test suites
- [ ] **Phase 7: Production Readiness** - Verify CORS, cookies, and API base URL in the deployed environment

## Phase Details

### Phase 1: Backend Hardening
**Goal**: The API enforces correct authorization, uses transactions where needed, validates inputs, and eliminates N+1 query patterns
**Depends on**: Nothing (first phase)
**Requirements**: AUTH-02, SEC-01, SEC-02, SEC-03, SEC-04, SEC-05, PERF-01, PERF-02, INFRA-02
**Success Criteria** (what must be TRUE):
  1. Creating a game as a non-pod-member returns 403
  2. Creating a deck ignores any player_id in the request body and uses the authenticated caller's player ID
  3. If result row insertion fails during game creation, no orphaned game row is left in the database
  4. Submitting a player name, pod name, or deck name exceeding max length returns 400, not 500
  5. Deck stats for a pod load via a single batch query, not one query per deck
**Plans:** 5/5 plans complete

Plans:
- [x] 01-01-PLAN.md — ErrForbidden sentinel + pod router error discrimination (INFRA-02)
- [x] 01-02-PLAN.md — Game auth + transaction (SEC-01, SEC-03)
- [x] 01-03-PLAN.md — Deck auth migration to router layer (SEC-02)
- [x] 01-04-PLAN.md — Batch deck stats + unfiltered endpoint block (PERF-01, PERF-02)
- [x] 01-05-PLAN.md — Validation guards: JWT secret, string lengths, invite limit (AUTH-02, SEC-04, SEC-05)

### Phase 2: Design Language
**Goal**: The app has a defined visual design system that all subsequent UI work is built against
**Depends on**: Phase 1
**Requirements**: DSNG-01, DSNG-02, DSNG-03
**Success Criteria** (what must be TRUE):
  1. A documented color palette, typography scale, and spacing tokens exist and are applied consistently across at least one representative view
  2. The chosen MUI component patterns are used consistently — no inline style overrides where MUI provides a pattern
  3. At least one view is verified usable on a phone-sized viewport (touch targets adequate, text readable without zooming)
**Plans:** 3 plans

Plans:
- [x] 02-01-PLAN.md — Create MUI dark theme + wire ThemeProvider globally (DSNG-01, DSNG-03)
- [x] 02-02-PLAN.md — Pod view mobile polish + visual verification (DSNG-02)
- [x] 02-03-PLAN.md — Remove monospace fontFamily override from AppBar title (DSNG-03 gap closure)

### Phase 3: Frontend Structure
**Goal**: The frontend codebase is organized into per-view subdirectories with shared components, routing bugs are fixed, and all views are individually polished against the Phase 2 design language
**Depends on**: Phase 2
**Requirements**: FEND-01, FEND-02, FEND-03, FEND-04, FEND-05, DSNG-04
**Success Criteria** (what must be TRUE):
  1. Route files live in per-view subdirectories (pod/, player/, deck/, game/ under app/src/routes/) with no large monolithic route files remaining
  2. Pod, Player, and Deck views all use the shared tab component with active tab persisted in the query string
  3. Shared tooltip icon and tooltip icon button components exist and are used where applicable
  4. Refreshing any page in the app does not produce a blank white screen
  5. The HomeView does not flash "No pods yet" before pod data has loaded
  6. Login, Home, Player, Deck, and Game views have each been individually audited and improved — layout, spacing, and typography are intentional and consistent with the Phase 2 design system, not just passively inherited via ThemeProvider
**Plans:** 8/8 plans complete

Plans:
- [x] 03-01-PLAN.md — Create shared components (TabbedLayout, TooltipIcon, SvgIconPlayingCards) + move utilities to components/ (FEND-02, FEND-03)
- [x] 03-02-PLAN.md — Extract HomeView + fix loading flash + FEND-04 blank screen + Login/Join UI polish (FEND-04, FEND-05, DSNG-04)
- [x] 03-03-PLAN.md — Split PodView into pod/ subdirectory + TabbedLayout migration + Pod UI-SPEC fixes + AppBar mobile fix (FEND-01, DSNG-04)
- [x] 03-04-PLAN.md — Split PlayerView into player/ subdirectory + TabbedLayout migration + Player UI-SPEC fixes (FEND-01, DSNG-04)
- [x] 03-05-PLAN.md — Split DeckView into deck/ subdirectory + TabbedLayout migration + Deck UI-SPEC fixes + TooltipIcon usage (FEND-01, DSNG-04)
- [x] 03-06-PLAN.md — Move GameView + NewGameView to subdirectories + Game UI-SPEC fixes + TooltipIconButton usage (FEND-01, DSNG-04)
- [x] 03-07-PLAN.md — Fix SPA handler static asset fallback causing blank screen on refresh (FEND-04 gap closure)
- [x] 03-08-PLAN.md — Login positioning, logout icon, tooltip placement, PlayersTab confirmation dialogs (DSNG-04 gap closure)

### Phase 4: Game Model Change
**Goal**: Games track decks only, the game form works cleanly on mobile, and the deck picker shows owner context
**Depends on**: Phase 3
**Requirements**: GAME-01, GAME-02, GAME-03, GAME-04
**Success Criteria** (what must be TRUE):
  1. Creating a game does not require selecting or entering a player — only decks are selected
  2. The deck picker in the game form displays the owner's name alongside the commander name (e.g., "Rakdos, Lord of Riots (Mike)")
  3. The game form is visually clean and navigable on a phone-sized screen without horizontal scrolling or overlapping elements
  4. The record display renders correctly for games with any number of players, not just 4
**Plans:** 5/5 plans complete

Plans:
- [x] 04-01-PLAN.md — Backend AddResult cleanup + frontend types + dynamic Record component (GAME-01, GAME-04)
- [x] 04-02-PLAN.md — GameView AddResultModal player removal + TooltipIconButton enhancement (GAME-01, GAME-02)
- [x] 04-03-PLAN.md — NewGameView full redesign with stacked card layout (GAME-01, GAME-02, GAME-03)
- [x] 04-04-PLAN.md — AddResultModal/EditResultModal deck label + Place/Kills bounds (GAME-01, GAME-02) [gap closure]
- [x] 04-05-PLAN.md — New Game card inline remove button + button moved to pod header (GAME-03) [gap closure]

### Phase 5: Pod & Deck UX
**Goal**: Pod and deck workflows are complete, onboarding is clear, and views show meaningful stats
**Depends on**: Phase 4
**Requirements**: POD-01, POD-02, POD-03, POD-04, DECK-01, DECK-02, DECK-03
**Success Criteria** (what must be TRUE):
  1. A user can create a new pod from the pod page — the action is not hidden in player settings
  2. A new user with no pods sees a clear prompt to create or join one — no confusing empty state
  3. The Pod Decks tab is sorted by record by default
  4. The Pod Players tab shows each player's record and points within that pod
  5. A player can create a new deck directly from the UI, and the deck commander update field shows the disambiguation tooltip
**Plans:** 8/8 plans complete

Plans:
- [x] 05-01-PLAN.md — Backend: API contract fix (pod/deck create return IDs) + pod-scoped player stats (POD-04)
- [x] 05-02-PLAN.md — Frontend quick fixes: Record min 4 places, icon link, DecksTab sort, retired filter, cancel button (POD-03, DECK-02, DECK-03)
- [x] 05-03-PLAN.md — Pod creation flow: HomeView onboarding + AppBar pod switcher + CreatePodDialog (POD-01, POD-02)
- [x] 05-04-PLAN.md — Pod Players tab redesign with card layout and pod-scoped stats (POD-04)
- [x] 05-05-PLAN.md — Deck creation route /deck/new + DeckSettingsTab freeSolo + Add Deck button (DECK-01)
- [x] 05-06-PLAN.md — pod.Create transaction: AddPlayerToPod + atomic 3-write tx (POD-01, POD-02) [gap closure]
- [x] 05-07-PLAN.md — Commander POST returns {id: N} body + PostCommander http.ts fix (DECK-01) [gap closure]
- [x] 05-08-PLAN.md — Pod Decks client-side sort + retired filter via DataGrid + Cancel copy (POD-03, DECK-02, DECK-03) [gap closure]

### Phase 6: Auth, Session & Test Coverage
**Goal**: Session expiry is handled gracefully in the UI, the JWT secret has a startup guard, and missing test coverage is in place
**Depends on**: Phase 5
**Requirements**: AUTH-01, TEST-01, TEST-02, TEST-03, TEST-04
**Success Criteria** (what must be TRUE):
  1. When the API returns a 401, the frontend automatically logs the user out and redirects to /login — no manual intervention needed
  2. commander.Repository.GetAll has repository-layer tests
  3. Commander and Format router handlers have handler-level tests
  4. Auth router Login, Logout, and Me handlers have tests
  5. A test verifies that game creation is rejected for a caller who is not a pod member
**Plans**: TBD

### Phase 7: Production Readiness
**Goal**: The app runs correctly in the deployed Docker environment — cookies are sent, CORS is not blocking requests, and the API base URL resolves correctly
**Depends on**: Phase 6
**Requirements**: INFRA-01, INFRA-03
**Success Criteria** (what must be TRUE):
  1. Login, cookie-authenticated API calls, and logout work end-to-end in the deployed environment (not just locally)
  2. The React app can reach the API in production — either via nginx reverse proxy (relative URLs) or a build-time env var — with no hardcoded localhost references
**Plans**: TBD

### Phase 8: User testing and iterative feedback resolution

**Goal:** [To be planned]
**Requirements**: TBD
**Depends on:** Phase 7
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8

| Phase                                             | Plans Complete | Status      | Completed  |
|---------------------------------------------------|----------------|-------------|------------|
| 1. Backend Hardening                              | 5/5            | Complete    | ✓          |
| 2. Design Language                                | 3/3            | Complete    | ✓          |
| 3. Frontend Structure                             | 8/8            | Complete    | 2026-03-24 |
| 4. Game Model Change                              | 5/5            | Complete    | 2026-03-24 |
| 5. Pod & Deck UX                                  | 8/8            | Complete    | 2026-03-27 |
| 6. Auth, Session & Test Coverage                  | 0/?            | Not started | -          |
| 7. Production Readiness                           | 0/?            | Not started | -          |
| 8. User testing and iterative feedback resolution | 0/?            | Not started | -          |

## Backlog

### Phase 999.1: Full rename of "edh-tracker" to "pod-tracker" (BACKLOG)

**Goal:** [Captured for future planning]
**Requirements:** TBD
**Plans:** 0 plans

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.2: Add Maybe support to Go code (BACKLOG)

**Goal:** Extract pointer-based optional values into a generic `Maybe[T]` type that provides methods for checking presence (`IsPresent`, `IsNil`, `IsZero`) and safe value access — replacing ad-hoc nil pointer checks throughout the codebase.
**Requirements:** TBD
**Plans:** 0 plans

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.3: Ability to run locally without needing to do full docker builds (BACKLOG)

**Goal:** [Captured for future planning]
**Requirements:** TBD
**Plans:** 0 plans

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.4: Game Tracker Feature (BACKLOG)

**Goal:** A live game tracker page that lets players track life totals, poison counters, and commander damage during a game, then convert the session into a saved game record on completion.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- Entry point: new CTA button to the left of "New Game" on the pod view — the primary way to record games
- Deck selection screen before game starts
- Life counters (start at 40): large tap targets on left/right of total for mobile; prompt to eliminate at 0
- Poison counters (hidden by default, start at 0): prompt to eliminate at 10
- Commander damage tracker (hidden by default): one counter per opponent, labeled by deck; prompt to eliminate at 21 from any single source
- Context menu per deck with Eliminate button → red "lost" highlight state
- Kill assignment modal: when a deck is eliminated, prompt which deck gets the kill point
- End game triggered when all but one deck are eliminated
- End game screen: CSS fireworks for the winner, game description field, confirm kills/places, save → persist game + results in DB and redirect to game page, or discard
- Session state maintained via cookie; reset options: restart current game or restart with new decks
- Button to end game immediately at any time

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.5: Ability to Tag Decks (BACKLOG)

**Goal:** Add the ability to tag a deck with strings describing Magic: The Gathering archetypes (e.g. "Voltron", "Aggro", "Control", "Stax"). Tags are stored in a new table with a Deck ↔ Tag relation. Show tags on Deck data grids and Deck Details view (sortable where applicable). Tag input uses autocomplete over existing tags with the ability to add new ones.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- New `tags` table and `deck_tag` join table (Deck ↔ Tag many-to-many)
- Tag input: autocomplete over existing tag strings; allows creating new tags inline
- Display tags on Deck data grid (sortable column) and Deck Details view
- Tags persist as reusable strings shared across decks

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.7: Saltiness Tracking for Winning Decks (BACKLOG)

**Goal:** Capture how "salty" (bitter/mad/annoyed) players are with a winning deck at the end of a game, and surface that data in deck and player stats.
**Requirements:** TBD
**Plans:** 0 plans

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.8: Ability to Capture Kudos for Winning Decks (BACKLOG)

**Goal:** Allow players to give "Kudos" (virtual thumbs-ups / pats on the back) to winning decks at the end of a game, and surface Kudos counts in deck and player stats.
**Requirements:** TBD
**Plans:** 0 plans

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.9: Deck Tiering (BACKLOG)

**Goal:** Automatically extract or infer deck power tiers from existing point and record statistics, so players can see how their decks rank in relative power level without manual input.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- Derive tiers from existing collected data (points, win rates, kill counts, placement records)
- No manual tier tagging — tiers should emerge from the numbers
- Surface tier assignments on deck listings, deck detail views, and player stats

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.13: Custom Scoring Schemas (BACKLOG)

**Goal:** Allow pods to create and select custom scoring schemas as alternatives to the default points system, so groups with different playstyles can tailor how wins and kills are rewarded.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- Default system (kills + placement bonuses) remains the out-of-the-box schema
- Pods should be able to define their own schema (e.g. different point weights, no kill points, winner-takes-all)
- Schema selection scoped to the pod level (different pods can use different schemas)
- High complexity — requires schema storage model, formula evaluation, and retroactive recalculation concerns

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.12: ChangedByUserID Audit Tracking (BACKLOG)

**Goal:** Add a `changed_by_user_id` column to Pod, Game, GameResult, and any other models that can be mutated by multiple users, so every write is attributed to the user who made it.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- Applies to models writable by more than one user (Pod, Game, GameResult at minimum)
- Column records the user_id of whoever performed the last mutation
- Useful for audit trails, debugging disputes, and future moderation features
- Schema migration required per affected table; GORM models and business layer need updating

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.11: Info/About Page — Scoring System (BACKLOG)

**Goal:** Add an info/about page that explains the scoring system and its philosophy, so new players understand how points are calculated and why.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- Explain the points formula (kills + placement bonuses)
- Communicate the philosophy behind the scoring design
- Accessible to unauthenticated users (no login required to read)

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.10: Sweeps Feature (BACKLOG)

**Goal:** Show "Sweep" statistics for decks and games. A Sweep is defined as a deck winning a game and receiving all kill points as well (e.g. in a 4-player game, 1st place && >= 3 kills).
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- A Sweep = 1st place finish AND all kills in the game (e.g. ≥ 3 kills in a 4-player game)
- Surface sweep counts on deck stat grids, player stat views, and game result listings
- Indicate sweeps with a broom icon in data grids

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)

### Phase 999.6: Fetch and Display Card Images (BACKLOG)

**Goal:** Display card artwork for commanders and deck cards by building Scryfall image URLs from MTGJSON data, so players can visually identify cards throughout the app.
**Requirements:** TBD
**Plans:** 0 plans

**Captured Context:**
- Approach: use MTGJSON (https://mtgjson.com/getting-started/) to establish card name truth and build Scryfall image links per https://mtgjson.com/faq/#how-do-i-access-a-card-s-imagery
- Borrow reddit cardfetcher bot's approach for constructing Scryfall URLs from card names
- Once MTGJSON/Scryfall pipeline is established, use it to validate/confirm commander names at input time (autocomplete or validation against Scryfall data)
- Display card images on commander entries, deck detail views, and anywhere a commander name is shown

Plans:
- [ ] TBD (promote with /gsd:review-backlog when ready)
