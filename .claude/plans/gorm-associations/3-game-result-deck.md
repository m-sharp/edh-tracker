# Phase 3 — Preloading: gameResult.Model → deck.Model (enrichment chain)

## Skill

Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Goal

Eliminate the per-result `getDeckName`, `getPlayerIDForDeck`, and `getCommanderEntry`
queries that still occur inside `gameResult.EnrichModels` (introduced in Phase 1).

Current pattern after Phase 1: O(N×P) queries for enrichment (where N=games, P=players
per game) — one per result for deck name, one for player ID, one for commander.
After: 3–4 additional batched queries regardless of result count (one for decks, one for
deckCommanders, one for commanders) via nested Preload on `gameResult.Model`.

## Dependencies

- **Phase 1** (game → gameResult preloading, `EnrichModels` constructor exists)
- **Phase 2** (deck.Model has `Commander *deckCommander.Model` with nested commander
  associations; `deck.Model` has `Player player.Model` for playerID/name access)

Both must be complete before implementing this phase.

## Scope

- `lib/repositories/gameResult/model.go` — add `Deck` BelongsTo field
- `lib/repositories/gameResult/repo.go` — add `GetByGameIDWithDeckInfo` method
- `lib/repositories/interfaces.go` — add new `GameResultRepository` signature
- `lib/repositories/game/repo.go` — extend all 4 `Get*WithResults` methods to chain deck preloads
- `lib/business/gameResult/functions.go` — update `EnrichModels` to take no closure deps
- `lib/business/business.go` — simplify wiring

## Association Declaration

Add to `lib/repositories/gameResult/model.go`:

```go
import deckRepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"

type Model struct {
    base.GormModelBase
    GameID    int
    DeckID    int
    Place     int
    KillCount int
    // Populated only when using GetByGameIDWithDeckInfo.
    Deck deckRepo.Model // BelongsTo via DeckID
}

func (Model) TableName() string { return "game_result" }
```

GORM BelongsTo convention: `Deck` + `DeckID` → `deck.id = game_result.deck_id`.

The `deck.Model.Commander`, `deck.Model.Player` associations from Phase 2 are nested and
will be transitively preloaded when specified.

## New Repository Method

Add to `lib/repositories/gameResult/repo.go`:

```go
func (r *Repository) GetByGameIDWithDeckInfo(ctx context.Context, gameID int) ([]Model, error) {
    var results []Model
    err := r.db.WithContext(ctx).
        Preload("Deck.Commander.Commander").
        Preload("Deck.Commander.PartnerCommander").
        Preload("Deck.Player").
        Where("game_id = ?", gameID).
        Find(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get GameResults with deck info for game %d: %w", gameID, err)
    }
    return results, nil
}
```

This issues ~4 queries regardless of result count:
1. `SELECT * FROM game_result WHERE game_id = ?`
2. `SELECT * FROM deck WHERE id IN (...)`
3. `SELECT * FROM deck_commander WHERE deck_id IN (...)`
4. `SELECT * FROM commander WHERE id IN (...)` (covers both commander + partner in one batch)

## Interface Update

Add to `GameResultRepository` in `lib/repositories/interfaces.go`:

```go
GetByGameIDWithDeckInfo(ctx context.Context, gameID int) ([]gameResult.Model, error)
```

## Repository Changes — Extend `Get*WithResults` (Option A)

Extend all four `Get*WithResults` methods in `lib/repositories/game/repo.go` to chain
the deck preload on the nested Results. This achieves ~5 total queries for game listing
regardless of scale:

```go
func (r *Repository) GetAllByPodWithResults(ctx context.Context, podID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).
        Preload("Results.Deck.Commander.Commander").
        Preload("Results.Deck.Commander.PartnerCommander").
        Preload("Results.Deck.Player").
        Where("pod_id = ?", podID).Find(&games).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Games with results for pod %d: %w", podID, err)
    }
    return games, nil
}
```

Apply the same three-line `Preload` chain to `GetAllByDeckWithResults`,
`GetAllByPlayerWithResults`, and `GetByIDWithResults`.

## Business Layer Changes

### Updated `EnrichModels` in `lib/business/gameResult/functions.go`

With deck data preloaded, `EnrichModels` no longer needs closure deps for DB lookups.
It reads directly from the preloaded `model.Deck` field:

```go
func EnrichModels() EnrichModelsFunc {
    return func(ctx context.Context, models []gameResultRepo.Model) ([]Entity, error) {
        numPlayers := len(models)
        results := make([]Entity, 0, len(models))
        for _, r := range models {
            entity := Entity{
                ID:       r.ID,
                GameID:   r.GameID,
                DeckID:   r.DeckID,
                PlayerID: r.Deck.PlayerID,
                DeckName: r.Deck.Name,
                Place:    r.Place,
                Kills:    r.KillCount,
                Points:   utils.GetPointsForPlace(r.KillCount, r.Place, numPlayers),
            }
            if r.Deck.Commander != nil {
                name := r.Deck.Commander.Commander.Name
                entity.CommanderName = &name
                if r.Deck.Commander.PartnerCommander != nil {
                    partnerName := r.Deck.Commander.PartnerCommander.Name
                    entity.PartnerCommanderName = &partnerName
                }
            }
            results = append(results, entity)
        }
        return results, nil
    }
}
```

The closure constructor now takes no parameters. The type signature of `EnrichModelsFunc`
is unchanged: `func(ctx context.Context, models []gameResultRepo.Model) ([]Entity, error)`.

### `lib/business/business.go` wiring simplification

```go
enrichGameResults := gameResult.EnrichModels()

// game.Functions unchanged structurally; repos already use extended Get*WithResults
GetAllByPod:    game.GetAllByPod(log, r.Games, enrichGameResults),
GetAllByDeck:   game.GetAllByDeck(log, r.Games, enrichGameResults),
GetAllByPlayer: game.GetAllByPlayer(log, r.Games, enrichGameResults),
GetByID:        game.GetByID(log, r.Games, enrichGameResults),
```

After this phase, check whether `getDeckName`, `getCommanderEntry`, and `getPlayerIDForDeck`
have any remaining callers. If `gameResult.GetByGameID` is the only user, it can be migrated
to use `GetByGameIDWithDeckInfo` + the dep-free `EnrichModels()`, then those three closures
can be removed from wiring entirely.

## Query Reduction (combined with Phase 1)

| Function | After Phase 1 only | After Phases 1 + 3 |
|----------|-------------------|-------------------|
| GetAllByPod (N games, P players) | 2 + N×P×(2–3) | ~5 queries |
| GetAllByDeck | 2 + N×P×(2–3) | ~5 queries |
| GetAllByPlayer | 2 + N×P×(2–3) | ~5 queries |
| GetByID | 2 + P×(2–3) | ~5 queries |

## Tests

- Integration test for `GetByGameIDWithDeckInfo`: verify `Deck`, `Deck.Commander`, and
  `Deck.Player` are populated; soft-deleted decks and results are excluded
- Integration tests for updated `Get*WithResults` methods: verify `Results[*].Deck` is
  populated with full commander and player data
- Business layer: verify `EnrichModels` builds correct entities from preloaded deck data
  (no mock closures needed — pure data transformation)

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/gameResult/...` passes (or skips)
3. `go test ./lib/business/gameResult/...` passes
4. Smoke test: game listing endpoints return correct deck names, player IDs, commander names
