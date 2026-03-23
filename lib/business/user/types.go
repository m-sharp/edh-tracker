package user

import "context"

type GetByIDFunc func(ctx context.Context, id int) (*Entity, error)
type GetByPlayerIDFunc func(ctx context.Context, playerID int) (*Entity, error)
type CreateFunc func(ctx context.Context, playerID, roleID int) (int, error)
type GetByOAuthFunc func(ctx context.Context, provider, subject string) (*Entity, error)
type CreateWithOAuthFunc func(ctx context.Context, playerName, provider, subject, email, displayName, avatarURL string) (*Entity, error)
type GetByEmailFunc func(ctx context.Context, email string) (*Entity, error)
type LinkOAuthFunc func(ctx context.Context, userID int, provider, subject, email, displayName, avatarURL string) (*Entity, error)

type Functions struct {
	GetByID         GetByIDFunc
	GetByPlayerID   GetByPlayerIDFunc
	Create          CreateFunc
	GetByOAuth      GetByOAuthFunc
	CreateWithOAuth CreateWithOAuthFunc
	GetByEmail      GetByEmailFunc
	LinkOAuth       LinkOAuthFunc
}
