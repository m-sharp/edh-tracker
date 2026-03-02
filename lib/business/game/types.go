package game

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib/business/gameresult"
)

type GetAllByPodFunc func(ctx context.Context, podID int) ([]Entity, error)
type GetAllByDeckFunc func(ctx context.Context, deckID int) ([]Entity, error)
type GetByIDFunc func(ctx context.Context, gameID int) (*Entity, error)
type CreateFunc func(ctx context.Context, description string, podID, formatID int, results []gameresult.InputEntity) error

type Functions struct {
	GetAllByPod  GetAllByPodFunc
	GetAllByDeck GetAllByDeckFunc
	GetByID      GetByIDFunc
	Create       CreateFunc
}
