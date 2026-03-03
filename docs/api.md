# API Reference

All routes are prefixed `/api/`. The pattern is: **plural path for GET-all, singular path for GET-one and POST**.

| Method | Path | Notes |
|---|---|---|
| GET | `/api/players` | list all players with stats |
| GET | `/api/player?player_id=X` | single player |
| POST | `/api/player` | create player — 201, no body |
| GET | `/api/decks` | list all active decks |
| GET | `/api/deck?deck_id=X` | single deck |
| POST | `/api/deck` | create deck |
| PATCH | `/api/deck?deck_id=X` | retire deck |
| GET | `/api/games` | list all games |
| GET | `/api/game?game_id=X` | single game |
| POST | `/api/game` | create game + results |
| GET | `/api/formats` | list formats (`commander`, `other`) |
| GET | `/api/commander` | list commanders |
| POST | `/api/commander` | create commander |
| GET | `/api/pod?pod_id=X` | single pod |
| POST | `/api/pod` | create pod |
| POST | `/api/pod/player` | add player to pod |

Query params are used for filtering (e.g. `GET /api/decks?player_id=X` for a player's decks).
