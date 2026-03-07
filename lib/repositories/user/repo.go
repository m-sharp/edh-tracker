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
	getUserByOAuth    = `SELECT id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url, created_at, updated_at, deleted_at FROM user WHERE oauth_provider = ? AND oauth_subject = ? AND deleted_at IS NULL LIMIT 1;`
	getRoleByName     = `SELECT id, name, created_at, updated_at, deleted_at FROM user_role WHERE name = ?;`
	insertUser        = `INSERT INTO user (player_id, role_id) VALUES (?, ?);`
	insertUserOAuth   = `INSERT INTO user (player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url) VALUES (?, ?, ?, ?, ?, ?, ?);`
	softDeleteUser    = `UPDATE user SET deleted_at = NOW() WHERE id = ?;`
	insertPlayerTx    = `INSERT INTO player (name) VALUES (?);`
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

func (r *Repository) GetByOAuth(ctx context.Context, provider, subject string) (*Model, error) {
	var users []Model
	if err := r.client.Db.SelectContext(ctx, &users, getUserByOAuth, provider, subject); err != nil {
		return nil, fmt.Errorf("failed to get User record for oauth %s/%s: %w", provider, subject, err)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

func (r *Repository) AddWithOAuth(
	ctx context.Context,
	playerID, roleID int,
	provider, subject, email, displayName, avatarURL string,
) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertUserOAuth, playerID, roleID, provider, subject, email, displayName, avatarURL)
	if err != nil {
		return 0, fmt.Errorf("failed to insert User record with OAuth: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new User with OAuth: %w", err)
	}
	return int(id), nil
}

func (r *Repository) CreatePlayerAndUser(
	ctx context.Context,
	playerName string,
	roleID int,
	provider, subject, email, displayName, avatarURL string,
) (*Model, error) {
	tx, err := r.client.Db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction for CreatePlayerAndUser: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	playerResult, err := tx.ExecContext(ctx, insertPlayerTx, playerName)
	if err != nil {
		return nil, fmt.Errorf("failed to insert player in CreatePlayerAndUser: %w", err)
	}
	playerID, err := playerResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get player insert ID in CreatePlayerAndUser: %w", err)
	}

	userResult, err := tx.ExecContext(ctx, insertUserOAuth, playerID, roleID, provider, subject, email, displayName, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user in CreatePlayerAndUser: %w", err)
	}
	userID, err := userResult.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user insert ID in CreatePlayerAndUser: %w", err)
	}

	var users []Model
	if err = tx.SelectContext(ctx, &users, getUserByID, userID); err != nil {
		return nil, fmt.Errorf("failed to fetch created user in CreatePlayerAndUser: %w", err)
	}
	if len(users) == 0 {
		err = fmt.Errorf("created user %d not found after insert", userID)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit CreatePlayerAndUser transaction: %w", err)
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
