package gameResult

import "github.com/m-sharp/edh-tracker/lib/utils"

// Aggregate is a computed summary of game_result rows via SQL aggregation.
type Aggregate struct {
	Record map[int]int // place → count
	Games  int
	Kills  int
	Points int
}

type gameStat struct {
	GameID      int
	Place       int
	KillCount   int
	PlayerCount int
}

type gameStatWithDeck struct {
	DeckID      int
	GameID      int
	Place       int
	KillCount   int
	PlayerCount int
}

type gameStatWithPlayer struct {
	PlayerID    int
	GameID      int
	Place       int
	KillCount   int
	PlayerCount int
}

type gameStats []gameStat

func (g gameStats) toAggregate() Aggregate {
	kills, points := 0, 0
	record := map[int]int{}
	for _, s := range g {
		kills += s.KillCount
		points += utils.GetPointsForPlace(s.KillCount, s.Place, s.PlayerCount)
		record[s.Place]++
	}
	return Aggregate{
		Record: record,
		Games:  len(g),
		Kills:  kills,
		Points: points,
	}
}
