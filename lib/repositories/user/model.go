package user

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

const (
	RoleAdmin  = "admin"
	RolePlayer = "player"
)

type Model struct {
	base.GormModelBase
	PlayerID      int
	RoleID        int
	OAuthProvider *string `gorm:"column:oauth_provider"`
	OAuthSubject  *string `gorm:"column:oauth_subject"`
	Email         *string
	DisplayName   *string
	AvatarURL     *string
}

func (Model) TableName() string { return "user" }

type RoleModel struct {
	base.GormModelBase
	Name string
}

func (RoleModel) TableName() string { return "user_role" }
