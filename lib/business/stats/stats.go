package stats

import (
	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

// Stats is the entity representation of aggregated game statistics.
type Stats struct {
	Record map[int]int `json:"record"`
	Games  int         `json:"games"`
	Kills  int         `json:"kills"`
	Points int         `json:"points"`
}

func FromAggregate(a *gameResult.Aggregate) Stats {
	if a == nil {
		return Stats{Record: map[int]int{}}
	}
	return Stats{
		Record: a.Record,
		Games:  a.Games,
		Kills:  a.Kills,
		Points: GetPointsForRecord(a.Kills, a.Record),
	}
}
