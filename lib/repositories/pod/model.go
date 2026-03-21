package pod

import (
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	playerPodRoleRepo "github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
)

type Model struct {
	base.GormModelBase
	Name    string
	Members []playerPodRoleRepo.Model `gorm:"foreignKey:PodID"`
}

func (Model) TableName() string { return "pod" }

type PlayerPodModel struct {
	base.GormModelBase
	PodID    int
	PlayerID int
}

func (PlayerPodModel) TableName() string { return "player_pod" }
