# Phase 3a — Game Repository

## Status
Pending

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
    Description string `gorm:"column:description"`
    PodID       int    `gorm:"column:pod_id"`
    FormatID    int    `gorm:"column:format_id"`
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

## TODO — Cascading Soft Delete

The existing code has a TODO comment:
```go
// TODO: Soft deleting a game should also delete all associated GameResult records
```

This is out of scope for the GORM migration phase. Do NOT implement it here — note it as a follow-up.

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

Add `testhelpers_test.go` with `newTestDB(t)`. Cleanup: truncate `game` (and `game_result`, `deck`, `player` if needed for FK deps in join tests).

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/game/...` passes (or skips)
3. Smoke test: `GET /api/games/pod/:id` and game creation endpoints work
