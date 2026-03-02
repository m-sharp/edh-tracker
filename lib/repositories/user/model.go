package user

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

const (
	RoleAdmin  = "admin"
	RolePlayer = "player"
)

type Model struct {
	base.ModelBase
	PlayerID      int     `db:"player_id"`
	RoleID        int     `db:"role_id"`
	OAuthProvider *string `db:"oauth_provider"`
	OAuthSubject  *string `db:"oauth_subject"`
	Email         *string `db:"email"`
	DisplayName   *string `db:"display_name"`
	AvatarURL     *string `db:"avatar_url"`
}

type RoleModel struct {
	base.ModelBase
	Name string `db:"name"`
}
