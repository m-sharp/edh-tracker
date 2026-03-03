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
}
