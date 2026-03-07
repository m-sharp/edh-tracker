# Phase 1b — Format Repository

## Status
Pending

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `format`
- Files: `lib/repositories/format/model.go`, `lib/repositories/format/repo.go`
- No tests exist currently — write new integration tests
- No business layer changes

## GORM Model

```go
package format

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    Name string `gorm:"column:name"`
}

func (Model) TableName() string { return "format" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetAll` | `db.Find(&formats)` | Soft-delete automatic |
| `GetById` | `db.First(&m, id)` | `ErrRecordNotFound` → nil,nil |
| `GetByName` | `db.Where("name = ?", name).First(&m)` | `ErrRecordNotFound` → nil,nil |

No write methods — format is read-only (seeded via migrations).

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Implementation

```go
func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
    var formats []Model
    if err := r.db.WithContext(ctx).Find(&formats).Error; err != nil {
        return nil, fmt.Errorf("failed to get Format records: %w", err)
    }
    if formats == nil {
        return []Model{}, nil
    }
    return formats, nil
}

func (r *Repository) GetById(ctx context.Context, id int) (*Model, error) {
    var m Model
    err := r.db.WithContext(ctx).First(&m, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get Format record for id %d: %w", id, err)
    }
    return &m, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
    var m Model
    err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get Format record for name %q: %w", name, err)
    }
    return &m, nil
}
```

## Test Migration

No existing tests. Write new integration tests:
- `TestGetAll` — format table is seeded by migrations; verify returns records
- `TestGetById_Found` / `TestGetById_NotFound`
- `TestGetByName_Found` / `TestGetByName_NotFound`

Note: Format table is seed-only. Tests can use known seeded values (e.g. "Commander") rather than inserting. Alternatively, the test DB should have formats seeded via migrations.

Add `testhelpers_test.go` with `newTestDB(t)` helper.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/format/...` passes (or skips if TEST_DBHOST unset)
3. Smoke test: `GET /api/formats` returns 200
