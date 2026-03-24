---
phase: 02-design-language
plan: 01
subsystem: ui
tags: [react, typescript, mui, theme, dark-mode, typography, google-fonts]

# Dependency graph
requires: []
provides:
  - MUI createTheme dark theme with palette (#0f1117 bg, #1a1a2e paper, #c9a227 gold accent)
  - Josefin Sans Google Font loaded globally for headings
  - ThemeProvider wrapping entire app in index.tsx with CssBaseline inside
  - Gradient background retired; flat dark background applied via CssBaseline
  - root.tsx Container uses theme token bgcolor: background.default instead of hardcoded hex
affects:
  - 02-02 (mobile polish — Pod view validation target)
  - Phase 3 (per-view restyling builds on this theme foundation)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "MUI theme in app/src/theme.ts — single source for palette, typography, component overrides"
    - "ThemeProvider > CssBaseline > app — CssBaseline must be inside ThemeProvider for dark bg to apply"
    - "Theme tokens in sx props (bgcolor: background.default) — no hardcoded hex in component props"

key-files:
  created:
    - app/src/theme.ts
  modified:
    - app/src/index.tsx
    - app/public/index.html
    - app/src/styles.css
    - app/src/routes/root.tsx

key-decisions:
  - "CssBaseline placed inside ThemeProvider — required for dark body background to apply (Pitfall 1)"
  - "ThemeProvider and createTheme both imported from @mui/material/styles (not @mui/material) per Pitfall 5"
  - "gradientBackground removed from both index.html class attribute and styles.css definition atomically"

patterns-established:
  - "Pattern: all theme config in app/src/theme.ts; index.tsx only imports and wraps"
  - "Pattern: theme tokens in sx props instead of hardcoded hex values"

requirements-completed:
  - DSNG-01
  - DSNG-03

# Metrics
duration: 8min
completed: 2026-03-23
---

# Phase 02 Plan 01: Design Language Theme Foundation Summary

**MUI dark theme with gold accent (#c9a227), Josefin Sans headings, and ThemeProvider wiring that applies the design system globally across the entire app**

## Performance

- **Duration:** 8 min
- **Started:** 2026-03-23T22:49:58Z
- **Completed:** 2026-03-23T22:57:00Z
- **Tasks:** 2
- **Files modified:** 5 (1 created, 4 modified)

## Accomplishments

- Created `app/src/theme.ts` with full MUI dark palette, Josefin Sans/Roboto typography scale, and component overrides (AppBar dark navy, Button border-radius 6px, Tab no textTransform)
- Wired ThemeProvider around the entire app with CssBaseline correctly placed inside it, so the dark body background applies globally
- Loaded Josefin Sans from Google Fonts CDN and retired the gradient background CSS class from both index.html and styles.css
- Replaced the one hardcoded hex bgcolor in root.tsx with the `background.default` theme token

## Task Commits

Each task was committed atomically:

1. **Task 1: Create MUI dark theme file** - `152bb09` (feat)
2. **Task 2: Wire ThemeProvider, load Josefin Sans, retire gradient, fix root bgcolor** - `92ca9ae` (feat)

**Plan metadata:** (docs commit — see below)

## Files Created/Modified

- `app/src/theme.ts` — NEW: createTheme with dark palette, Josefin Sans/Roboto typography, MuiButton/MuiAppBar/MuiChip/MuiTab component overrides
- `app/src/index.tsx` — ThemeProvider wrapper added with CssBaseline inside it; theme imported from ./theme
- `app/public/index.html` — Josefin Sans Google Fonts link added; gradientBackground class removed from body
- `app/src/styles.css` — .gradientBackground CSS definition removed
- `app/src/routes/root.tsx` — Container bgcolor changed from `"#f0f5fa"` to `"background.default"`

## Decisions Made

- CssBaseline placed inside ThemeProvider (not outside) — this is required for the dark body background to apply; placing it outside means no theme context and browser default applies
- Both ThemeProvider and createTheme imported from `@mui/material/styles` not `@mui/material` — canonical sub-path avoids potential reconciliation issues
- gradientBackground retired in both locations atomically — removing only one would leave partial state

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Theme foundation is in place; all MUI components now inherit the dark palette and typography scale globally
- Plan 02-02 (mobile polish) validates the Pod view at 375px viewport with Tabs scrollable behavior
- Phase 3 per-view restyling can build against this theme; individual views will need audit per DSNG-04

## Self-Check: PASSED

- app/src/theme.ts: FOUND
- .planning/phases/02-design-language/02-01-SUMMARY.md: FOUND
- Task commit 152bb09: FOUND
- Task commit 92ca9ae: FOUND

---
*Phase: 02-design-language*
*Completed: 2026-03-23*
