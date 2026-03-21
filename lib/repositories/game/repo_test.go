package game_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/game"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetAllByPod(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	gameID, err := repo.Add(ctx, "Game With Results", podID, formatID)
	require.NoError(t, err)
	testDeck := testHelpers.CreateTestDeck(t, db)
	testHelpers.CreateTestGameResult(t, db, gameID, testDeck.ID, 1, 2)

	games, err := repo.GetAllByPod(ctx, podID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
	require.Len(t, games[0].Results, 1)
	assert.Equal(t, gameID, games[0].Results[0].GameID)
	assert.Equal(t, testDeck.ID, games[0].Results[0].DeckID)
	assert.Equal(t, 1, games[0].Results[0].Place)
	assert.Equal(t, 2, games[0].Results[0].KillCount)
	assert.Equal(t, testDeck.Name, games[0].Results[0].Deck.Name)
	assert.Equal(t, testDeck.PlayerID, games[0].Results[0].Deck.PlayerID)

	// Add another deck result
	deck2ID := testHelpers.CreateTestDeck(t, db).ID
	testHelpers.CreateTestGameResult(t, db, gameID, deck2ID, 2, 1)
	games, err = repo.GetAllByPod(ctx, podID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	require.Len(t, games[0].Results, 2)
}

func TestGetAllByDeck(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	testDeck := testHelpers.CreateTestDeck(t, db)
	testHelpers.CreateTestGameResult(t, db, gameID, testDeck.ID, 1, 0)

	games, err := repo.GetAllByDeck(ctx, testDeck.ID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
	require.Len(t, games[0].Results, 1)
	assert.Equal(t, gameID, games[0].Results[0].GameID)
	assert.Equal(t, testDeck.ID, games[0].Results[0].DeckID)
	assert.Equal(t, testDeck.Name, games[0].Results[0].Deck.Name)
	assert.Equal(t, testDeck.PlayerID, games[0].Results[0].Deck.PlayerID)
}

func TestGetAllByPlayerID(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	testDeck := testHelpers.CreateTestDeck(t, db)
	testHelpers.CreateTestGameResult(t, db, gameID, testDeck.ID, 1, 0)

	games, err := repo.GetAllByPlayerID(ctx, testDeck.PlayerID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
	require.Len(t, games[0].Results, 1)
	assert.Equal(t, gameID, games[0].Results[0].GameID)
	assert.Equal(t, testDeck.ID, games[0].Results[0].DeckID)
	assert.Equal(t, testDeck.Name, games[0].Results[0].Deck.Name)
	assert.Equal(t, testDeck.PlayerID, games[0].Results[0].Deck.PlayerID)
}

func TestGetByID_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	gameID, err := repo.Add(ctx, "Friday Night", podID, formatID)
	require.NoError(t, err)
	testDeck := testHelpers.CreateTestDeck(t, db)
	testHelpers.CreateTestGameResult(t, db, gameID, testDeck.ID, 2, 1)

	got, err := repo.GetByID(ctx, gameID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, gameID, got.ID)
	assert.Equal(t, "Friday Night", got.Description)
	require.Len(t, got.Results, 1)
	assert.Equal(t, gameID, got.Results[0].GameID)
	assert.Equal(t, testDeck.ID, got.Results[0].DeckID)
	assert.Equal(t, testDeck.Name, got.Results[0].Deck.Name)
	assert.Equal(t, testDeck.PlayerID, got.Results[0].Deck.PlayerID)
}

func TestGetByID_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)

	got, err := repo.GetByID(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	id, err := repo.Add(ctx, "Test Game", podID, formatID)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Test Game", got.Description)
	assert.Equal(t, podID, got.PodID)
	assert.Equal(t, formatID, got.FormatID)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	games := []game.Model{
		{Description: fmt.Sprintf("Bulk Game %d", 1), PodID: podID, FormatID: formatID},
		{Description: fmt.Sprintf("Bulk Game %d", 2), PodID: podID, FormatID: formatID},
		{Description: fmt.Sprintf("Bulk Game %d", 3), PodID: podID, FormatID: formatID},
	}

	ids, err := repo.BulkAdd(ctx, games)
	require.NoError(t, err)
	require.Len(t, ids, 3)

	seen := make(map[int]bool)
	for _, id := range ids {
		assert.Greater(t, id, 0)
		assert.False(t, seen[id], "duplicate ID returned: %d", id)
		seen[id] = true
	}
}

func TestUpdate(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	id, err := repo.Add(ctx, "Old Description", podID, formatID)
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, "New Description"))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "New Description", got.Description)
}

func TestUpdate_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)

	err := repo.Update(context.Background(), 999999, "New Description")
	assert.ErrorContains(t, err, "unexpected number of rows")
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	id, err := repo.Add(ctx, "To Delete", podID, formatID)
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}
