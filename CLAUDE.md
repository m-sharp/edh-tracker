# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

EDH Tracker is a Magic: The Gathering Commander (EDH) game tracking app. It tracks players, decks, and game results with a points system (kills + place-based bonuses).

## Tech Stack

- **Backend**: Go 1.26, Gorilla Mux, Uber Zap logger, GORM
- **Database**: MySQL (`pod_tracker` database)
- **Frontend**: React 18 + TypeScript, React Router v6, Material-UI (MUI) v5
- **Deployment**: Docker (separate images for API, React app, and MySQL)

## Commands

### Backend
```bash
go mod vendor      # Vendor dependencies
go run main.go     # Run API server locally (requires DB env vars)
go vet ./lib/...   # Compile check (no binary output)
go test ./lib/...  # Run tests
```

### Frontend (from `app/`)
```bash
npm start        # Dev server
npm run build    # Production build
npm test         # Run tests
```

### Docker
```bash
# Build images
docker build -t edh-tracker .
docker build -f app/Dockerfile -t edh-tracker-app .

# Run API (requires DB + auth env vars)
docker run -p 8080:8081 \
  --env DBHOST=host.docker.internal \
  --env DBUSER=root \
  --env DBPASSWORD=<pass> \
  --env DBPORT=3306 \
  --env GOOGLE_CLIENT_ID=<client_id> \
  --env GOOGLE_CLIENT_SECRET=<client_secret> \
  --env OAUTH_REDIRECT_URL=<redirect_url> \
  --env JWT_SECRET=<secret> \
  --env FRONTEND_URL=http://localhost:8081 \
  --env DEV=1 edh-tracker

# Run React web app
docker run -p 8081:8081 -it edh-tracker-app:latest
```

## Architecture

See [`docs/architecture.md`](docs/architecture.md) for full detail on the 4-layer backend (routers → business → repositories → DB), functional DI pattern, and frontend structure.

### Quick Summary

```
lib/routers/       HTTP handlers
lib/business/      Domain logic + entity construction (functional DI closures)
lib/repositories/  Pure DB access returning Model types
lib/migrations/    Auto-run numbered schema migrations
lib/seeder/        Optional seed data (triggered by SEED env var)
```

`api.go` wires all routes. `lib/config.go` reads env vars.

**Required env vars**: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `JWT_SECRET`, `FRONTEND_URL`

**Optional env vars**: `SEED` (triggers data seeder), `DEV` (development mode — disables secure cookies)

## Data Model

See [`docs/data-model.md`](docs/data-model.md) for entity relationships, table schemas, Model/Entity split, points formula, and format validation rules.

## API Routes

See [`docs/api.md`](docs/api.md) for the full route table.

Route pattern: **plural path for GET-all, singular path for GET-one and POST** (e.g. `GET /api/players` vs `POST /api/player`).

## After Making API Changes

After modifying routes, handlers, repositories, or migrations, run the `/smoke-test` skill to rebuild the Docker image and verify that core endpoints are responding correctly.

<!-- GSD:project-start source:PROJECT.md -->
## Project

**EDH Tracker — Launch Preparation**

EDH Tracker is a Magic: The Gathering Commander (EDH) game tracking app for small playgroups. It tracks players, decks, commanders, and game results with a points system based on kills and finish position. The app is being prepared for soft launch with a small friend group, with the intent to eventually open it more broadly.

**Core Value:** A pod can sit down, record a game in under a minute, and immediately see accurate standings — on their phones.

### Constraints

- **Tech stack**: Go + Gorilla Mux + GORM + MySQL backend; React + MUI + React Router v6 frontend — no framework changes
- **Auth**: Google OAuth only — no email/password auth
- **Deployment**: Docker (separate images for API, React app, MySQL) — deployment shape must remain compatible
- **Compatibility**: No breaking changes to existing game/player/deck data already in the database
<!-- GSD:project-end -->

<!-- GSD:stack-start source:codebase/STACK.md -->
## Technology Stack

