# Phase 4c — Preloading: deck.Model → player.Model + format.Model

## Status
Needs Review

## Skill
Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Goal

Eliminate per-deck `getPlayerName` and `getFormat` queries in deck listing functions
(`GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID`).

Current pattern: `getFormat(ctx, d.FormatID)` called per deck (1 query each); `getPlayerName`
also called per deck with a local cache that still issues 1 query per unique player.
After: 2 additional batched queries (one for all players referenced by the returned decks,
one for all formats).

## Scope

- `lib/repositories/deck/model.go` — add `Player` and `Format` BelongsTo fields
- `lib/repositories/deck/repo.go` — extend `Get*WithCommanders` helper to also Preload Player + Format
  (or rename to `Get*WithAll` if preferred for clarity)
- `lib/repositories/interfaces.go` — rename/add new signature if method names change
- `lib/business/deck/functions.go` — drop `getPlayerName` + `getFormat` deps from listing functions

Depends on:
- Phase 1a (player repo GORM-migrated)
- Phase 1b (format repo GORM-migrated)
- Phase 2a (deck repo GORM-migrated)
- Phase 4b is NOT a hard prerequisite, but implementing 4c after 4b avoids touching the same
  methods twice. Recommended order: 4b then 4c.

## Association Declarations

Add to `lib/repositories/deck/model.go`:

```go
import (
    formatRepo "github.com/m-sharp/edh-tracker/lib/repositories/format"
    playerRepo "github.com/m-sharp/edh-tracker/lib/repositories/player"
    deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
)

type Model struct {
    base.GormModelBase
    PlayerID int
    Name     string
    FormatID int
    Retired  bool
    // Association fields (populated only when using Get*WithAll methods)
    Commander *deckCommanderRepo.Model `gorm:"foreignKey:DeckID"`
    Player    playerRepo.Model         // BelongsTo via PlayerID
    Format    formatRepo.Model         // BelongsTo via FormatID
}

func (Model) TableName() string { return "deck" }
```

GORM BelongsTo convention: `Player` field + `PlayerID` FK column → GORM resolves
`player.id = deck.player_id`. Same for `Format` / `FormatID`.

## Repository Method Changes

If implementing after Phase 4b, extend (or rename) the `preloadCommanders` helper in
`lib/repositories/deck/repo.go` to chain Player and Format preloads:

```go
func (r *Repository) preloadAll(db *gorm.DB) *gorm.DB {
    return db.
        Preload("Commander.Commander").
        Preload("Commander.PartnerCommander").
        Preload("Player").
        Preload("Format")
}
```

Rename existing `Get*WithCommanders` methods to `Get*WithAll` (or keep both):

```go
func (r *Repository) GetAllWithAll(ctx context.Context) ([]Model, error) { ... }
func (r *Repository) GetAllForPlayerWithAll(ctx context.Context, playerID int) ([]Model, error) { ... }
func (r *Repository) GetAllByPlayerIDsWithAll(ctx context.Context, playerIDs []int) ([]Model, error) { ... }
func (r *Repository) GetByIDWithAll(ctx context.Context, deckID int) (*Model, error) { ... }
```

If NOT implementing after 4b, use `preloadAll` without the Commander chain and add
commander preloads later.

## Interface Updates

Replace (or add alongside) `Get*WithCommanders` in `DeckRepository`:

```go
GetAllWithAll(ctx context.Context) ([]deck.Model, error)
GetAllForPlayerWithAll(ctx context.Context, playerID int) ([]deck.Model, error)
GetAllByPlayerIDsWithAll(ctx context.Context, playerIDs []int) ([]deck.Model, error)
GetByIDWithAll(ctx context.Context, deckID int) (*deck.Model, error)
```

## Business Layer Changes

### Updated listing functions

`GetAll`, `GetAllForPlayer`, `GetAllByPod`, `GetByID` drop `getPlayerName` and
`getFormat` parameters. They use `d.Player.Name` and `d.Format.Name` from preloaded fields:

```go
func GetAll(
    deckRepo repos.DeckRepository,
    gameResultRepo repos.GameResultRepository,
    // getPlayerName removed
    // getFormat removed
    // getCommanderEntry removed (from Phase 4b)
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

The local `playerCache` map in `GetAll` and `GetAllByPod` is removed — GORM's batch
preload replaces it.

### `lib/business/business.go` wiring update

After 4b + 4c, the four deck listing constructors have no enrichment closure deps:

```go
Decks: deck.Functions{
    GetAll:          deck.GetAll(r.Decks, r.GameResults),
    GetAllForPlayer: deck.GetAllForPlayer(r.Decks, r.GameResults),
    GetAllByPod:     deck.GetAllByPod(r.Decks, r.Pods, r.GameResults),
    GetByID:         deck.GetByID(r.Decks, r.GameResults),
    // ... Create, Update, SoftDelete, Retire unchanged
    GetDeckName:        getDeckName,         // still needed by gameResult enrichment
    GetCommanderEntry:  getCommanderEntry,   // still needed by gameResult enrichment (until 4d)
    GetPlayerIDForDeck: getPlayerIDForDeck,  // still needed by gameResult enrichment (until 4d)
},
```

`getPlayerName` and `getFormat` closures may still be needed by other callers
(e.g., player business layer) — check before removing from wiring.

## Query Reduction (deck listing, combined with 4b)

| Function | Before 4b+4c | After 4b+4c |
|----------|-------------|-------------|
| GetAll (N decks) | 2–3 per-deck + 2 per-deck = 4–5 per-deck | ~4 batched total |
| GetAllForPlayer (N decks) | 3–4 per-deck | ~4 batched total |
| GetAllByPod (N decks) | 4–5 per-deck | ~4 batched total |
| GetByID | 4–5 queries | ~4 queries |

The ~4 batched queries: deck fetch + deckCommander batch + commander batch + player batch
+ format batch. GORM may combine some depending on the query structure.

## Tests

- Integration tests for each `Get*WithAll` repo method
- Verify `Player.Name` and `Format.Name` are populated on returned models
- Business layer: remove `getPlayerName` and `getFormat` mocks from listing function tests

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/deck/...` passes (or skips)
3. `go test ./lib/business/deck/...` passes
4. Smoke test: `GET /api/decks` returns correct player names and format names
