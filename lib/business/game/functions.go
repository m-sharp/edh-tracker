package game

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business/deck"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/gameresult"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	gamerepo "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameResultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

func GetAllByPod(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	gameResultRepo repos.GameResultRepository,
	getDeckName deck.GetDeckNameFunc,
	getCommanderEntry deck.GetCommanderEntryFunc,
) GetAllByPodFunc {
	return func(ctx context.Context, podID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByPod(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for pod %d: %w", podID, err)
		}

		result := make([]Entity, 0, len(games))
		for _, g := range games {
			entity, err := buildGameEntity(ctx, g, gameResultRepo, getDeckName, getCommanderEntry)
			if err != nil {
				log.Warn("Failed to build game entity, dropping from results",
					zap.Int("game_id", g.ID), zap.Error(err))
				continue
			}
			result = append(result, entity)
		}

		return result, nil
	}
}

func GetAllByDeck(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	gameResultRepo repos.GameResultRepository,
	getDeckName deck.GetDeckNameFunc,
	getCommanderEntry deck.GetCommanderEntryFunc,
) GetAllByDeckFunc {
	return func(ctx context.Context, deckID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByDeck(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for deck %d: %w", deckID, err)
		}

		result := make([]Entity, 0, len(games))
		for _, g := range games {
			entity, err := buildGameEntity(ctx, g, gameResultRepo, getDeckName, getCommanderEntry)
			if err != nil {
				log.Warn("Failed to build game entity, dropping from results",
					zap.Int("game_id", g.ID), zap.Error(err))
				continue
			}
			result = append(result, entity)
		}

		return result, nil
	}
}

func GetByID(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	gameResultRepo repos.GameResultRepository,
	getDeckName deck.GetDeckNameFunc,
	getCommanderEntry deck.GetCommanderEntryFunc,
) GetByIDFunc {
	return func(ctx context.Context, gameID int) (*Entity, error) {
		g, err := gameRepo.GetById(ctx, gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get game %d: %w", gameID, err)
		}
		if g == nil {
			return nil, nil
		}

		entity, err := buildGameEntity(ctx, *g, gameResultRepo, getDeckName, getCommanderEntry)
		if err != nil {
			return nil, err
		}

		return &entity, nil
	}
}

func Create(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	gameResultRepo repos.GameResultRepository,
	deckRepo repos.DeckRepository,
	getFormat format.GetByIDFunc,
) CreateFunc {
	return func(ctx context.Context, description string, podID, formatID int, inputs []gameresult.InputEntity) error {
		for _, input := range inputs {
			if err := input.Validate(); err != nil {
				return fmt.Errorf("invalid game result: %w", err)
			}
		}

		f, err := getFormat(ctx, formatID)
		if err != nil {
			return fmt.Errorf("failed to look up format %d: %w", formatID, err)
		}
		if f == nil {
			return fmt.Errorf("format %d not found", formatID)
		}

		if f.Name != "other" {
			for _, input := range inputs {
				d, err := deckRepo.GetById(ctx, input.DeckID)
				if err != nil {
					return fmt.Errorf("failed to look up deck %d: %w", input.DeckID, err)
				}
				if d == nil {
					return fmt.Errorf("deck %d not found", input.DeckID)
				}
				if d.FormatID != formatID {
					return fmt.Errorf("deck %d format does not match game format", input.DeckID)
				}
			}
		}

		gameID, err := gameRepo.Add(ctx, description, podID, formatID)
		if err != nil {
			return fmt.Errorf("failed to create game: %w", err)
		}

		results := make([]gameResultrepo.Model, 0, len(inputs))
		for _, input := range inputs {
			results = append(results, gameResultrepo.Model{
				GameID:    gameID,
				DeckID:    input.DeckID,
				Place:     input.Place,
				KillCount: input.Kills,
			})
		}

		if err := gameResultRepo.BulkAdd(ctx, results); err != nil {
			return fmt.Errorf("failed to create game results: %w", err)
		}

		return nil
	}
}

func buildGameEntity(
	ctx context.Context,
	g gamerepo.Model,
	gameResultRepo repos.GameResultRepository,
	getDeckName deck.GetDeckNameFunc,
	getCommanderEntry deck.GetCommanderEntryFunc,
) (Entity, error) {
	// TODO: Getting GameResults by GameID should be a Business Function under lib/business/gameresult.
	// TODO: This method should take GameResult Entities to associate with the Game Entity being returned
	resultModels, err := gameResultRepo.GetByGameId(ctx, g.ID)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to get results for game %d: %w", g.ID, err)
	}

	deckNameCache := map[int]string{}

	results := make([]gameresult.Entity, 0, len(resultModels))
	for _, r := range resultModels {
		deckName, err := cachedDeckName(ctx, r.DeckID, deckNameCache, getDeckName)
		if err != nil {
			return Entity{}, err
		}

		entity := gameresult.Entity{
			ID:       r.ID,
			GameID:   r.GameID,
			DeckID:   r.DeckID,
			DeckName: deckName,
			Place:    r.Place,
			Kills:    r.KillCount,
			Points:   gameresult.GetPointsForPlace(r.KillCount, r.Place),
		}

		commanders, err := getCommanderEntry(ctx, r.DeckID)
		if err != nil {
			return Entity{}, fmt.Errorf("failed to get commander for deck %d: %w", r.DeckID, err)
		}
		if commanders != nil {
			name := commanders.CommanderName
			entity.CommanderName = &name
			entity.PartnerCommanderName = commanders.PartnerCommanderName
		}

		results = append(results, entity)
	}

	return Entity{
		ID:          g.ID,
		Description: g.Description,
		PodID:       g.PodID,
		FormatID:    g.FormatID,
		Results:     results,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}, nil
}

// TODO: Not much savings to be found here, remove for now
func cachedDeckName(ctx context.Context, deckID int, cache map[int]string, getDeckName deck.GetDeckNameFunc) (string, error) {
	if name, ok := cache[deckID]; ok {
		return name, nil
	}
	name, err := getDeckName(ctx, deckID)
	if err != nil {
		return "", err
	}
	cache[deckID] = name
	return name, nil
}
