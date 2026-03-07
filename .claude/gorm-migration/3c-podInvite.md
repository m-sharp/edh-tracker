# Phase 3c — PodInvite Repository

## Scope

- Table: `pod_invite`
- Files: `lib/repositories/podInvite/model.go`, `lib/repositories/podInvite/repo.go`, `lib/repositories/podInvite/repo_test.go`
- Depends on: pod (Phase 2c), player (Phase 1a)
- No business layer changes

## GORM Model

```go
package podInvite

import (
    "time"

    "github.com/m-sharp/edh-tracker/lib/repositories/base"
)

type Model struct {
    base.GormModelBase
    PodID             int        `gorm:"column:pod_id"`
    InviteCode        string     `gorm:"column:invite_code"`
    CreatedByPlayerID int        `gorm:"column:created_by_player_id"`
    ExpiresAt         *time.Time `gorm:"column:expires_at"`
    UsedCount         int        `gorm:"column:used_count"`
}

func (Model) TableName() string { return "pod_invite" }
```

`ExpiresAt` is nullable (`*time.Time`) — GORM handles this correctly.

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByCode` | `db.Where("invite_code = ?", code).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `Add` | `db.Create(&m)` | No ID needed; return error only |
| `IncrementUsedCount` | Raw `UPDATE` via `db.Exec` | See below |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — Add (no returned ID)

The `Add` interface returns only `error`. Build and create a model:

```go
func (r *Repository) Add(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
    m := Model{
        PodID:             podID,
        InviteCode:        code,
        CreatedByPlayerID: createdByPlayerID,
        ExpiresAt:         expiresAt,
    }
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return fmt.Errorf("failed to insert pod invite: %w", err)
    }
    return nil
}
```

## Special Pattern — IncrementUsedCount (raw UPDATE)

`used_count = used_count + 1` can't be expressed safely as a GORM `Update` (it would send the literal value, not an expression). Use `db.Exec` for this:

```go
func (r *Repository) IncrementUsedCount(ctx context.Context, code string) error {
    err := r.db.WithContext(ctx).
        Exec("UPDATE pod_invite SET used_count = used_count + 1 WHERE invite_code = ?", code).
        Error
    if err != nil {
        return fmt.Errorf("failed to increment used_count for invite %q: %w", code, err)
    }
    return nil
}
```

Alternatively, use GORM's expression support:
```go
r.db.WithContext(ctx).Model(&Model{}).
    Where("invite_code = ?", code).
    Update("used_count", gorm.Expr("used_count + 1"))
```

Either approach is acceptable. The `gorm.Expr` approach is more idiomatic GORM.

## Test Migration

Remove existing sqlmock tests. Replace with integration tests.

Tests to write:
- `TestGetByCode_Found` — with and without expiresAt
- `TestGetByCode_NotFound`
- `TestAdd_WithExpiry` / `TestAdd_NoExpiry`
- `TestIncrementUsedCount` — verify used_count goes from 0 to 1 to 2

Add `testhelpers_test.go` with `newTestDB(t)`. Cleanup: truncate `pod_invite` (and `pod`/`player` if needed for FK).

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/podInvite/...` passes (or skips)
3. Smoke test: `POST /api/pod/:id/invite` generates invite; `POST /api/pod/join` uses invite code
