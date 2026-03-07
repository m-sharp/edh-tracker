package podInvite

import (
	"context"
	"fmt"
	"time"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getByCode          = `SELECT id, pod_id, invite_code, created_by_player_id, expires_at, used_count, created_at, updated_at, deleted_at FROM pod_invite WHERE invite_code = ? AND deleted_at IS NULL LIMIT 1;`
	insertInvite       = `INSERT INTO pod_invite (pod_id, invite_code, created_by_player_id, expires_at) VALUES (?, ?, ?, ?);`
	incrementUsedCount = `UPDATE pod_invite SET used_count = used_count + 1 WHERE invite_code = ?;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*Model, error) {
	var rows []Model
	if err := r.client.Db.SelectContext(ctx, &rows, getByCode, code); err != nil {
		return nil, fmt.Errorf("failed to get pod invite for code %q: %w", code, err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (r *Repository) Add(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
	if _, err := r.client.Db.ExecContext(ctx, insertInvite, podID, code, createdByPlayerID, expiresAt); err != nil {
		return fmt.Errorf("failed to insert pod invite: %w", err)
	}
	return nil
}

func (r *Repository) IncrementUsedCount(ctx context.Context, code string) error {
	if _, err := r.client.Db.ExecContext(ctx, incrementUsedCount, code); err != nil {
		return fmt.Errorf("failed to increment used_count for invite %q: %w", code, err)
	}
	return nil
}
