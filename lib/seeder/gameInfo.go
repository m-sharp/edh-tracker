package seeder

import "time"

type Game struct {
	Date    time.Time `json:"date"`
	Format  string    `json:"format"`
	Results []Result  `json:"results"`
}

type Result struct {
	Player           string `json:"player"`
	Name             string `json:"name"`             // deck display name; defaults to Commander
	Commander        string `json:"commander"`        // primary commander; defaults to Name
	PartnerCommander string `json:"partnerCommander"` // optional partner commander name
	Place            int    `json:"place"`
	Kills            int    `json:"kills"`
}

// DeckName returns the deck's display name (Name if set, else Commander).
func (r Result) DeckName() string {
	if r.Name != "" {
		return r.Name
	}
	return r.Commander
}

// CommanderName returns the primary commander name (Commander if set, else Name).
func (r Result) CommanderName() string {
	if r.Commander != "" {
		return r.Commander
	}
	return r.Name
}
