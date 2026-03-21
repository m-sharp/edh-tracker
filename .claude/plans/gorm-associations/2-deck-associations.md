# Phase 2 — Preloading: deck.Model → deckCommander + commander + player + format

## Skill

Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Goal

Eliminate per-deck `getCommanderEntry`, `getPlayerName`, and `getFormat` queries in deck
listing functions (`GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID`).

Current pattern: 4–5 queries per deck (commander entry + commander name + optional partner
name + player name + format). After: ~5 batched queries regardless of deck count.

This phase merges what were originally two separate sub-phases (4b + 4c). The implementation
goes straight to `Get*WithAll` methods — there is no intermediate `Get*WithCommanders` step.

## Dependencies

None — standalone, can run before or in parallel with Phase 1.

## Scope

- `lib/repositories/deckCommander/model.go` — add `Commander` + `PartnerCommander` BelongsTo fields
- `lib/repositories/deck/model.go` — add `Commander`, `Player`, `Format` association fields
- `lib/repositories/deck/repo.go` — add `preloadAll` helper + 4 `Get*WithAll` methods
- `lib/repositories/interfaces.go` — add 4 `Get*WithAll` `DeckRepository` signatures
- `lib/business/deck/functions.go` — add `commanderInfoFromModel` helper; drop `getCommanderEntry`,
  `getPlayerName`, `getFormat` deps from `GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID`
- `lib/business/business.go` — update deck wiring

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

Add all three association fields in one pass:

```go
import (
    deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
    formatRepo        "github.com/m-sharp/edh-tracker/lib/repositories/format"
    playerRepo        "github.com/m-sharp/edh-tracker/lib/repositories/player"
)

type Model struct {
    base.GormModelBase
    PlayerID int
    Name     string
    FormatID int
    Retired  bool
    // Association fields — populated only when using Get*WithAll methods.
    Commander *deckCommanderRepo.Model `gorm:"foreignKey:DeckID"`
    Player    playerRepo.Model         // BelongsTo via PlayerID
    Format    formatRepo.Model         // BelongsTo via FormatID
}

func (Model) TableName() string { return "deck" }
```

GORM BelongsTo convention: `Player` field + `PlayerID` FK column → `player.id = deck.player_id`.
Same for `Format` / `FormatID`.

## New Repository Methods

Add to `lib/repositories/deck/repo.go` (additive — do not replace existing flat methods):

```go
func (r *Repository) preloadAll(db *gorm.DB) *gorm.DB {
    return db.
        Preload("Commander.Commander").
        Preload("Commander.PartnerCommander").
        Preload("Player").
        Preload("Format")
}

func (r *Repository) GetAllWithAll(ctx context.Context) ([]Model, error) {
    var decks []Model
    err := r.preloadAll(r.db.WithContext(ctx)).
        Where("retired = ?", false).Find(&decks).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Decks with all associations: %w", err)
    }
    return decks, nil
}

func (r *Repository) GetAllForPlayerWithAll(ctx context.Context, playerID int) ([]Model, error) {
    var decks []Model
    err := r.preloadAll(r.db.WithContext(ctx)).
        Where("player_id = ?", playerID).Find(&decks).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Decks with all associations for player %d: %w", playerID, err)
    }
    return decks, nil
}

func (r *Repository) GetAllByPlayerIDsWithAll(ctx context.Context, playerIDs []int) ([]Model, error) {
    if len(playerIDs) == 0 {
        return []Model{}, nil
    }
    var decks []Model
    err := r.preloadAll(r.db.WithContext(ctx)).
        Where("player_id IN ?", playerIDs).Find(&decks).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Decks with all associations for player IDs: %w", err)
    }
    return decks, nil
}

func (r *Repository) GetByIDWithAll(ctx context.Context, deckID int) (*Model, error) {
    var m Model
    err := r.preloadAll(r.db.WithContext(ctx)).First(&m, deckID).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get Deck with all associations for id %d: %w", deckID, err)
    }
    return &m, nil
}
```

