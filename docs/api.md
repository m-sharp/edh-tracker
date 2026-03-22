# API Reference

All routes are prefixed `/api/`. The pattern is: **plural path for GET-all, singular path for GET-one and mutating operations**.

Routes marked **Auth** require a valid JWT cookie (set via OAuth login). Routes marked **Manager** additionally require the caller to be a pod manager.

## Auth

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/auth/google` | — | Initiate Google OAuth login; accepts optional `?redirect=<path>` |
| GET | `/api/auth/google/callback` | — | OAuth callback handler |
| POST | `/api/auth/logout` | — | Clear JWT cookie |
| GET | `/api/auth/me` | Auth | Get current authenticated user |

## Players

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/players` | — | List all players with stats; accepts `?pod_id=X` to filter by pod |
| GET | `/api/player?player_id=X` | — | Single player with stats |
| PATCH | `/api/player?player_id=X` | — | Update player name |

## Decks

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/decks` | — | List all active decks; accepts `?player_id=X` to filter by player |
| GET | `/api/deck?deck_id=X` | — | Single deck with commander info |
| POST | `/api/deck` | — | Create deck |
| PATCH | `/api/deck?deck_id=X` | — | Update deck (retire/un-retire or rename) |

## Games

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/games` | — | List games; accepts `?pod_id=X`, `?deck_id=X`, or `?player_id=X` |
| GET | `/api/game?game_id=X` | — | Single game with results |
| POST | `/api/game` | — | Create game + initial results |
| PATCH | `/api/game?game_id=X` | — | Update game description |
| DELETE | `/api/game?game_id=X` | Manager | Soft delete game (caller must be pod manager) |
| POST | `/api/game/result` | — | Add a result to an existing game |
| PATCH | `/api/game/result?result_id=X` | — | Update result (place, kills, deck) |
| DELETE | `/api/game/result?result_id=X` | Manager | Soft delete game result (caller must be pod manager) |

## Pods

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/pod?pod_id=X` | Auth | Get pod with members and roles |
| POST | `/api/pod` | Auth | Create new pod |
| PATCH | `/api/pod?pod_id=X` | Auth | Update pod name |
| DELETE | `/api/pod?pod_id=X` | Auth | Soft delete pod |
| POST | `/api/pod/player` | Auth | Add player to pod |
| PATCH | `/api/pod/player?pod_id=X&player_id=Y` | Auth | Promote player to manager |
| DELETE | `/api/pod/player?pod_id=X&player_id=Y` | Auth | Remove (kick) player from pod |
| POST | `/api/pod/invite` | Auth | Generate invite code for pod |
| POST | `/api/pod/join` | Auth | Join pod using invite code |
| POST | `/api/pod/leave` | Auth | Leave pod (caller leaves their own membership) |

## Formats

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/formats` | — | List formats (seeded: `commander`, `other`) |

## Commanders

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/api/commander` | — | List all commanders |
| POST | `/api/commander` | — | Create commander |
