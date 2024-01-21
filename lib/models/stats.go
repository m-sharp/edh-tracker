package models

type Stats struct {
	Record map[int]int `json:"record"`
	Games  int         `json:"games"`
	Kills  int         `json:"kills"`
	Points int         `json:"points"`
}

type GameStat struct {
	GameID    int `db:"game_id"`
	Place     int `db:"place"`
	KillCount int `db:"kill_count"`
}

type GameStats []GameStat

func (g GameStats) ToStats() Stats {
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

	return Stats{
		Record: record,
		Games:  len(g),
		Kills:  kills,
		Points: getPointsForRecord(kills, record),
	}
}

var (
	placeMultipliers = map[int]int{
		1: 3,
		2: 2,
		3: 1,
	}
)

func getPointsForRecord(kills int, record map[int]int) int {
	points := kills
	for place, count := range record {
		if multiplier, ok := placeMultipliers[place]; ok {
			points += count * multiplier
		}
	}

	return points
}

func getPointsForPlace(kills, place int) int {
	points := kills
	if placePoints, ok := placeMultipliers[place]; ok {
		points += placePoints
	}

	return points
}
