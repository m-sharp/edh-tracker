# Phase 0 — Prerequisites

## Goal

Add GORM as a dependency, extend `DBClient` with a `*gorm.DB` field, define `GormModelBase`, and establish integration test infrastructure. No repository logic changes yet.

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

Add `GormDb *gorm.DB` to `DBClient` and open the GORM connection in `NewDBClient`, reusing the same DSN.

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

In `NewDBClient`, after opening the sqlx DB and passing the connection check, open GORM:

```go
gormDB, err := gorm.Open(gormmysql.Open(config.FormatDSN()), &gorm.Config{})
if err != nil {
    return nil, fmt.Errorf("error opening gorm connection: %w", err)
}
// Reuse pool settings on the underlying *sql.DB
sqlDB, _ := gormDB.DB()
sqlDB.SetConnMaxLifetime(maxConnTTL)
sqlDB.SetMaxOpenConns(maxConnCount)
sqlDB.SetMaxIdleConns(maxConnCount)

inst := &DBClient{log: log, Db: sqlx.NewDb(db, "mysql"), GormDb: gormDB}
```

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
    ID        int            `gorm:"primaryKey;column:id"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

## Step 4 — Integration Test Infrastructure

For each repo package that is migrated, create a `testhelpers_test.go` file with a `newTestDB(t)` helper. This pattern is the same across all packages:

```go
package <domain>

import (
    "fmt"
    "os"
    "testing"

    gormmysql "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    host := os.Getenv("TEST_DBHOST")
    if host == "" {
        t.Skip("TEST_DBHOST not set; skipping integration test")
    }
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("TEST_DBUSER"),
        os.Getenv("TEST_DBPASSWORD"),
        host,
        os.Getenv("TEST_DBPORT"),
        os.Getenv("TEST_DBNAME"),
    )
    db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
    require.NoError(t, err)
    return db
}
```

Tests that require a real DB call `t.Skip` when `TEST_DBHOST` is unset, so they are safe to run in CI without a DB.

Required env vars:
- `TEST_DBHOST`
- `TEST_DBUSER`
- `TEST_DBPASSWORD`
- `TEST_DBPORT`
- `TEST_DBNAME` (recommended: `pod_tracker_test` — a dedicated test DB with the same schema)

## Verification

After completing all steps:
1. `go vet ./lib/...` — must pass with no errors
2. `go run main.go` (with DB env vars set) — API starts cleanly; GORM connection log appears
3. No changes to any repository or business logic — all existing sqlmock tests still pass
