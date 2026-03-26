---
phase: 05-pod-deck-ux
plan: "04"
subsystem: frontend
tags: [ui, players-tab, pod-view, card-layout, stats]
dependency_graph:
  requires: [05-01]
  provides: [POD-04]
  affects: [app/src/routes/pod/PlayersTab.tsx]
tech_stack:
  added: []
  patterns: [Paper card layout, TooltipIconButton for actions, pod-scoped stats display]
key_files:
  modified:
    - app/src/routes/pod/PlayersTab.tsx
decisions:
  - "Member role shows no chip (manager-only chip) per UI-SPEC — cleaner than showing both roles"
  - "bullet separator using unicode \\u2022 between stats items inline with Typography body2"
metrics:
  duration: "7min"
  completed: "2026-03-26T03:00:27Z"
  tasks_completed: 1
  files_modified: 1
---

# Phase 05 Plan 04: Pod Players Tab Card Layout Summary

Redesigned Pod Players tab from List/ListItem layout to Paper card layout with pod-scoped stats. Each player now displays their record, points, and kills within the specific pod, with icon-based promote/remove actions and UI-SPEC-compliant dialog copy.

## What Was Built

**PlayersTab.tsx rewrite** — replaced the `List` / `ListItem` / `ListItemText` layout with a vertical stack of `Paper` cards (elevation=2). Each card has two rows:

1. **Header row:** Player name as a React Router `Link`, optional "Manager" `Chip` (managers only; members show nothing), and action buttons (TooltipIconButton with PersonAdd/PersonOff icons, manager-only, hidden for self).
2. **Stats row:** `<Record>` component for pod-scoped W/L record, then points and kills separated by bullet characters (•).

Dialog copy updated to match UI-SPEC contract:
- Confirm titles include the player name: "Promote {name}?" / "Remove {name}?"
- Cancel button text: "Never mind"
- Promote action: "Make Manager"
- Remove action: "Remove" (color=error)

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — `p.stats.record`, `p.stats.points`, and `p.stats.kills` are wired to real pod-scoped data from the backend (GET /api/players?pod_id=N implemented in Plan 01).

## Self-Check: PASSED

- FOUND: app/src/routes/pod/PlayersTab.tsx
- FOUND: commit 37622d5
