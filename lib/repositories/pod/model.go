package pod

import "time"

type Model struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// TODO: Consider if a specific repository subpackage is needed for PlayerPod
type PlayerPodModel struct {
	ID        int        `db:"id"`
	PodID     int        `db:"pod_id"`
	PlayerID  int        `db:"player_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
