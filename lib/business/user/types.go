package user

import "context"

type GetByIDFunc func(ctx context.Context, id int) (*Entity, error)
type GetByPlayerIDFunc func(ctx context.Context, playerID int) (*Entity, error)
type CreateFunc func(ctx context.Context, playerID, roleID int) (int, error)

type Functions struct {
	GetByID       GetByIDFunc
	GetByPlayerID GetByPlayerIDFunc
	Create        CreateFunc
}
