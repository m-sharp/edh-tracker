package commander

import (
	"time"

	repo "github.com/m-sharp/edh-tracker/lib/repositories/commander"
)

type Entity struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToEntity(m repo.Model) Entity {
	return Entity{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
