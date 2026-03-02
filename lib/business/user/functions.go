package user

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

func GetByID(userRepo repos.UserRepository) GetByIDFunc {
	return func(ctx context.Context, id int) (*Entity, error) {
		m, err := userRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user %d: %w", id, err)
		}
		if m == nil {
			return nil, nil
		}
		e := ToEntity(*m)
		return &e, nil
	}
}

func GetByPlayerID(userRepo repos.UserRepository) GetByPlayerIDFunc {
	return func(ctx context.Context, playerID int) (*Entity, error) {
		m, err := userRepo.GetByPlayerID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user for player %d: %w", playerID, err)
		}
		if m == nil {
			return nil, nil
		}
		e := ToEntity(*m)
		return &e, nil
	}
}

func Create(userRepo repos.UserRepository) CreateFunc {
	return func(ctx context.Context, playerID, roleID int) (int, error) {
		return userRepo.Add(ctx, playerID, roleID)
	}
}
