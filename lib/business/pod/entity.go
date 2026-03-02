package pod

import (
	"fmt"
	"time"

	repo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
)

type Entity struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e Entity) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// PlayerPodInputEntity is the input for adding a player to a pod.
type PlayerPodInputEntity struct {
	PodID    int `json:"pod_id"`
	PlayerID int `json:"player_id"`
}

func (p PlayerPodInputEntity) Validate() error {
	if p.PodID == 0 {
		return fmt.Errorf("pod_id is required")
	}
	if p.PlayerID == 0 {
		return fmt.Errorf("player_id is required")
	}
	return nil
}

func ToEntity(m repo.Model) Entity {
	return Entity{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
