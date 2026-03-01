package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getUserByID       = `SELECT id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url, created_at, updated_at, deleted_at FROM user WHERE id = ? AND deleted_at IS NULL;`
	getUserByPlayerID = `SELECT id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url, created_at, updated_at, deleted_at FROM user WHERE player_id = ? AND deleted_at IS NULL;`
	getRoleByName     = `SELECT id, name, created_at, updated_at, deleted_at FROM user_role WHERE name = ?;`
	insertUser        = `INSERT INTO user (player_id, role_id) VALUES (?, ?);`
	softDeleteUser    = `UPDATE user SET deleted_at = NOW() WHERE id = ?;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Model, error) {
	var users []Model
	if err := r.client.Db.SelectContext(ctx, &users, getUserByID, id); err != nil {
		return nil, fmt.Errorf("failed to get User record for id %d: %w", id, err)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

func (r *Repository) GetByPlayerID(ctx context.Context, playerID int) (*Model, error) {
	var users []Model
	if err := r.client.Db.SelectContext(ctx, &users, getUserByPlayerID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get User record for player_id %d: %w", playerID, err)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (*RoleModel, error) {
	var roles []RoleModel
	if err := r.client.Db.SelectContext(ctx, &roles, getRoleByName, name); err != nil {
		return nil, fmt.Errorf("failed to get UserRole record for name %q: %w", name, err)
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("no UserRole found with name %q", name)
	}
	return &roles[0], nil
}

func (r *Repository) Add(ctx context.Context, playerID, roleID int) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertUser, playerID, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert User record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by User insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new User: %w", err)
	}

	return int(id), nil
}

func (r *Repository) BulkAdd(ctx context.Context, playerIDs []int, roleID int) error {
	if len(playerIDs) == 0 {
		return nil
	}

	query := "INSERT INTO user (player_id, role_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?),", len(playerIDs)), ",")
	args := make([]interface{}, 0, len(playerIDs)*2)
	for _, id := range playerIDs {
		args = append(args, id, roleID)
	}
	if _, err := r.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert User records: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	result, err := r.client.Db.ExecContext(ctx, softDeleteUser, id)
	if err != nil {
		return fmt.Errorf("failed to soft-delete User record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by soft-delete: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by User soft-delete: got %d, expected 1", numAffected)
	}

	return nil
}
