package pod

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/m-sharp/edh-tracker/lib/utils"

	"github.com/m-sharp/edh-tracker/lib/errs"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	"github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
)

// maxInviteUses is the maximum number of times an invite code can be used.
// The pod_invite table has no max_used_count column, so this is a hardcoded limit.
const maxInviteUses = 25

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

func Create(podRepo repos.PodRepository, roleRepo repos.PlayerPodRoleRepository) CreateFunc {
	return func(ctx context.Context, name string, creatorPlayerID int) (int, error) {
		podID, err := podRepo.Add(ctx, name)
		if err != nil {
			return 0, fmt.Errorf("failed to add pod: %w", err)
		}

		if err = roleRepo.SetRole(ctx, podID, creatorPlayerID, playerPodRole.RoleManager); err != nil {
			return 0, fmt.Errorf("failed to set creator as manager: %w", err)
		}

		return podID, nil
	}
}

func AddPlayer(podRepo repos.PodRepository, roleRepo repos.PlayerPodRoleRepository) AddPlayerFunc {
	return func(ctx context.Context, podID, playerID int) error {
		if err := podRepo.AddPlayerToPod(ctx, podID, playerID); err != nil {
			return err
		}
		return roleRepo.SetRole(ctx, podID, playerID, playerPodRole.RoleMember)
	}
}

func GetRole(roleRepo repos.PlayerPodRoleRepository) GetRoleFunc {
	return func(ctx context.Context, podID, playerID int) (string, error) {
		m, err := roleRepo.GetRole(ctx, podID, playerID)
		if err != nil {
			return "", fmt.Errorf("failed to get role for player %d in pod %d: %w", playerID, podID, err)
		}
		if m == nil {
			return "", nil
		}
		return m.Role, nil
	}
}

func PromoteToManager(roleRepo repos.PlayerPodRoleRepository) PromoteToManagerFunc {
	return func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error {
		callerRole, err := roleRepo.GetRole(ctx, podID, callerPlayerID)
		if err != nil {
			return fmt.Errorf("failed to check caller role: %w", err)
		}
		if callerRole == nil || callerRole.Role != playerPodRole.RoleManager {
			return fmt.Errorf("forbidden: caller is not a manager of pod %d: %w", podID, errs.ErrForbidden)
		}

		return roleRepo.SetRole(ctx, podID, targetPlayerID, playerPodRole.RoleManager)
	}
}

func GenerateInvite(inviteRepo repos.PodInviteRepository) GenerateInviteFunc {
	return func(ctx context.Context, podID, callerPlayerID int) (string, error) {
		code := uuid.New().String()
		expiresAt := time.Now().Add(7 * 24 * time.Hour)

		if err := inviteRepo.Add(ctx, podID, callerPlayerID, code, &expiresAt); err != nil {
			return "", fmt.Errorf("failed to create invite for pod %d: %w", podID, err)
		}

		return code, nil
	}
}

func JoinByInvite(inviteRepo repos.PodInviteRepository, podRepo repos.PodRepository, roleRepo repos.PlayerPodRoleRepository) JoinByInviteFunc {
	return func(ctx context.Context, inviteCode string, playerID int) (*Entity, error) {
		invite, err := inviteRepo.GetByCode(ctx, inviteCode)
		if err != nil {
			return nil, fmt.Errorf("failed to look up invite code: %w", err)
		}
		if invite == nil {
			return nil, fmt.Errorf("invite code not found or expired")
		}
		if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("invite code has expired")
		}
		if invite.UsedCount >= maxInviteUses {
			return nil, fmt.Errorf("invite code has reached its maximum number of uses")
		}

		if err = podRepo.AddPlayerToPod(ctx, invite.PodID, playerID); err != nil {
			return nil, fmt.Errorf("failed to add player to pod: %w", err)
		}
		if err = roleRepo.SetRole(ctx, invite.PodID, playerID, playerPodRole.RoleMember); err != nil {
			return nil, fmt.Errorf("failed to set member role: %w", err)
		}
		if err = inviteRepo.IncrementUsedCount(ctx, inviteCode); err != nil {
			return nil, fmt.Errorf("failed to increment invite used count: %w", err)
		}

		m, err := podRepo.GetByID(ctx, invite.PodID)
		if err != nil {
			return nil, fmt.Errorf("failed to get pod after join: %w", err)
		}
		if m == nil {
			return nil, fmt.Errorf("pod not found after join")
		}
		e := ToEntity(*m)
		return &e, nil
	}
}

func Leave(podRepo repos.PodRepository, roleRepo repos.PlayerPodRoleRepository) LeaveFunc {
	return func(ctx context.Context, podID, playerID int) error {
		role, err := roleRepo.GetRole(ctx, podID, playerID)
		if err != nil {
			return fmt.Errorf("failed to get role: %w", err)
		}
		if role != nil && role.Role == playerPodRole.RoleManager {
			members, err := roleRepo.GetMembersWithRoles(ctx, podID)
			if err != nil {
				return fmt.Errorf("failed to get pod members: %w", err)
			}
			managerCount := 0
			for _, m := range members {
				if m.Role == playerPodRole.RoleManager {
					managerCount++
				}
			}
			if managerCount <= 1 {
				return fmt.Errorf("forbidden: cannot leave pod as the only manager; promote another member first: %w", errs.ErrForbidden)
			}
		}

		if err = podRepo.RemovePlayer(ctx, podID, playerID); err != nil {
			return fmt.Errorf("failed to remove player from pod: %w", err)
		}
		return nil
	}
}

func SoftDelete(podRepo repos.PodRepository) SoftDeleteFunc {
	return func(ctx context.Context, podID, callerPlayerID int) error {
		return podRepo.SoftDelete(ctx, podID)
	}
}

func Update(podRepo repos.PodRepository) UpdateFunc {
	return func(ctx context.Context, podID int, name string) error {
		return podRepo.Update(ctx, podID, name)
	}
}

func GetMembersWithRoles(roleRepo repos.PlayerPodRoleRepository) GetMembersWithRolesFunc {
	return func(ctx context.Context, podID int) ([]PlayerWithRole, error) {
		models, err := roleRepo.GetMembersWithRoles(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get members with roles for pod %d: %w", podID, err)
		}

		result := make([]PlayerWithRole, 0, len(models))
		for _, m := range models {
			result = append(result, PlayerWithRole{
				PlayerID: m.PlayerID,
				Role:     utils.TitleCase(m.Role),
			})
		}
		return result, nil
	}
}

func RemovePlayer(podRepo repos.PodRepository, roleRepo repos.PlayerPodRoleRepository) RemovePlayerFunc {
	return func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error {
		callerRole, err := roleRepo.GetRole(ctx, podID, callerPlayerID)
		if err != nil {
			return fmt.Errorf("failed to check caller role: %w", err)
		}
		if callerRole == nil || callerRole.Role != playerPodRole.RoleManager {
			return fmt.Errorf("forbidden: caller is not a manager of pod %d: %w", podID, errs.ErrForbidden)
		}

		return podRepo.RemovePlayer(ctx, podID, targetPlayerID)
	}
}
