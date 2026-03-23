# Coding Conventions

**Analysis Date:** 2026-03-22

## Package Naming

- One package per domain directory: `player`, `deck`, `game`, `gameResult`, `pod`, `format`, `commander`, `user`
- Multi-word domain names use camelCase for the package name: `gameResult`, `deckCommander`, `playerPodRole`, `podInvite`
- Test helpers live in their own packages: `lib/business/testHelpers`, `lib/repositories/testHelpers`
- The base GORM model lives in `lib/repositories/base`

## File Layout Per Domain

Each domain package in `lib/business/<domain>/` contains exactly:
- `entity.go` — `Entity` struct, `ToEntity(model, ...)` conversion, `Validate()` method
- `functions.go` — Free constructors returning typed closures
- `types.go` — `*Func` type aliases and the `Functions` struct

Each domain package in `lib/repositories/<domain>/` contains:
- `model.go` — raw DB scan struct embedding `base.GormModelBase`
- `repo.go` — `Repository` struct with methods, `NewRepository(client)` constructor

## Naming Conventions

**Files:** lowercase, matching their package name (`entity.go`, `functions.go`, `types.go`, `repo.go`, `model.go`)

**Go Types:**
- Repository raw DB structs: `<domain>.Model` (e.g., `player.Model`, `gameResult.Model`)
- Business enriched objects: `<domain>.Entity` (e.g., `player.Entity`, `deck.Entity`)
- Business function type aliases: `<Operation>Func` (e.g., `GetAllFunc`, `CreateFunc`, `GetByIDFunc`)
- Business function containers: `Functions` struct in each domain package
- Interface names: `<Domain>Repository` defined in `lib/repositories/interfaces.go`

**Go Functions:**
- Business free constructors match operation name exactly: `GetAll`, `GetByID`, `Create`, `Update`, `SoftDelete`
- Repository methods are lowercase verbs: `GetAll`, `GetById`, `GetByName`, `Add`, `BulkAdd`, `Update`, `SoftDelete`
- Note the inconsistency: repositories use `GetById` (lowercase `d`), business layer uses `GetByID` (uppercase `D`)

**Variables/Parameters:** camelCase throughout

**TypeScript:**
- Interfaces use PascalCase: `Player`, `Deck`, `Game`, `Pod`, `GameResult`
- API functions use PascalCase verbs: `GetPlayers`, `PostGame`, `PatchDeck`, `DeletePod`
- Verb prefix convention: `Get*` for reads, `Post*` for creates, `Patch*` for updates, `Delete*` for deletes

## Model vs Entity Split

- `<domain>.Model` — struct used to scan raw DB rows. Embeds `base.GormModelBase` which provides `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`. Lives in `lib/repositories/<domain>/model.go`.
- `<domain>.Entity` — enriched business object with computed fields, cross-domain data, and JSON tags. Lives in `lib/business/<domain>/entity.go`.
- Conversion is always done by `ToEntity(m Model, ...)` in the business layer entity file.
- Entity structs define a `Validate() error` method for input validation.

Example from `lib/business/player/entity.go`:
```go
type Entity struct {
    ID        int         `json:"id"`
    Name      string      `json:"name"`
    Stats     stats.Stats `json:"stats"`
    PodIDs    []int       `json:"pod_ids"`
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`
}

func ToEntity(m repo.Model, agg *gameResult.Aggregate, podIDs []int) Entity { ... }
func (e Entity) Validate() error { ... }
```

## Business Layer Conventions (Functional DI)

Each business operation is a free constructor that captures its repository dependencies and returns a typed closure.

Pattern from `lib/business/player/functions.go`:
```go
func GetAll(playerRepo repos.PlayerRepository, gameResultRepo repos.GameResultRepository, podRepo repos.PodRepository) GetAllFunc {
    return func(ctx context.Context) ([]Entity, error) {
        // implementation
    }
}
```

- Constructor function names match the operation: `GetAll`, `GetByID`, `Create`, `Update`, `SoftDelete`
- The closure signature is typed via `*Func` aliases defined in `types.go`
- All closures accept `context.Context` as first argument
- `lib/business/business.go` wires all constructors and exposes them as a `Business` struct of `Functions` structs

## Repository Interface Conventions

Defined in `lib/repositories/interfaces.go`. Rules enforced across all repos:

- `Add(ctx, ...)` always returns `(int, error)` where `int` is `LastInsertId()`
- `GetByName(ctx, name)` returns `(nil, nil)` when not found — no error for missing rows
- `GetById(ctx, id)` returns `(nil, nil)` when not found
- `SoftDelete(ctx, id)` sets `deleted_at`; all `GetAll`/`GetById` queries filter `deleted_at IS NULL`
- `BulkAdd` methods exist alongside singular `Add` for batch inserts
- Compile-time interface checks in `lib/repositories/repositories.go`:
  ```go
  var _ PlayerRepository = (*player.Repository)(nil)
  ```

## Route Naming Conventions

- **Plural path for GET-all:** `GET /api/players`, `GET /api/decks`, `GET /api/games`
- **Singular path for GET-one and POST:** `GET /api/player?player_id=1`, `POST /api/player`
- Sub-resources use path nesting: `POST /api/pod/player`, `DELETE /api/pod/player`, `POST /api/pod/invite`, `POST /api/pod/join`
- `POST` returns `201 Created` with no body for creates
- `PATCH` returns `200 OK` or `204 No Content` depending on whether a body is returned
- Query param names use `snake_case`: `player_id`, `pod_id`, `deck_id`, `game_id`

## Error Handling

**Backend (Go):**
- Errors are wrapped with context using `fmt.Errorf("failed to get player %d: %w", id, err)`
- Error messages include the entity type and ID for traceability
- Business functions propagate repo errors upward; routers map them to HTTP status codes
- "Not found" is not an error — `nil, nil` signals absence; routers return `404` explicitly
- Sentinel strings (not typed errors) are used for some domain errors: `"forbidden: ..."`, `"unexpected number of rows"`
- Router error responses use `trackerHttp.WriteError(log, w, statusCode, err, logMsg, clientMsg)`

**Frontend (TypeScript):**
- Async functions `throw new Error(...)` on bad HTTP status
- Some functions attach `.status` to the error object for caller inspection

## Logging

**Framework:** Uber Zap (`go.uber.org/zap`)

**Patterns:**
- Logger passed as `*zap.Logger` to constructors that need it (primarily `game` business functions)
- In tests, replace with `zap.NewNop()` to suppress output
- Router errors logged via `trackerHttp.WriteError(log, w, ...)` which calls `log.Error(logMsg, zap.Error(err))`
- No structured fields beyond `zap.Error(err)` in most log calls

## Import Organization

Imports are grouped in this order:
1. Standard library
2. Third-party packages
3. Internal packages (`github.com/m-sharp/edh-tracker/lib/...`)

Internal cross-domain imports use aliased package names when there would be collisions:
```go
import (
    repos "github.com/m-sharp/edh-tracker/lib/repositories"
    deckRepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
    gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)
```

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
  ```go
  // CreatePlayerAndUser atomically inserts a player row and a linked user row in one transaction.
  ```
- Compile-time checks get a comment: `// Compile-time interface satisfaction checks.`

---

*Convention analysis: 2026-03-22*
