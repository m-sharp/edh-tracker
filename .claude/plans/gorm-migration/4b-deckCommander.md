# Phase 4b — Preloading: deck.Model → deckCommander.Model → commander.Model

## Status
Approved

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Goal

Eliminate the per-deck commander lookups in deck listing functions
(`deck.GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID`).

Current pattern: per-deck `getCommanderEntry` calls make 1 query to `deck_commander`
plus 1–2 queries to `commander` (for name + optional partner name) = 2–3 queries per deck.
After: up to 3 batched queries regardless of deck count (deck_commander rows, commander rows,
and a separate batch for partner commander rows when any deck has a partner), loaded via
nested GORM `Preload`.

## Scope

- `lib/repositories/deckCommander/model.go` — add `Commander` + `PartnerCommander` BelongsTo fields
- `lib/repositories/deck/model.go` — add `Commander` HasOne field
- `lib/repositories/deck/repo.go` — add `Get*WithCommanders` methods
- `lib/repositories/interfaces.go` — add new `DeckRepository` signatures
- `lib/business/deck/functions.go` — add `commanderInfoFromModel` helper; update `GetAll`,
  `GetAllForPlayer`, `GetAllByPod`, `GetByID` to drop `getCommanderEntry` dep

Depends on: Phase 2a (deck), Phase 2b (deckCommander), Phase 1c (commander) GORM-migrated

## Association Declarations

### `lib/repositories/deckCommander/model.go`

Add BelongsTo associations for both commander slots. GORM convention matches field names
to existing FK columns (`CommanderID` → `commander.id`, `PartnerCommanderID` → `commander.id`):

```go
import commanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/commander"

type Model struct {
    base.GormModelBase
    DeckID             int
    CommanderID        int
    PartnerCommanderID *int
    // Associations — populated when preloaded.
    Commander        commanderRepo.Model  `gorm:"foreignKey:CommanderID"`        // BelongsTo via CommanderID
    PartnerCommander *commanderRepo.Model `gorm:"foreignKey:PartnerCommanderID"` // BelongsTo via PartnerCommanderID (nil if no partner)
}

func (Model) TableName() string { return "deck_commander" }
```

When `PartnerCommanderID` is nil, GORM leaves `PartnerCommander` as nil — no query issued.

### `lib/repositories/deck/model.go`

Add HasOne for the deck's commander record:

```go
import deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"

type Model struct {
    base.GormModelBase
    PlayerID int
    Name     string
    FormatID int
    Retired  bool
    // Populated only when using Get*WithCommanders methods.
    Commander *deckCommanderRepo.Model `gorm:"foreignKey:DeckID"`
}

func (Model) TableName() string { return "deck" }
```

## New Repository Methods

Add to `lib/repositories/deck/repo.go` (additive — do not replace existing methods):

```go
// preloadCommanders is the shared Preload chain for commander associations.
func (r *Repository) preloadCommanders(db *gorm.DB) *gorm.DB {
    return db.Preload("Commander.Commander").Preload("Commander.PartnerCommander")
}

func (r *Repository) GetAllWithCommanders(ctx context.Context) ([]Model, error) {
    var decks []Model
    err := r.preloadCommanders(r.db.WithContext(ctx)).
        Where("retired = ?", false).Find(&decks).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Decks with commanders: %w", err)
    }
    return decks, nil
}

func (r *Repository) GetAllForPlayerWithCommanders(ctx context.Context, playerID int) ([]Model, error) {
    var decks []Model
    err := r.preloadCommanders(r.db.WithContext(ctx)).
        Where("player_id = ?", playerID).Find(&decks).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Decks with commanders for player %d: %w", playerID, err)
    }
    return decks, nil
}

func (r *Repository) GetAllByPlayerIDsWithCommanders(ctx context.Context, playerIDs []int) ([]Model, error) {
    if len(playerIDs) == 0 {
        return []Model{}, nil
    }
    var decks []Model
    err := r.preloadCommanders(r.db.WithContext(ctx)).
        Where("player_id IN ?", playerIDs).Find(&decks).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Decks with commanders for player IDs: %w", err)
    }
    return decks, nil
}

func (r *Repository) GetByIDWithCommanders(ctx context.Context, deckID int) (*Model, error) {
    var m Model
    err := r.preloadCommanders(r.db.WithContext(ctx)).First(&m, deckID).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get Deck with commanders for id %d: %w", deckID, err)
    }
    return &m, nil
}
```

The `Preload("Commander.Commander")` call asks GORM to:
1. Load `deck_commander` rows where `deck_id IN (...)` — one batch query
2. For each loaded `deckCommander`, load `commander` rows where `id IN (...)` — one batch query
3. Same for `PartnerCommander` (GORM issues a separate batch for the partner FK)

