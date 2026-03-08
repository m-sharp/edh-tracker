# Phase 4a — Preloading: game.Model → []gameResult.Model

## Status
Needs Review

## Skill
Load `.claude/skills/gorm.md` at the start of each implementation session for this phase.

## Goal

Eliminate the per-game DB round-trip for `game_result` records that currently occurs
inside `game.GetAllByPod`, `game.GetAllByDeck`, and `game.GetAllByPlayer`.

Current pattern: 1 query for N games, then 1 query per game for its results = 1+N queries.
After: 2 queries total (games + all results batched via GORM Preload).

Per-result enrichment (deck name, player, commander) is unchanged in this phase;
that is addressed in Phase 4d.

## Scope

- `lib/repositories/game/model.go` — add `Results` association field
- `lib/repositories/game/repo.go` — add four `Get*WithResults` methods
- `lib/repositories/interfaces.go` — add new `GameRepository` method signatures
- `lib/business/gameResult/functions.go` — add new `EnrichModels` constructor
- `lib/business/gameResult/types.go` — add `EnrichModelsFunc` type
- `lib/business/game/functions.go` — replace `getGameResults` dep with `enrichGameResults`
- `lib/business/business.go` — update `NewBusiness` wiring

Depends on: Phase 3a (game repo GORM-migrated), Phase 3b (gameResult repo GORM-migrated)

## Association Declaration

Add to `lib/repositories/game/model.go`:

```go
import gameResultRepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"

type Model struct {
    base.GormModelBase
    Description string
    PodID       int
    FormatID    int
    // Populated only when using Get*WithResults repo methods.
    Results []gameResultRepo.Model `gorm:"foreignKey:GameID"`
}

func (Model) TableName() string { return "game" }
```

GORM automatically filters soft-deleted `game_result` rows (`deleted_at IS NULL`) when
preloading because `GormModelBase` embeds `gorm.DeletedAt`.

## New Repository Methods

Add to `lib/repositories/game/repo.go` alongside the existing flat methods (do not replace):

```go
func (r *Repository) GetAllByPodWithResults(ctx context.Context, podID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).Preload("Results").Where("pod_id = ?", podID).Find(&games).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Games with results for pod %d: %w", podID, err)
    }
    return games, nil
}

func (r *Repository) GetAllByDeckWithResults(ctx context.Context, deckID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).
        Preload("Results").
        Joins("INNER JOIN game_result ON game.id = game_result.game_id").
        Where("game_result.deck_id = ? AND game_result.deleted_at IS NULL", deckID).
        Find(&games).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Games with results for deck %d: %w", deckID, err)
    }
    return games, nil
}

func (r *Repository) GetAllByPlayerWithResults(ctx context.Context, playerID int) ([]Model, error) {
    var games []Model
    err := r.db.WithContext(ctx).
        Preload("Results").
        Distinct("game.id, game.description, game.pod_id, game.format_id, game.created_at, game.updated_at, game.deleted_at").
        Joins("INNER JOIN game_result ON game.id = game_result.game_id").
        Joins("INNER JOIN deck ON game_result.deck_id = deck.id").
        Where("deck.player_id = ? AND game_result.deleted_at IS NULL AND deck.deleted_at IS NULL", playerID).
        Find(&games).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get Games with results for player %d: %w", playerID, err)
    }
    return games, nil
}

func (r *Repository) GetByIDWithResults(ctx context.Context, gameID int) (*Model, error) {
    var m Model
    err := r.db.WithContext(ctx).Preload("Results").First(&m, gameID).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get Game with results for id %d: %w", gameID, err)
    }
    return &m, nil
}
```

Note on Preload + Joins: the `Joins` clause filters which games are returned (games that have
a result for the given deck/player). The `Preload("Results")` then fetches ALL non-deleted
results for the returned games in a second batched query — correct behavior for game entities.

## Interface Updates

Add to the `GameRepository` interface in `lib/repositories/interfaces.go`:

```go
GetAllByPodWithResults(ctx context.Context, podID int) ([]game.Model, error)
GetAllByDeckWithResults(ctx context.Context, deckID int) ([]game.Model, error)
GetAllByPlayerWithResults(ctx context.Context, playerID int) ([]game.Model, error)
GetByIDWithResults(ctx context.Context, gameID int) (*game.Model, error)
```

## Business Layer Changes

