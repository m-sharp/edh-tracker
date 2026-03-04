package user

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	userrepo "github.com/m-sharp/edh-tracker/lib/repositories/user"
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

func GetByOAuth(userRepo repos.UserRepository) GetByOAuthFunc {
	return func(ctx context.Context, provider, subject string) (*Entity, error) {
		m, err := userRepo.GetByOAuth(ctx, provider, subject)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by oauth %s/%s: %w", provider, subject, err)
		}
		if m == nil {
			return nil, nil
		}
		e := ToEntity(*m)
		return &e, nil
	}
}

func CreateWithOAuth(userRepo repos.UserRepository) CreateWithOAuthFunc {
	return func(ctx context.Context, playerName, provider, subject, email, displayName, avatarURL string) (*Entity, error) {
		role, err := userRepo.GetRoleByName(ctx, userrepo.RolePlayer)
		if err != nil {
			return nil, fmt.Errorf("failed to get player role: %w", err)
		}

		m, err := userRepo.CreatePlayerAndUser(ctx, playerName, role.ID, provider, subject, email, displayName, avatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create player and user with oauth: %w", err)
		}

		e := ToEntity(*m)
		return &e, nil
	}
}
