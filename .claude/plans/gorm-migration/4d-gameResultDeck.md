# Phase 4d — Preloading: gameResult.Model → deck.Model (enrichment chain)

## Status
Needs Review

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Goal

Eliminate the per-result `getDeckName`, `getPlayerIDForDeck`, and `getCommanderEntry`
queries that still occur inside `gameResult.EnrichModels` (introduced in Phase 4a).

Current pattern after Phase 4a: O(N×P) queries for enrichment (where N=games, P=players
per game) — one per result for deck name, one for player ID, one for commander.
After: 3–4 additional batched queries regardless of result count (one for decks, one for
deckCommanders, one for commanders) via nested Preload on gameResult.Model.

## Scope

- `lib/repositories/gameResult/model.go` — add `Deck` BelongsTo field
- `lib/repositories/gameResult/repo.go` — add `GetByGameIDWithDeckInfo` method
- `lib/repositories/interfaces.go` — add new `GameResultRepository` signature
- `lib/business/gameResult/functions.go` — update `EnrichModels` to read from preloaded Deck
- `lib/business/game/functions.go` — update game listing functions to use fully hydrated methods
- `lib/business/business.go` — simplify wiring (remove enrichment closure deps)

Depends on:
- Phase 4a (game → gameResult preloading, `EnrichModels` exists)
- Phase 4b (deck.Model has `Commander *deckCommander.Model` with nested commander associations)
- Phase 4c (deck.Model has `Player player.Model` for playerID/name access)

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

The `deck.Model.Commander`, `deck.Model.Player` associations from Phases 4b and 4c
are nested and will be transitively preloaded when specified.

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

### Updated game business functions

The game listing functions from Phase 4a use `gameRepo.GetAllByPodWithResults` which
preloads `game.Results []gameResult.Model`. After Phase 4d, these results also need
to have `Deck` preloaded. Options:

**Option A (preferred):** Extend the game repo `Get*WithResults` methods to also chain the
deck preload on the nested Results:

```go
func (r *Repository) GetAllByPodWithResults(ctx context.Context, podID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).
        Preload("Results.Deck.Commander.Commander").
        Preload("Results.Deck.Commander.PartnerCommander").
        Preload("Results.Deck.Player").
        Where("pod_id = ?", podID).Find(&games).Error
    // ...
}
```

This reduces `GetAllByPod` for a pod with N games and P players/game to approximately:
- 1 query for games
- 1 batched query for all results
- 1 batched query for all decks (across all results)
- 1 batched query for all deckCommanders
- 1 batched query for all commanders
Total: ~5 queries, regardless of N or P.

**Option B:** Keep game repo `Get*WithResults` unchanged; `enrichGameResults` uses
`gameResultRepo.GetByGameIDWithDeckInfo` per game (still O(N) queries for deck info).

Option A is recommended since it achieves the full N+1 elimination.

### `lib/business/business.go` wiring simplification

After Phase 4d, `enrichGameResults` needs no closure deps:

```go
enrichGameResults := gameResult.EnrichModels()

// getDeckName, getCommanderEntry, getPlayerIDForDeck may now be removable
// from wiring if no other callers remain. Check before removing.
```

The `gameResult.GetByGameID` function (used in `GameResults.GetByGameID`) still uses
the old enrichment path and can be migrated separately or kept as-is.

## Query Reduction (game listing, combined with 4a)

| Function | After 4a only | After 4a + 4d |
|----------|--------------|---------------|
| GetAllByPod (N games, P players) | 2 + N×P×(2–3) queries | ~5 queries |
| GetAllByDeck | 2 + N×P×(2–3) queries | ~5 queries |
| GetAllByPlayer | 2 + N×P×(2–3) queries | ~5 queries |
| GetByID | 2 + P×(2–3) queries | ~5 queries |

## Tests

- Integration test for `GetByGameIDWithDeckInfo`: verify Deck, Deck.Commander, and Deck.Player
  are populated; soft-deleted decks and results are excluded
- Business layer: verify `EnrichModels` builds correct entities from preloaded deck data
  (no mock closures needed — pure data transformation)

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/gameResult/...` passes (or skips)
3. `go test ./lib/business/gameResult/...` passes
4. Smoke test: game listing endpoints return correct deck names, player IDs, commander names
