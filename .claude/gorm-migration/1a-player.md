# Phase 1a — Player Repository

## Status
Pending

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `player`
- Files: `lib/repositories/player/model.go`, `lib/repositories/player/repo.go`, `lib/repositories/player/repo_test.go`
- No business layer changes

## GORM Model

```go
package player

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    Name string `gorm:"column:name"`
}

func (Model) TableName() string { return "player" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetAll` | `db.Find(&players)` | Soft-delete filter automatic |
| `GetById` | `db.First(&m, id)` | `ErrRecordNotFound` → nil,nil |
| `GetByName` | `db.Where("name = ?", name).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `GetByNames` | `db.Where("name IN ?", names).Find(&players)` | Replaces sqlx.In + Rebind |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM after create |
| `BulkAdd` | `db.CreateInBatches(&models, 100)` then return via `GetByNames` | See below |
| `Update` | `db.Model(&Model{}).Where("id = ?", id).Update("name", name)` | Check RowsAffected |
| `SoftDelete` | `db.Delete(&Model{}, id)` | Sets `deleted_at` automatically |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — BulkAdd

`BulkAdd` inserts then returns the created models via `GetByNames`. With GORM, use `CreateInBatches` which populates IDs, then call `GetByNames`:

```go
func (r *Repository) BulkAdd(ctx context.Context, names []string) ([]Model, error) {
    if len(names) == 0 {
        return []Model{}, nil
    }
    models := make([]Model, len(names))
    for i, n := range names {
        models[i] = Model{Name: n}
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&models, 100).Error; err != nil {
        return nil, fmt.Errorf("failed to bulk insert Player records: %w", err)
    }
    return r.GetByNames(ctx, names)
}
```

## Special Pattern — GetByNames (IN clause)

```go
func (r *Repository) GetByNames(ctx context.Context, names []string) ([]Model, error) {
    if len(names) == 0 {
        return []Model{}, nil
    }
    var players []Model
    if err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&players).Error; err != nil {
        return nil, fmt.Errorf("failed to get Player records by names: %w", err)
    }
    return players, nil
}
```

## Test Migration

Remove all sqlmock tests in `repo_test.go`. Replace with integration tests using `newTestDB(t)`.

Tests to write:
- `TestGetAll` — insert 2 players, GetAll returns both; insert soft-deleted player, not returned
- `TestGetById_Found` / `TestGetById_NotFound`
- `TestGetByName_Found` / `TestGetByName_NotFound`
- `TestGetByNames`
- `TestAdd` — returns correct ID
- `TestBulkAdd` — inserts and returns models with IDs
- `TestUpdate_Found` / `TestUpdate_NotFound`
- `TestSoftDelete` — player not returned by GetAll after delete

Add `testhelpers_test.go` with `newTestDB(t)` helper (see Phase 0 template). Cleanup: truncate `player` table in `t.Cleanup`.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/player/...` passes (or skips if TEST_DBHOST unset)
3. Smoke test: `GET /api/players` returns 200
