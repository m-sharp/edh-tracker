# Phase 5 — Final Cleanup

## Status
Done

## Skill
Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Goal

Remove all sqlx dependencies after all 11 repository domains have been migrated to GORM. Project should compile cleanly with no sqlx references remaining.

## Prerequisites

All phases complete:
- Phase 1a — player
- Phase 1b — format
- Phase 1c — commander
- Phase 2a — deck
- Phase 2b — deckCommander
- Phase 2c — pod
- Phase 2d — playerPodRole
- Phase 3a — game
- Phase 3b — gameResult
- Phase 3c — podInvite
- Phase 3d — user

> Note: `client.Db` is used by all 18 migration files. Migrations must be updated to use
> `client.GormDb` (or its underlying `*sql.DB`) before `Db *sqlx.DB` can be removed.

## Step 0 — Migrate `lib/migrations/` Off sqlx

All 18 migration files (e.g. `1.go` through `18.go`) and `migrate.go` use `client.Db` directly
via `QueryRowContext`, `ExecContext`, `QueryContext`, and `BeginTxx`. The `Db *sqlx.DB` field
cannot be removed until migrations no longer reference it.

Options (choose one before starting Phase 5):
1. **Run migrations via GORM raw SQL** — replace `client.Db.ExecContext(ctx, sql)` with
   `r.db.WithContext(ctx).Exec(sql)` (or `r.db.Exec(sql)` for non-contextual calls).
   Transactions: use `r.db.Transaction(func(tx *gorm.DB) error { ... })`.
2. **Run migrations via `sqlDB`** — extract the underlying `*sql.DB` from GORM:
   `sqlDB, _ := client.GormDb.DB()` and use `sqlDB.ExecContext(ctx, sql)` etc.
   This avoids touching migration SQL while still removing the sqlx dependency.

Recommended: Option 2 (minimal diff — only the field access changes, no migration logic touched).

After this step, `client.Db` has zero remaining callers and Step 1 can proceed.

## Step 1 — Remove `Db *sqlx.DB` from DBClient

In `lib/db.go`:
- Remove `Db *sqlx.DB` field from `DBClient`
- Remove sqlx import
- Remove `sql.Open` + `sqlx.NewDb` wiring (only GORM connection remains)
- Remove `CheckConnection` use of `Db.Ping()` — replace with GORM's `sqlDB.Ping()`:

```go
func (d *DBClient) CheckConnection() error {
    d.log.Debug("Pinging DB for health check...")
    sqlDB, err := d.GormDb.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}
```

## Step 2 — Remove sqlx from repositories.go

In `lib/repositories/repositories.go`:
- The `New(log, client)` signature is unchanged (client is still `*lib.DBClient`)
- No sqlx imports remain in this file after repo migrations

## Step 3 — Remove ModelBase from base.go

Once all repos use `GormModelBase`, remove `ModelBase` from `lib/repositories/base/base.go`:

```go
package base

import (
    "time"
    "gorm.io/gorm"
)

// GormModelBase is the base struct for all GORM repository models.
type GormModelBase struct {
    ID        int            `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## Step 4 — Remove sqlx from go.mod

```bash
go mod tidy
go mod vendor
```

Verify `github.com/jmoiron/sqlx` and `github.com/go-sql-driver/mysql` are removed from `go.mod` (the MySQL driver is now pulled in transitively by `gorm.io/driver/mysql` so it may remain — that's fine).

Also verify `github.com/DATA-DOG/go-sqlmock` is removed (all tests now use integration test helpers).

## Step 5 — Grep for Remaining sqlx / db: References

```bash
grep -r "sqlx" lib/
grep -r '"db:"' lib/
grep -r "go-sqlmock" lib/
```

All results should be empty after cleanup.

## Step 6 — Remove DBError type (if unused)

Check if `lib.DBError` / `lib.NewDBError` are still used anywhere after the migration. If not, remove from `lib/db.go`.

```bash
grep -r "DBError\|NewDBError" .
```

## Step 7 — Run Full Test Suite

```bash
go test ./lib/...
```

All integration tests pass (or skip if TEST_DBHOST is unset). No sqlmock tests remain.

## Step 8 — Smoke Test

Run `/smoke-test` to rebuild the Docker image and verify all core API endpoints respond correctly:
- `GET /api/players`
- `GET /api/formats`
- `GET /api/pods`
- `POST /api/player`
- `POST /api/deck`
- Auth endpoints

## Verification Checklist

- [ ] `go vet ./lib/...` passes (no binary output)
- [ ] `go mod tidy` produces no changes
- [ ] `grep -r "sqlx" lib/` — empty
- [ ] `grep -r '"db:"' lib/` — empty
- [ ] `grep -r "go-sqlmock" lib/` — empty
- [ ] All integration tests pass or skip
- [ ] Smoke test passes for all core endpoints
- [ ] Docker image builds cleanly
