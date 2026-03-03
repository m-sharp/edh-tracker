package user

import (
	"time"

	repo "github.com/m-sharp/edh-tracker/lib/repositories/user"
)

type Entity struct {
	ID            int       `json:"id"`
	PlayerID      int       `json:"player_id"`
	RoleID        int       `json:"role_id"`
	OAuthProvider *string   `json:"oauth_provider"`
	OAuthSubject  *string   `json:"oauth_subject"`
	Email         *string   `json:"email"`
	DisplayName   *string   `json:"display_name"`
	AvatarURL     *string   `json:"avatar_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToEntity(m repo.Model) Entity {
	return Entity{
		ID:            m.ID,
		PlayerID:      m.PlayerID,
		RoleID:        m.RoleID,
		OAuthProvider: m.OAuthProvider,
		OAuthSubject:  m.OAuthSubject,
		Email:         m.Email,
		DisplayName:   m.DisplayName,
		AvatarURL:     m.AvatarURL,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
