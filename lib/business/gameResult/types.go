package gameResult

import "context"

type GetByGameIDFunc func(ctx context.Context, gameID int) ([]Entity, error)

type Functions struct {
	GetByGameID GetByGameIDFunc
}
