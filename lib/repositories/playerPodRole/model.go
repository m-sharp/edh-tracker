package playerPodRole

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

const (
	RoleManager = "manager"
	RoleMember  = "member"
)

type Model struct {
	base.GormModelBase
	PodID    int
	PlayerID int
	Role     string
}

func (Model) TableName() string { return "player_pod_role" }
