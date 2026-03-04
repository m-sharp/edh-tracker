package playerPodRole

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getRole            = `SELECT id, pod_id, player_id, role, created_at, updated_at, deleted_at FROM player_pod_role WHERE pod_id = ? AND player_id = ? AND deleted_at IS NULL LIMIT 1;`
	getMembersWithRole = `SELECT id, pod_id, player_id, role, created_at, updated_at, deleted_at FROM player_pod_role WHERE pod_id = ? AND deleted_at IS NULL;`
	setRole            = `INSERT INTO player_pod_role (pod_id, player_id, role) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE role = VALUES(role), deleted_at = NULL;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetRole(ctx context.Context, podID, playerID int) (*Model, error) {
	var rows []Model
	if err := r.client.Db.SelectContext(ctx, &rows, getRole, podID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get role for player %d in pod %d: %w", playerID, podID, err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (r *Repository) SetRole(ctx context.Context, podID, playerID int, role string) error {
	if _, err := r.client.Db.ExecContext(ctx, setRole, podID, playerID, role); err != nil {
		return fmt.Errorf("failed to set role %q for player %d in pod %d: %w", role, playerID, podID, err)
	}
	return nil
}

func (r *Repository) GetMembersWithRoles(ctx context.Context, podID int) ([]Model, error) {
	var rows []Model
	if err := r.client.Db.SelectContext(ctx, &rows, getMembersWithRole, podID); err != nil {
		return nil, fmt.Errorf("failed to get members with roles for pod %d: %w", podID, err)
	}
	if rows == nil {
		return []Model{}, nil
	}
	return rows, nil
}

func (r *Repository) BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error {
	if len(playerIDs) == 0 {
		return nil
	}

	query := "INSERT INTO player_pod_role (pod_id, player_id, role) VALUES " +
		strings.TrimSuffix(strings.Repeat("(?,?,?),", len(playerIDs)), ",")
	args := make([]interface{}, 0, len(playerIDs)*3)
	for _, id := range playerIDs {
		args = append(args, podID, id, role)
	}

	if _, err := r.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert player_pod_role records: %w", err)
	}
	return nil
}
