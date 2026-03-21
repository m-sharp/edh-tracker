package pod

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.GormModelBase
	Name string
}

func (Model) TableName() string { return "pod" }

type PlayerPodModel struct {
	base.GormModelBase
	PodID    int
	PlayerID int
}

func (PlayerPodModel) TableName() string { return "player_pod" }
