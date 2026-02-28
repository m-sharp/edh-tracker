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

**Index convention**: Any column used in a WHERE clause that is not a primary key or foreign key must have an explicit index. Add indexes in a new migration whenever adding queries that filter on non-PK/FK columns.

### Frontend (`app/src/`)

- `index.tsx` — React Router setup with all routes and their loaders
- `http.ts` — All API client functions (fetch wrappers); API base URL is hardcoded to `http://localhost:8080/api/`
- `types.ts` — TypeScript interfaces for all domain types
- `routes/` — One component per page (Players, Decks, Games, NewGame, etc.)

React Router data loaders (`loader` functions in `index.tsx`) fetch data before rendering. The frontend at `app/Dockerfile` serves the built React app.

### Data Model

#### Entity Relationships

```
Player ──→ many Decks
       ──→ one User (optional)
       ──→ many PlayerPods ──→ Pod

Pod ──→ many Games
    ──→ many PlayerPods ──→ Player

Format ──→ many Decks
       ──→ many Games

Deck ──→ one Player
     ──→ one Format
     ──→ one DeckCommander (optional, via deck_commander join table)
     ──→ many GameResults

DeckCommander ──→ one Commander
              ──→ one Commander (partner_commander_id, nullable)

Game ──→ one Pod
     ──→ one Format
     ──→ many GameResults

GameResult ──→ one Game
           ──→ one Deck (includes commander info via JOIN)
```

#### Tables

| Table | Key Columns |
|---|---|
| `player` | id, name, deleted_at |
| `deck` | id, player_id, name, format_id, retired, deleted_at |
| `game` | id, description, pod_id, format_id, deleted_at |
| `game_result` | id, game_id, deck_id, place, kill_count, deleted_at |
| `format` | id, name — seeded with `commander` and `other` |
| `commander` | id, name (unique card name) |
| `deck_commander` | id, deck_id, commander_id, partner_commander_id (nullable) |
| `pod` | id, name, deleted_at |
| `player_pod` | id, pod_id, player_id (join table, unique constraint) |
| `user` | id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url |
| `user_role` | id, name — seeded with `admin` and `player` |

All tables use soft deletes (`deleted_at` nullable DATETIME) and track `created_at`/`updated_at`.

#### Key Model Structs (`lib/models/`)

- **`Deck`** — includes `Name`, `FormatID`, `FormatName`, `Retired`, and `Commanders *DeckCommanderEntry` (nullable, populated via LEFT JOIN on `deck_commander`)
- **`DeckCommanderEntry`** — `CommanderID`, `CommanderName`, `PartnerCommanderID *int`, `PartnerCommanderName *string`
- **`GameResult`** — includes `DeckName`, `CommanderName *string`, `PartnerCommanderName *string`, and computed `Points`
- **`Game`** — includes `PodID` and `FormatID`
- **`PlayerInfo`** — embeds `Player` + `Stats` + `PodIDs []int`
- **`DeckWithStats`** — embeds `Deck` + `Stats`
- **`Stats`** — `Record map[int]int` (place → count), `Games`, `Kills`, `Points` — computed at query time in `lib/models/stats.go`, never stored

#### Points Formula

`points = kills + place_bonus` where place bonuses are **3/2/1/0** for 1st–4th place.

#### Format Validation

When creating a game result, the deck's `format_id` must match the game's `format_id`, **except** when the format name is `"other"` (which skips the check).

#### Nullable Fields & SQL Scanning

Deck queries use LEFT JOIN to fetch commander data. `deckRow` is an internal scan struct using `sql.NullInt64`/`sql.NullString` for nullable commander fields, which are then converted to the `*DeckCommanderEntry` pointer on `Deck`.

### API Routes

All routes are prefixed `/api/`:
- `GET/POST /players`, `GET /player?player_id=X`
- `GET/POST /decks`, `GET /deck?deck_id=X`, `PATCH /deck?deck_id=X` (retire)
- `GET/POST /games`, `GET /game?game_id=X`
- `GET /formats`
- `GET/POST /commander`

Query params are used for filtering (e.g. `GET /decks?player_id=X` for player's decks).