### 1. New `EnrichModels` in `lib/business/gameResult/functions.go`

Extracts the inner enrichment loop from `GetByGameID` into a standalone function that
operates on already-fetched `[]gameResultRepo.Model`:

```go
func EnrichModels(
    getDeckName deck.GetDeckNameFunc,
    getCommanderEntry deck.GetCommanderEntryFunc,
    getPlayerIDForDeck deck.GetPlayerIDForDeckFunc,
) EnrichModelsFunc {
    return func(ctx context.Context, models []gameResultRepo.Model) ([]Entity, error) {
        numPlayers := len(models)
        deckNameCache := map[int]string{}
        playerIDCache := map[int]int{}
        results := make([]Entity, 0, len(models))
        for _, r := range models {
            // Same enrichment logic as GetByGameID inner loop:
            // - resolve deck name (cached)
            // - resolve playerID (cached)
            // - resolve commander entry
            // - compute points via utils.GetPointsForPlace
        }
        return results, nil
    }
}
```

The same deps (`getDeckName`, `getCommanderEntry`, `getPlayerIDForDeck`) are kept here
until Phase 4d replaces them with preloaded deck data.

### 2. New type in `lib/business/gameResult/types.go`

```go
type EnrichModelsFunc func(ctx context.Context, models []gameResultRepo.Model) ([]Entity, error)
```

Add to `gameResult.Functions` struct:
```go
EnrichModels EnrichModelsFunc
```

### 3. Updated game business constructors in `lib/business/game/functions.go`

Replace `getGameResults gameResult.GetByGameIDFunc` with
`enrichGameResults gameResult.EnrichModelsFunc` in `GetAllByPod`, `GetAllByDeck`,
`GetAllByPlayer`, and `GetByID`:

```go
func GetAllByPod(
    log *zap.Logger,
    gameRepo repos.GameRepository,
    enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByPodFunc {
    return func(ctx context.Context, podID int) ([]Entity, error) {
        games, err := gameRepo.GetAllByPodWithResults(ctx, podID)
        // ...
        for _, g := range games {
            results, err := enrichGameResults(ctx, g.Results)
            if err != nil {
                log.Warn("Failed to enrich results for game", ...)
                continue
            }
            result = append(result, buildGameEntity(g, results))
        }
        return result, nil
    }
}
```

Apply the same pattern to `GetAllByDeck`, `GetAllByPlayer`, and `GetByID`.

### 4. `lib/business/business.go` wiring

```go
enrichGameResults := gameResult.EnrichModels(getDeckName, getCommanderEntry, getPlayerIDForDeck)

// game.Functions updated:
GetAllByPod:    game.GetAllByPod(log, r.Games, enrichGameResults),
GetAllByDeck:   game.GetAllByDeck(log, r.Games, enrichGameResults),
GetAllByPlayer: game.GetAllByPlayer(log, r.Games, enrichGameResults),
GetByID:        game.GetByID(log, r.Games, enrichGameResults),

// gameResult.Functions gains:
EnrichModels: enrichGameResults,

// getGameResults still wired for GameResults.GetByGameID:
getGameResults := gameResult.GetByGameID(r.GameResults, getDeckName, getCommanderEntry, getPlayerIDForDeck)
GameResults: gameResult.Functions{
    GetByGameID:        getGameResults,
    GetGameIDForResult: gameResult.GetGameIDForResult(r.GameResults),
    EnrichModels:       enrichGameResults,
},
```

## Query Reduction

| Function | Before | After |
|----------|--------|-------|
| GetAllByPod (N games) | 1 + N queries | 2 queries |
| GetAllByDeck (N games) | 1 + N queries | 2 queries |
| GetAllByPlayer (N games) | 1 + N queries | 2 queries |
| GetByID | 2 queries | 2 queries (unchanged) |

Per-result enrichment (deck name, commander, player ID) unchanged — addressed in Phase 4d.

## Tests

- Integration tests for each `Get*WithResults` repo method
- Verify preloaded `Results` is populated and soft-deleted results are excluded
- Business layer tests: update `GetAllByPod` etc. to inject mock `EnrichModelsFunc`
  instead of `GetByGameIDFunc`

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/game/...` passes (or skips)
3. `go test ./lib/business/game/...` passes
4. Smoke test: game listing endpoints return correct data with all results populated
