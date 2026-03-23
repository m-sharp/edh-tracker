package game

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	gameRepository "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameResultRepository "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

func GetAllByPod(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByPodFunc {
	return func(ctx context.Context, podID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByPod(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for pod %d: %w", podID, err)
		}
		return enrichGameModels(ctx, log, games, enrichGameResults), nil
	}
}

func GetAllByPodPaginated(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByPodPaginatedFunc {
	return func(ctx context.Context, podID, limit, offset int) ([]Entity, int, error) {
		games, total, err := gameRepo.GetAllByPodPaginated(ctx, podID, limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get games for pod %d: %w", podID, err)
		}
		return enrichGameModels(ctx, log, games, enrichGameResults), total, nil
	}
}

func GetAllByDeck(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByDeckFunc {
	return func(ctx context.Context, deckID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByDeck(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for deck %d: %w", deckID, err)
		}
		return enrichGameModels(ctx, log, games, enrichGameResults), nil
	}
}

func GetAllByDeckPaginated(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByDeckPaginatedFunc {
	return func(ctx context.Context, deckID, limit, offset int) ([]Entity, int, error) {
		games, total, err := gameRepo.GetAllByDeckPaginated(ctx, deckID, limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get games for deck %d: %w", deckID, err)
		}
		return enrichGameModels(ctx, log, games, enrichGameResults), total, nil
	}
}

func GetByID(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetByIDFunc {
	return func(ctx context.Context, gameID int) (*Entity, error) {
		g, err := gameRepo.GetByID(ctx, gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get game %d: %w", gameID, err)
		}
		if g == nil {
			return nil, nil
		}

		results, err := enrichGameResults(ctx, g.Results)
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
	client *lib.DBClient,
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

		err = client.GormDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			txGameRepo := gameRepository.NewRepository(&lib.DBClient{GormDb: tx})
			txGameResultRepo := gameResultRepository.NewRepository(&lib.DBClient{GormDb: tx})

			gameID, err := txGameRepo.Add(ctx, description, podID, formatID)
			if err != nil {
				return fmt.Errorf("failed to create game: %w", err)
			}

			results := make([]gameResultRepository.Model, 0, len(inputs))
			for _, input := range inputs {
				results = append(results, gameResultRepository.Model{
					GameID:    gameID,
					DeckID:    input.DeckID,
					Place:     input.Place,
					KillCount: input.Kills,
				})
			}

			if err := txGameResultRepo.BulkAdd(ctx, results); err != nil {
				return fmt.Errorf("failed to create game results: %w", err)
			}

			return nil
		})
		return err
	}
}

func GetAllByPlayer(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByPlayerFunc {
	return func(ctx context.Context, playerID int) ([]Entity, error) {
		games, err := gameRepo.GetAllByPlayerID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get games for player %d: %w", playerID, err)
		}
		return enrichGameModels(ctx, log, games, enrichGameResults), nil
	}
}

func GetAllByPlayerIDPaginated(
	log *zap.Logger,
	gameRepo repos.GameRepository,
	enrichGameResults gameResult.EnrichModelsFunc,
) GetAllByPlayerIDPaginatedFunc {
	return func(ctx context.Context, playerID, limit, offset int) ([]Entity, int, error) {
		games, total, err := gameRepo.GetAllByPlayerIDPaginated(ctx, playerID, limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get games for player %d: %w", playerID, err)
		}
		return enrichGameModels(ctx, log, games, enrichGameResults), total, nil
	}
}

func Update(gameRepo repos.GameRepository) UpdateFunc {
	return func(ctx context.Context, gameID int, description string) error {
		return gameRepo.Update(ctx, gameID, description)
	}
}

func SoftDelete(gameRepo repos.GameRepository) SoftDeleteFunc {
	return func(ctx context.Context, gameID int) error {
		return gameRepo.SoftDelete(ctx, gameID)
	}
}

func AddResult(
	gameResultRepo repos.GameResultRepository,
) AddResultFunc {
	return func(ctx context.Context, gameID, deckID, playerID, place, killCount int) (int, error) {
		return gameResultRepo.Add(ctx, gameResultRepository.Model{
			GameID:    gameID,
			DeckID:    deckID,
			Place:     place,
			KillCount: killCount,
		})
	}
}

func UpdateResult(gameResultRepo repos.GameResultRepository) UpdateResultFunc {
	return func(ctx context.Context, resultID, place, killCount, deckID int) error {
		return gameResultRepo.Update(ctx, resultID, place, killCount, deckID)
	}
}

func DeleteResult(gameResultRepo repos.GameResultRepository) DeleteResultFunc {
	return func(ctx context.Context, resultID int) error {
		return gameResultRepo.SoftDelete(ctx, resultID)
	}
}

// enrichGameModels converts []gameRepository.Model → []Entity.
// Games whose result enrichment fails are silently dropped (warning logged).
func enrichGameModels(
	ctx context.Context,
	log *zap.Logger,
	games []gameRepository.Model,
	enrich gameResult.EnrichModelsFunc,
) []Entity {
	result := make([]Entity, 0, len(games))
	for _, g := range games {
		results, err := enrich(ctx, g.Results)
		if err != nil {
			log.Warn("Failed to get results for game, dropping from results",
				zap.Int("game_id", g.ID), zap.Error(err))
			continue
		}
		result = append(result, buildGameEntity(g, results))
	}
	return result
}

func buildGameEntity(g gameRepository.Model, results []gameResult.Entity) Entity {
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
