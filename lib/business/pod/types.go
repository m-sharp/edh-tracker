package pod

import "context"

type GetByIDFunc func(ctx context.Context, podID int) (*Entity, error)
type GetByPlayerIDFunc func(ctx context.Context, playerID int) ([]Entity, error)
type CreateFunc func(ctx context.Context, name string) (int, error)
type AddPlayerFunc func(ctx context.Context, podID, playerID int) error

type Functions struct {
	GetByID       GetByIDFunc
	GetByPlayerID GetByPlayerIDFunc
	Create        CreateFunc
	AddPlayer     AddPlayerFunc
}
