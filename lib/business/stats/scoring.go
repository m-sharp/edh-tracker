package stats

// placeMultipliers defines the bonus points awarded per finishing position.
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

// GetPointsForRecord computes total points from an aggregated record (place → win count).
func GetPointsForRecord(kills int, record map[int]int) int {
	points := kills
	for place, count := range record {
		if multiplier, ok := placeMultipliers[place]; ok {
			points += count * multiplier
		}
	}
	return points
}
