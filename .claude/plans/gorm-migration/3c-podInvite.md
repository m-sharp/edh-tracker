# Phase 3c — PodInvite Repository

## Status
Approved

## Skill
Use the `/gorm` skill tool at the start of each implementation session for this phase.

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
    PodID             int
    InviteCode        string
    CreatedByPlayerID int
    ExpiresAt         *time.Time
    UsedCount         int
}

func (Model) TableName() string { return "pod_invite" }
```

`ExpiresAt` is nullable (`*time.Time`) — GORM handles this correctly.

All field names infer correct column names via GORM's snake_case convention (`PodID` → `pod_id`, `InviteCode` → `invite_code`, `CreatedByPlayerID` → `created_by_player_id`, `UsedCount` → `used_count`). No explicit `gorm:"column:..."` tags needed.

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByCode` | `db.Where("invite_code = ?", code).First(&m)` | `ErrRecordNotFound` → nil,nil; soft-delete scope automatic |
| `Add` | `db.Create(&m)` | No ID needed; return error only |
| `IncrementUsedCount` | `gorm.Expr("used_count + 1")` via `Model(&Model{}).Where(...).Update(...)` | See below |

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

## Special Pattern — IncrementUsedCount (gorm.Expr)

`used_count = used_count + 1` must use a SQL expression rather than a Go literal to avoid GORM sending the pre-incremented value. Use `gorm.Expr`:

```go
func (r *Repository) IncrementUsedCount(ctx context.Context, code string) error {
    result := r.db.WithContext(ctx).Model(&Model{}).
        Where("invite_code = ?", code).
        Update("used_count", gorm.Expr("used_count + 1"))
    if result.Error != nil {
        return fmt.Errorf("failed to increment used_count for invite %q: %w", code, result.Error)
    }
    return nil
}
```

`Model(&Model{})` applies the soft-delete scope automatically (`AND deleted_at IS NULL`). In practice this is unreachable — the business layer always calls `GetByCode` first (which won't return deleted invites) before calling `IncrementUsedCount`. No `RowsAffected` check needed.

## Behavior Changes from sqlx Migration

**`GetByCode` — behavior equivalent, no change**

The original SQL explicitly has `AND deleted_at IS NULL` and `LIMIT 1`. GORM's `First` on a model embedding `gorm.DeletedAt` applies the soft-delete scope automatically, and `First` always adds `LIMIT 1`. Behavior is identical.

**`IncrementUsedCount` — soft-delete scope now applied**

The original SQL has no `deleted_at IS NULL` guard — it could increment `used_count` on a soft-deleted invite. The `gorm.Expr` implementation uses `Model(&Model{})` which adds the soft-delete scope automatically, making updates to deleted invites a silent no-op. This is an intentional minor change: in practice the business layer always calls `GetByCode` first (which returns nil for deleted invites), so `IncrementUsedCount` is never called with a deleted code.

## Test Migration

Remove existing sqlmock tests. Replace with integration tests.

**FK prerequisites:** `pod_invite` has foreign keys to `pod(id)` and `player(id)`. Each test that inserts a `pod_invite` row must first insert prerequisite `player` and `pod` rows within the same transaction. The tx rollback cleanup handles all of these automatically.

Tests to write:
- `TestGetByCode_Found_WithExpiry` — invite with non-nil `ExpiresAt`; verify all fields
- `TestGetByCode_Found_NoExpiry` — invite with nil `ExpiresAt`; verify `ExpiresAt` is nil
- `TestGetByCode_NotFound`
- `TestAdd_WithExpiry` / `TestAdd_NoExpiry` — insert then `GetByCode` to verify
- `TestIncrementUsedCount` — insert invite (used_count = 0), call once → verify 1, call again → verify 2; read back via `GetByCode` after each call

Use `base.NewTestDB(t)` from `lib/repositories/base/testHelpers.go`. Define a `newRepo(t)` helper in `repo_test.go` (see Phase 1a pattern). No `testhelpers_test.go` needed.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/podInvite/...` passes (or skips)
3. Smoke test: `POST /api/pod/:id/invite` generates invite; `POST /api/pod/join` uses invite code
