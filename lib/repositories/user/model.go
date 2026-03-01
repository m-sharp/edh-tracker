package user

import "time"

const (
	RoleAdmin  = "admin"
	RolePlayer = "player"
)

type Model struct {
	ID            int        `db:"id"`
	PlayerID      int        `db:"player_id"`
	RoleID        int        `db:"role_id"`
	OAuthProvider *string    `db:"oauth_provider"`
	OAuthSubject  *string    `db:"oauth_subject"`
	Email         *string    `db:"email"`
	DisplayName   *string    `db:"display_name"`
	AvatarURL     *string    `db:"avatar_url"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

// TODO: Consider if a specific repository subpackage is needed for Role
type RoleModel struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
