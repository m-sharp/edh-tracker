package gameResult

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAggregate_Empty(t *testing.T) {
	gs := gameStats{}
	agg := gs.toAggregate()
	assert.Equal(t, 0, agg.Games)
	assert.Equal(t, 0, agg.Kills)
	assert.Equal(t, 0, agg.Points)
	assert.NotNil(t, agg.Record)
	assert.Len(t, agg.Record, 0)
}

func TestToAggregate_WithData(t *testing.T) {
	gs := gameStats{
		{GameID: 1, Place: 1, KillCount: 2},
		{GameID: 2, Place: 2, KillCount: 0},
		{GameID: 3, Place: 1, KillCount: 1},
	}
	agg := gs.toAggregate()
	assert.Equal(t, 3, agg.Games)
	assert.Equal(t, 3, agg.Kills)
	assert.Equal(t, map[int]int{1: 2, 2: 1}, agg.Record)
	// 3 kills + 2 first-place wins * 3 + 1 second-place win * 2 = 3 + 6 + 2 = 11
	assert.Equal(t, 11, agg.Points)
}

func TestGetPointsForRecord(t *testing.T) {
	// 3 kills, record: {1:2, 2:1} => 3 + 2*3 + 1*2 = 11
	got := getPointsForRecord(3, map[int]int{1: 2, 2: 1})
	assert.Equal(t, 11, got)
}
