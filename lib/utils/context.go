package utils

import "context"

type contextKey string

const (
	userIDKey   contextKey = "userID"
	playerIDKey contextKey = "playerID"
)

// TODO: Is player even necessary?
func ContextWithUserInfo(ctx context.Context, userID, playerID int) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, playerIDKey, playerID)
	return ctx
}

// UserFromContext extracts userID and playerID from context
func UserFromContext(ctx context.Context) (int, int, bool) {
	userID, ok1 := ctx.Value(userIDKey).(int)
	playerID, ok2 := ctx.Value(playerIDKey).(int)
	return userID, playerID, ok1 && ok2
}
