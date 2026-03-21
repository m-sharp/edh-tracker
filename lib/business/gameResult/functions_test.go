package gameResult

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	commanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/commander"
	deckRepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

// TODO: These mock repos are being repeated all over in @lib/business/ packages. We should have a single @lib/business/testHelpers/ package in which all of these are declared and used from in a unified manner.
// mockGameResultRepo implements repos.GameResultRepository
type mockGameResultRepo struct {
	GetByGameIDWithDeckInfoFn func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error)
	GetByIDFn                 func(ctx context.Context, resultID int) (*gameresultrepo.Model, error)
}

func (m *mockGameResultRepo) GetByGameId(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	panic("unexpected call to GetByGameId")
}
func (m *mockGameResultRepo) GetByGameIDWithDeckInfo(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	if m.GetByGameIDWithDeckInfoFn != nil {
		return m.GetByGameIDWithDeckInfoFn(ctx, gameID)
	}
	panic("unexpected call to GetByGameIDWithDeckInfo")
}
func (m *mockGameResultRepo) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForPlayer")
}
func (m *mockGameResultRepo) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForDeck")
}
func (m *mockGameResultRepo) GetByID(ctx context.Context, resultID int) (*gameresultrepo.Model, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, resultID)
	}
	panic("unexpected call to GetByID")
}
func (m *mockGameResultRepo) Add(ctx context.Context, model gameresultrepo.Model) (int, error) {
	panic("unexpected call to Add")
}
func (m *mockGameResultRepo) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	panic("unexpected call to Update")
}
func (m *mockGameResultRepo) BulkAdd(ctx context.Context, results []gameresultrepo.Model) error {
	panic("unexpected call to BulkAdd")
}
func (m *mockGameResultRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}

// makeModel builds a gameresultrepo.Model with an inline Deck preloaded.
func makeModel(id, gameID, deckID, place, killCount int, deckName string, playerID int, commander *deckCommanderRepo.Model) gameresultrepo.Model {
	return gameresultrepo.Model{
		GormModelBase: base.GormModelBase{ID: id},
		GameID:        gameID,
		DeckID:        deckID,
		Place:         place,
		KillCount:     killCount,
		Deck: deckRepo.Model{
			GormModelBase: base.GormModelBase{ID: deckID},
			Name:          deckName,
			PlayerID:      playerID,
			Commander:     commander,
		},
	}
}

func TestGetByGameID_NoResults(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIDWithDeckInfoFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{}, nil
		},
	}
	fn := GetByGameID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func TestGetByGameID_NoCommander(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIDWithDeckInfoFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				makeModel(1, 1, 20, 1, 2, "My Deck", 7, nil),
			}, nil
		},
	}
	fn := GetByGameID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "My Deck", got[0].DeckName)
	assert.Equal(t, 7, got[0].PlayerID)
	assert.Nil(t, got[0].CommanderName)
	// Points: 2 kills + max(0, 1-1)=0 bonus = 2 (1-player game)
	assert.Equal(t, 2, got[0].Points)
}

func TestGetByGameID_WithCommander(t *testing.T) {
	cmdr := &deckCommanderRepo.Model{
		Commander: commanderRepo.Model{Name: "Krenko"},
	}
	repo := &mockGameResultRepo{
		GetByGameIDWithDeckInfoFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				makeModel(1, 1, 20, 1, 0, "Krenko Deck", 5, cmdr),
			}, nil
		},
	}
	fn := GetByGameID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].CommanderName)
	assert.Equal(t, "Krenko", *got[0].CommanderName)
	assert.Nil(t, got[0].PartnerCommanderName)
}

func TestGetByGameID_WithPartner(t *testing.T) {
	partner := &commanderRepo.Model{Name: "Goblin Chieftain"}
	cmdr := &deckCommanderRepo.Model{
		Commander:        commanderRepo.Model{Name: "Krenko"},
		PartnerCommander: partner,
	}
	repo := &mockGameResultRepo{
		GetByGameIDWithDeckInfoFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				makeModel(1, 1, 20, 2, 1, "Partner Deck", 5, cmdr),
			}, nil
		},
	}
	fn := GetByGameID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].CommanderName)
	assert.Equal(t, "Krenko", *got[0].CommanderName)
	require.NotNil(t, got[0].PartnerCommanderName)
	assert.Equal(t, "Goblin Chieftain", *got[0].PartnerCommanderName)
}

