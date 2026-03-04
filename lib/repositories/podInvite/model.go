package podInvite

import (
	"time"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

type Model struct {
	base.ModelBase
	PodID             int        `db:"pod_id"`
	InviteCode        string     `db:"invite_code"`
	CreatedByPlayerID int        `db:"created_by_player_id"`
	ExpiresAt         *time.Time `db:"expires_at"`
	UsedCount         int        `db:"used_count"`
}
