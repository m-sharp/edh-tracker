package playerPodRole

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

const (
	RoleManager = "manager"
	RoleMember  = "member"
)

type Model struct {
	base.ModelBase
	PodID    int    `db:"pod_id"`
	PlayerID int    `db:"player_id"`
	Role     string `db:"role"`
}