## Languages
- Go 1.26 - Backend API server (`main.go`, `api.go`, `lib/`)
- TypeScript 4.9.5 - Frontend React app (`app/src/`)
- SQL - Database schema migrations (`lib/migrations/`)
## Runtime
- Go 1.26.0 (alpine-based Docker image `golang:1.26.0-alpine3.23`)
- Node.js 18 (alpine-based Docker image `node:18-alpine`) — build only; runtime is a Go static file server
- Backend: Go modules with vendoring (`go.mod`, `go.sum`, `vendor/`)
- Frontend: npm (`app/package.json`, `app/package-lock.json`)
- Lockfile: `app/package-lock.json` present
## Frameworks
- `gorilla/mux` v1.8.1 - HTTP router and URL parameter extraction
- `gorm.io/gorm` v1.31.1 - ORM for MySQL access (`lib/db.go`, `lib/repositories/`)
- `gorm.io/driver/mysql` v1.6.0 - GORM MySQL driver
- React 18.0.0 - UI component framework (`app/src/`)
- React Router DOM 6.21.1 - Client-side routing
- Material-UI (MUI) v5.15.2 + `@mui/icons-material` v5.15.3 - UI component library
- `@mui/x-data-grid` v6.18.6 - Data grid component
- `@emotion/react` v11.11.3 + `@emotion/styled` v11.11.0 - CSS-in-JS (MUI peer dependency)
- `react-scripts` v5.0.1 (Create React App) - Dev server, production build, test runner
- `stretchr/testify` v1.11.1 - Assertions and `require` in Go tests (`assert`, `require` packages)
## Key Dependencies
- `go-sql-driver/mysql` v1.8.1 - Low-level MySQL driver (used through GORM); `lib/db.go`
- `golang-jwt/jwt/v5` v5.3.1 - JWT signing and validation; `lib/trackerHttp/auth.go`
- `golang.org/x/oauth2` v0.35.0 - OAuth2 client; `lib/routers/auth.go`
- `google/uuid` v1.6.0 - UUID generation (invite codes, nonces)
- `go.uber.org/zap` v1.26.0 - Structured logging throughout all layers
- `cloud.google.com/go/compute/metadata` v0.3.0 - Indirect; pulled in by `golang.org/x/oauth2` for Google endpoint support
## Configuration
- All config read from OS environment variables at startup via `lib/config.go`
- Required vars: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `JWT_SECRET`, `FRONTEND_URL`
- Optional vars: `DEV` (disables secure cookies in development), `SEED` (triggers data seeder on startup)
- Config struct: `lib.Config` with typed key constants (e.g., `lib.DBHost`, `lib.JWTSecret`)
- Backend: `Dockerfile` at project root — multi-stage not used, copies `vendor/` for reproducible builds
- Frontend + static server: `app/Dockerfile` — multi-stage: React build via `node:18-alpine`, Go static server via `golang:1.26.0-alpine3.23`
- MySQL: `mysql/Dockerfile` — `mysql:latest` base image with `mysql/custom.cnf`
## Platform Requirements
- Go 1.26+
- Node.js 18+
- MySQL instance (or Docker)
- All required env vars set
- Docker (three separate images: API, React app, MySQL)
- API listens on port 8081 (mapped to 8080 externally in example runs)
- Frontend static server listens on port 8081
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

