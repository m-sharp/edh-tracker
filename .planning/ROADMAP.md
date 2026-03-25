# Roadmap: EDH Tracker — Launch Preparation

## Overview

The backend is mature and auth is complete, but several correctness gaps (missing authorization checks, non-transactional game creation, input validation), frontend structural debt, and a missing design language block the soft launch. This roadmap hardens the backend first, establishes the design language, restructures the frontend, ships the game model change, rounds out pod and deck UX, closes auth and test gaps, and finally verifies deployment — delivering a mobile-ready app that a small friend group can use at the table.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Backend Hardening** - Close authorization, transaction, validation, and performance gaps in the API
- [ ] **Phase 2: Design Language** - Define and apply the visual design system before any UI work begins
- [x] **Phase 3: Frontend Structure** - Refactor route files, build shared components, and fix routing bugs (completed 2026-03-24)
- [ ] **Phase 4: Game Model Change** - Remove player field from game entry and redesign the game form
- [ ] **Phase 5: Pod & Deck UX** - Complete pod and deck feature gaps and onboarding flow
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
**Plans:** 7/8 plans executed

Plans:
- [x] 03-01-PLAN.md — Create shared components (TabbedLayout, TooltipIcon, SvgIconPlayingCards) + move utilities to components/ (FEND-02, FEND-03)
- [x] 03-02-PLAN.md — Extract HomeView + fix loading flash + FEND-04 blank screen + Login/Join UI polish (FEND-04, FEND-05, DSNG-04)
- [x] 03-03-PLAN.md — Split PodView into pod/ subdirectory + TabbedLayout migration + Pod UI-SPEC fixes + AppBar mobile fix (FEND-01, DSNG-04)
- [x] 03-04-PLAN.md — Split PlayerView into player/ subdirectory + TabbedLayout migration + Player UI-SPEC fixes (FEND-01, DSNG-04)
- [x] 03-05-PLAN.md — Split DeckView into deck/ subdirectory + TabbedLayout migration + Deck UI-SPEC fixes + TooltipIcon usage (FEND-01, DSNG-04)
- [x] 03-06-PLAN.md — Move GameView + NewGameView to subdirectories + Game UI-SPEC fixes + TooltipIconButton usage (FEND-01, DSNG-04)
- [x] 03-07-PLAN.md — Fix SPA handler static asset fallback causing blank screen on refresh (FEND-04 gap closure)
- [ ] 03-08-PLAN.md — Login positioning, logout icon, tooltip placement, PlayersTab confirmation dialogs (DSNG-04 gap closure)

### Phase 4: Game Model Change
**Goal**: Games track decks only, the game form works cleanly on mobile, and the deck picker shows owner context
**Depends on**: Phase 3
**Requirements**: GAME-01, GAME-02, GAME-03, GAME-04
**Success Criteria** (what must be TRUE):
  1. Creating a game does not require selecting or entering a player — only decks are selected
  2. The deck picker in the game form displays the owner's name alongside the commander name (e.g., "Rakdos, Lord of Riots (Mike)")
  3. The game form is visually clean and navigable on a phone-sized screen without horizontal scrolling or overlapping elements
  4. The record display renders correctly for games with any number of players, not just 4
**Plans**: TBD
**UI hint**: yes

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
**Plans**: TBD
**UI hint**: yes

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

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Backend Hardening | 5/5 | Complete | ✓ |
| 2. Design Language | 3/3 | Complete | ✓ |
| 3. Frontend Structure | 7/8 | In Progress|  |
| 4. Game Model Change | 0/? | Not started | - |
| 5. Pod & Deck UX | 0/? | Not started | - |
| 6. Auth, Session & Test Coverage | 0/? | Not started | - |
| 7. Production Readiness | 0/? | Not started | - |

### Phase 8: User testing and iterative feedback resolution

**Goal:** [To be planned]
**Requirements**: TBD
**Depends on:** Phase 7
**Plans:** 0 plans

Plans:
- [ ] TBD (run /gsd:plan-phase 8 to break down)
