package game

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	gamerepo "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameResultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

func GetAllByPod(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	getGameResults gameResult.GetByGameIDFunc,
) GetAllByPodFunc {
	return func(ctx context.Context, podID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByPod(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for pod %d: %w", podID, err)
		}

		result := make([]Entity, 0, len(games))
		for _, g := range games {
			results, err := getGameResults(ctx, g.ID)
			if err != nil {
				log.Warn("Failed to get results for game, dropping from results",
					zap.Int("game_id", g.ID), zap.Error(err))
				continue
			}
			result = append(result, buildGameEntity(g, results))
		}

		return result, nil
	}
}

func GetAllByDeck(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	getGameResults gameResult.GetByGameIDFunc,
) GetAllByDeckFunc {
	return func(ctx context.Context, deckID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByDeck(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for deck %d: %w", deckID, err)
		}

		result := make([]Entity, 0, len(games))
		for _, g := range games {
			results, err := getGameResults(ctx, g.ID)
			if err != nil {
				log.Warn("Failed to get results for game, dropping from results",
					zap.Int("game_id", g.ID), zap.Error(err))
				continue
			}
			result = append(result, buildGameEntity(g, results))
		}

		return result, nil
	}
}

func GetByID(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	getGameResults gameResult.GetByGameIDFunc,
) GetByIDFunc {
	return func(ctx context.Context, gameID int) (*Entity, error) {
		g, err := gameRepo.GetById(ctx, gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get game %d: %w", gameID, err)
		}
		if g == nil {
			return nil, nil
		}

		results, err := getGameResults(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		entity := buildGameEntity(*g, results)
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
	return func(ctx context.Context, description string, podID, formatID int, inputs []gameResult.InputEntity) error {
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

func buildGameEntity(g gamerepo.Model, results []gameResult.Entity) Entity {
	return Entity{
		ID:          g.ID,
		Description: g.Description,
		PodID:       g.PodID,
		FormatID:    g.FormatID,
		Results:     results,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}
