# Testing Patterns

**Analysis Date:** 2026-03-22

## Test Framework

**Runner:** `react-scripts test` (frontend), `go test` (backend)

**Backend assertion libraries:**
- `github.com/stretchr/testify/assert` — non-fatal assertions
- `github.com/stretchr/testify/require` — fatal assertions (stops test on failure)

**Run Commands:**
```bash
go test ./lib/...          # Run all backend tests
go vet ./lib/...           # Compile check without binary output
npm test                   # Run frontend tests (from app/)
```

Note: `go build ./...` and `go build ./lib/...` must NOT be used — the former crawls `app/node_modules/` and the latter leaves a binary in the project root.

## Test File Organization

**Co-location:** All test files live in the same directory as the code under test.

**File naming:** `<file>_test.go` co-located with the source file. Examples:
- `lib/business/player/functions_test.go` tests `lib/business/player/functions.go`
- `lib/routers/player_test.go` tests `lib/routers/player.go`
- `lib/repositories/player/repo_test.go` tests `lib/repositories/player/repo.go`

**Package declarations:**
- Business layer tests (`lib/business/<domain>/`) use the **same package** as the code: `package player`, `package game`
- Repository integration tests use `package <domain>_test` (external test package): `package player_test`
- Router tests use the **same package** as routers: `package routers`
- This means business tests have access to unexported identifiers; repo tests do not

**Export files:** When an unexported function must be accessible to the `_test` package (e.g., resetting package-level cache state), an `export_test.go` file in the same package exposes it:
```go
// lib/business/format/export_test.go
package format

func resetCache() {
    cache.Lock()
    defer cache.Unlock()
    cache.m = nil
}
```

## Three Distinct Test Types

### 1. Router Tests (`lib/routers/*_test.go`)

Test HTTP handler behaviour. No DB or real business logic — inject mock `Functions` structs directly.

**Pattern:**
```go
func newTestPlayerRouter(players player.Functions) *PlayerRouter {
    return &PlayerRouter{
        log:     zap.NewNop(),
        players: players,
    }
}

func TestPlayerRouter_GetAll_Success(t *testing.T) {
    router := newTestPlayerRouter(player.Functions{
        GetAll: func(ctx context.Context) ([]player.Entity, error) {
            return []player.Entity{{Name: "Alice"}}, nil
        },
    })

    req := httptest.NewRequest(http.MethodGet, "/api/players", nil)
    rr := httptest.NewRecorder()
    router.GetAll(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    var got []player.Entity
    require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
    assert.Len(t, got, 1)
}
```

Key points:
- Always use `httptest.NewRecorder()` for the response writer
- Always use `httptest.NewRequest()` for the request
- Use `zap.NewNop()` for the logger
- Inject only the `Functions` fields needed for the test — others stay nil (zero value)
- Calling an un-injected function panics: `panic("unexpected call to GetAll")` — this is intentional

**Simulating auth:** Use the `withAuth` helper defined in `lib/routers/pod_test.go`:
```go
func withAuth(r *http.Request, playerID int) *http.Request {
    return r.WithContext(utils.ContextWithUserInfo(r.Context(), 1, playerID))
}

req := withAuth(httptest.NewRequest(http.MethodPost, "/api/pod", body), 10)
```

**Test naming:** `Test<RouterType>_<HandlerName>_<Scenario>` e.g., `TestPodRouter_PodCreate_Unauthenticated`

### 2. Business Layer Tests (`lib/business/<domain>/*_test.go`)

Test domain logic in isolation. Use mock repository structs from `lib/business/testHelpers/mocks.go`.

**Pattern:**
```go
func TestGetByID_Success(t *testing.T) {
    playerRepo := &testHelpers.MockPlayerRepo{
        GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
            return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: 5}, Name: "Alice"}, nil
        },
    }
    gameResultRepo := &testHelpers.MockGameResultRepo{
        GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
            return &gameresultrepo.Aggregate{Games: 3, Record: map[int]int{1: 1}}, nil
        },
    }

    fn := GetByID(playerRepo, gameResultRepo, nil) // nil for unused deps
    got, err := fn(context.Background(), 5)
    require.NoError(t, err)
    assert.Equal(t, "Alice", got.Name)
}
```

Key points:
- Call the free constructor to get the closure, then invoke the closure directly
- Pass `nil` for repo dependencies not exercised in the test
- Mock structs are defined in `lib/business/testHelpers/mocks.go` — one struct per interface
- Each mock struct field is named `<MethodName>Fn`; if nil when called, panics with `"unexpected call to <Method>"`

**Mock struct example:**
```go
type MockPlayerRepo struct {
    GetAllFn    func(ctx context.Context) ([]playerRepo.Model, error)
    GetByIdFn   func(ctx context.Context, playerID int) (*playerRepo.Model, error)
    GetByNameFn func(ctx context.Context, name string) (*playerRepo.Model, error)
    UpdateFn    func(ctx context.Context, playerID int, name string) error
}
```

**Compile-time checks in mocks file:**
```go
var (
    _ repos.PlayerRepository = (*MockPlayerRepo)(nil)
    _ repos.PodRepository    = (*MockPodRepo)(nil)
    // ...
)
```

### 3. Repository Integration Tests (`lib/repositories/<domain>/repo_test.go`)

Test against a real MySQL database. Require the Docker DB to be running.

