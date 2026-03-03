package gameResult

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

// mockGameResultRepo implements repos.GameResultRepository
type mockGameResultRepo struct {
	GetByGameIdFn func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error)
}

func (m *mockGameResultRepo) GetByGameId(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	if m.GetByGameIdFn != nil {
		return m.GetByGameIdFn(ctx, gameID)
	}
	panic("unexpected call to GetByGameId")
}
func (m *mockGameResultRepo) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForPlayer")
}
func (m *mockGameResultRepo) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForDeck")
}
func (m *mockGameResultRepo) BulkAdd(ctx context.Context, results []gameresultrepo.Model) error {
	panic("unexpected call to BulkAdd")
}
func (m *mockGameResultRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}

func TestGetByGameID_NoResults(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIdFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{}, nil
		},
	}
	getDeckName := func(ctx context.Context, deckID int) (string, error) {
		panic("should not be called")
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*deck.CommanderInfo, error) {
		panic("should not be called")
	}

	fn := GetByGameID(repo, getDeckName, getCommanderEntry)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func TestGetByGameID_NoCommander(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIdFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, GameID: 1, DeckID: 20, Place: 1, KillCount: 2},
			}, nil
		},
	}
	getDeckName := func(ctx context.Context, deckID int) (string, error) {
		return "My Deck", nil
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*deck.CommanderInfo, error) {
		return nil, nil
	}

	fn := GetByGameID(repo, getDeckName, getCommanderEntry)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "My Deck", got[0].DeckName)
	assert.Nil(t, got[0].CommanderName)
	// Points: 2 kills + max(0, 1-1)=0 bonus = 2 (1-player game)
	assert.Equal(t, 2, got[0].Points)
}

func TestGetByGameID_WithCommander(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIdFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, GameID: 1, DeckID: 20, Place: 1, KillCount: 0},
			}, nil
		},
	}
	getDeckName := func(ctx context.Context, deckID int) (string, error) {
		return "Krenko Deck", nil
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*deck.CommanderInfo, error) {
		return &deck.CommanderInfo{CommanderID: 5, CommanderName: "Krenko"}, nil
	}

	fn := GetByGameID(repo, getDeckName, getCommanderEntry)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].CommanderName)
	assert.Equal(t, "Krenko", *got[0].CommanderName)
	assert.Nil(t, got[0].PartnerCommanderName)
}

func TestGetByGameID_WithPartner(t *testing.T) {
	partnerID := 6
	partnerName := "Goblin Chieftain"
	repo := &mockGameResultRepo{
		GetByGameIdFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, GameID: 1, DeckID: 20, Place: 2, KillCount: 1},
			}, nil
		},
	}
	getDeckName := func(ctx context.Context, deckID int) (string, error) {
		return "Partner Deck", nil
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*deck.CommanderInfo, error) {
		return &deck.CommanderInfo{
			CommanderID:          5,
			CommanderName:        "Krenko",
			PartnerCommanderID:   &partnerID,
			PartnerCommanderName: &partnerName,
		}, nil
	}

	fn := GetByGameID(repo, getDeckName, getCommanderEntry)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].CommanderName)
	assert.Equal(t, "Krenko", *got[0].CommanderName)
	require.NotNil(t, got[0].PartnerCommanderName)
	assert.Equal(t, "Goblin Chieftain", *got[0].PartnerCommanderName)
}

func TestGetByGameID_DeckNameCached(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIdFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, GameID: 1, DeckID: 20, Place: 1, KillCount: 0},
				{ModelBase: base.ModelBase{ID: 2}, GameID: 1, DeckID: 20, Place: 2, KillCount: 0},
			}, nil
		},
	}
	getDeckNameCallCount := 0
	getDeckName := func(ctx context.Context, deckID int) (string, error) {
		getDeckNameCallCount++
		return "Shared Deck", nil
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*deck.CommanderInfo, error) {
		return nil, nil
	}

	fn := GetByGameID(repo, getDeckName, getCommanderEntry)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 1, getDeckNameCallCount, "getDeckName should be called exactly once due to cache")
}

func TestGetByGameID_PointsCalculation(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIdFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, GameID: 1, DeckID: 10, Place: 1, KillCount: 2},
				{ModelBase: base.ModelBase{ID: 2}, GameID: 1, DeckID: 11, Place: 2, KillCount: 0},
				{ModelBase: base.ModelBase{ID: 3}, GameID: 1, DeckID: 12, Place: 4, KillCount: 1},
			}, nil
		},
	}
	getDeckName := func(ctx context.Context, deckID int) (string, error) {
		return "Deck", nil
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*deck.CommanderInfo, error) {
		return nil, nil
	}

	fn := GetByGameID(repo, getDeckName, getCommanderEntry)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 3)
	assert.Equal(t, 4, got[0].Points) // 2 kills + max(0,3-1)=2 (3-player game, 1st place)
	assert.Equal(t, 1, got[1].Points) // 0 kills + max(0,3-2)=1 (3-player game, 2nd place)
	assert.Equal(t, 1, got[2].Points) // 1 kill + max(0,3-4)=0 (3-player game, 4th place)
}
