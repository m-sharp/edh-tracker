# Phase 2d — PlayerPodRole Repository

## Status
Pending

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `player_pod_role`
- Files: `lib/repositories/playerPodRole/model.go`, `lib/repositories/playerPodRole/repo.go`, `lib/repositories/playerPodRole/repo_test.go`
- Depends on: pod (Phase 2c), player (Phase 1a)
- No business layer changes

## GORM Model

```go
package playerPodRole

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

const (
    RoleManager = "manager"
    RoleMember  = "member"
)

type Model struct {
    base.GormModelBase
    PodID    int    `gorm:"column:pod_id"`
    PlayerID int    `gorm:"column:player_id"`
    Role     string `gorm:"column:role"`
}

func (Model) TableName() string { return "player_pod_role" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetRole` | `db.Where("pod_id = ? AND player_id = ?", ...).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `SetRole` | `db.Clauses(clause.OnConflict{...}).Create(&m)` | Upsert — see below |
| `GetMembersWithRoles` | `db.Where("pod_id = ?", podID).Find(&rows)` | Soft-delete automatic |
| `BulkAdd` | `db.CreateInBatches(&entries, 100)` | No upsert needed for initial bulk |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — SetRole (Upsert via ON DUPLICATE KEY UPDATE)

The current SQL is:
```sql
INSERT INTO player_pod_role (pod_id, player_id, role) VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE role = VALUES(role), deleted_at = NULL;
```

GORM equivalent using `clause.OnConflict`:

```go
import "gorm.io/gorm/clause"

func (r *Repository) SetRole(ctx context.Context, podID, playerID int, role string) error {
    m := Model{PodID: podID, PlayerID: playerID, Role: role}
    err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "pod_id"}, {Name: "player_id"}},
        DoUpdates: clause.Assignments(map[string]any{
            "role":       role,
            "deleted_at": nil,
        }),
    }).Create(&m).Error
    if err != nil {
        return fmt.Errorf("failed to set role %q for player %d in pod %d: %w", role, playerID, podID, err)
    }
    return nil
}
```

This relies on the UNIQUE constraint on `(pod_id, player_id)` in the table schema. Verify this constraint exists in the migration.

## Special Pattern — BulkAdd

BulkAdd does not need upsert — it's called only when initially populating roles (pod creation). Plain `CreateInBatches` is sufficient:

```go
func (r *Repository) BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error {
    if len(playerIDs) == 0 {
        return nil
    }
    entries := make([]Model, len(playerIDs))
    for i, id := range playerIDs {
        entries[i] = Model{PodID: podID, PlayerID: id, Role: role}
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
        return fmt.Errorf("failed to bulk insert player_pod_role records: %w", err)
    }
    return nil
}
```

## Test Migration

Remove existing sqlmock tests in `repo_test.go`. Replace with integration tests.

Tests to write:
- `TestGetRole_Found` / `TestGetRole_NotFound`
- `TestSetRole_Insert` — new role created
- `TestSetRole_Update` — existing role updated (upsert)
- `TestSetRole_Restore` — previously soft-deleted role restored (deleted_at → NULL)
- `TestGetMembersWithRoles` — returns all non-deleted roles for pod
- `TestBulkAdd`

Add `testhelpers_test.go` with `newTestDB(t)`. Cleanup: truncate `player_pod_role` (and `pod`/`player` if needed for FK).

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/playerPodRole/...` passes (or skips)
3. Smoke test: pod role promotion and `GET /api/pod/:id/members` endpoints work
