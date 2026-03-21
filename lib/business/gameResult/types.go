package gameResult

import (
	"context"

	gameResultRepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

type GetByGameIDFunc func(ctx context.Context, gameID int) ([]Entity, error)
type GetGameIDForResultFunc func(ctx context.Context, resultID int) (int, error)
type EnrichModelsFunc func(ctx context.Context, models []gameResultRepo.Model) ([]Entity, error)

type Functions struct {
	GetByGameID        GetByGameIDFunc
	GetGameIDForResult GetGameIDForResultFunc
	EnrichModels       EnrichModelsFunc
}
