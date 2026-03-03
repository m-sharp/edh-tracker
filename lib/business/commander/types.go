package commander

import "context"

type GetByIDFunc func(ctx context.Context, id int) (*Entity, error)
type CreateFunc func(ctx context.Context, name string) (int, error)
type GetCommanderNameFunc func(ctx context.Context, id int) (string, error)

type Functions struct {
	GetByID          GetByIDFunc
	Create           CreateFunc
	GetCommanderName GetCommanderNameFunc
}
