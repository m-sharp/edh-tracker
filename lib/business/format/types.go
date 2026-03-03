package format

import "context"

type GetAllFunc func(ctx context.Context) ([]Entity, error)
type GetByIDFunc func(ctx context.Context, id int) (*Entity, error)

type Functions struct {
	GetAll  GetAllFunc
	GetByID GetByIDFunc
}