## Package Naming
- One package per domain directory: `player`, `deck`, `game`, `gameResult`, `pod`, `format`, `commander`, `user`
- Multi-word domain names use camelCase for the package name: `gameResult`, `deckCommander`, `playerPodRole`, `podInvite`
- Test helpers live in their own packages: `lib/business/testHelpers`, `lib/repositories/testHelpers`
- The base GORM model lives in `lib/repositories/base`
## File Layout Per Domain
- `entity.go` — `Entity` struct, `ToEntity(model, ...)` conversion, `Validate()` method
- `functions.go` — Free constructors returning typed closures
- `types.go` — `*Func` type aliases and the `Functions` struct
- `model.go` — raw DB scan struct embedding `base.GormModelBase`
- `repo.go` — `Repository` struct with methods, `NewRepository(client)` constructor
## Naming Conventions
- Repository raw DB structs: `<domain>.Model` (e.g., `player.Model`, `gameResult.Model`)
- Business enriched objects: `<domain>.Entity` (e.g., `player.Entity`, `deck.Entity`)
- Business function type aliases: `<Operation>Func` (e.g., `GetAllFunc`, `CreateFunc`, `GetByIDFunc`)
- Business function containers: `Functions` struct in each domain package
- Interface names: `<Domain>Repository` defined in `lib/repositories/interfaces.go`
- Business free constructors match operation name exactly: `GetAll`, `GetByID`, `Create`, `Update`, `SoftDelete`
- Repository methods are lowercase verbs: `GetAll`, `GetById`, `GetByName`, `Add`, `BulkAdd`, `Update`, `SoftDelete`
- Note the inconsistency: repositories use `GetById` (lowercase `d`), business layer uses `GetByID` (uppercase `D`)
- Interfaces use PascalCase: `Player`, `Deck`, `Game`, `Pod`, `GameResult`
- API functions use PascalCase verbs: `GetPlayers`, `PostGame`, `PatchDeck`, `DeletePod`
- Verb prefix convention: `Get*` for reads, `Post*` for creates, `Patch*` for updates, `Delete*` for deletes
## Model vs Entity Split
- `<domain>.Model` — struct used to scan raw DB rows. Embeds `base.GormModelBase` which provides `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`. Lives in `lib/repositories/<domain>/model.go`.
- `<domain>.Entity` — enriched business object with computed fields, cross-domain data, and JSON tags. Lives in `lib/business/<domain>/entity.go`.
- Conversion is always done by `ToEntity(m Model, ...)` in the business layer entity file.
- Entity structs define a `Validate() error` method for input validation.
## Business Layer Conventions (Functional DI)
- Constructor function names match the operation: `GetAll`, `GetByID`, `Create`, `Update`, `SoftDelete`
- The closure signature is typed via `*Func` aliases defined in `types.go`
- All closures accept `context.Context` as first argument
- `lib/business/business.go` wires all constructors and exposes them as a `Business` struct of `Functions` structs
## Repository Interface Conventions
- `Add(ctx, ...)` always returns `(int, error)` where `int` is `LastInsertId()`
- `GetByName(ctx, name)` returns `(nil, nil)` when not found — no error for missing rows
- `GetById(ctx, id)` returns `(nil, nil)` when not found
- `SoftDelete(ctx, id)` sets `deleted_at`; all `GetAll`/`GetById` queries filter `deleted_at IS NULL`
- `BulkAdd` methods exist alongside singular `Add` for batch inserts
- Compile-time interface checks in `lib/repositories/repositories.go`:
## Route Naming Conventions
- **Plural path for GET-all:** `GET /api/players`, `GET /api/decks`, `GET /api/games`
- **Singular path for GET-one and POST:** `GET /api/player?player_id=1`, `POST /api/player`
- Sub-resources use path nesting: `POST /api/pod/player`, `DELETE /api/pod/player`, `POST /api/pod/invite`, `POST /api/pod/join`
- `POST` returns `201 Created` with no body for creates
- `PATCH` returns `200 OK` or `204 No Content` depending on whether a body is returned
- Query param names use `snake_case`: `player_id`, `pod_id`, `deck_id`, `game_id`
## Error Handling
- Errors are wrapped with context using `fmt.Errorf("failed to get player %d: %w", id, err)`
- Error messages include the entity type and ID for traceability
- Business functions propagate repo errors upward; routers map them to HTTP status codes
- "Not found" is not an error — `nil, nil` signals absence; routers return `404` explicitly
- Sentinel strings (not typed errors) are used for some domain errors: `"forbidden: ..."`, `"unexpected number of rows"`
- Router error responses use `trackerHttp.WriteError(log, w, statusCode, err, logMsg, clientMsg)`
- Async functions `throw new Error(...)` on bad HTTP status
- Some functions attach `.status` to the error object for caller inspection
## Logging
- Logger passed as `*zap.Logger` to constructors that need it (primarily `game` business functions)
- In tests, replace with `zap.NewNop()` to suppress output
- Router errors logged via `trackerHttp.WriteError(log, w, ...)` which calls `log.Error(logMsg, zap.Error(err))`
- No structured fields beyond `zap.Error(err)` in most log calls
## Import Organization
## Frontend Component Conventions
- Components are default exports from `app/src/routes/<name>.tsx`
- Component functions return `ReactElement`
- Route components consume loader data via `useLoaderData() as <Type>`
- Auth state accessed via `useAuth()` hook from `app/src/auth.tsx`
- All HTTP calls centralized in `app/src/http.ts`; never call `fetch` directly from components
- All TypeScript interfaces defined in `app/src/types.ts`; imported as named imports
- MUI components used throughout; `@mui/material` for layout/UI, `@mui/x-data-grid` for tables
- `credentials: "include"` must be set on every `fetch` call
## JSON Field Names
- Go structs use `snake_case` JSON tags matching DB column names: `json:"player_id"`, `json:"pod_ids"`, `json:"kill_count"`
- TypeScript interfaces mirror these snake_case field names
- No camelCase JSON in the API
## Comments
- Comments on exported identifiers are rare — reserved for non-obvious behaviour
- Interface methods with special semantics get inline doc comments:
- Compile-time checks get a comment: `// Compile-time interface satisfaction checks.`
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

