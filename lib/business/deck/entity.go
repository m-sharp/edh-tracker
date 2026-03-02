package deck

import (
	"fmt"
	"time"

	"github.com/m-sharp/edh-tracker/lib/business/stats"
	repo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

type Entity struct {
	ID         int            `json:"id"`
	PlayerID   int            `json:"player_id"`
	PlayerName string         `json:"player_name,omitempty"`
	Name       string         `json:"name"`
	FormatID   int            `json:"format_id"`
	FormatName string         `json:"format_name"`
	Retired    bool           `json:"retired"`
	Commanders *CommanderInfo `json:"commanders,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// CommanderInfo is the entity representation of deck commander assignments.
type CommanderInfo struct {
	CommanderID          int     `json:"commander_id"`
	CommanderName        string  `json:"commander_name"`
	PartnerCommanderID   *int    `json:"partner_commander_id,omitempty"`
	PartnerCommanderName *string `json:"partner_commander_name,omitempty"`
}

type EntityWithStats struct {
	Entity
	Stats stats.Stats `json:"stats"`
}

// ValidateCreate checks required fields on an incoming create request.
func ValidateCreate(playerID int, name string, formatID int) error {
	if playerID == 0 {
		return fmt.Errorf("player_id is required")
	}
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if formatID == 0 {
		return fmt.Errorf("format_id is required")
	}

	return nil
}

func ToEntity(m repo.Model, playerName, formatName string, commanders *CommanderInfo) Entity {
	return Entity{
		ID:         m.ID,
		PlayerID:   m.PlayerID,
		PlayerName: playerName,
		Name:       m.Name,
		FormatID:   m.FormatID,
		FormatName: formatName,
		Retired:    m.Retired,
		Commanders: commanders,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ToEntityWithStats(e Entity, agg *gameResult.Aggregate) EntityWithStats {
	return EntityWithStats{
		Entity: e,
		Stats:  stats.FromAggregate(agg),
	}
}
