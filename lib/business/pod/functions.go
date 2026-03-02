package pod

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

func GetByID(podRepo repos.PodRepository) GetByIDFunc {
	return func(ctx context.Context, podID int) (*Entity, error) {
		m, err := podRepo.GetByID(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pod %d: %w", podID, err)
		}
		if m == nil {
			return nil, nil
		}
		e := ToEntity(*m)
		return &e, nil
	}
}

func GetByPlayerID(podRepo repos.PodRepository) GetByPlayerIDFunc {
	return func(ctx context.Context, playerID int) ([]Entity, error) {
		models, err := podRepo.GetByPlayerID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pods for player %d: %w", playerID, err)
		}

		entities := make([]Entity, 0, len(models))
		for _, m := range models {
			entities = append(entities, ToEntity(m))
		}

		return entities, nil
	}
}

func Create(podRepo repos.PodRepository) CreateFunc {
	return func(ctx context.Context, name string) (int, error) {
		return podRepo.Add(ctx, name)
	}
}

func AddPlayer(podRepo repos.PodRepository) AddPlayerFunc {
	return func(ctx context.Context, podID, playerID int) error {
		return podRepo.AddPlayerToPod(ctx, podID, playerID)
	}
}
