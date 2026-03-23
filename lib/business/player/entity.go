package player

import (
	"fmt"
	"time"

	"github.com/m-sharp/edh-tracker/lib/business/stats"
	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	repo "github.com/m-sharp/edh-tracker/lib/repositories/player"
)

type Entity struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Stats     stats.Stats `json:"stats"`
	PodIDs    []int       `json:"pod_ids"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (e Entity) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(e.Name) > 256 {
		return fmt.Errorf("name must be 256 characters or fewer")
	}
	return nil
}

type PlayerWithRoleEntity struct {
	Entity
	Role string `json:"role"`
}

func ToEntity(m repo.Model, agg *gameResult.Aggregate, podIDs []int) Entity {
	if podIDs == nil {
		podIDs = []int{}
	}
	return Entity{
		ID:        m.ID,
		Name:      m.Name,
		Stats:     stats.FromAggregate(agg),
		PodIDs:    podIDs,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
