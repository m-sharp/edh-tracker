# Phase 3b — GameResult Repository

## Status
Pending

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `game_result`
- Files: `lib/repositories/gameResult/model.go`, `lib/repositories/gameResult/repo.go`, `lib/repositories/gameResult/repo_test.go`, `lib/repositories/gameResult/stats.go`
- Depends on: game (Phase 3a), deck (Phase 2a)
- Complex: correlated subquery for `player_count` in stats queries; `Aggregate` type preserved
- No business layer changes

## GORM Model

```go
package gameResult

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    GameID    int `gorm:"column:game_id"`
    DeckID    int `gorm:"column:deck_id"`
    Place     int `gorm:"column:place"`
    KillCount int `gorm:"column:kill_count"`
}

func (Model) TableName() string { return "game_result" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByGameId` | `db.Where("game_id = ?", gameID).Find(&results)` | Soft-delete automatic |
| `GetByID` | `db.First(&m, resultID)` | `ErrRecordNotFound` → nil,nil |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM |
| `BulkAdd` | `db.CreateInBatches(&results, 100)` | No return needed |
| `Update` | `db.Model(&Model{}).Where("id = ?", resultID).Updates(map)` | Check RowsAffected |
| `SoftDelete` | `db.Delete(&Model{}, id)` | Sets deleted_at |
| `GetStatsForPlayer` | Raw SQL kept | Correlated subquery; see below |
| `GetStatsForDeck` | Raw SQL kept | Correlated subquery; see below |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — Stats Queries (Raw SQL)

The stats queries use a correlated subquery for `player_count`. Keep these as raw SQL since GORM can't express correlated subqueries cleanly:

```go
const getStatsForPlayer = `SELECT game_result.game_id, game_result.place, game_result.kill_count,
        (SELECT COUNT(*) FROM game_result gr2
          WHERE gr2.game_id = game_result.game_id
            AND gr2.deleted_at IS NULL) AS player_count
   FROM game_result INNER JOIN deck ON game_result.deck_id = deck.id
  WHERE deck.player_id = ?
    AND deck.deleted_at IS NULL
    AND game_result.deleted_at IS NULL;`

const getStatsForDeck = `SELECT game_result.game_id, game_result.place, game_result.kill_count,
      (SELECT COUNT(*) FROM game_result gr2
        WHERE gr2.game_id = game_result.game_id
          AND gr2.deleted_at IS NULL) AS player_count
 FROM game_result INNER JOIN deck ON game_result.deck_id = deck.id
WHERE deck.id = ? AND game_result.deleted_at IS NULL;`
```

Scan into `gameStat` structs using `gorm:"column:..."` tags:

```go
// stats.go — update gameStat to use gorm tags instead of db tags
type gameStat struct {
    GameID      int `gorm:"column:game_id"`
    Place       int `gorm:"column:place"`
    KillCount   int `gorm:"column:kill_count"`
    PlayerCount int `gorm:"column:player_count"`
}

func (r *Repository) GetStatsForPlayer(ctx context.Context, playerID int) (*Aggregate, error) {
    var stats gameStats
    if err := r.db.WithContext(ctx).Raw(getStatsForPlayer, playerID).Scan(&stats).Error; err != nil {
        return nil, fmt.Errorf("failed to get stats for player %d: %w", playerID, err)
    }
    agg := stats.toAggregate()
    return &agg, nil
}
```

`Scan` on raw queries does not error on empty results (unlike `First`) — the `toAggregate` function handles empty `gameStats` correctly by returning zero values.

## Special Pattern — Update (multiple fields)

```go
func (r *Repository) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
    result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", resultID).
        Updates(map[string]any{
            "place":      place,
            "kill_count": killCount,
            "deck_id":    deckID,
        })
    if result.Error != nil {
        return fmt.Errorf("failed to update GameResult record: %w", result.Error)
    }
    if result.RowsAffected != 1 {
        return fmt.Errorf("unexpected rows affected by GameResult update: got %d, expected 1", result.RowsAffected)
    }
    return nil
}
```

Using a map avoids GORM skipping zero integer values (e.g. `place = 0` is valid).

## stats.go Changes

Update `gameStat` struct tags from `db:` to `gorm:` — the `toAggregate()` logic is unchanged:

```go
type gameStat struct {
    GameID      int `gorm:"column:game_id"`
    Place       int `gorm:"column:place"`
    KillCount   int `gorm:"column:kill_count"`
    PlayerCount int `gorm:"column:player_count"`
}
```

The `Aggregate` type and `gameStats.toAggregate()` function are unchanged.

## Test Migration

Remove existing sqlmock tests. Replace with integration tests.

Tests to write:
- `TestGetByGameId`
- `TestGetByID_Found` / `TestGetByID_NotFound`
- `TestAdd`
- `TestBulkAdd`
- `TestUpdate`
- `TestSoftDelete`
- `TestGetStatsForPlayer` — requires game + deck + player rows; verify kills, points, record
- `TestGetStatsForDeck` — similar

Add `testhelpers_test.go` with `newTestDB(t)`. Cleanup: truncate `game_result` (and upstream tables as needed).

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/gameResult/...` passes (or skips)
3. Smoke test: stats endpoints (`GET /api/player/:id/stats`, `GET /api/deck/:id/stats`) return correct data
