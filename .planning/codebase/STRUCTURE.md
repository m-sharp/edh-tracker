# Codebase Structure

**Analysis Date:** 2026-03-22

## Directory Layout

```
edh-tracker/
├── main.go                      # Entry point: wires config → DB → migrations → repos → business → API
├── api.go                       # ApiRouter: registers all sub-routers and middleware chains
├── go.mod / go.sum              # Go module definition
├── vendor/                      # Vendored Go dependencies (committed)
├── lib/                         # All backend Go source code
│   ├── config.go                # Config struct + env var key constants
│   ├── db.go                    # DBClient (wraps *gorm.DB), DBError type
│   ├── log.go                   # Zap logger factory
│   ├── web.go                   # Server struct wrapping Gorilla Mux (port 8081)
│   ├── business/                # Domain logic layer
│   │   ├── business.go          # Business struct + NewBusiness() wiring
│   │   ├── <domain>/            # One subdirectory per domain
│   │   │   ├── entity.go        # Entity type + ToEntity() converter + Validate()
│   │   │   ├── functions.go     # Free function constructors returning closures
│   │   │   ├── types.go         # Func type aliases + Functions struct
│   │   │   └── *_test.go        # Business logic tests
│   │   ├── stats/               # Shared Stats type + FromAggregate() helper
│   │   └── testHelpers/         # Shared test fixtures for business layer tests
│   ├── repositories/            # DB access layer
│   │   ├── repositories.go      # Repositories struct + New() + compile-time interface checks
│   │   ├── interfaces.go        # One interface per repository (used for DI + testing)
│   │   ├── <domain>/            # One subdirectory per domain
│   │   │   ├── model.go         # DB scan struct (embeds base.GormModelBase)
│   │   │   └── repo.go          # Repository struct + GORM query methods
│   │   ├── base/
│   │   │   └── base.go          # GormModelBase (ID, CreatedAt, UpdatedAt, DeletedAt)
│   │   └── testHelpers/         # Shared test helpers (mock DB, sqlmock setup)
│   ├── routers/                 # HTTP handler layer
│   │   ├── auth.go              # OAuth login/callback/logout/me (Google OAuth2 + JWT)
│   │   ├── player.go            # GET /api/players, GET /api/player, PATCH /api/player
│   │   ├── deck.go              # Deck CRUD endpoints
│   │   ├── game.go              # Game + game result endpoints
│   │   ├── pod.go               # Pod CRUD + membership + invite endpoints
│   │   ├── format.go            # GET /api/formats
│   │   ├── commander.go         # Commander CRUD endpoints
│   │   └── *_test.go            # Router handler tests (httptest + mock closures)
│   ├── migrations/              # Schema migration runner
│   │   ├── migrate.go           # RunAll(), Migration interface, getAllMigrations() map
│   │   ├── 1.go … 19.go         # Numbered migration structs (Upgrade + Downgrade)
│   ├── seeder/                  # Optional seed data (triggered by SEED env var)
│   │   └── seeder.go            # Seeder struct + Run(); uses repositories.Repositories directly
│   ├── trackerHttp/             # HTTP infrastructure helpers
│   │   ├── http.go              # Route struct, ApiRouter interface, WriteError/WriteJson, CallerPlayerID, middleware types
│   │   ├── auth.go              # RequireAuth / OptionalAuth middleware + IssueJWT
│   │   └── cookie.go            # SetCookie / ClearCookie helpers + cookie name constants
│   └── utils/
│       └── context.go           # ContextWithUserInfo / UserFromContext helpers
├── app/                         # React frontend (separate Docker image)
│   ├── Dockerfile               # Frontend Docker build
│   ├── package.json             # npm dependencies
│   ├── tsconfig.json            # TypeScript config
│   ├── public/                  # Static assets (index.html)
│   └── src/
│       ├── index.tsx            # Entry point: router config, AuthProvider, root render
│       ├── auth.tsx             # AuthProvider + useAuth() hook + AuthUser interface
│       ├── http.ts              # All API fetch functions (one per endpoint)
│       ├── types.ts             # TypeScript interfaces matching API Entity JSON shapes
│       ├── stats.tsx            # Reusable stat column definitions for MUI DataGrid
│       ├── common.ts            # Shared utility functions
│       ├── styles.css           # Global CSS
│       └── routes/              # One file per page view
│           ├── root.tsx         # Layout shell: AppBar + pod selector + <Outlet>
│           ├── RequireAuth.tsx  # Auth guard wrapper component
│           ├── login.tsx        # Login page
│           ├── join.tsx         # Pod join-by-invite page
│           ├── pod.tsx          # Pod detail view (players/decks/games tabs)
│           ├── player.tsx       # Player detail view
│           ├── deck.tsx         # Deck detail view
│           ├── game.tsx         # Game detail view
│           ├── new.tsx          # New game creation form
│           └── error.tsx        # Error boundary page
├── docs/                        # Architecture and API reference docs
│   ├── architecture.md          # Layer overview and wiring diagram
│   ├── api.md                   # Full API route table
│   └── data-model.md            # Entity relationships, table schemas, points formula
├── mysql/                       # MySQL Docker configuration
├── scripts/                     # Utility scripts
├── data/                        # Seed data files
└── .claude/                     # Claude Code skills and plans
    └── skills/                  # Runnable skill definitions
```

