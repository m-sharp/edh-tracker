---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
status: Phase 4 planned — ready to execute
stopped_at: Phase 4 plans created (3 plans, 2 waves)
last_updated: "2026-03-24T18:00:00.000Z"
progress:
  total_phases: 8
  completed_phases: 3
  total_plans: 20
  completed_plans: 17
---

# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-03-22)

**Core value:** A pod can sit down, record a game in under a minute, and immediately see accurate standings — on their phones.
**Current focus:** Phase 04 — game-model-change

## Current Position

Phase: 04
Plan: Ready to execute (3 plans created)

## Performance Metrics

**Velocity:**

- Total plans completed: 0
- Average duration: -
- Total execution time: 0 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| - | - | - | - |

**Recent Trend:**

- Last 5 plans: -
- Trend: -

*Updated after each plan completion*
| Phase 01 P05 | 3 | 3 tasks | 6 files |
| Phase 01 P04 | 8 | 2 tasks | 10 files |
| Phase 01-backend-hardening P01 | 9min | 2 tasks | 5 files |
| Phase 01-backend-hardening P03 | 15 | 1 tasks | 5 files |
| Phase 01-backend-hardening P02 | 14min | 2 tasks | 6 files |
| Phase 01-backend-hardening P06 | 4min | 2 tasks | 3 files |
| Phase 02-design-language P01 | 8min | 2 tasks | 5 files |
| Phase 03 P01 | 10min | 2 tasks | 10 files |
| Phase 03 P02 | 4min | 2 tasks | 5 files |
| Phase 03-frontend-structure P03 | 15min | 2 tasks | 7 files |
| Phase 03 P04 | 7min | 2 tasks | 6 files |
| Phase 03 P05 | 8min | 1 tasks | 5 files |
| Phase 03 P06 | 7min | 2 tasks | 4 files |
| Phase 03 P07 | 1 | 1 tasks | 1 files |
| Phase 03 P08 | 3min | 2 tasks | 4 files |

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- [Pre-roadmap]: Games track decks only, not players — player is implicit via deck ownership
- [Pre-roadmap]: Frontend design language to be defined before UI implementation begins
- [Pre-roadmap]: Soft launch to small friend group before broader rollout
- [Phase 01]: maxInviteUses hardcoded to 25 — pod_invite table has no max_used_count column; constant placed in functions.go
- [Phase 01]: Removed deck.GetAll entirely — no other callers existed once router default path was removed (plan 04)
- [Phase 01]: GetStatsForDecks added to GameResultRepository interface to maintain functional DI pattern (plan 04)
- [Phase 01-backend-hardening]: ErrForbidden placed in lib/errs (not lib/business) to avoid circular import — sub-packages cannot import their parent package in Go
- [Phase 01-backend-hardening]: errors.Is used at router layer to discriminate 403 (forbidden) vs 500 (DB error) — plain errors without ErrForbidden wrapper now correctly return 500
- [Phase 01-backend-hardening]: assertCallerOwnsDeck placed on DeckRouter (router layer owns auth) — business layer Update/SoftDelete/Retire no longer take callerPlayerID
- [Phase 01-backend-hardening]: DeckCreate ignores body player_id, uses JWT callerPlayerID exclusively (SEC-02)
- [Phase 01-backend-hardening]: Integration tests used for Create success path — transaction wrapper bypasses interface mocks; nil client safe for error-path unit tests
- [Phase 02-design-language]: CssBaseline placed inside ThemeProvider — required for dark body background to apply globally
- [Phase 02-design-language]: ThemeProvider/createTheme imported from @mui/material/styles (not @mui/material) — canonical sub-path per MUI v5 conventions
- [Phase 02-design-language]: DSNG-02 partially met — dark theme verified on Pod view at 375px; 3 mobile usability issues deferred to Phase 3 gap closure (touch tab scroll, AppBar title clipping, DataGrid narrow viewport)
- [Phase 03]: SvgIconPlayingCards extracted to components/ with optional fontSize prop; root.tsx wraps usage in Box for layout margin
- [Phase 03]: app/src/components/ established as canonical shared frontend code location; original utilities deleted from app/src/
- [Phase 03]: HomeView loading state initialized to true — CircularProgress renders before fetch starts; empty state only shown after fetch confirms zero pods (FEND-05)
- [Phase 03]: RequireAuth spinner wrapped in centered Box — FEND-04 blank screen was caused by unpositioned invisible spinner during auth check on refresh
- [Phase 03]: Button component=Link pattern used in JoinView for Go home — MUI styling on React Router navigation without anchor tag
- [Phase 03]: PodView passes all data as props to tab components; data loading stays in index.tsx, tabs are pure display
- [Phase 03]: AppBar title uses xs:none/sm:flex to prevent crowding with PodSelector + Avatar + Logout on 375px viewports (P-07)
- [Phase 03]: PlayerDecksTab/PlayerGamesTab accept playerId not player object — minimal prop footprint for data-fetching tabs
- [Phase 03]: NewGameView moved as straight file move only — Phase 4 redesigns it entirely (plan excluded from DSNG-04 audit per D-20)
- [Phase 03]: homepage set to '/' not '.' in CRA config — absolute asset paths prevent sub-route refresh blank screen
- [Phase 03]: Login page uses flex-start + top padding (pt xs:4/sm:8) instead of center alignment — closes UAT gap #4
- [Phase 03]: Logout replaced with LogoutIcon in IconButton wrapped in Tooltip — closes UAT gap #5
- [Phase 03]: TooltipIcon/TooltipIconButton both default placement='top' via optional prop — closes UAT gap #6
- [Phase 03]: Pod PlayersTab uses single confirmAction state to drive shared Dialog for Promote/Remove confirmation — closes UAT gap #8

### Roadmap Evolution

- Phase 8 added: User testing and iterative feedback resolution

### Pending Todos

- None currently

### Blockers/Concerns

- [Phase 7 risk]: CORS / cookie behavior in deployed environment is unverified — must be confirmed before launch
- [Phase 4 risk]: Game model change (remove player field) requires both API and UI changes — coordinate carefully

## Session Continuity

Last session: 2026-03-24T17:18:42.470Z
Stopped at: Phase 4 UI-SPEC approved
Resume file: .planning/phases/04-game-model-change/04-UI-SPEC.md
