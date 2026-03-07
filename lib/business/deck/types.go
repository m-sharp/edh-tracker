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

type GetAllFunc func(ctx context.Context) ([]EntityWithStats, error)
type GetAllForPlayerFunc func(ctx context.Context, playerID int) ([]EntityWithStats, error)
type GetAllByPodFunc func(ctx context.Context, podID int) ([]EntityWithStats, error)
type GetByIDFunc func(ctx context.Context, deckID int) (*EntityWithStats, error)
type CreateFunc func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error)
type UpdateFunc func(ctx context.Context, deckID int, callerPlayerID int, fields UpdateFields) error
type SoftDeleteFunc func(ctx context.Context, deckID int, callerPlayerID int) error
type RetireFunc func(ctx context.Context, deckID int) error
type GetDeckNameFunc func(ctx context.Context, deckID int) (string, error)
type GetCommanderEntryFunc func(ctx context.Context, deckID int) (*CommanderInfo, error)

type Functions struct {
	GetAll            GetAllFunc
	GetAllForPlayer   GetAllForPlayerFunc
	GetAllByPod       GetAllByPodFunc
	GetByID           GetByIDFunc
	Create            CreateFunc
	Update            UpdateFunc
	SoftDelete        SoftDeleteFunc
	Retire            RetireFunc
	GetDeckName       GetDeckNameFunc
	GetCommanderEntry GetCommanderEntryFunc
}
