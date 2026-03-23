package deck

import "context"

// UpdateFields holds optional fields that may be updated on a deck.
// CommanderID and PartnerCommanderID, when non-nil, replace the existing deck_commander rows.
type UpdateFields struct {
	Name               *string
	FormatID           *int
	CommanderID        *int
	PartnerCommanderID *int
	Retired            *bool
}

type GetAllForPlayerFunc func(ctx context.Context, playerID int) ([]EntityWithStats, error)
type GetAllByPodFunc func(ctx context.Context, podID int) ([]EntityWithStats, error)
type GetAllByPodPaginatedFunc func(ctx context.Context, podID, limit, offset int) ([]EntityWithStats, int, error)
type GetAllByPlayerPaginatedFunc func(ctx context.Context, playerID, limit, offset int) ([]EntityWithStats, int, error)
type GetByIDFunc func(ctx context.Context, deckID int) (*EntityWithStats, error)
type CreateFunc func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error)
type UpdateFunc func(ctx context.Context, deckID int, fields UpdateFields) error
type SoftDeleteFunc func(ctx context.Context, deckID int) error
type RetireFunc func(ctx context.Context, deckID int) error
type GetDeckNameFunc func(ctx context.Context, deckID int) (string, error)

type Functions struct {
	GetAllForPlayer         GetAllForPlayerFunc
	GetAllByPod             GetAllByPodFunc
	GetAllByPodPaginated    GetAllByPodPaginatedFunc
	GetAllByPlayerPaginated GetAllByPlayerPaginatedFunc
	GetByID                 GetByIDFunc
	Create                  CreateFunc
	Update                  UpdateFunc
	SoftDelete              SoftDeleteFunc
	Retire                  RetireFunc
	GetDeckName             GetDeckNameFunc
}
