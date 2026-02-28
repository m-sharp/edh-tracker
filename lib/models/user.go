package models

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	RoleAdmin  = "admin"
	RolePlayer = "player"

	GetUserByID       = `SELECT id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url, created_at, updated_at, deleted_at FROM user WHERE id = ? AND deleted_at IS NULL;`
	GetUserByPlayerID = `SELECT id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url, created_at, updated_at, deleted_at FROM user WHERE player_id = ? AND deleted_at IS NULL;`
	InsertUser        = `INSERT INTO user (player_id, role_id) VALUES (?, ?);`
	SoftDeleteUser    = `UPDATE user SET deleted_at = NOW() WHERE id = ?;`
)

// UserRole maps to the user_role lookup table.
type UserRole struct {
	Model
	Name string `json:"name" db:"name"`
}

// User maps to the user table. OAuth fields are nullable to support
// stub records created before any real login has occurred.
type User struct {
	Model
	PlayerID      int     `json:"player_id"      db:"player_id"`
	RoleID        int     `json:"role_id"        db:"role_id"`
	OAuthProvider *string `json:"oauth_provider" db:"oauth_provider"`
	OAuthSubject  *string `json:"oauth_subject"  db:"oauth_subject"`
	Email         *string `json:"email"          db:"email"`
	DisplayName   *string `json:"display_name"   db:"display_name"`
	AvatarURL     *string `json:"avatar_url"     db:"avatar_url"`
}

type UserProvider struct {
	client *lib.DBClient
}

func NewUserProvider(client *lib.DBClient) *UserProvider {
	return &UserProvider{
		client: client,
	}
}

func (u *UserProvider) GetByID(ctx context.Context, id int) (*User, error) {
	var users []User
	if err := u.client.Db.SelectContext(ctx, &users, GetUserByID, id); err != nil {
		return nil, fmt.Errorf("failed to get User record for id %d: %w", id, err)
	}

	if len(users) == 0 || len(users) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of users returned for ID %d: got %d, expected 1",
			id, len(users),
		)
	}

	return &users[0], nil
}

func (u *UserProvider) GetByPlayerID(ctx context.Context, playerID int) (*User, error) {
	var users []User
	if err := u.client.Db.SelectContext(ctx, &users, GetUserByPlayerID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get User record for player_id %d: %w", playerID, err)
	}

	if len(users) == 0 || len(users) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of users returned for player_id %d: got %d, expected 1",
			playerID, len(users),
		)
	}

	return &users[0], nil
}

func (u *UserProvider) Add(ctx context.Context, playerID, roleID int) error {
	result, err := u.client.Db.ExecContext(ctx, InsertUser, playerID, roleID)
	if err != nil {
		return fmt.Errorf("failed to insert User record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by User insert: got %d, expected 1", numAffected)
	}

	return nil
}

func (u *UserProvider) SoftDelete(ctx context.Context, id int) error {
	result, err := u.client.Db.ExecContext(ctx, SoftDeleteUser, id)
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
