# Architecture

**Analysis Date:** 2026-03-22

## Pattern Overview

**Overall:** Layered monolith (backend) + SPA (frontend), communicating over a REST API

**Key Characteristics:**
- Strict 4-layer backend: routers → business → repositories → DB
- Functional dependency injection — closures replace struct-based services
- No shared model layer; domain types split between raw DB scan structs and enriched entities
- Frontend is a single-page React app; all routing is client-side

## Backend Layers

**Routers (`lib/routers/`):**
- Purpose: HTTP handlers only — parse requests, call business functions, write responses
- Location: `lib/routers/`
- Contains: One `*Router` struct per domain (e.g., `PlayerRouter`, `DeckRouter`, `PodRouter`, `AuthRouter`)
- Depends on: `lib/business/` (via domain `Functions` structs), `lib/trackerHttp/` helpers
- Used by: `api.go` (route registration)
- Example: `lib/routers/player.go`, `lib/routers/auth.go`

**Business (`lib/business/`):**
- Purpose: Domain logic and entity construction; no HTTP knowledge
- Location: `lib/business/<domain>/`
- Contains: `functions.go` (constructors returning closures), `types.go` (function type aliases + `Functions` struct), `entity.go` (enriched output type + `ToEntity()` converter)
- Depends on: `lib/repositories/` interfaces (injected as constructor arguments)
- Used by: `lib/routers/` (holds `Functions` struct fields), `lib/business/business.go` (wired at startup)
- Example: `lib/business/player/functions.go`

**Repositories (`lib/repositories/<domain>/`):**
- Purpose: Pure DB access; return raw Model types with no business logic
- Location: `lib/repositories/<domain>/`
- Contains: `repo.go` (GORM queries), `model.go` (DB scan struct embedding `base.GormModelBase`)
- Depends on: `lib/db.go` (`*lib.DBClient`)
- Used by: `lib/business/` functions (via interfaces in `lib/repositories/interfaces.go`)
- Example: `lib/repositories/player/`, `lib/repositories/deck/`

**Infrastructure (`lib/`):**
- Purpose: Shared primitives — config, logging, DB client, web server
- Location: `lib/config.go`, `lib/db.go`, `lib/log.go`, `lib/web.go`
- Contains: `lib.Config`, `lib.DBClient` (wraps `*gorm.DB`), `lib.Server` (Gorilla Mux wrapper)

## Functional DI Pattern

Each business operation is implemented as a free constructor that captures dependencies and returns a typed closure:

```go
// types.go — type aliases
type GetAllFunc func(ctx context.Context) ([]Entity, error)

type Functions struct {
    GetAll    GetAllFunc
    GetByID   GetByIDFunc
    Create    CreateFunc
    // ...
}

// functions.go — constructor returns closure
func GetAll(playerRepo repos.PlayerRepository, ...) GetAllFunc {
    return func(ctx context.Context) ([]Entity, error) {
        // uses captured repos
    }
}
```

`business.NewBusiness()` in `lib/business/business.go` wires all constructors, passing repository interfaces. The resulting `*Business` struct holds one `Functions` struct per domain. Routers extract the domain `Functions` struct at construction and call closures directly.

**Testing benefit:** Router tests inject mock closures directly into `Functions` struct fields — no mock struct type is needed.

## Cross-Domain Dependency Graph

```
player, format, commander, pod, user   (leaf — no cross-domain deps)
  ↓
deck  (uses GetPlayerNameFunc, GetByIDFunc/format, GetCommanderNameFunc)
  ↓
game  (uses GetDeckNameFunc, GetCommanderEntryFunc, GetByGameIDFunc/gameResult)
```

Shared function types are constructed once in `business.NewBusiness()` and passed to multiple domain constructors (e.g., `getFormat` and `getCommanderName` are reused across `deck` and `game`).

## Data Flow — Request Lifecycle

**GET request:**
1. Gorilla Mux matches route → dispatches to router handler
2. `api.go:SetupRoutes` has applied middleware chain: `GNUMiddleware → CORSMiddleware → [RequireAuth] → handler`
3. `RequireAuth` middleware (`lib/trackerHttp/auth.go`) validates `edh_session` JWT cookie; injects `userID`/`playerID` into context via `lib/utils/context.go`
4. Router handler calls business closure (e.g., `p.players.GetAll(ctx)`)
5. Business closure calls one or more repository methods, receiving `<domain>.Model` values
6. Business layer converts Models to Entities via `ToEntity()`, computing enrichments (stats, name lookups, etc.)
7. Router marshals Entity to JSON and writes response via `trackerHttp.WriteJson()`

