# Phase 2c — Pod Repository

## Status
Approved

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Tables: `pod`, `player_pod`
- Files: `lib/repositories/pod/model.go`, `lib/repositories/pod/repo.go`
- No existing tests — write new integration tests
- No FK deps on other migrated domains (leaf for this table)
- No business layer changes

## GORM Models

```go
package pod

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    Name string
}

func (Model) TableName() string { return "pod" }

// PlayerPodModel represents the player_pod junction table.
type PlayerPodModel struct {
    base.GormModelBase
    PodID    int
    PlayerID int
}

func (PlayerPodModel) TableName() string { return "player_pod" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetAll` | `db.Find(&pods)` | Soft-delete automatic on pod |
| `GetByID` | `db.First(&m, podID)` | `ErrRecordNotFound` → nil,nil |
| `GetByPlayerID` | JOIN query — see below | `player_pod` JOIN |
| `GetByName` | `db.Where("name = ?", name).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `GetIDsByPlayerID` | `db.Model(&PlayerPodModel{}).Where(...).Pluck("pod_id", &ids)` | Pluck scalar list |
| `GetPlayerIDs` | `db.Model(&PlayerPodModel{}).Where(...).Pluck("player_id", &ids)` | Pluck scalar list |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM |
| `BulkAddPlayers` | `db.CreateInBatches(&entries, 100)` | Junction table bulk insert |
| `AddPlayerToPod` | `db.Create(&PlayerPodModel{...})` | Error-only; RowsAffected check dropped (see below) |
| `SoftDelete` | `db.Delete(&Model{}, podID)` | Sets deleted_at |
| `Update` | `db.Model(&Model{}).Where("id = ?", podID).Update("name", name)` | |
| `RemovePlayer` | `db.Where("pod_id = ? AND player_id = ?", ...).Delete(&PlayerPodModel{})` | Soft-delete junction row |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — GetByPlayerID (JOIN)

```go
func (r *Repository) GetByPlayerID(ctx context.Context, playerID int) ([]Model, error) {
    var pods []Model
    err := r.db.WithContext(ctx).
        Joins("INNER JOIN player_pod ON pod.id = player_pod.pod_id").
        Where("player_pod.player_id = ? AND player_pod.deleted_at IS NULL", playerID).
        Find(&pods).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Pod records for player %d: %w", playerID, err)
    }
    if pods == nil {
        return []Model{}, nil
    }
    return pods, nil
}
```

Note: `pod.deleted_at IS NULL` is handled automatically because `Model` embeds `GormModelBase` with `gorm.DeletedAt`. The explicit `player_pod.deleted_at IS NULL` filter is needed because GORM doesn't know to apply it to the joined table.

## Special Pattern — GetIDsByPlayerID / GetPlayerIDs (Pluck)

```go
func (r *Repository) GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error) {
    var ids []int
    err := r.db.WithContext(ctx).
        Model(&PlayerPodModel{}).
        Where("player_id = ? AND deleted_at IS NULL", playerID).
        Pluck("pod_id", &ids).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Pod IDs for player %d: %w", playerID, err)
    }
    if ids == nil {
        return []int{}, nil
    }
    return ids, nil
}

func (r *Repository) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
    var ids []int
    err := r.db.WithContext(ctx).
        Model(&PlayerPodModel{}).
        Where("pod_id = ? AND deleted_at IS NULL", podID).
        Pluck("player_id", &ids).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get player IDs for pod %d: %w", podID, err)
    }
    if ids == nil {
        return []int{}, nil
    }
    return ids, nil
}
```

Note: `PlayerPodModel` embeds `GormModelBase` (which contains `gorm.DeletedAt`), so GORM applies `AND deleted_at IS NULL` automatically when `Model(&PlayerPodModel{})` is used. The explicit `AND deleted_at IS NULL` in the WHERE clause is redundant but retained for readability.

## Special Pattern — BulkAddPlayers

```go
func (r *Repository) BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error {
    if len(playerIDs) == 0 {
        return nil
    }
    entries := make([]PlayerPodModel, len(playerIDs))
    for i, id := range playerIDs {
        entries[i] = PlayerPodModel{PodID: podID, PlayerID: id}
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
        return fmt.Errorf("failed to bulk insert PlayerPod records: %w", err)
    }
    return nil
}
```

## Special Pattern — AddPlayerToPod

The current implementation checks `numAffected != 1` after `ExecContext`. With GORM's `Create`, the insert either succeeds (exactly 1 row) or returns an error — the RowsAffected check adds nothing:

```go
func (r *Repository) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
    m := PlayerPodModel{PodID: podID, PlayerID: playerID}
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return fmt.Errorf("failed to insert PlayerPod record: %w", err)
    }
    return nil
}
```

## Behavior Changes from sqlx Migration

**`SoftDelete` (pod) — GORM adds soft-delete scope**

The original SQL (`UPDATE pod SET deleted_at = NOW() WHERE id = ?`) has no `deleted_at IS NULL` guard — it would re-stamp `deleted_at` on an already-deleted pod. GORM's `db.Delete(&Model{}, podID)` adds the soft-delete scope automatically, so it has no effect on already-deleted pods. Strictly safer behavior.

**`AddPlayerToPod` — RowsAffected check removed**

The original implementation checks `numAffected != 1` and returns an error. The GORM version relies solely on `Create`'s error return, which is equivalent — `Create` always inserts exactly 1 row or errors.

**`Update` (pod) — no RowsAffected check in either version**

Neither the original SQL (`UPDATE pod SET name = ? WHERE id = ? AND deleted_at IS NULL`) nor the GORM version checks RowsAffected. If the pod ID doesn't exist or is already deleted, the update is a silent no-op. This is consistent behavior, not a regression.

## Test Migration

No existing tests. Write new integration tests:
- `TestGetAll` — returns non-deleted pods
- `TestGetByID_Found` / `TestGetByID_NotFound`
- `TestGetByPlayerID` — pod returned via player_pod join
- `TestGetByName_Found` / `TestGetByName_NotFound`
- `TestGetIDsByPlayerID` — pluck pod IDs
- `TestGetPlayerIDs` — pluck player IDs
- `TestAdd`
- `TestBulkAddPlayers`
- `TestAddPlayerToPod`
- `TestSoftDelete`
- `TestUpdate`
- `TestRemovePlayer` — player_pod row soft-deleted; not returned by GetPlayerIDs

Use `base.NewTestDB(t)` from `lib/repositories/base/testHelpers.go`. Define a `newRepo(t)` helper in `repo_test.go` (see Phase 1a pattern). No `testhelpers_test.go` needed.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/pod/...` passes (or skips)
3. Smoke test: `GET /api/pods` and pod creation/membership flows work
