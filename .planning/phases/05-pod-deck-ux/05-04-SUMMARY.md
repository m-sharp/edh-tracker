---
plan: 05-04
phase: 05-pod-deck-ux
status: complete
completed: 2026-03-26
commits:
  - c2ac214
---

# Plan 05-04: Pod Players Tab Redesign — Summary

## What Was Built

**PlayersTab redesigned** (`app/src/routes/pod/PlayersTab.tsx`) — Replaced the flat List/ListItem layout with a Paper card per player. Each card shows:

- Player name as a link to their profile page
- "Manager" chip for pod managers
- TooltipIconButton actions (PersonAdd for promote, PersonOff for remove) — manager-only
- Pod-scoped stats row: record (W-L-D), points, kills — all from the new backend GetStatsForPlayersInPod query from Plan 05-01
- Promote/remove dialogs updated to match UI-SPEC copywriting (player name in title, "Never mind" cancel, "Make Manager" / "Remove")

## Key Files

- `app/src/routes/pod/PlayersTab.tsx` — Full rewrite to card layout with pod-scoped stats

## Self-Check: PASSED

- Card-per-player layout renders ✓
- Pod-scoped stats (kills, record, points) displayed ✓
- Manager chip shown for managers ✓
- Owner-only promote/remove actions (tooltip icon buttons) ✓
- Dialog copy matches UI-SPEC ✓
- TypeScript: exit 0 ✓
