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
    GameID    int
    DeckID    int
    Place     int
    KillCount int
}

func (Model) TableName() string { return "game_result" }
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByGameId` | `db.Where("game_id = ?", gameID).Find(&results)` | Soft-delete automatic |
| `GetByID` | `db.First(&m, resultID)` | `ErrRecordNotFound` → nil,nil |
| `Add` | `db.Create(&m)` | Takes full `Model`; `m.ID` set by GORM — see below |
| `BulkAdd` | `db.CreateInBatches(&results, 100)` | Returns `error` only — see below |
| `Update` | `db.Model(&Model{}).Where("id = ?", resultID).Updates(map)` | Check RowsAffected |
| `SoftDelete` | `db.Delete(&Model{}, id)` | Sets deleted_at; RowsAffected check dropped — see below |
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

## Special Pattern — Add (takes full Model)

Unlike other repos where `Add` takes individual fields, `gameResult.Add` takes a full `Model` struct (interface requirement — no business layer changes). GORM populates `m.ID` after create:

```go
func (r *Repository) Add(ctx context.Context, m Model) (int, error) {
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return 0, fmt.Errorf("failed to insert GameResult record: %w", err)
    }
    return m.ID, nil
}
```

## Special Pattern — BulkAdd (returns error only)

Unlike game's `BulkAdd` which returns `[]int` IDs, gameResult's interface returns only `error`. No callers need the inserted IDs:

```go
func (r *Repository) BulkAdd(ctx context.Context, results []Model) error {
    if len(results) == 0 {
        return nil
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&results, 100).Error; err != nil {
        return fmt.Errorf("failed to bulk insert GameResult records: %w", err)
    }
    return nil
}
```

## Special Pattern — Stats Queries (Raw SQL)

The stats queries use a correlated subquery for `player_count`. Keep these as raw SQL since GORM can't express correlated subqueries cleanly. The constants remain in `repo.go` alongside the methods that use them:

```go
// repo.go — constants stay here alongside GetStatsForPlayer / GetStatsForDeck
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

Scan into `gameStat` structs using GORM's naming convention (all fields infer correctly):

```go
func (r *Repository) GetStatsForPlayer(ctx context.Context, playerID int) (*Aggregate, error) {
    var stats gameStats
    if err := r.db.WithContext(ctx).Raw(getStatsForPlayer, playerID).Scan(&stats).Error; err != nil {
        return nil, fmt.Errorf("failed to get stats for player %d: %w", playerID, err)
    }
    agg := stats.toAggregate()
    return &agg, nil
}

func (r *Repository) GetStatsForDeck(ctx context.Context, deckID int) (*Aggregate, error) {
    var stats gameStats
    if err := r.db.WithContext(ctx).Raw(getStatsForDeck, deckID).Scan(&stats).Error; err != nil {
        return nil, fmt.Errorf("failed to get stats for deck %d: %w", deckID, err)
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

Update `gameStat` struct tags from `db:` to GORM — all fields infer the correct column names via GORM's snake_case naming convention (`GameID` → `game_id`, `KillCount` → `kill_count`, `PlayerCount` → `player_count`), so no tags are needed:

```go
type gameStat struct {
    GameID      int
    Place       int
    KillCount   int
    PlayerCount int
}
```

The `Aggregate` type and `gameStats.toAggregate()` function are unchanged.

## Behavior Changes from sqlx Migration

**`SoftDelete` — RowsAffected check dropped**

The original `SoftDelete` runs `UPDATE game_result SET deleted_at = NOW() WHERE id = ?` with no soft-delete guard, then checks `RowsAffected != 1` and errors. GORM's `db.Delete(&Model{}, id)` adds the soft-delete scope automatically (`AND deleted_at IS NULL`), making re-deletion a no-op with `RowsAffected = 0` and no error. The RowsAffected check is dropped, consistent with the approach used in earlier phases.

There is no business-layer existence guard for `SoftDelete` (the business layer calls the repo directly). A delete of a non-existent or already-deleted game result silently succeeds.

## Test Migration

Remove existing sqlmock tests from `repo_test.go`. Replace with integration tests.

**Keep `stats_test.go` unchanged** — it tests `toAggregate()` as a pure unit test with no DB access. It does not use sqlmock and requires no migration.

Tests to write in `repo_test.go`:
- `TestGetByGameId`
- `TestGetByID_Found` / `TestGetByID_NotFound`
- `TestAdd`
- `TestBulkAdd`
- `TestUpdate` / `TestUpdate_NotFound`
- `TestSoftDelete`
- `TestGetStatsForPlayer` — requires game + deck + player rows; verify kills, points, record
- `TestGetStatsForDeck` — similar

Use `base.NewTestDB(t)` from `lib/repositories/base/testHelpers.go`. Define a `newRepo(t)` helper in `repo_test.go` (see Phase 1a pattern). No `testhelpers_test.go` needed.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/gameResult/...` passes (or skips)
3. Smoke test: stats endpoints (`GET /api/player/:id/stats`, `GET /api/deck/:id/stats`) return correct data
