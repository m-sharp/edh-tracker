package utils

// GetPointsForPlace computes points for a single game result.
func GetPointsForPlace(kills, place, numPlayers int) int {
	if place <= 0 || numPlayers <= 0 {
		return kills
	}
	return kills + max(0, numPlayers-place)
}
