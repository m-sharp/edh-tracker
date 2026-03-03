package gameResult

// Aggregate is a computed summary of game_result rows via SQL aggregation.
type Aggregate struct {
	Record map[int]int // place → count
	Games  int
	Kills  int
}

type gameStat struct {
	GameID    int `db:"game_id"`
	Place     int `db:"place"`
	KillCount int `db:"kill_count"`
}

type gameStats []gameStat

func (g gameStats) toAggregate() Aggregate {
	kills := 0
	record := map[int]int{}
	for _, stat := range g {
		kills += stat.KillCount

		if _, ok := record[stat.Place]; !ok {
			record[stat.Place] = 1
		} else {
			record[stat.Place] += 1
		}
	}

	return Aggregate{
		Record: record,
		Games:  len(g),
		Kills:  kills,
	}
}
