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

func GetAllByPod(
	playerRepo repos.PlayerRepository,
	gameResultRepo repos.GameResultRepository,
	podRepo repos.PodRepository,
	roleRepo repos.PlayerPodRoleRepository,
) GetAllByPodFunc {
	return func(ctx context.Context, podID int) ([]PlayerWithRoleEntity, error) {
		members, err := roleRepo.GetMembersWithRoles(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get members for pod %d: %w", podID, err)
		}

		roleByPlayerID := make(map[int]string, len(members))
		for _, m := range members {
			roleByPlayerID[m.PlayerID] = m.Role
		}

		result := make([]PlayerWithRoleEntity, 0, len(members))
		for _, m := range members {
			p, err := playerRepo.GetById(ctx, m.PlayerID)
			if err != nil {
				return nil, fmt.Errorf("failed to get player %d: %w", m.PlayerID, err)
			}
			if p == nil {
				continue
			}

			agg, err := gameResultRepo.GetStatsForPlayer(ctx, m.PlayerID)
			if err != nil {
				return nil, fmt.Errorf("failed to get stats for player %d: %w", m.PlayerID, err)
			}

			podIDs, err := podRepo.GetIDsByPlayerID(ctx, m.PlayerID)
			if err != nil {
				return nil, fmt.Errorf("failed to get pod IDs for player %d: %w", m.PlayerID, err)
			}

			result = append(result, PlayerWithRoleEntity{
				Entity: ToEntity(*p, agg, podIDs),
				Role:   roleByPlayerID[m.PlayerID],
			})
		}

		return result, nil
	}
}

func Update(playerRepo repos.PlayerRepository) UpdateFunc {
	return func(ctx context.Context, playerID int, name string) error {
		return playerRepo.Update(ctx, playerID, name)
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
