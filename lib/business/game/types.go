package game

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
)

type GetAllByPodFunc func(ctx context.Context, podID int) ([]Entity, error)
type GetAllByDeckFunc func(ctx context.Context, deckID int) ([]Entity, error)
type GetAllByPlayerFunc func(ctx context.Context, playerID int) ([]Entity, error)
type GetByIDFunc func(ctx context.Context, gameID int) (*Entity, error)
type CreateFunc func(ctx context.Context, description string, podID, formatID int, results []gameResult.InputEntity) error
type UpdateFunc func(ctx context.Context, gameID int, description string) error
type SoftDeleteFunc func(ctx context.Context, gameID int) error
type AddResultFunc func(ctx context.Context, gameID, deckID, playerID, place, killCount int) (int, error)
type UpdateResultFunc func(ctx context.Context, resultID, place, killCount, deckID int) error
type DeleteResultFunc func(ctx context.Context, resultID int) error

type Functions struct {
	GetAllByPod    GetAllByPodFunc
	GetAllByDeck   GetAllByDeckFunc
	GetAllByPlayer GetAllByPlayerFunc
	GetByID        GetByIDFunc
	Create         CreateFunc
	Update         UpdateFunc
	SoftDelete     SoftDeleteFunc
	AddResult      AddResultFunc
	UpdateResult   UpdateResultFunc
	DeleteResult   DeleteResultFunc
}
