package player

import "context"

type GetAllFunc func(ctx context.Context) ([]Entity, error)
type GetAllByPodFunc func(ctx context.Context, podID int) ([]PlayerWithRoleEntity, error)
type GetByIDFunc func(ctx context.Context, playerID int) (*Entity, error)
type CreateFunc func(ctx context.Context, name string) (int, error)
type UpdateFunc func(ctx context.Context, playerID int, name string) error
type GetPlayerNameFunc func(ctx context.Context, playerID int) (string, error)

type Functions struct {
	GetAll        GetAllFunc
	GetAllByPod   GetAllByPodFunc
	GetByID       GetByIDFunc
	Create        CreateFunc
	Update        UpdateFunc
	GetPlayerName GetPlayerNameFunc
}