## Pattern Overview
- Strict 4-layer backend: routers → business → repositories → DB
- Functional dependency injection — closures replace struct-based services
- No shared model layer; domain types split between raw DB scan structs and enriched entities
- Frontend is a single-page React app; all routing is client-side
## Backend Layers
- Purpose: HTTP handlers only — parse requests, call business functions, write responses
- Location: `lib/routers/`
- Contains: One `*Router` struct per domain (e.g., `PlayerRouter`, `DeckRouter`, `PodRouter`, `AuthRouter`)
- Depends on: `lib/business/` (via domain `Functions` structs), `lib/trackerHttp/` helpers
- Used by: `api.go` (route registration)
- Example: `lib/routers/player.go`, `lib/routers/auth.go`
- Purpose: Domain logic and entity construction; no HTTP knowledge
- Location: `lib/business/<domain>/`
- Contains: `functions.go` (constructors returning closures), `types.go` (function type aliases + `Functions` struct), `entity.go` (enriched output type + `ToEntity()` converter)
- Depends on: `lib/repositories/` interfaces (injected as constructor arguments)
- Used by: `lib/routers/` (holds `Functions` struct fields), `lib/business/business.go` (wired at startup)
- Example: `lib/business/player/functions.go`
- Purpose: Pure DB access; return raw Model types with no business logic
- Location: `lib/repositories/<domain>/`
- Contains: `repo.go` (GORM queries), `model.go` (DB scan struct embedding `base.GormModelBase`)
- Depends on: `lib/db.go` (`*lib.DBClient`)
- Used by: `lib/business/` functions (via interfaces in `lib/repositories/interfaces.go`)
- Example: `lib/repositories/player/`, `lib/repositories/deck/`
- Purpose: Shared primitives — config, logging, DB client, web server
- Location: `lib/config.go`, `lib/db.go`, `lib/log.go`, `lib/web.go`
- Contains: `lib.Config`, `lib.DBClient` (wraps `*gorm.DB`), `lib.Server` (Gorilla Mux wrapper)
## Functional DI Pattern
```go
```
## Cross-Domain Dependency Graph
```
```
## Data Flow — Request Lifecycle
- `api.go` auto-applies `RequireAuth` to all mutating methods unless `Route.NoAuth` is set
- CORS preflight (OPTIONS) handlers are registered automatically for each mutating path
- Router handler reads body, validates inputs, calls business closure, returns 201 (no body) for creates
## Model/Entity Split
- Embeds `base.GormModelBase` (`ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`)
- Maps directly to a database table (GORM `TableName()` method)
- May embed related models for GORM preloading (e.g., `deck.Model` embeds `player.Model`, `format.Model`)
- Example: `lib/repositories/player/model.go`
- Contains resolved names, computed stats, and other cross-domain data
- Serialized directly to JSON for API responses
- Constructed via `ToEntity(model, ...)` in the business layer
- Example: `lib/business/player/entity.go` — `Entity` includes `Stats` (computed from `gameResult.Aggregate`) and `PodIDs` (joined from pod repo)
## Authentication
- `RequireAuth` (`lib/trackerHttp/auth.go`): validates JWT, injects `userID`/`playerID` into context
- `OptionalAuth`: reads JWT if present but never rejects
- `trackerHttp.CallerPlayerID()`: helper used in handlers to extract playerID from context or write 401
## Frontend Architecture
| Path | Component | Loader/Action |
|---|---|---|
| `/` | `HomeView` (inline `index.tsx`) | none |
| `/login` | `LoginPage` | none |
| `/join` | `JoinView` | none |
| `/pod/:podId` | `PodView` | `podLoader` (parallel fetches) |
| `/pod/:podId/new-game` | `NewGameView` | `newGameLoader` + `createGame` action |
| `/pod/:podId/game/:gameId` | `GameView` | `gameLoader` |
| `/player/:playerId` | `PlayerView` | `GetPlayer` |
| `/player/:playerId/deck/:deckId` | `DeckView` | `GetDeck` |
## Error Handling
- Business functions wrap repo errors with `fmt.Errorf("...: %w", err)` to preserve stack context
- Routers use internal 500 for unexpected errors, 400 for bad inputs
- `lib.DBError` wraps failing SQL queries with the query string for debugging
- GORM `ErrRecordNotFound` is suppressed in logs via `quietLogger`
## Cross-Cutting Concerns
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd:profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