## Interface Updates

Add to `DeckRepository` in `lib/repositories/interfaces.go`:

```go
GetAllWithCommanders(ctx context.Context) ([]deck.Model, error)
GetAllForPlayerWithCommanders(ctx context.Context, playerID int) ([]deck.Model, error)
GetAllByPlayerIDsWithCommanders(ctx context.Context, playerIDs []int) ([]deck.Model, error)
GetByIDWithCommanders(ctx context.Context, deckID int) (*deck.Model, error)
```

## Business Layer Changes

### New helper `commanderInfoFromModel` in `lib/business/deck/functions.go`

Reads commander data from a preloaded `deck.Model.Commander` field instead of making
a DB call:

```go
func commanderInfoFromModel(d deckRepo.Model) *CommanderInfo {
    if d.Commander == nil {
        return nil
    }
    entry := &CommanderInfo{
        CommanderID:   d.Commander.CommanderID,
        CommanderName: d.Commander.Commander.Name,
    }
    if d.Commander.PartnerCommanderID != nil {
        entry.PartnerCommanderID = d.Commander.PartnerCommanderID
        // PartnerCommander is guaranteed non-nil here: GORM always loads a BelongsTo
        // association when its FK is non-nil and referential integrity holds.
        name := d.Commander.PartnerCommander.Name
        entry.PartnerCommanderName = &name
    }
    return entry
}
```

### Updated listing functions

`GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID` switch to the `*WithCommanders`
repo methods and call `commanderInfoFromModel(d)` instead of `getCommanderEntry(ctx, d.ID)`.

The `getCommanderEntry GetCommanderEntryFunc` parameter is **removed** from these four
constructors. `GetCommanderEntry` (the closure constructor) remains in place for callers
that still need point lookups by deck ID (e.g., `gameResult.GetByGameID` until Phase 4d).

Note: `GetAllByPod`'s internal deck fetch changes from `deckRepo.GetAllByPlayerIDs` to
`deckRepo.GetAllByPlayerIDsWithCommanders` — this is the only non-obvious repo call change
in the business layer (the enrichment loop structure is otherwise unchanged).

Example diff for `GetAll`:

```go
func GetAll(
    deckRepo repos.DeckRepository,
    gameResultRepo repos.GameResultRepository,
    getPlayerName player.GetPlayerNameFunc,
    getFormat format.GetByIDFunc,
    // getCommanderEntry removed
) GetAllFunc {
    return func(ctx context.Context) ([]EntityWithStats, error) {
        decks, err := deckRepo.GetAllWithCommanders(ctx)
        // ...
        for _, d := range decks {
            commanders := commanderInfoFromModel(d)   // reads from d.Commander
            // getFormat and getPlayerName calls unchanged (addressed in Phase 4c)
            entity := ToEntity(d, playerName, f.Name, commanders)
            // ...
        }
    }
}
```

### `lib/business/business.go` wiring update

The four deck constructors no longer receive `getCommanderEntry`:

```go
Decks: deck.Functions{
    GetAll:          deck.GetAll(r.Decks, r.GameResults, getPlayerName, getFormat),
    GetAllForPlayer: deck.GetAllForPlayer(r.Decks, r.GameResults, getPlayerName, getFormat),
    GetAllByPod:     deck.GetAllByPod(r.Decks, r.Pods, r.GameResults, getPlayerName, getFormat),
    GetByID:         deck.GetByID(r.Decks, r.GameResults, getPlayerName, getFormat),
    // ... other fields unchanged
    GetCommanderEntry: getCommanderEntry,   // still present for gameResult enrichment
},
```

## Query Reduction (deck listing)

| Function | Before | After |
|----------|--------|-------|
| GetAll (N decks) | 2–3 per-deck | up to 3 batched (deck_commander + commander + partner commander) |
| GetAllForPlayer (N decks) | 2–3 per-deck | up to 3 batched |
| GetAllByPod (N decks) | 2–3 per-deck | up to 3 batched |
| GetByID | 2–3 queries | up to 3 queries |

`getPlayerName` and `getFormat` per-deck calls unchanged — addressed in Phase 4c.

## Tests

- Integration tests for each `Get*WithCommanders` repo method
- Verify that `Commander`, `Commander.Commander`, and `Commander.PartnerCommander` are
  populated correctly; nil for decks without commanders
- Business layer: update listing function tests to remove `getCommanderEntry` mock

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/deck/...` passes (or skips)
3. `go test ./lib/business/deck/...` passes
4. Smoke test: `GET /api/decks` and `GET /api/deck/:id` return correct commander names
