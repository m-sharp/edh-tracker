package deck

import "context"

type GetAllFunc func(ctx context.Context) ([]EntityWithStats, error)
type GetAllForPlayerFunc func(ctx context.Context, playerID int) ([]EntityWithStats, error)
type GetByIDFunc func(ctx context.Context, deckID int) (*EntityWithStats, error)
type CreateFunc func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error)
type RetireFunc func(ctx context.Context, deckID int) error
type GetDeckNameFunc func(ctx context.Context, deckID int) (string, error)
type GetCommanderEntryFunc func(ctx context.Context, deckID int) (*CommanderInfo, error)

type Functions struct {
	GetAll            GetAllFunc
	GetAllForPlayer   GetAllForPlayerFunc
	GetByID           GetByIDFunc
	Create            CreateFunc
	Retire            RetireFunc
	GetDeckName       GetDeckNameFunc
	GetCommanderEntry GetCommanderEntryFunc
}
