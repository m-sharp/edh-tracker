package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPointsForPlace(t *testing.T) {
	tests := []struct {
		kills, place, want int
	}{
		{kills: 2, place: 1, want: 5}, // 2 kills + 3 bonus
		{kills: 0, place: 2, want: 2}, // 0 kills + 2 bonus
		{kills: 1, place: 3, want: 2}, // 1 kill + 1 bonus
		{kills: 3, place: 4, want: 3}, // 3 kills + 0 bonus
		{kills: 0, place: 0, want: 0}, // 0 kills + 0 bonus
	}
	for _, tt := range tests {
		got := GetPointsForPlace(tt.kills, tt.place)
		assert.Equal(t, tt.want, got)
	}
}

func TestGetPointsForRecord(t *testing.T) {
	// 3 kills, record: {1:2, 2:1} => 3 + 2*3 + 1*2 = 11
	got := GetPointsForRecord(3, map[int]int{1: 2, 2: 1})
	assert.Equal(t, 11, got)
}
