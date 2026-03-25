---
phase: 03-frontend-structure
plan: "08"
subsystem: frontend
tags: [ui-fix, uat-closure, login, appbar, tooltip, pod-players]
dependency_graph:
  requires: []
  provides: [UAT-gap-04, UAT-gap-05, UAT-gap-06, UAT-gap-08]
  affects: [app/src/routes/login.tsx, app/src/routes/root.tsx, app/src/components/TooltipIcon.tsx, app/src/routes/pod/PlayersTab.tsx]
tech_stack:
  added: []
  patterns: [confirmation-dialog, icon-button, tooltip-placement]
key_files:
  created: []
  modified:
    - app/src/routes/login.tsx
    - app/src/routes/root.tsx
    - app/src/components/TooltipIcon.tsx
    - app/src/routes/pod/PlayersTab.tsx
decisions:
  - Login page uses flex-start + top padding (pt xs:4/sm:8) instead of center alignment
  - Logout replaced with LogoutIcon in an IconButton wrapped in Tooltip to prevent AppBar clipping
  - TooltipIcon and TooltipIconButton both default placement to 'top' via optional prop
  - Pod PlayersTab uses single confirmAction state object to drive one shared Dialog for both Promote and Remove actions
metrics:
  duration: "3 minutes"
  completed_date: "2026-03-24"
  tasks_completed: 2
  files_modified: 4
---

# Phase 03 Plan 08: UAT UI Fix Closure Summary

**One-liner:** Four UAT gaps closed — login repositioned with flex-start + top padding, logout replaced with icon button, tooltips default to top placement, and pod player actions use contained buttons with confirmation dialogs.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Fix login positioning, logout icon, and tooltip placement | 64a4731 | login.tsx, root.tsx, TooltipIcon.tsx |
| 2 | Add contained buttons and confirmation dialogs to Pod PlayersTab | e4dea16 | PlayersTab.tsx |

## Changes Made

### Task 1 — Login, Logout Icon, Tooltip Placement

**login.tsx:** Changed outer Box `justifyContent` from `"center"` to `"flex-start"` and added `pt: { xs: 4, sm: 8 }`. Closes UAT gap #4.

**root.tsx:** Replaced text `Button` logout with `IconButton` containing `LogoutIcon` wrapped in a `Tooltip title="Logout"`. Added `IconButton`, `Tooltip` to MUI imports; added `LogoutIcon` from `@mui/icons-material/Logout`. Closes UAT gap #5.

**TooltipIcon.tsx:** Added optional `placement?: "top" | "bottom" | "left" | "right"` prop (defaulting to `"top"`) to both `TooltipIconProps` and `TooltipIconButtonProps`. Both components now pass `placement` to their `<Tooltip>` element. Closes UAT gap #6.

### Task 2 — Pod PlayersTab Contained Buttons + Confirmation Dialogs

**PlayersTab.tsx:** Added `variant="contained"` to both Promote and Remove buttons. Added `confirmAction` state (`{ type: "promote" | "remove"; player: PlayerWithRole } | null`). Both button `onClick` handlers now set `confirmAction` instead of directly calling handlers. A single shared `<Dialog>` at the component bottom renders contextually: title ("Promote player?" / "Remove player?"), body with player name and action description, Cancel + Confirm buttons. Confirm button color is `"primary"` for promote and `"error"` for remove. Closes UAT gap #8.

## Decisions Made

- **Single shared dialog state:** Rather than separate `promoting` and `removing` state variables, a single `confirmAction` with a `type` discriminant drives one Dialog — reduces state count and mirrors a pattern that scales cleanly to additional action types.
- **Tooltip on logout icon:** Required for discoverability since text label was removed; `Tooltip title="Logout"` provides the label accessibly.

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.
