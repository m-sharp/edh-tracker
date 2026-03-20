---
name: gorm
description: Use this skill to load GORM patterns, conventions, and code examples for the EDH Tracker. Invoke at the start of any GORM repository implementation or migration session to get import paths, model definitions, query patterns, bulk insert, upsert, transactions, soft-delete, and test infrastructure.
version: 1.0.0
---

# GORM Patterns — EDH Tracker Reference

## Import Paths

```go
"gorm.io/gorm"
"gorm.io/driver/mysql"
"gorm.io/gorm/clause"
```

## Model Definition

All GORM models in this project embed `GormModelBase` from `lib/repositories/base/`:

```go
// lib/repositories/base/base.go
type GormModelBase struct {
    ID        int            `gorm:"primaryKey;column:id"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

Domain model example:
```go
type Model struct {
    base.GormModelBase
    Name string `gorm:"column:name"`
}
```

**Soft delete is automatic** when `gorm.DeletedAt` is embedded — GORM appends `WHERE deleted_at IS NULL` to all queries and sets `deleted_at` on `db.Delete()`.

**Table name**: GORM pluralises by default (`Model` → `models`). Override with:
```go
func (Model) TableName() string { return "player" }
```

## Opening a GORM DB (main.go)

```go
import (
    gormmysql "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

gormDB, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
```

The DSN is the same MySQL DSN used for sqlx (`config.FormatDSN()`). Connection pool settings must be applied to the underlying `*sql.DB`:
```go
sqlDB, _ := gormDB.DB()
sqlDB.SetConnMaxLifetime(maxConnTTL)
sqlDB.SetMaxOpenConns(maxConnCount)
sqlDB.SetMaxIdleConns(maxConnCount)
```

## DBClient

```go
type DBClient struct {
    log    *zap.Logger
    Db     *sqlx.DB   // remove after full migration
    GormDb *gorm.DB
}
```

Repository constructors receive `*lib.DBClient` and use `client.GormDb`.

## Context Propagation

Pass context to every GORM call via `.WithContext(ctx)`:
```go
db.WithContext(ctx).Where("id = ?", id).First(&m)
```

## Common Query Patterns

### Find all (with soft-delete filter automatic)
```go
var results []Model
if err := r.db.WithContext(ctx).Find(&results).Error; err != nil { ... }
if results == nil { results = []Model{} }
```

### Find one by ID
```go
var m Model
err := r.db.WithContext(ctx).First(&m, id).Error
if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
```

### Find one by field — return nil,nil if not found
```go
var m Model
err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
```

### IN clause (replaces sqlx.In + Rebind)
```go
var results []Model
r.db.WithContext(ctx).Where("name IN ?", names).Find(&results)
```

### Single insert — get LastInsertId
```go
m := Model{Name: name}
if err := r.db.WithContext(ctx).Create(&m).Error; err != nil { ... }
return m.ID, nil  // GORM sets ID after Create
```

### Soft delete (sets deleted_at)
```go
r.db.WithContext(ctx).Delete(&Model{}, id)
```

### Update single field
```go
r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Update("name", name)
```

### Update multiple fields (map avoids zero-value skipping)
```go
r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Updates(map[string]any{
    "name":      newName,
    "format_id": newFormatID,
})
```

### Dynamic partial update (UpdateFields pattern)
```go
updates := map[string]any{}
if f.Name != nil     { updates["name"] = *f.Name }
if f.FormatID != nil { updates["format_id"] = *f.FormatID }
if len(updates) == 0 { return nil }
return r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Updates(updates).Error
```

### Retire / toggle bool field
```go
r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Update("retired", true)
```

## Bulk Insert (CreateInBatches)

```go
records := []Model{{Name: "a"}, {Name: "b"}}
if err := r.db.WithContext(ctx).CreateInBatches(&records, 100).Error; err != nil { ... }
// records[i].ID is populated after insert
```

For domains that need IDs back after BulkAdd:
```go
r.db.WithContext(ctx).CreateInBatches(&games, 100)
ids := make([]int, len(games))
for i, g := range games { ids[i] = g.ID }
```

## Pluck (scalar list queries — replaces SelectContext into []int)

