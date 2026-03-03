package commander

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

func GetByID(commanderRepo repos.CommanderRepository) GetByIDFunc {
	return func(ctx context.Context, id int) (*Entity, error) {
		m, err := commanderRepo.GetById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get commander %d: %w", id, err)
		}
		if m == nil {
			return nil, nil
		}
		e := ToEntity(*m)
		return &e, nil
	}
}

func Create(commanderRepo repos.CommanderRepository) CreateFunc {
	return func(ctx context.Context, name string) (int, error) {
		return commanderRepo.Add(ctx, name)
	}
}

func GetCommanderName(commanderRepo repos.CommanderRepository) GetCommanderNameFunc {
	return func(ctx context.Context, id int) (string, error) {
		m, err := commanderRepo.GetById(ctx, id)
		if err != nil {
			return "", fmt.Errorf("failed to look up commander %d: %w", id, err)
		}
		if m == nil {
			return "", fmt.Errorf("commander %d not found", id)
		}
		return m.Name, nil
	}
}