The `Preload("Commander.Commander")` call asks GORM to:
1. Load `deck_commander` rows where `deck_id IN (...)` — one batch query
2. For each loaded `deckCommander`, load `commander` rows where `id IN (...)` — one batch query
3. Same for `PartnerCommander` (GORM issues a separate batch for the partner FK)
4. Load `player` rows where `id IN (...)` — one batch query
5. Load `format` rows where `id IN (...)` — one batch query

## Interface Updates

Add to `DeckRepository` in `lib/repositories/interfaces.go`:

```go
GetAllWithAll(ctx context.Context) ([]deck.Model, error)
GetAllForPlayerWithAll(ctx context.Context, playerID int) ([]deck.Model, error)
GetAllByPlayerIDsWithAll(ctx context.Context, playerIDs []int) ([]deck.Model, error)
GetByIDWithAll(ctx context.Context, deckID int) (*deck.Model, error)
```

## Business Layer Changes

### New helper `commanderInfoFromModel` in `lib/business/deck/functions.go`

Reads commander data from a preloaded `deck.Model.Commander` field instead of making DB calls:

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

`GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID` switch to the `*WithAll` repo methods
and use preloaded field values directly:

```go
func GetAll(
    deckRepo repos.DeckRepository,
    gameResultRepo repos.GameResultRepository,
    // getPlayerName removed
    // getFormat removed
    // getCommanderEntry removed
) GetAllFunc {
    return func(ctx context.Context) ([]EntityWithStats, error) {
        decks, err := deckRepo.GetAllWithAll(ctx)
        // ...
        for _, d := range decks {
            commanders := commanderInfoFromModel(d)
            entity := ToEntity(d, d.Player.Name, d.Format.Name, commanders)
            // ...
        }
    }
}
```

Apply the same pattern to `GetAllForPlayer`, `GetAllByPod`, `GetByID`.

Note for `GetAllByPod`: the internal deck fetch changes from `deckRepo.GetAllByPlayerIDs`
to `deckRepo.GetAllByPlayerIDsWithAll`. The local `playerCache` map in `GetAll` and
`GetAllByPod` is removed — GORM batch preload replaces it.

### `lib/business/business.go` wiring

After this phase, the four deck listing constructors have no enrichment closure deps:

```go
Decks: deck.Functions{
    GetAll:          deck.GetAll(r.Decks, r.GameResults),
    GetAllForPlayer: deck.GetAllForPlayer(r.Decks, r.GameResults),
    GetAllByPod:     deck.GetAllByPod(r.Decks, r.Pods, r.GameResults),
    GetByID:         deck.GetByID(r.Decks, r.GameResults),
    // Create, Update, SoftDelete, Retire unchanged
    GetDeckName:        getDeckName,         // still needed by gameResult enrichment until Phase 3
    GetCommanderEntry:  getCommanderEntry,   // still needed by gameResult enrichment until Phase 3
    GetPlayerIDForDeck: getPlayerIDForDeck,  // still needed by gameResult enrichment until Phase 3
},
```

`getPlayerName` and `getFormat` closures: check whether other callers remain before
removing from wiring (e.g., player business layer may still use `getFormat`).

## Query Reduction

| Function | Before | After |
|----------|--------|-------|
| GetAll (N decks) | 4–5 per-deck | ~5 batched total |
| GetAllForPlayer (N decks) | 3–4 per-deck | ~5 batched total |
| GetAllByPod (N decks) | 4–5 per-deck | ~5 batched total |
| GetByID | 4–5 queries | ~5 queries |

The ~5 batched queries: deck fetch + deckCommander batch + commander batch + player batch
+ format batch.

## Tests

- Integration tests for each `Get*WithAll` repo method
- Verify `Commander.Commander`, `Commander.PartnerCommander`, `Player.Name`, `Format.Name`
  are populated correctly; `Commander` is nil for decks without commanders
- Business layer: remove `getCommanderEntry`, `getPlayerName`, `getFormat` mocks from
  listing function tests

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/deck/...` passes (or skips)
3. `go test ./lib/business/deck/...` passes
4. Smoke test: `GET /api/decks` and `GET /api/deck/:id` return correct player names,
   format names, and commander names