## Domain Directories

The following domains exist as parallel subdirectories under both `lib/repositories/` and `lib/business/`:

| Domain | Repositories | Business |
|---|---|---|
| `player` | `lib/repositories/player/` | `lib/business/player/` |
| `deck` | `lib/repositories/deck/` | `lib/business/deck/` |
| `game` | `lib/repositories/game/` | `lib/business/game/` |
| `gameResult` | `lib/repositories/gameResult/` | `lib/business/gameResult/` |
| `pod` | `lib/repositories/pod/` | `lib/business/pod/` |
| `user` | `lib/repositories/user/` | `lib/business/user/` |
| `format` | `lib/repositories/format/` | `lib/business/format/` |
| `commander` | `lib/repositories/commander/` | `lib/business/commander/` |
| `deckCommander` | `lib/repositories/deckCommander/` | (no business layer — used by deck) |
| `playerPodRole` | `lib/repositories/playerPodRole/` | (no business layer — used by pod/player) |
| `podInvite` | `lib/repositories/podInvite/` | (no business layer — used by pod) |

## Key Files and Their Roles

**Entry Points:**
- `main.go` — startup sequence: config validation → logger → DB client → migrations → repos → optional seeder → business → API router → HTTP server
- `api.go` — `NewApiRouter()` collects all sub-routers; `SetupRoutes()` registers routes with middleware (CORS, RequireAuth, CORS preflight)
- `app/src/index.tsx` — React entry point; defines router config and wraps app in `<AuthProvider>`

**Configuration:**
- `lib/config.go` — `lib.Config` struct; maps env var names to internal key constants; required vs optional configs enforced at startup
- `lib/db.go` — `lib.DBClient` wrapping `*gorm.DB`; database name constant `DBName = "pod_tracker"`
- `lib/web.go` — HTTP server on port 8081

**Wiring:**
- `lib/repositories/repositories.go` — `Repositories` struct bundling all concrete `*Repository` types; compile-time interface satisfaction checks
- `lib/business/business.go` — `Business` struct + `NewBusiness()` — constructs all domain `Functions` structs with repo dependencies injected

**HTTP Infrastructure:**
- `lib/trackerHttp/http.go` — `Route` struct (Path, Method, Handler, MiddleWare, RequireAuth, NoAuth), `ApiRouter` interface, `WriteError`, `WriteJson`, `CallerPlayerID`, `GetQueryId`
- `lib/trackerHttp/auth.go` — JWT middleware (`RequireAuth`, `OptionalAuth`), `IssueJWT`
- `lib/utils/context.go` — `ContextWithUserInfo` / `UserFromContext` for passing auth identity through request context

