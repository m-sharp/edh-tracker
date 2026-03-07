package pod

import "context"

type GetByIDFunc func(ctx context.Context, podID int) (*Entity, error)
type GetByPlayerIDFunc func(ctx context.Context, playerID int) ([]Entity, error)
type CreateFunc func(ctx context.Context, name string, creatorPlayerID int) (int, error)
type AddPlayerFunc func(ctx context.Context, podID, playerID int) error
type GetRoleFunc func(ctx context.Context, podID, playerID int) (string, error)
type PromoteToManagerFunc func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error
type GenerateInviteFunc func(ctx context.Context, podID, callerPlayerID int) (string, error)
type JoinByInviteFunc func(ctx context.Context, inviteCode string, playerID int) (*Entity, error)
type LeaveFunc func(ctx context.Context, podID, playerID int) error
type SoftDeleteFunc func(ctx context.Context, podID, callerPlayerID int) error
type UpdateFunc func(ctx context.Context, podID int, name string) error
type GetMembersWithRolesFunc func(ctx context.Context, podID int) ([]PlayerWithRole, error)
type RemovePlayerFunc func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error

type Functions struct {
	GetByID             GetByIDFunc
	GetByPlayerID       GetByPlayerIDFunc
	Create              CreateFunc
	AddPlayer           AddPlayerFunc
	GetRole             GetRoleFunc
	PromoteToManager    PromoteToManagerFunc
	GenerateInvite      GenerateInviteFunc
	JoinByInvite        JoinByInviteFunc
	Leave               LeaveFunc
	SoftDelete          SoftDeleteFunc
	Update              UpdateFunc
	GetMembersWithRoles GetMembersWithRolesFunc
	RemovePlayer        RemovePlayerFunc
}
