package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPointsForPlace(t *testing.T) {
	tests := []struct {
		kills, place, numPlayers, want int
	}{
		{kills: 2, place: 1, numPlayers: 4, want: 5}, // 2 kills + 3 bonus
		{kills: 0, place: 2, numPlayers: 4, want: 2}, // 0 kills + 2 bonus
		{kills: 1, place: 3, numPlayers: 4, want: 2}, // 1 kill + 1 bonus
		{kills: 3, place: 4, numPlayers: 4, want: 3}, // 3 kills + 0 bonus
		{kills: 0, place: 0, numPlayers: 4, want: 0}, // place=0 guard fires
		{kills: 2, place: 1, numPlayers: 0, want: 2}, // numPlayers=0 guard fires
	}
	for _, tt := range tests {
		got := GetPointsForPlace(tt.kills, tt.place, tt.numPlayers)
		assert.Equal(t, tt.want, got)
	}
}