**Frontend:**
- `app/src/http.ts` — single source of truth for all API calls; typed return values matching `app/src/types.ts`
- `app/src/auth.tsx` — auth state management; `AuthProvider` fetches `/api/auth/me` on mount
- `app/src/routes/root.tsx` — layout component with `AppBar`, pod selector dropdown, and `<Outlet>`

## Migration Numbering Scheme

- Migrations live in `lib/migrations/` as `1.go` through `19.go` (currently 19 migrations)
- Each file defines a `MigrationN` struct implementing the `Migration` interface (`Upgrade`, `Downgrade`)
- `migrate.go:getAllMigrations()` returns a `map[int]Migration` keyed by migration number
- `RunAll()` runs only migrations not yet applied (tracked via a `migration` table in MySQL)
- Migrations run in ascending numeric order at server startup
- To add a new migration: create `lib/migrations/N.go` with `MigrationN` struct; add entry to `getAllMigrations()`
- Note: migration 19 exists as a file even though it was registered after skipping (no gap in numbering — the number-to-count mapping must stay sequential)

## Naming Conventions

**Backend files:**
- Domain directories: camelCase (`gameResult`, `deckCommander`, `playerPodRole`)
- Go files: lowercase (`functions.go`, `entity.go`, `types.go`, `model.go`, `repo.go`)
- Migration files: numeric (`1.go`, `2.go`, … `19.go`)

**Frontend files:**
- Route components: camelCase `.tsx` (`pod.tsx`, `player.tsx`, `new.tsx`)
- Shared modules: camelCase `.ts` or `.tsx` (`http.ts`, `types.ts`, `auth.tsx`)

## Where to Add New Code

**New domain (backend):**
1. Repository: `lib/repositories/<domain>/model.go` + `repo.go`
2. Interface: add to `lib/repositories/interfaces.go`
3. Register concrete type: `lib/repositories/repositories.go` (`Repositories` struct + `New()`)
4. Business: `lib/business/<domain>/types.go` + `functions.go` + `entity.go`
5. Wire: `lib/business/business.go` (`Business` struct + `NewBusiness()`)
6. Router: `lib/routers/<domain>.go` implementing `trackerHttp.ApiRouter`
7. Register router: `api.go:NewApiRouter()`

**New endpoint on existing domain:**
1. Add function type to `lib/business/<domain>/types.go`
2. Add constructor in `lib/business/<domain>/functions.go`
3. Add field to `Functions` struct in `lib/business/<domain>/types.go`
4. Wire in `lib/business/business.go`
5. Add handler method + route in `lib/routers/<domain>.go`

**New migration:**
- Create `lib/migrations/N.go` (next number in sequence)
- Add `N: &MigrationN{}` to `getAllMigrations()` in `lib/migrations/migrate.go`

**New frontend page:**
- Add component file to `app/src/routes/<name>.tsx`
- Register route in `app/src/index.tsx` router config
- Add API calls to `app/src/http.ts`
- Add TypeScript types to `app/src/types.ts` if needed

**New API function (frontend):**
- Add to `app/src/http.ts` with typed parameters and return value

## Special Directories

**`vendor/`:**
- Purpose: Vendored Go dependencies
- Generated: Yes (via `go mod vendor`)
- Committed: Yes

**`app/build/`:**
- Purpose: Production React build output (served by frontend Docker image)
- Generated: Yes (via `npm run build`)
- Committed: No (should not be committed)

**`lib/business/testHelpers/`:**
- Purpose: Shared fixtures and helpers for business layer tests
- Generated: No
- Committed: Yes

**`lib/repositories/testHelpers/`:**
- Purpose: Shared sqlmock helpers for repository tests (e.g., `newMockDB`)
- Generated: No
- Committed: Yes

**`.planning/`:**
- Purpose: GSD planning documents (codebase maps, phase plans)
- Generated: By Claude Code GSD commands
- Committed: Yes

---

*Structure analysis: 2026-03-22*
