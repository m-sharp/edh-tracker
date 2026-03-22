package podInvite

import (
	"time"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

type Model struct {
	base.GormModelBase
	PodID             int
	InviteCode        string
	CreatedByPlayerID int
	ExpiresAt         *time.Time
	UsedCount         int
}

func (Model) TableName() string { return "pod_invite" }