```go
var ids []int
r.db.WithContext(ctx).Model(&PlayerPodModel{}).
    Where("player_id = ? AND deleted_at IS NULL", playerID).
    Pluck("pod_id", &ids)
```

Note: Pluck does NOT apply soft-delete automatically unless `DeletedAt` is embedded in the model being queried. Add explicit `AND deleted_at IS NULL` when using Pluck on junction tables whose model may not embed `gorm.DeletedAt`.

## Upsert — ON DUPLICATE KEY UPDATE

```go
import "gorm.io/gorm/clause"

r.db.WithContext(ctx).Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "pod_id"}, {Name: "player_id"}},
    DoUpdates: clause.Assignments(map[string]any{
        "role":       role,
        "deleted_at": nil,
    }),
}).Create(&record)
```

## JOIN Queries

Simple JOIN — scan into a flat struct:
```go
var results []Model
r.db.WithContext(ctx).
    Joins("INNER JOIN game_result ON game.id = game_result.game_id").
    Where("game_result.deck_id = ? AND game_result.deleted_at IS NULL", deckID).
    Find(&results)
```

Multi-table JOIN with DISTINCT:
```go
r.db.WithContext(ctx).
    Distinct().
    Joins("INNER JOIN game_result ON game.id = game_result.game_id").
    Joins("INNER JOIN deck ON game_result.deck_id = deck.id").
    Where("deck.player_id = ? AND game_result.deleted_at IS NULL AND deck.deleted_at IS NULL", playerID).
    Find(&results)
```

## Raw SQL (for correlated subqueries)

When GORM can't express a query cleanly (e.g. correlated subquery for player_count):
```go
type gameStat struct {
    GameID      int `gorm:"column:game_id"`
    Place       int `gorm:"column:place"`
    KillCount   int `gorm:"column:kill_count"`
    PlayerCount int `gorm:"column:player_count"`
}

var stats []gameStat
r.db.WithContext(ctx).Raw(getStatsForPlayer, playerID).Scan(&stats)
```

Keep SQL constants for raw queries; they remain in `repo.go`.

## Transactions

```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    player := playerModel{Name: playerName}
    if err := tx.Create(&player).Error; err != nil {
        return err  // triggers rollback
    }
    u := userModel{PlayerID: player.ID, ...}
    if err := tx.Create(&u).Error; err != nil {
        return err
    }
    return nil  // commits
})
```

## Soft-Delete on FK-Related Rows (DeleteByDeckID pattern)

When deleting all rows matching a FK (not by primary key):
```go
// This sets deleted_at on all matching rows
r.db.WithContext(ctx).Where("deck_id = ?", deckID).Delete(&Model{})
```

## Testing with go-sqlmock + GORM

For unit tests that still need sqlmock:
```go
import (
    "github.com/DATA-DOG/go-sqlmock"
    gormmysql "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

db, mock, _ := sqlmock.New()
gormDB, _ := gorm.Open(gormmysql.New(gormmysql.Config{Conn: db}), &gorm.Config{})
```

## Integration Test Infrastructure

The preferred test approach after GORM migration is integration tests against a real MySQL DB.

Required env vars:
- `TEST_DBHOST`
- `TEST_DBUSER`
- `TEST_DBPASSWORD`
- `TEST_DBPORT`
- `TEST_DBNAME` (use a dedicated test DB, e.g. `pod_tracker_test`)

Helper pattern (place in `testhelpers_test.go` in each package):
```go
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
    t.Cleanup(func() { /* truncate test tables */ })
    return db
}
```

Tests skip gracefully when env vars are absent, so CI without a DB doesn't fail.

## RowsAffected Check

GORM doesn't panic on 0 rows affected — check explicitly when required:
```go
result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Update("name", name)
if result.Error != nil { return result.Error }
if result.RowsAffected != 1 {
    return fmt.Errorf("unexpected rows affected: got %d, expected 1", result.RowsAffected)
}
```

For soft-deletes and updates where "not found" should silently succeed (e.g. RemovePlayer), skip the RowsAffected check.

## Error Wrapping Convention

Keep the same `fmt.Errorf("failed to ...: %w", err)` pattern as existing repos for consistency.
