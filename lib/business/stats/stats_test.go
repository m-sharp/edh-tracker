package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

func TestFromAggregate_Nil(t *testing.T) {
	s := FromAggregate(nil)
	assert.NotNil(t, s.Record)
	assert.Equal(t, 0, s.Games)
	assert.Equal(t, 0, s.Kills)
	assert.Equal(t, 0, s.Points)
}

func TestFromAggregate_NonNil(t *testing.T) {
	agg := &gameresultrepo.Aggregate{
		Games:  3,
		Kills:  4,
		Points: 12,
		Record: map[int]int{1: 2},
	}
	s := FromAggregate(agg)
	assert.Equal(t, 3, s.Games)
	assert.Equal(t, 4, s.Kills)
	assert.Equal(t, 12, s.Points)
	assert.Equal(t, map[int]int{1: 2}, s.Record)
}
