# Phase 2b — DeckCommander Repository

## Status
Pending

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `deck_commander`
- Files: `lib/repositories/deckCommander/model.go`, `lib/repositories/deckCommander/repo.go`
- No existing tests — write new integration tests
- Depends on: deck (Phase 2a), commander (Phase 1c)
- No business layer changes

## GORM Model

```go
package deckCommander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    DeckID             int  `gorm:"column:deck_id"`
    CommanderID        int  `gorm:"column:commander_id"`
    PartnerCommanderID *int `gorm:"column:partner_commander_id"`
}

func (Model) TableName() string { return "deck_commander" }
```

`PartnerCommanderID` is nullable — GORM handles `*int` correctly (stores NULL when nil).

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByDeckId` | `db.Where("deck_id = ?", deckID).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM |
| `BulkAdd` | `db.CreateInBatches(&entries, 100)` | No return value needed |
| `DeleteByDeckID` | `db.Where("deck_id = ?", deckID).Delete(&Model{})` | Soft-delete all rows for deck |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}
```

## Special Pattern — Nullable PartnerCommanderID

GORM handles `*int` fields transparently — no special handling needed. Pass nil to store NULL:

```go
func (r *Repository) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
    m := Model{
        DeckID:             deckID,
        CommanderID:        commanderID,
        PartnerCommanderID: partnerCommanderID,
    }
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return 0, fmt.Errorf("failed to insert DeckCommander record: %w", err)
    }
    return m.ID, nil
}
```

## Special Pattern — DeleteByDeckID (soft-delete by FK)

`db.Delete()` with a `Where` clause sets `deleted_at` on all matching rows:

```go
func (r *Repository) DeleteByDeckID(ctx context.Context, deckID int) error {
    if err := r.db.WithContext(ctx).Where("deck_id = ?", deckID).Delete(&Model{}).Error; err != nil {
        return fmt.Errorf("failed to soft-delete DeckCommander records for deck %d: %w", deckID, err)
    }
    return nil
}
```

Note: No RowsAffected check here — it's valid for a deck to have 0 entries (e.g. after a previous delete).

## Special Pattern — BulkAdd

```go
func (r *Repository) BulkAdd(ctx context.Context, entries []Model) error {
    if len(entries) == 0 {
        return nil
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
        return fmt.Errorf("failed to bulk insert DeckCommander records: %w", err)
    }
    return nil
}
```

## Test Migration

No existing tests. Write new integration tests:
- `TestGetByDeckId_Found` — with and without partner commander
- `TestGetByDeckId_NotFound`
- `TestAdd_WithPartner` / `TestAdd_WithoutPartner`
- `TestBulkAdd`
- `TestDeleteByDeckID` — soft-deletes all entries; GetByDeckId returns nil after

Add `testhelpers_test.go` with `newTestDB(t)`. Cleanup: truncate `deck_commander` (and upstream tables needed for FK constraints).

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/deckCommander/...` passes (or skips)
3. Smoke test: `POST /api/deck` (creates deck + commander) returns 201; `GET /api/deck/:id` returns commander info
