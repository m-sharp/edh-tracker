---
name: smoke-test
description: Use this skill to build and smoke test the edh-tracker API after making changes to routes, handlers, models, or migrations. Triggers when the user asks to "smoke test", "test the API", "verify API changes", or after completing any modification to files in lib/routers/, lib/models/, lib/migrations/, or api.go.
version: 1.0.0
---

# EDH Tracker API Smoke Test

Run this after any API change to verify the build and core endpoints are working.

## Step 1 — Build

```bash
build-edh-tracker-api
```

## Step 2 — Start (detached)

```bash
stop-edh-tracker-api 2>/dev/null || true
run-edh-tracker-api-bg
```

## Step 3 — Wait for ready

```bash
until curl -sf http://localhost:8080/api/formats > /dev/null; do sleep 1; done && echo "API ready"
```

## Step 4 — Assert endpoints

```bash
# Formats: must return commander + other
curl -s http://localhost:8080/api/formats | jq .

# Players: must return a non-error response
curl -s http://localhost:8080/api/players | jq length
```

## Step 5 — Stop

```bash
stop-edh-tracker-api
```

---

## Route Reference

All routes are prefixed `/api/`. The pattern is: **plural path for GET-all, singular path for GET-one and POST**.

| Method | Path | Notes |
|---|---|---|
| GET | `/api/players` | list all players with stats |
| GET | `/api/player?player_id=X` | single player |
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

## Aliases Reference

These aliases are defined in `~/.zshrc`:

| Alias | What it does |
|---|---|
| `build-edh-tracker-api` | `docker build -t edh-tracker .` |
| `run-edh-tracker-api-bg` | Run API container detached on port 8080 |
| `run-edh-tracker-api-seed-bg` | Same as above with `SEED=1` |
| `stop-edh-tracker-api` | `docker stop edh-tracker-api && docker rm edh-tracker-api` |
