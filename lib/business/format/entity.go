package format

import (
	"time"

	repo "github.com/m-sharp/edh-tracker/lib/repositories/format"
	"github.com/m-sharp/edh-tracker/lib/utils"
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
		Name:      utils.TitleCase(m.Name),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
