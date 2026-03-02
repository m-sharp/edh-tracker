package game

import (
	"time"

	"github.com/m-sharp/edh-tracker/lib/business/gameresult"
)

type Entity struct {
	ID          int                 `json:"id"`
	Description string              `json:"description"`
	PodID       int                 `json:"pod_id"`
	FormatID    int                 `json:"format_id"`
	Results     []gameresult.Entity `json:"results"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
