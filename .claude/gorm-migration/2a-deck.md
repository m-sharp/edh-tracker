# Phase 2a — Deck Repository

## Status
Pending

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `deck`
- Files: `lib/repositories/deck/model.go`, `lib/repositories/deck/repo.go`, `lib/repositories/deck/repo_test.go`
- Depends on: player (Phase 1a), format (Phase 1b)
- No business layer changes

## GORM Model

```go
package deck

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    PlayerID int    `gorm:"column:player_id"`
    Name     string `gorm:"column:name"`
    FormatID int    `gorm:"column:format_id"`
    Retired  bool   `gorm:"column:retired"`
}

func (Model) TableName() string { return "deck" }

// UpdateFields holds the optional fields that may be updated on a deck.
// Only non-nil fields are applied.
type UpdateFields struct {
    Name     *string
    FormatID *int
    Retired  *bool
}
```

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetAll` | `db.Where("retired = ?", false).Find(&decks)` | retired=0 filter preserved; soft-delete automatic |
| `GetAllForPlayer` | `db.Where("player_id = ?", playerID).Find(&decks)` | Includes retired decks (intentional — matches current) |
| `GetAllByPlayerIDs` | `db.Where("player_id IN ?", playerIDs).Find(&decks)` | Replaces sqlx.In + Rebind |
| `GetById` | `db.First(&m, deckID)` | `ErrRecordNotFound` → nil,nil |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM |
| `BulkAdd` | `db.CreateInBatches(&models, 100)` + select-back | See below |
| `Update` | dynamic map-based `db.Updates(map)` | See UpdateFields pattern |
| `Retire` | `db.Model(&Model{}).Where("id = ?", id).Update("retired", true)` | Check RowsAffected |
| `SoftDelete` | `db.Delete(&Model{}, id)` | Sets deleted_at automatically |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — BulkAdd (select-back after insert)

The current implementation selects back inserted decks using player_id IN + name IN. With GORM/CreateInBatches, the IDs are populated on the structs directly — no select-back needed:

```go
func (r *Repository) BulkAdd(ctx context.Context, decks []Model) ([]Model, error) {
    if len(decks) == 0 {
        return []Model{}, nil
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&decks, 100).Error; err != nil {
        return nil, fmt.Errorf("failed to bulk insert Deck records: %w", err)
    }
    return decks, nil  // IDs populated by GORM
}
```

This simplifies the business layer slightly — verify callers only need the returned models (they do: `lib/business/deck/functions.go` uses them for entity construction).

## Special Pattern — Dynamic Update (UpdateFields)

Use a `map[string]any` to avoid GORM skipping zero values:

```go
func (r *Repository) Update(ctx context.Context, deckID int, fields UpdateFields) error {
    updates := map[string]any{}
    if fields.Name != nil     { updates["name"] = *fields.Name }
    if fields.FormatID != nil { updates["format_id"] = *fields.FormatID }
    if fields.Retired != nil  { updates["retired"] = *fields.Retired }
    if len(updates) == 0 {
        return nil
    }
    result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", deckID).Updates(updates)
    if result.Error != nil {
        return fmt.Errorf("failed to update Deck record: %w", result.Error)
    }
    if result.RowsAffected != 1 {
        return fmt.Errorf("unexpected rows affected by Deck update: got %d, expected 1", result.RowsAffected)
    }
    return nil
}
```

## Special Pattern — Retire

```go
func (r *Repository) Retire(ctx context.Context, deckID int) error {
    result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", deckID).Update("retired", true)
    if result.Error != nil {
        return fmt.Errorf("failed to retire Deck: %w", result.Error)
    }
    if result.RowsAffected != 1 {
        return fmt.Errorf("unexpected rows affected by Deck retirement: got %d, expected 1", result.RowsAffected)
    }
    return nil
}
```

## Test Migration

Remove all sqlmock tests in `repo_test.go`. Replace with integration tests.

Tests to write:
- `TestGetAll` — only non-retired, non-deleted decks returned
- `TestGetAllForPlayer` — includes retired decks for that player
- `TestGetAllByPlayerIDs` — multiple player IDs
- `TestGetById_Found` / `TestGetById_NotFound`
- `TestAdd` — returns correct ID
- `TestBulkAdd` — IDs populated, all returned
- `TestUpdate_PartialFields` — only specified fields change
- `TestUpdate_NoFields` — no-op, no error
- `TestRetire` — deck.Retired set to true
- `TestSoftDelete` — deck not returned after delete

Add `testhelpers_test.go` with `newTestDB(t)` (tx rollback pattern — see Phase 0). No explicit cleanup needed: `t.Cleanup` rolls back the transaction automatically.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/deck/...` passes (or skips)
3. Smoke test: `GET /api/decks` and `POST /api/deck` work correctly