**POST/PATCH/DELETE:**
- `api.go` auto-applies `RequireAuth` to all mutating methods unless `Route.NoAuth` is set
- CORS preflight (OPTIONS) handlers are registered automatically for each mutating path
- Router handler reads body, validates inputs, calls business closure, returns 201 (no body) for creates

## Model/Entity Split

**`<domain>.Model`** — raw DB scan struct located at `lib/repositories/<domain>/model.go`
- Embeds `base.GormModelBase` (`ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`)
- Maps directly to a database table (GORM `TableName()` method)
- May embed related models for GORM preloading (e.g., `deck.Model` embeds `player.Model`, `format.Model`)
- Example: `lib/repositories/player/model.go`

**`<domain>.Entity`** — enriched business object located at `lib/business/<domain>/entity.go`
- Contains resolved names, computed stats, and other cross-domain data
- Serialized directly to JSON for API responses
- Constructed via `ToEntity(model, ...)` in the business layer
- Example: `lib/business/player/entity.go` — `Entity` includes `Stats` (computed from `gameResult.Aggregate`) and `PodIDs` (joined from pod repo)

## Authentication

**Flow:**
1. Frontend redirects to `GET /api/auth/google` — server stores CSRF nonce in cookie, redirects to Google OAuth
2. Google redirects to `GET /api/auth/google/callback` — server validates CSRF, exchanges code for token, fetches Google profile
3. Server looks up or creates user (linking to existing seeded player by email if present); issues HMAC-SHA256 JWT
4. JWT stored as `edh_session` HTTP-only cookie; reissued (sliding 24h expiry) on every authenticated request
5. `GET /api/auth/me` returns the authenticated user's `user.Entity` (used by frontend `AuthProvider` on load)

**Middleware:**
- `RequireAuth` (`lib/trackerHttp/auth.go`): validates JWT, injects `userID`/`playerID` into context
- `OptionalAuth`: reads JWT if present but never rejects
- `trackerHttp.CallerPlayerID()`: helper used in handlers to extract playerID from context or write 401

## Frontend Architecture

**Framework:** React 18 + TypeScript, React Router v6 (`createBrowserRouter`), Material-UI v5

**Entry point:** `app/src/index.tsx` — creates the router, wraps app in `<AuthProvider>`, renders with `<RouterProvider>`

**Auth context:** `app/src/auth.tsx` — `AuthProvider` calls `GET /api/auth/me` on mount; exposes `useAuth()` hook returning `{ user, loading, logout }`

**Route protection:** `app/src/routes/RequireAuth.tsx` — wrapper component that redirects to `/login` when `user` is null

**Layout:** `app/src/routes/root.tsx` — `Root` component renders `<AppBar>` with pod selector + auth controls, and `<Outlet>` for nested route content

**HTTP client:** `app/src/http.ts` — one typed async function per API endpoint; all calls include `credentials: "include"` for cookie auth

**Types:** `app/src/types.ts` — TypeScript interfaces matching API Entity JSON shapes

**Route→component mapping:**
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

**Backend strategy:** Wrap-and-return errors up through layers; routers write HTTP errors using `trackerHttp.WriteError(log, w, statusCode, err, logMsg, clientMsg)`

**Patterns:**
- Business functions wrap repo errors with `fmt.Errorf("...: %w", err)` to preserve stack context
- Routers use internal 500 for unexpected errors, 400 for bad inputs
- `lib.DBError` wraps failing SQL queries with the query string for debugging
- GORM `ErrRecordNotFound` is suppressed in logs via `quietLogger`

## Cross-Cutting Concerns

**Logging:** Uber Zap (`go.uber.org/zap`); structured, named loggers per component (e.g., `log.Named("PlayerRouter")`)
**Validation:** Done in business entity layer (`Entity.Validate()`) or inline in router handlers
**Authentication:** JWT cookie middleware applied per-route in `api.go:SetupRoutes`
**CORS:** `CORSMiddleware` applied to all routes; CORS preflight (OPTIONS) registered automatically for mutating methods
**Soft deletes:** All major entities use GORM `DeletedAt` (from `base.GormModelBase`); cascade soft-deletes handled via GORM `AfterDelete` hooks in model files

---

*Architecture analysis: 2026-03-22*
