package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPointsForPlace(t *testing.T) {
	tests := []struct {
		name  string
		kills int
		place int
		want  int
	}{
		{"1st place, 0 kills", 0, 1, 3},
		{"2nd place, 0 kills", 0, 2, 2},
		{"3rd place, 0 kills", 0, 3, 1},
		{"4th place, 0 kills", 0, 4, 0},
		{"1st place, 2 kills", 2, 1, 5},
		{"2nd place, 3 kills", 3, 2, 5},
		{"3rd place, 1 kill", 1, 3, 2},
		{"4th place, 2 kills", 2, 4, 2},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			assert.Equal(tt, testCase.want, getPointsForPlace(testCase.kills, testCase.place))
		})
	}
}

func TestGetPointsForRecord(t *testing.T) {
	tests := []struct {
		name   string
		kills  int
		record map[int]int
		want   int
	}{
		{"no games", 0, map[int]int{}, 0},
		{"1st place win with 2 kills", 2, map[int]int{1: 1}, 5},
		{"two 2nd place finishes, 0 kills", 0, map[int]int{2: 2}, 4},
		{"mix of places with kills", 3, map[int]int{1: 1, 3: 2}, 8},
		{"only 4th place, no kills", 0, map[int]int{4: 3}, 0},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			assert.Equal(tt, testCase.want, getPointsForRecord(testCase.kills, testCase.record))
		})
	}
}

func TestGameStats_ToStats(t *testing.T) {
	tests := []struct {
		name  string
		input GameStats
		want  Stats
	}{
		{
			name:  "empty",
			input: GameStats{},
			want:  Stats{Record: map[int]int{}, Games: 0, Kills: 0, Points: 0},
		},
		{
			name:  "single game, 1st place, 2 kills",
			input: GameStats{{GameID: 1, Place: 1, KillCount: 2}},
			want:  Stats{Record: map[int]int{1: 1}, Games: 1, Kills: 2, Points: 5},
		},
		{
			name: "multiple games, mixed places",
			input: GameStats{
				{GameID: 1, Place: 1, KillCount: 1},
				{GameID: 2, Place: 2, KillCount: 0},
				{GameID: 3, Place: 1, KillCount: 2},
			},
			// kills=3, record={1:2,2:1} → 3 + 2*3 + 1*2 = 11
			want: Stats{Record: map[int]int{1: 2, 2: 1}, Games: 3, Kills: 3, Points: 11},
		},
		{
			name:  "4th place only, no kills",
			input: GameStats{{GameID: 1, Place: 4, KillCount: 0}},
			want:  Stats{Record: map[int]int{4: 1}, Games: 1, Kills: 0, Points: 0},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(tt *testing.T) {
			assert.Equal(tt, testCase.want, testCase.input.ToStats())
		})
	}
}
