# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

EDH Tracker is a Magic: The Gathering Commander (EDH) game tracking app. It tracks players, decks, and game results with a points system (kills + place-based bonuses).

## Tech Stack

- **Backend**: Go 1.26, Gorilla Mux, Uber Zap logger
- **Database**: MySQL (`pod_tracker` database)
- **Frontend**: React 18 + TypeScript, React Router v6, Material-UI (MUI) v5
- **Deployment**: Docker (separate images for API, React app, and MySQL)

## Commands

### Frontend (from `app/`)
```bash
npm start        # Dev server
npm run build    # Production build
npm test         # Run tests
```

### Backend
```bash
go mod vendor    # Vendor dependencies
go run main.go   # Run API server locally (requires DB env vars)
```

### Docker
```bash
# Build images
docker build -t edh-tracker .
docker build -f app/Dockerfile -t edh-tracker-app .

# Run API (requires DB env vars)
docker run -p 8080:8081 \
  --env DBHOST=host.docker.internal \
  --env DBUSER=root \
  --env DBPASSWORD=<pass> \
  --env DBPORT=3306 \
  --env DEV=1 edh-tracker
  
# Run React web app
docker run -p 8081:8081 -it edh-tracker-app:latest
```

## Architecture

### Backend (`lib/`)

The Go server is structured around three layers per entity:

1. **Routers** (`lib/routers/`): HTTP handlers — parse request, call model methods, write JSON response
2. **Models** (`lib/models/`): Business logic + DB queries — each entity has a `Provider` struct with a `*sql.DB`
3. **Migrations** (`lib/migrations/`): Numbered migration files run automatically on startup via `migrate.go`

`api.go` wires all routes via Gorilla Mux. `lib/http.go` contains CORS middleware. `lib/config.go` reads all configs from environment variables into a string map.

**Required env vars**: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT`

### Frontend (`app/src/`)

- `index.tsx` — React Router setup with all routes and their loaders
- `http.ts` — All API client functions (fetch wrappers); API base URL is hardcoded to `http://localhost:8080/api/`
- `types.ts` — TypeScript interfaces for all domain types
- `routes/` — One component per page (Players, Decks, Games, NewGame, etc.)

React Router data loaders (`loader` functions in `index.tsx`) fetch data before rendering. The frontend at `app/Dockerfile` serves the built React app.

### Data Model

```
Player → has many Decks
Deck → belongs to Player, has many GameResults
Game → has many GameResults
GameResult → belongs to Game + Deck; has place, kill_count, points
```

Points formula: `points = kills + place_bonus` where place bonuses are 3/2/1/0 for 1st–4th.

Stats (record, games, kills, points) are computed at query time in `lib/models/stats.go`, not stored in the DB.

### API Routes

All routes are prefixed `/api/`:
- `GET/POST /players`, `GET /player?player_id=X`
- `GET/POST /decks`, `GET /deck?deck_id=X`, `PATCH /deck?deck_id=X` (retire)
- `GET/POST /games`, `GET /game?game_id=X`

Query params are used for filtering (e.g. `GET /decks?player_id=X` for player's decks).