func TestGetByGameID_MultipleSameDeck(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIDWithDeckInfoFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				makeModel(1, 1, 20, 1, 0, "Shared Deck", 3, nil),
				makeModel(2, 1, 20, 2, 0, "Shared Deck", 3, nil),
			}, nil
		},
	}
	fn := GetByGameID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "Shared Deck", got[0].DeckName)
	assert.Equal(t, "Shared Deck", got[1].DeckName)
}

func TestGetByGameID_PointsCalculation(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByGameIDWithDeckInfoFn: func(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
			return []gameresultrepo.Model{
				makeModel(1, 1, 10, 1, 2, "Deck A", 1, nil),
				makeModel(2, 1, 11, 2, 0, "Deck B", 2, nil),
				makeModel(3, 1, 12, 4, 1, "Deck C", 3, nil),
			}, nil
		},
	}
	fn := GetByGameID(repo)
	got, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, got, 3)
	assert.Equal(t, 4, got[0].Points) // 2 kills + max(0,3-1)=2 (3-player game, 1st place)
	assert.Equal(t, 1, got[1].Points) // 0 kills + max(0,3-2)=1 (3-player game, 2nd place)
	assert.Equal(t, 1, got[2].Points) // 1 kill + max(0,3-4)=0 (3-player game, 4th place)
}

// TestEnrichModels_* tests the EnrichModels closure directly with pre-populated models.

func TestEnrichModels_NoCommander(t *testing.T) {
	models := []gameresultrepo.Model{
		makeModel(1, 1, 20, 1, 2, "My Deck", 7, nil),
	}
	fn := EnrichModels()
	got, err := fn(context.Background(), models)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 1, got[0].ID)
	assert.Equal(t, 1, got[0].GameID)
	assert.Equal(t, 20, got[0].DeckID)
	assert.Equal(t, 7, got[0].PlayerID)
	assert.Equal(t, "My Deck", got[0].DeckName)
	assert.Equal(t, 1, got[0].Place)
	assert.Equal(t, 2, got[0].Kills)
	assert.Nil(t, got[0].CommanderName)
}

func TestEnrichModels_WithCommander(t *testing.T) {
	cmdr := &deckCommanderRepo.Model{
		Commander: commanderRepo.Model{Name: "Krenko"},
	}
	models := []gameresultrepo.Model{
		makeModel(1, 1, 20, 1, 0, "Krenko Deck", 5, cmdr),
	}
	fn := EnrichModels()
	got, err := fn(context.Background(), models)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].CommanderName)
	assert.Equal(t, "Krenko", *got[0].CommanderName)
	assert.Nil(t, got[0].PartnerCommanderName)
}

func TestEnrichModels_WithPartner(t *testing.T) {
	partner := &commanderRepo.Model{Name: "Goblin Chieftain"}
	cmdr := &deckCommanderRepo.Model{
		Commander:        commanderRepo.Model{Name: "Krenko"},
		PartnerCommander: partner,
	}
	models := []gameresultrepo.Model{
		makeModel(1, 1, 20, 2, 1, "Partner Deck", 5, cmdr),
	}
	fn := EnrichModels()
	got, err := fn(context.Background(), models)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].CommanderName)
	assert.Equal(t, "Krenko", *got[0].CommanderName)
	require.NotNil(t, got[0].PartnerCommanderName)
	assert.Equal(t, "Goblin Chieftain", *got[0].PartnerCommanderName)
}

func TestGetGameIDForResult_Found(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByIDFn: func(ctx context.Context, resultID int) (*gameresultrepo.Model, error) {
			return &gameresultrepo.Model{GormModelBase: base.GormModelBase{ID: resultID}, GameID: 7}, nil
		},
	}
	fn := GetGameIDForResult(repo)
	gameID, err := fn(context.Background(), 5)
	require.NoError(t, err)
	assert.Equal(t, 7, gameID)
}

func TestGetGameIDForResult_NotFound(t *testing.T) {
	repo := &mockGameResultRepo{
		GetByIDFn: func(ctx context.Context, resultID int) (*gameresultrepo.Model, error) {
			return nil, nil
		},
	}
	fn := GetGameIDForResult(repo)
	_, err := fn(context.Background(), 99)
	assert.ErrorContains(t, err, "not found")
}
