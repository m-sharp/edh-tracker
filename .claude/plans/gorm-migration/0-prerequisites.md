# Phase 0 — Prerequisites

## Status
Done

## Skill
Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Goal

Add GORM as a dependency, extend `DBClient` with a `*gorm.DB` field, define `GormModelBase`, and establish integration test infrastructure. No repository logic changes yet.

## Key Context

- **Migrations** (`lib/migrations/*.go`) access `client.Db` directly — all 18 migration files use
  `client.Db.QueryRowContext`, `ExecContext`, etc. They are unaffected by Phase 0 since `Db *sqlx.DB`
  is retained until Phase 4.
- **All 11 repositories** currently use `client.Db` — unaffected until their individual phase.
- The **vendor directory** exists — `go mod vendor` is required after adding deps.
- **`DBError` / `NewDBError`** remain in `lib/db.go` through all phases; removed only in Phase 4 if unused.

## Files to Modify / Create

| File | Change |
|---|---|
| `go.mod` / `go.sum` | Add `gorm.io/gorm`, `gorm.io/driver/mysql` |
| `lib/db.go` | Add `GormDb *gorm.DB` to `DBClient`; open GORM connection in `NewDBClient` |
| `lib/repositories/base/base.go` | Add `GormModelBase` alongside existing `ModelBase` |

## Step 1 — Add Dependencies

```bash
go get gorm.io/gorm gorm.io/driver/mysql
go mod vendor
```

## Step 2 — Update `lib/db.go`

Add `GormDb *gorm.DB` to `DBClient` and open the GORM connection in `NewDBClient`, sharing the existing `*sql.DB` (no second pool).

```go
import (
    gormmysql "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type DBClient struct {
    log    *zap.Logger
    Db     *sqlx.DB   // keep until full migration
    GormDb *gorm.DB
}
```

In `NewDBClient`, after the pool settings are applied to `db` and before `inst` construction, wrap the same `*sql.DB` with GORM:

```go
// GORM wraps the same *sql.DB — no second pool
gormDB, err := gorm.Open(gormmysql.New(gormmysql.Config{Conn: db}), &gorm.Config{})
if err != nil {
    return nil, fmt.Errorf("error opening gorm connection: %w", err)
}

inst := &DBClient{log: log, Db: sqlx.NewDb(db, "mysql"), GormDb: gormDB}
```

**Note:** `gormmysql.Config{Conn: db}` (not `gormmysql.Open(dsn)`) is required to share the existing `*sql.DB`. The pool settings are already applied to `db` — do not re-apply them to GORM's `sqlDB`.

Import: same alias as the DSN form — `gormmysql "gorm.io/driver/mysql"` — no change needed.

Note: `CheckConnection` continues to use `Db.Ping()` — no change needed there.

## Step 3 — Update `lib/repositories/base/base.go`

Add `GormModelBase` while keeping the existing `ModelBase` for any repos not yet migrated:

```go
package base

import (
    "time"

    "gorm.io/gorm"
)

// ModelBase is used by sqlx-based repositories (pre-GORM).
type ModelBase struct {
    ID        int        `db:"id"`
    CreatedAt time.Time  `db:"created_at"`
    UpdatedAt time.Time  `db:"updated_at"`
    DeletedAt *time.Time `db:"deleted_at"`
}

// GormModelBase is used by GORM-based repositories.
type GormModelBase struct {
    ID        int            `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## Step 4 — Integration Test Infrastructure

The shared test helper package lives at `lib/repositories/testHelpers/`:

- `testHelpers.NewTestDB(t)` — opens a connection to `host.docker.internal:3306/pod_tracker`, begins a transaction, and registers a `t.Cleanup` rollback — no explicit table cleanup needed.
- `testHelpers.New<Domain>Repo(db)` — thin wrappers over each repo's `NewRepositoryFromDB`.
- `testHelpers.CreateTest*` helpers — insert FK prerequisites and return the generated ID.

Each repo test file uses an **external test package** (`package <domain>_test`) and sets up per-test with two lines — no shared `newRepo` helper:

```go
db := testHelpers.NewTestDB(t)
repo := testHelpers.NewPlayerRepo(db)
```

`base/testHelpers.go` has been deleted; `base.NewTestDB` no longer exists. All new and migrated tests import `testHelpers` directly.

## Verification

After completing all steps:
1. `go vet ./lib/...` — must pass with no errors
2. `go run main.go` (with DB env vars set) — API starts cleanly; GORM connection log appears
3. No changes to any repository or business logic — all existing sqlmock tests still pass
