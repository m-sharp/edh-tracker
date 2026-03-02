package pod

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.ModelBase
	Name string `db:"name"`
}

type PlayerPodModel struct {
	base.ModelBase
	PodID    int `db:"pod_id"`
	PlayerID int `db:"player_id"`
}
