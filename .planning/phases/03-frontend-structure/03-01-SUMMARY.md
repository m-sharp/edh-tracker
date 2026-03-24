---
phase: 03-frontend-structure
plan: 01
subsystem: ui
tags: [react, typescript, mui, react-router]

# Dependency graph
requires:
  - phase: 02-design-language
    provides: MUI theme and component patterns established
provides:
  - TabbedLayout component with query-string-persisted tab state
  - TooltipIcon and TooltipIconButton shared components
  - SvgIconPlayingCards extracted as standalone component
  - app/src/components/ directory as canonical shared frontend code location
  - AsyncComponentHelper, StatColumns, Record, CommanderColumn, MatchesDisplay moved to components/
affects:
  - 03-02 and later plans that import from components/

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "TabbedLayout: query-string tab persistence via useSearchParams + useNavigate with replace:true"
    - "TooltipIcon: span wrapper for MUI Tooltip ref-forwarding compatibility"
    - "Shared component directory: app/src/components/ as canonical location"

key-files:
  created:
    - app/src/components/TabbedLayout.tsx
    - app/src/components/TooltipIcon.tsx
    - app/src/components/SvgIconPlayingCards.tsx
    - app/src/components/common.ts
    - app/src/components/stats.tsx
    - app/src/components/matches.tsx
  modified:
    - app/src/routes/root.tsx
    - app/src/routes/pod.tsx
    - app/src/routes/player.tsx
    - app/src/routes/deck.tsx

key-decisions:
  - "SvgIconPlayingCards extracted to components/ as default export with optional fontSize prop; root.tsx wraps usage in Box for layout margin"
  - "Original utility files (common.ts, stats.tsx, matches.tsx) deleted from app/src/ after move to components/"

patterns-established:
  - "TabbedLayout: query key per view (podTab/playerTab/deckTab), hidden prop filters tabs before index calculation"
  - "TooltipIcon uses span wrapper for MUI Tooltip ref-forwarding; TooltipIconButton uses IconButton directly (already forwards ref)"

requirements-completed: [FEND-02, FEND-03]

# Metrics
duration: 10min
completed: 2026-03-24
---

# Phase 03 Plan 01: Create Shared Component Library Summary

**TabbedLayout (query-string tabs), TooltipIcon/TooltipIconButton, SvgIconPlayingCards added to app/src/components/; three existing utilities moved from app/src/ to components/ with all import paths updated**

## Performance

- **Duration:** 10 min
- **Started:** 2026-03-24T02:52:46Z
- **Completed:** 2026-03-24T03:02:46Z
- **Tasks:** 2
- **Files modified:** 10

## Accomplishments
- Created app/src/components/ as canonical shared frontend code directory with 6 files
- Implemented TabbedLayout with query-string tab persistence (useSearchParams + useNavigate, replace:true, hidden tab filtering)
- Implemented TooltipIcon (enterTouchDelay=0 for mobile toggle) and TooltipIconButton (tap=click)
- Extracted SvgIconPlayingCards from root.tsx into standalone reusable component
- Moved AsyncComponentHelper, StatColumns/Record/CommanderColumn, MatchesDisplay to components/ and deleted originals
- Updated all downstream imports in pod.tsx, player.tsx, deck.tsx — TypeScript compiles cleanly

## Task Commits

Each task was committed atomically:

1. **Task 1: Create new shared components (TabbedLayout, TooltipIcon, SvgIconPlayingCards)** - `2373f7f` (feat)
2. **Task 2: Move shared utilities to components/ and update all downstream imports** - `b1c036e` (feat)

**Plan metadata:** (docs commit follows)

## Files Created/Modified
- `app/src/components/TabbedLayout.tsx` - Shared tab component with query-string-persisted active tab; supports hidden tabs and loading state
- `app/src/components/TooltipIcon.tsx` - TooltipIcon (info icon with tap-to-toggle tooltip) and TooltipIconButton (icon button with hover tooltip)
- `app/src/components/SvgIconPlayingCards.tsx` - Playing cards SVG icon extracted from root.tsx; accepts optional fontSize prop
- `app/src/components/common.ts` - AsyncComponentHelper moved from app/src/common.ts
- `app/src/components/stats.tsx` - Record, RecordComparator, StatColumns, CommanderColumn moved from app/src/stats.tsx
- `app/src/components/matches.tsx` - MatchesDisplay, MatchUpDisplay moved from app/src/matches.tsx
- `app/src/routes/root.tsx` - Replaced inline SvgIconPlayingCards with import from components/; added Box wrapper for layout margin
- `app/src/routes/pod.tsx` - Updated StatColumns import to ../components/stats
- `app/src/routes/player.tsx` - Updated common/matches/stats imports to ../components/
- `app/src/routes/deck.tsx` - Updated common/matches/stats imports to ../components/

## Decisions Made
- SvgIconPlayingCards extracted with optional `fontSize` prop (layout-specific `mr: 2` removed from component; root.tsx wraps usage in `<Box sx={{ display: "flex", mr: 2 }}>`)
- Original utility files deleted from app/src/ after copy to components/ (not kept as re-exports, clean break)

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- app/src/components/ established with all 6 planned files
- TabbedLayout ready for use in plans 03-02 through 03-04 (pod/player/deck restructure)
- SvgIconPlayingCards ready for use in login.tsx and join.tsx per UI-SPEC issues L-04 and J-01

## Self-Check: PASSED

All 6 component files confirmed present in app/src/components/. All 3 original files confirmed deleted from app/src/. Both task commits (2373f7f, b1c036e) verified in git log.

---
*Phase: 03-frontend-structure*
*Completed: 2026-03-24*
