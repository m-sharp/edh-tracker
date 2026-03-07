# Phase 2b — DeckCommander Repository

## Status
Approved

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Scope

- Table: `deck_commander`
- Files: `lib/repositories/deckCommander/model.go`, `lib/repositories/deckCommander/repo.go`
- No existing tests — write new integration tests
- Depends on: deck (Phase 2a), commander (Phase 1c)
- No business logic changes; `lib/business/deck/functions_test.go` requires a minor collateral fix (see below)

## GORM Model

```go
package deckCommander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
    base.GormModelBase
    DeckID             int
    CommanderID        int
    PartnerCommanderID *int
}

func (Model) TableName() string { return "deck_commander" }
```

`PartnerCommanderID` is nullable — GORM handles `*int` correctly (stores NULL when nil).

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByDeckId` | `db.Where("deck_id = ?", deckID).First(&m)` | `ErrRecordNotFound` → nil,nil; `.First()` adds `ORDER BY id ASC LIMIT 1` (original SQL had no ORDER BY — safe, a deck has at most one active commander entry) |
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

## Collateral Fix — `lib/business/deck/functions_test.go`

Two test helper closures construct `deckCommanderrepo.Model` using the old embedded field name `ModelBase: base.ModelBase{ID: 1}`. After the model switches to `GormModelBase`, this fails to compile. Drop the base field entirely — the ID is never asserted on in these tests:

```go
// Before (two locations, lines ~222 and ~245):
&deckCommanderrepo.Model{ModelBase: base.ModelBase{ID: 1}, DeckID: 7, CommanderID: 5}

// After:
&deckCommanderrepo.Model{DeckID: 7, CommanderID: 5}
```

If, after removing the `base.ModelBase` references, the `base` import is no longer used elsewhere in the file, remove it from the import block.

## Test Migration

No existing tests. Write new integration tests:
- `TestGetByDeckId_Found` — with and without partner commander
- `TestGetByDeckId_NotFound`
- `TestAdd_WithPartner` / `TestAdd_WithoutPartner`
- `TestBulkAdd`
- `TestDeleteByDeckID` — soft-deletes all entries; GetByDeckId returns nil after

Add `testhelpers_test.go` with `newTestDB(t)` (tx rollback pattern — see Phase 0). No explicit cleanup needed: `t.Cleanup` rolls back the transaction automatically.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/deckCommander/...` passes (or skips)
3. Smoke test: `POST /api/deck` (creates deck + commander) returns 201; `GET /api/deck/:id` returns commander info
