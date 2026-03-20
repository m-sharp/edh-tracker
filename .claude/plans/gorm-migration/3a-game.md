# Phase 3a — Game Repository

## Status
Approved

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `game`
- Files: `lib/repositories/game/model.go`, `lib/repositories/game/repo.go`, `lib/repositories/game/repo_test.go`
- Depends on: pod (Phase 2c), format (Phase 1b)
- Complex: JOIN queries for GetAllByDeck and GetAllByPlayerID; BulkAdd returning []int IDs
- No business layer changes

## GORM Model

```go
package game

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    Description string
    PodID       int
    FormatID    int
}

func (Model) TableName() string { return "game" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetAllByPod` | `db.Where("pod_id = ?", podID).Find(&games)` | Soft-delete automatic |
| `GetAllByDeck` | JOIN with `game_result` — see below | |
| `GetAllByPlayerID` | Double JOIN with DISTINCT — see below | |
| `GetById` | `db.First(&m, gameID)` | `ErrRecordNotFound` → nil,nil |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM |
| `BulkAdd` | `db.CreateInBatches(&games, 100)` | Return IDs from structs — see below |
| `Update` | `db.Model(&Model{}).Where("id = ?", gameID).Update("description", desc)` | Check RowsAffected |
| `SoftDelete` | `db.Delete(&Model{}, id)` | Sets deleted_at |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — GetAllByDeck (JOIN)

Old SQL JOINs `game` and `game_result` on `game.id = game_result.game_id`:

```go
func (r *Repository) GetAllByDeck(ctx context.Context, deckID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).
        Joins("INNER JOIN game_result ON game.id = game_result.game_id").
        Where("game_result.deck_id = ? AND game_result.deleted_at IS NULL", deckID).
        Find(&games).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Game records for deck %d: %w", deckID, err)
    }
    if games == nil {
        return []Model{}, nil
    }
    return games, nil
}
```

`game.deleted_at IS NULL` is applied automatically by GORM on the `game` model.

## Special Pattern — GetAllByPlayerID (DISTINCT double JOIN)

Old SQL JOINs `game` → `game_result` → `deck` with DISTINCT:

```go
func (r *Repository) GetAllByPlayerID(ctx context.Context, playerID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).
        Distinct("game.id, game.description, game.pod_id, game.format_id, game.created_at, game.updated_at, game.deleted_at").
        Joins("INNER JOIN game_result ON game.id = game_result.game_id").
        Joins("INNER JOIN deck ON game_result.deck_id = deck.id").
        Where("deck.player_id = ? AND game_result.deleted_at IS NULL AND deck.deleted_at IS NULL", playerID).
        Find(&games).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Game records for player %d: %w", playerID, err)
    }
    if games == nil {
        return []Model{}, nil
    }
    return games, nil
}
```

## Special Pattern — BulkAdd returning []int IDs

The old implementation uses the `LastInsertId` + sequential offset trick (assumes auto-increment is contiguous). GORM populates IDs on each struct after `CreateInBatches`, so IDs can be read directly — no sequential assumption needed:

```go
func (r *Repository) BulkAdd(ctx context.Context, games []Model) ([]int, error) {
    if len(games) == 0 {
        return []int{}, nil
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&games, 100).Error; err != nil {
        return nil, fmt.Errorf("failed to bulk insert Game records: %w", err)
    }
    ids := make([]int, len(games))
    for i, g := range games {
        ids[i] = g.ID
    }
    return ids, nil
}
```

This is strictly more correct than the old sequential-offset approach.

## Special Pattern — Add

```go
func (r *Repository) Add(ctx context.Context, description string, podID, formatID int) (int, error) {
    m := Model{Description: description, PodID: podID, FormatID: formatID}
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return 0, fmt.Errorf("failed to insert Game record: %w", err)
    }
    return m.ID, nil
}
```

## Special Pattern — Update

```go
func (r *Repository) Update(ctx context.Context, gameID int, description string) error {
    result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", gameID).Update("description", description)
    if result.Error != nil {
        return fmt.Errorf("failed to update Game record: %w", result.Error)
    }
    if result.RowsAffected != 1 {
        return fmt.Errorf("unexpected number of rows affected by Game update: got %d, expected 1", result.RowsAffected)
    }
    return nil
}
```

`db.Model(&Model{})` applies the soft-delete scope automatically (`AND deleted_at IS NULL`), matching the original SQL's explicit `WHERE id = ? AND deleted_at IS NULL`.

## TODO — Cascading Soft Delete

The existing code has a TODO comment:
```go
// TODO: Soft deleting a game should also delete all associated GameResult records
```

This is out of scope for the GORM migration phase. Do NOT implement it here — note it as a follow-up.

## Behavior Changes from sqlx Migration

**`SoftDelete` — GORM adds soft-delete scope; RowsAffected check dropped**

The original SQL (`UPDATE game SET deleted_at = NOW() WHERE id = ?`) has no `deleted_at IS NULL` guard — re-deleting an already-deleted game would re-stamp `deleted_at` and return `numAffected = 1`. The original implementation then checks `numAffected != 1` and errors on 0.

GORM's `db.Delete(&Model{}, id)` adds the soft-delete scope automatically, so re-deleting an already-deleted game is a no-op (`RowsAffected = 0`, no error). The RowsAffected check is dropped, consistent with the approach used for deck and pod in earlier phases.

Unlike the deck repository, there is no business-layer existence guard for game (`SoftDelete` in `lib/business/game/functions.go` calls the repo directly). A delete of a non-existent or already-deleted game silently succeeds.

**`Update` — behavior equivalent, no change**

The original SQL has an explicit `WHERE id = ? AND deleted_at IS NULL`. GORM's `db.Model(&Model{}).Where("id = ?", gameID)` adds the same scope automatically. The RowsAffected check is preserved — updating a non-existent or already-deleted game still returns an error.

## Test Migration

Remove existing sqlmock tests. Replace with integration tests.

Tests to write:
- `TestGetAllByPod`
- `TestGetAllByDeck` — requires game_result rows
- `TestGetAllByPlayerID` — requires game_result + deck rows
- `TestGetById_Found` / `TestGetById_NotFound`
- `TestAdd`
- `TestBulkAdd` — IDs are non-sequential, all correct
- `TestUpdate` / `TestUpdate_NotFound`
- `TestSoftDelete`

Add `testhelpers_test.go` with `newTestDB(t)` (tx rollback pattern — see Phase 0). No explicit cleanup needed: `t.Cleanup` rolls back the transaction automatically.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/game/...` passes (or skips)
3. Smoke test: `GET /api/games/pod/:id` and game creation endpoints work
