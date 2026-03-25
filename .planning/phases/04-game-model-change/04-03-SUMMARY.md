---
phase: 04-game-model-change
plan: 03
status: complete
commits:
  - 2fc7958 (Task 1 — initial rewrite)
  - 9b60e42 (Task 2 — UAT feedback fixes)
files_modified:
  - app/src/routes/new/index.tsx
  - lib/business/game/types.go
  - lib/business/game/functions.go
  - lib/business/game/functions_test.go
  - lib/routers/game.go
  - lib/routers/game_test.go
---

# 04-03 Summary: NewGameView Redesign

## What Was Built

Complete rewrite of `app/src/routes/new/index.tsx` as a stacked card layout with deck-only entry. All player picker references removed.

## Decisions Made

- **CardState[]** replaces `ResultsMap` + `numPlayers` — each card holds `{ key, deckId, place, kills }` as strings for clean TextField UX
- **FormControl + InputLabel** wrapper required for MUI Select label to render correctly with outlined variant and notched border
- **Place/Kills bounds**: `min=1/max=N` for Place, `min=0/max=N` for Kills, where N = `cards.length` (dynamic)
- **CreateFunc returns `(int, error)`**: game ID propagated from transaction through business → router → frontend for post-submit redirect
- **Redirect to game page**: `createGame` action reads `{ id }` from 201 JSON body and navigates to `/pod/:podId/game/:id`
- Card titles ("Deck 1", "Deck 2") removed — redundant with card layout
- Outer `px: 2` padding removed for better mobile width utilization

## UAT Result

Approved after feedback round addressing: Format label rendering, Place/Kills bounds, card title removal, padding, and post-submit redirect.
