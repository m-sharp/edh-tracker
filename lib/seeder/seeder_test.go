package seeder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestSeeder() *Seeder {
	return NewSeeder(zap.NewNop(), nil)
}

func TestCollectEntities_Empty(t *testing.T) {
	s := newTestSeeder()
	players, commanders, entries, err := s.collectEntities([]Game{}, map[string]int{"commander": 1})
	require.NoError(t, err)
	assert.Empty(t, players)
	assert.Empty(t, commanders)
	assert.Empty(t, entries)
}

func TestCollectEntities_Dedup(t *testing.T) {
	s := newTestSeeder()
	formatIDs := map[string]int{"commander": 1}
	games := []Game{
		{
			Format: "commander",
			Results: []Result{
				{Player: "Alice", Name: "Atraxa", Place: 1, Kills: 2},
				{Player: "Bob", Name: "Najeela", Place: 2, Kills: 0},
			},
		},
		{
			Format: "commander",
			Results: []Result{
				{Player: "Alice", Name: "Atraxa", Place: 2, Kills: 1},
				{Player: "Bob", Name: "Najeela", Place: 1, Kills: 0},
			},
		},
	}
	players, commanders, entries, err := s.collectEntities(games, formatIDs)
	require.NoError(t, err)
	assert.Len(t, players, 2)
	assert.Len(t, commanders, 2)
	assert.Len(t, entries, 2, "duplicate player+commander pairs should be deduped")
}

func TestCollectEntities_UnknownFormat(t *testing.T) {
	s := newTestSeeder()
	formatIDs := map[string]int{"commander": 1}
	games := []Game{
		{
			Format:  "unknown-format",
			Results: []Result{{Player: "Alice", Name: "Atraxa", Place: 1}},
		},
	}
	_, _, _, err := s.collectEntities(games, formatIDs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format")
}
