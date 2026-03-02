package player

import "context"

type GetAllFunc func(ctx context.Context) ([]Entity, error)
type GetByIDFunc func(ctx context.Context, playerID int) (*Entity, error)
type CreateFunc func(ctx context.Context, name string) (int, error)
type GetPlayerNameFunc func(ctx context.Context, playerID int) (string, error)

type Functions struct {
	GetAll        GetAllFunc
	GetByID       GetByIDFunc
	Create        CreateFunc
	GetPlayerName GetPlayerNameFunc
}
