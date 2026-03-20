# Phase 1c — Commander Repository

## Status
Done

## Skill
Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Scope

- Table: `commander`
- Files: `lib/repositories/commander/model.go`, `lib/repositories/commander/repo.go`
- No existing tests — write new integration tests
- No business layer changes

## GORM Model

```go
package commander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    Name string
}

func (Model) TableName() string { return "commander" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetById` | `db.First(&m, id)` | `ErrRecordNotFound` → nil,nil |
| `GetByName` | `db.Where("name = ?", name).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `GetByNames` | `db.Where("name IN ?", names).Find(&results)` | Replaces sqlx.In + Rebind |
| `Add` | `db.Create(&m)` | `m.ID` populated by GORM |
| `BulkAdd` | `db.CreateInBatches(&models, 100)` then `GetByNames` | See below |

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

Same pattern as player BulkAdd — create models, use CreateInBatches, then select back via GetByNames:

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
        return nil, fmt.Errorf("failed to bulk insert Commander records: %w", err)
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
    var commanders []Model
    if err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&commanders).Error; err != nil {
        return nil, fmt.Errorf("failed to get Commander records by names: %w", err)
    }
    return commanders, nil
}
```

## Test Migration

No existing tests. Write new integration tests:
- `TestGetById_Found` / `TestGetById_NotFound`
- `TestGetByName_Found` / `TestGetByName_NotFound`
- `TestGetByNames` — multiple names, partial match
- `TestAdd` — returns correct ID
- `TestBulkAdd` — inserts and returns models

Use `base.NewTestDB(t)` from `lib/repositories/base/testHelpers.go`. Define a `newRepo(t)` helper in `repo_test.go` (see Phase 1a pattern). No `testhelpers_test.go` needed.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/commander/...` passes (or skips)
3. Smoke test: `POST /api/deck` (which creates commanders) returns 201
