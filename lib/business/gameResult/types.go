package gameResult

import "context"

type GetByGameIDFunc func(ctx context.Context, gameID int) ([]Entity, error)
type GetGameIDForResultFunc func(ctx context.Context, resultID int) (int, error)

type Functions struct {
	GetByGameID        GetByGameIDFunc
	GetGameIDForResult GetGameIDForResultFunc
}
