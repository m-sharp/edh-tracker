---
created: 2026-03-24T16:07:54.537Z
title: Security review all API and frontend route authorization
area: api
files:
  - lib/routers/player.go
  - lib/routers/deck.go
  - lib/routers/game.go
  - lib/routers/pod.go
  - app/src/routes/Player.tsx
  - app/src/routes/Deck.tsx
  - app/src/routes/Game.tsx
---

## Problem

Endpoints and frontend routes may be accessible to authenticated users who don't have a legitimate relationship to the requested resource. Known example: `GET /player/:id` (and the frontend route `/player/:id`) can be accessed by any authenticated user, even if that player is not in any of the caller's pods. This leaks player data across pod boundaries.

The same pattern likely affects:
- Deck detail routes (`/player/:id/deck/:deckId`)
- Game detail routes (`/pod/:podId/game/:gameId`)
- Any API endpoint that accepts an ID query param without verifying pod membership

## Solution

1. Audit every API endpoint in `lib/routers/` — for each, determine what relationship check (if any) is needed
2. Audit frontend route loaders in `app/src/` — ensure they only fetch data the caller is authorized to see
3. Define the authorization model: likely "caller must share at least one pod with the target player/deck/game"
4. Implement missing checks, likely by adding pod-membership lookups in business layer functions or as middleware
5. Write tests covering the unauthorized access cases
