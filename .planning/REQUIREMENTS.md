# Requirements: EDH Tracker — Launch Preparation

**Defined:** 2026-03-22
**Core Value:** A pod can sit down, record a game in under a minute, and immediately see accurate standings — on their phones.

## v1 Requirements

### Design

- [x] **DSNG-01**: App has a defined visual design language (color palette, typography, spacing tokens) implemented consistently across all views
- [ ] **DSNG-02**: All views are usable on mobile (phone-sized viewport) — layout adapts, touch targets are adequate, text is readable without zooming
- [x] **DSNG-03**: MUI components used properly and consistently — no ad-hoc styling where MUI has a clear pattern
- [x] **DSNG-04**: All views (Login, Home, Player, Deck, Game) are individually audited against the Phase 2 design language and receive view-specific layout, spacing, and typography improvements beyond what the global ThemeProvider provides automatically — no view is left with structural layout issues, misaligned spacing, or typography that conflicts with the design system

### Frontend Structure

- [x] **FEND-01**: Large route files refactored into per-view subdirectories (`pod/`, `player/`, `deck/`, `game/` under `app/src/routes/`)
- [x] **FEND-02**: Shared tab component used across Pod, Player, and Deck views — active tab persisted via query string, shared loading/error state
- [x] **FEND-03**: Shared tooltip icon and tooltip icon button components available and used where applicable
- [x] **FEND-04**: Page refresh no longer causes blank white screen
- [x] **FEND-05**: HomeView no longer flashes "No pods yet" before data loads

### Pods

- [ ] **POD-01**: Pod creation is accessible from the pod page — not buried in player settings
- [ ] **POD-02**: New users with no pods are guided to create or join one (no confusing empty state)
- [x] **POD-03**: Pod → Decks tab is sorted by record (win rate or points) by default
- [x] **POD-04**: Pod → Players tab shows each player's record and points within that pod

### Games

- [x] **GAME-01**: Games do not require a player field — decks are the unit of game entry; player is implicit via deck ownership
- [x] **GAME-02**: Deck picker in game form displays owner name alongside commander name (e.g., "Rakdos, Lord of Riots (Mike)")
- [x] **GAME-03**: New game form is visually clean and easy to use on mobile
- [x] **GAME-04**: Record renderer supports any number of places (not hardcoded to 4-player games)

### Decks

- [ ] **DECK-01**: Player can create a new deck from the UI (not just via data import)
- [x] **DECK-02**: Commander update field has tooltip: "This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead."
- [x] **DECK-03**: Retired deck visibility is defined and consistent across pod, player, and game views

### Authentication & Session

- [ ] **AUTH-01**: 401 responses from the API call `logout()` from auth context and redirect to `/login`
- [x] **AUTH-02**: Server rejects startup if JWT secret is shorter than 32 bytes

### Security

- [x] **SEC-01**: `POST /api/game` verifies the caller is a member of the target pod before creating the game
- [x] **SEC-02**: `POST /api/deck` uses the caller's player ID from the JWT context — ignores any `player_id` in the request body
- [x] **SEC-03**: Game creation (game row + result rows) is wrapped in a single DB transaction
- [x] **SEC-04**: Pod invite join validates max use count in addition to expiry
- [x] **SEC-05**: String field inputs (player name, pod name, deck name, etc.) validated for max length — returns 400 not 500 on overflow

### Performance

- [x] **PERF-01**: Deck stats are fetched in a single batch query, not one-per-deck in a loop
- [x] **PERF-02**: `GET /api/decks` requires at least one filter (`pod_id` or `player_id`) — unfiltered path removed or returns 400

### Test Coverage

- [ ] **TEST-01**: `commander.Repository.GetAll` has repository-layer tests
- [ ] **TEST-02**: Commander and Format routers have handler-level tests
- [ ] **TEST-03**: Auth router tests cover Login, Logout, and Me handlers
- [ ] **TEST-04**: Game creation router test verifies pod-membership authorization

### Production Readiness

- [ ] **INFRA-01**: CORS / cookie behavior verified in deployed environment; nginx reverse proxy added if needed
- [x] **INFRA-02**: `PromotePlayer` and `KickPlayer` return correct HTTP codes — 403 for authorization failures, 500 for infrastructure errors
- [ ] **INFRA-03**: `API_BASE_URL` in `app/src/http.ts` works correctly in production — either relative URLs (if nginx proxy) or configurable via build-time env var (`REACT_APP_API_BASE_URL`)

## v2 Requirements

### Scaling

- **SCALE-01**: Pagination on player, commander, and format list endpoints
- **SCALE-02**: Rate limiting on invite code generation and pod join attempts
- **SCALE-03**: Request body size limits applied to all handlers

### Observability

- **OBS-01**: Structured error context for silent failures (e.g., enrichGameModels drops no games silently)
- **OBS-02**: JWT re-issue failure logged at warn level

### Polish

- **POL-01**: `sitemap.xml`, `robots.txt`, `favicon.ico` served by frontend server
- **POL-02**: `pod.SoftDelete` parameter cleanup (unused `callerPlayerID`)

## Out of Scope

| Feature | Reason |
|---------|--------|
| Public marketing or advertising | Soft launch only — growth via word of mouth |
| Email/password auth | Google OAuth is sufficient for this use case |
| Mobile native app | Web-first; mobile browser must work well enough |
| Multi-origin CORS | App is single-frontend only; no external API callers |
| Framework migration (CRA → Vite) | Not needed for launch; future optimization |
| Frontend automated tests | High effort for existing code; post-launch investment |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| DSNG-01 | Phase 2 | Complete |
| DSNG-02 | Phase 2 | Pending |
| DSNG-03 | Phase 2 | Complete |
| DSNG-04 | Phase 3 | Complete |
| FEND-01 | Phase 3 | Complete |
| FEND-02 | Phase 3 | Complete |
| FEND-03 | Phase 3 | Complete |
| FEND-04 | Phase 3 | Complete |
| FEND-05 | Phase 3 | Complete |
| POD-01 | Phase 5 | Pending |
| POD-02 | Phase 5 | Pending |
| POD-03 | Phase 5 | Complete |
| POD-04 | Phase 5 | Complete |
| GAME-01 | Phase 4 | Complete |
| GAME-02 | Phase 4 | Complete |
| GAME-03 | Phase 4 | Complete |
| GAME-04 | Phase 4 | Complete |
| DECK-01 | Phase 5 | Pending |
| DECK-02 | Phase 5 | Complete |
| DECK-03 | Phase 5 | Complete |
| AUTH-01 | Phase 6 | Pending |
| AUTH-02 | Phase 1 | Complete |
| SEC-01 | Phase 1 | Complete |
| SEC-02 | Phase 1 | Complete |
| SEC-03 | Phase 1 | Complete |
| SEC-04 | Phase 1 | Complete |
| SEC-05 | Phase 1 | Complete |
| PERF-01 | Phase 1 | Complete |
| PERF-02 | Phase 1 | Complete |
| TEST-01 | Phase 6 | Pending |
| TEST-02 | Phase 6 | Pending |
| TEST-03 | Phase 6 | Pending |
| TEST-04 | Phase 6 | Pending |
| INFRA-01 | Phase 7 | Pending |
| INFRA-02 | Phase 1 | Complete |
| INFRA-03 | Phase 7 | Pending |

**Coverage:**
- v1 requirements: 36 total
- Mapped to phases: 36
- Unmapped: 0

---
*Requirements defined: 2026-03-22*
*Last updated: 2026-03-22 after roadmap creation*
