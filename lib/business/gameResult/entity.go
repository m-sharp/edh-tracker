package gameResult

import "fmt"

var placeMultipliers = map[int]int{
	1: 3,
	2: 2,
	3: 1,
}

// GetPointsForPlace computes points for a single game result.
func GetPointsForPlace(kills, place int) int {
	points := kills
	if bonus, ok := placeMultipliers[place]; ok {
		points += bonus
	}
	return points
}

// InputEntity is used by the game creation flow — what the router receives for each result.
type InputEntity struct {
	DeckID int `json:"deck_id"`
	Place  int `json:"place"`
	Kills  int `json:"kill_count"`
}

func (i InputEntity) Validate() error {
	if i.DeckID == 0 {
		return fmt.Errorf("deck_id is required")
	}
	if i.Place < 1 {
		return fmt.Errorf("place must be >= 1")
	}
	if i.Kills < 0 {
		return fmt.Errorf("kill_count cannot be negative")
	}
	return nil
}

// Entity is the enriched game result returned by GET endpoints.
type Entity struct {
	ID                   int     `json:"id"`
	GameID               int     `json:"game_id"`
	DeckID               int     `json:"deck_id"`
	DeckName             string  `json:"deck_name"`
	CommanderName        *string `json:"commander_name,omitempty"`
	PartnerCommanderName *string `json:"partner_commander_name,omitempty"`
	Place                int     `json:"place"`
	Kills                int     `json:"kill_count"`
	Points               int     `json:"points"`
}