**DB helpers in `lib/repositories/testHelpers/`:**
- `testDB.go` — `NewTestDB(t)` wraps the connection in a transaction rolled back on `t.Cleanup`
- `testDB.go` — `NewTestDBNoTx(t)` for operations that open their own transaction
- `helpers.go` — `NewPlayerRepo(db)`, `NewDeckRepo(db)`, etc. — construct repos from `*gorm.DB`
- `helpers.go` — fixture creators: `CreateTestPlayer`, `CreateTestDeck`, `CreateTestGame`, `CreateTestPod`, etc.

**Pattern:**
```go
func TestGetById_Found(t *testing.T) {
    db := testHelpers.NewTestDB(t)
    repo := testHelpers.NewPlayerRepo(db)
    ctx := context.Background()

    id, err := repo.Add(ctx, "Alice")
    require.NoError(t, err)

    got, err := repo.GetById(ctx, id)
    require.NoError(t, err)
    require.NotNil(t, got)
    assert.Equal(t, id, got.ID)
    assert.Equal(t, "Alice", got.Name)
}
```

Key points:
- `NewTestDB` wraps everything in a transaction; `t.Cleanup` rolls back automatically — no manual teardown needed
- Use `NewTestDBNoTx` only when the code under test opens its own transaction (e.g., `CreatePlayerAndUser`)
- When using `NewTestDBNoTx`, you MUST specify manual cleanup actions
- IDs for fixtures generated with `atomic.AddInt64(&fixtureCounter, 1)` to avoid collisions in parallel tests
- `require.NoError` for operations that must succeed for the test to have meaning
- `assert.Nil` / `assert.NotNil` for pointer return value assertions

**Not-found pattern in repo tests:**
```go
func TestGetByName_NotFound(t *testing.T) {
    db := testHelpers.NewTestDB(t)
    repo := testHelpers.NewPlayerRepo(db)
    got, err := repo.GetByName(context.Background(), "NoSuchPlayer")
    require.NoError(t, err) // no error
    assert.Nil(t, got)      // nil result
}
```

**Error string assertions for expected errors:**
```go
err := repo.Update(context.Background(), 999999, "Ghost")
assert.ErrorContains(t, err, "unexpected number of rows")
```

## Table-Driven Tests

Required for any function with more than 2 code paths. Used consistently in entity validation tests and utility tests.

**Standard pattern:**
```go
func TestInputEntityValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   InputEntity
        wantErr bool
    }{
        {name: "valid", input: InputEntity{DeckID: 1, Place: 1, Kills: 0}, wantErr: false},
        {name: "zero deck id", input: InputEntity{DeckID: 0, Place: 1, Kills: 0}, wantErr: true},
        {name: "place below 1", input: InputEntity{DeckID: 1, Place: 0, Kills: 0}, wantErr: true},
        {name: "negative kills", input: InputEntity{DeckID: 1, Place: 1, Kills: -1}, wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

Flat table-driven tests (without `t.Run`) are also used for pure functions:
```go
func TestGetPointsForPlace(t *testing.T) {
    tests := []struct {
        kills, place, numPlayers, want int
    }{
        {kills: 2, place: 1, numPlayers: 4, want: 5},
        // ...
    }
    for _, tt := range tests {
        got := GetPointsForPlace(tt.kills, tt.place, tt.numPlayers)
        assert.Equal(t, tt.want, got)
    }
}
```

## Mocking Strategy

**Two separate mock locations:**

1. `lib/business/testHelpers/mocks.go` — mock implementations of all repository interfaces. Used by business layer tests.

2. `lib/repositories/testHelpers/` — real repo constructors pointing at the test DB. Used by repository integration tests.

**No third-party mocking library** is used. Mocks are hand-written structs with `Fn` fields. Unset `Fn` fields panic on call — intentional design to catch unexpected invocations.

**What to mock:** Repository interfaces (`PlayerRepository`, `DeckRepository`, etc.)

**What NOT to mock:** Business `Functions` structs in business layer tests — call the real constructor with mock repos instead.

## Assertion Usage

- Use `require.NoError(t, err)` for setup operations that must succeed for the test to proceed
- Use `assert.NoError(t, err)` for the actual assertion under test when you still want to continue
- Use `require.NotNil(t, ptr)` before dereferencing a pointer
- Use `assert.ErrorContains(t, err, "substring")` to check partial error messages
- Use `assert.Equal(t, expected, actual)` — expected value always first
- Use `assert.Len(t, slice, n)` to check slice lengths
- Use `assert.GreaterOrEqual(t, len(slice), n)` when testing a lower bound (integration tests may have pre-existing rows)

## Test Naming

**Backend:**
- `Test<Type>_<Method>_<Scenario>` for router and business tests
- `Test<Method>_<Scenario>` for repository tests (no type prefix since package is domain-specific)
- Examples: `TestPlayerRouter_GetById_MissingParam`, `TestGetByID_NotFound`, `TestSoftDelete_CascadesToAllPlayerRows`

## Frontend Testing

No test files exist in `app/src/`. The frontend has no automated tests.

**Type checking** is used as a substitute for unit tests:
```bash
./node_modules/.bin/tsc --noEmit   # from app/
```

This is the only frontend verification command. Do NOT use `npm run build` or `npx tsc` for type checking.

## Coverage

No coverage targets are enforced. No coverage configuration exists in any config file.

Run coverage manually:
```bash
go test ./lib/... -cover
```

## Key Test Infrastructure Files

- `lib/business/testHelpers/mocks.go` — all mock repository implementations
- `lib/repositories/testHelpers/testDB.go` — `NewTestDB(t)` and `NewTestDBNoTx(t)`
- `lib/repositories/testHelpers/helpers.go` — `New*Repo(db)` constructors and `Create*` fixture helpers

---

*Testing analysis: 2026-03-22*
