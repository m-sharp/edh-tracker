package player

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

func GetAll(playerRepo repos.PlayerRepository, gameResultRepo repos.GameResultRepository, podRepo repos.PodRepository) GetAllFunc {
	return func(ctx context.Context) ([]Entity, error) {
		players, err := playerRepo.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get players: %w", err)
		}

		entities := make([]Entity, 0, len(players))
		for _, p := range players {
			agg, err := gameResultRepo.GetStatsForPlayer(ctx, p.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get stats for player %d: %w", p.ID, err)
			}

			podIDs, err := podRepo.GetIDsByPlayerID(ctx, p.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get pod IDs for player %d: %w", p.ID, err)
			}

			entities = append(entities, ToEntity(p, agg, podIDs))
		}

		return entities, nil
	}
}

func GetByID(playerRepo repos.PlayerRepository, gameResultRepo repos.GameResultRepository, podRepo repos.PodRepository) GetByIDFunc {
	return func(ctx context.Context, playerID int) (*Entity, error) {
		p, err := playerRepo.GetById(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get player %d: %w", playerID, err)
		}
		if p == nil {
			return nil, nil
		}

		agg, err := gameResultRepo.GetStatsForPlayer(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for player %d: %w", playerID, err)
		}

		podIDs, err := podRepo.GetIDsByPlayerID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pod IDs for player %d: %w", playerID, err)
		}

		e := ToEntity(*p, agg, podIDs)
		return &e, nil
	}
}

func Create(playerRepo repos.PlayerRepository) CreateFunc {
	return func(ctx context.Context, name string) (int, error) {
		return playerRepo.Add(ctx, name)
	}
}

func GetPlayerName(playerRepo repos.PlayerRepository) GetPlayerNameFunc {
	return func(ctx context.Context, playerID int) (string, error) {
		p, err := playerRepo.GetById(ctx, playerID)
		if err != nil {
			return "", fmt.Errorf("failed to look up player %d: %w", playerID, err)
		}
		if p == nil {
			return "", fmt.Errorf("player %d not found", playerID)
		}
		return p.Name, nil
	}
}
