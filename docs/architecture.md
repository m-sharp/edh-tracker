# Backend Architecture

## Layer Overview

The Go server uses a 4-layer architecture:

```
lib/routers/       HTTP handlers (parse request → call business → write JSON)
lib/business/      Domain logic: functional DI closures, entity construction
lib/repositories/  Pure DB access: SQL queries, return Model types
lib/migrations/    Numbered migration files, run automatically on startup
```

`api.go` wires all routes via Gorilla Mux. `lib/trackerHttp/http.go` contains CORS middleware and shared HTTP helpers. `lib/trackerHttp/auth.go` contains JWT middleware. `lib/config.go` reads env vars into a `*lib.Config` map.

**Required env vars**: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `JWT_SECRET`, `FRONTEND_URL`

**Optional env vars**: `SEED` (triggers data seeder), `DEV` (development mode — disables secure cookies)

## Repositories (`lib/repositories/`)

- One sub-package per domain: `player`, `deck`, `game`, `gameResult`, `pod`, `playerPodRole`, `podInvite`, `user`, `format`, `commander`, `deckCommander`
- Each has a `Repository` struct with a `*lib.DBClient` and methods that return `<domain>.Model` types
- `repositories.go` — `Repositories` struct bundles all concrete `*Repository` types
- `interfaces.go` — one interface per repository; used by the business layer for DI and testing

## Business Layer (`lib/business/`)

Functional DI pattern: each operation is a **free constructor function** that captures its dependencies via closure and returns a typed function.

```go
// Constructor (captures deps)
func GetAll(repo repositories.PlayerRepository, ...) GetAllFunc

// Typed function alias
type GetAllFunc func(ctx context.Context) ([]Entity, error)

// Functions struct groups all ops for a domain
type Functions struct {
    GetAll  GetAllFunc
    GetByID GetByIDFunc
    Create  CreateFunc
    ...
}
```

`business.go` — `Business` struct holds a `Functions` struct per domain; `NewBusiness(log, repos)` wires all closures.

### Cross-Domain Dependencies

Business functions receive peer domain functions as parameters (not repos), keeping the dependency graph acyclic:

```
player, format, commander, pod, user, gameResult  (no deps)
    ↓
  deck  (uses player.GetPlayerName, format.GetByID, commander.GetCommanderName)
    ↓
  game  (uses deck.GetDeckName, gameResult.EnrichModels)
```

## Routers (`lib/routers/`)

- Each router takes `*business.Business` and holds the relevant `<domain>.Functions` field
- `NewPlayerRouter(log, biz)`, `NewDeckRouter(log, biz)`, etc.
- Tests inject mock functions directly into `Functions` struct fields — no mock struct needed

## Dependency Wiring (`main.go` → `api.go`)

```
main.go:
  client := lib.NewDBClient(...)
  repos  := repositories.New(log, client)
  biz    := business.NewBusiness(log, repos)
  api    := NewApiRouter(cfg, log, biz)

api.go:
  NewApiRouter(cfg, log, biz *business.Business) — wires all routers
```

## Authentication & Authorization (`lib/trackerHttp/`, `lib/routers/auth.go`)

- **JWT middleware** (`lib/trackerHttp/auth.go`) — validates JWT cookie on each request; populates user info into request context
- **OAuth 2.0 flow** (`lib/routers/auth.go`) — Google OAuth login/callback/logout/me routes
- **Per-route flags** on `trackerHttp.Route`:
  - `RequireAuth: true` — enforces authentication; returns 401 if not logged in
  - `NoAuth: true` — skips JWT check for state-changing routes that don't need auth (e.g. logout)
  - Neither flag — JWT is parsed if present but not required (public read routes)
- **`trackerHttp.CallerPlayerID(w, r)`** — free function that extracts the authenticated player ID from context; writes 401 and returns false if missing
- **Pod role enforcement** — game and pod handlers call `requirePodManager` inline; no framework-level role middleware

## Seeder (`lib/seeder/`)

- Triggered by `SEED` env var (non-empty value)
- Guards against re-runs by checking if default pod (`"OG EDH Pod"`) exists
- Reads `./data/gameInfos.json`; caches player/deck/commander/format IDs in-memory
- Uses `repositories.Repositories` directly (not business layer)

## Migrations (`lib/migrations/`)

Numbered migration files run automatically on startup using **GORM**. Migrations always record themselves; there is no opt-out. Migration numbering is sequential — add a new file for any schema change. Currently 19 migrations are registered (1–19).

## Frontend (`app/src/`)

- `index.tsx` — React Router setup with all routes and their loaders
- `http.ts` — All API client functions (fetch wrappers); API base URL is hardcoded to `http://localhost:8080/api/`
- `types.ts` — TypeScript interfaces for all domain types
- `routes/` — One component per page (Players, Decks, Games, NewGame, etc.)

React Router data loaders (`loader` functions in `index.tsx`) fetch data before rendering.
