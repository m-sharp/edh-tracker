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
	_, err := repo.Add(ctx, "Game A", podID, formatID)
	require.NoError(t, err)
	_, err = repo.Add(ctx, "Game B", podID, formatID)
	require.NoError(t, err)

	games, err := repo.GetAllByPod(ctx, podID)
	require.NoError(t, err)
	assert.Len(t, games, 2)
}

func TestGetAllByDeck(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID := testHelpers.CreateTestDeck(t, db).ID

	testHelpers.CreateTestGameResult(t, db, gameID, deckID, 1, 0)

	games, err := repo.GetAllByDeck(ctx, deckID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
}

func TestGetAllByPlayerID(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	testDeck := testHelpers.CreateTestDeck(t, db)

	testHelpers.CreateTestGameResult(t, db, gameID, testDeck.ID, 1, 0)

	playerID := testDeck.PlayerID

	games, err := repo.GetAllByPlayerID(ctx, playerID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
}

func TestGetById_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	id, err := repo.Add(ctx, "Friday Night", podID, formatID)
	require.NoError(t, err)

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "Friday Night", got.Description)
	assert.Equal(t, podID, got.PodID)
}

func TestGetById_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)

	got, err := repo.GetById(context.Background(), 999999)
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

	got, err := repo.GetById(ctx, id)
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

	got, err := repo.GetById(ctx, id)
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

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetAllByPodWithResults(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	gameID, err := repo.Add(ctx, "Game With Results", podID, formatID)
	require.NoError(t, err)
	deckID := testHelpers.CreateTestDeck(t, db).ID
	testHelpers.CreateTestGameResult(t, db, gameID, deckID, 1, 2)

	games, err := repo.GetAllByPodWithResults(ctx, podID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
	require.Len(t, games[0].Results, 1)
	assert.Equal(t, gameID, games[0].Results[0].GameID)
	assert.Equal(t, deckID, games[0].Results[0].DeckID)
	assert.Equal(t, 1, games[0].Results[0].Place)
	assert.Equal(t, 2, games[0].Results[0].KillCount)

	// Add another deck result
	deck2ID := testHelpers.CreateTestDeck(t, db).ID
	testHelpers.CreateTestGameResult(t, db, gameID, deck2ID, 2, 1)
	games, err = repo.GetAllByPodWithResults(ctx, podID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	require.Len(t, games[0].Results, 2)
}

func TestGetAllByDeckWithResults(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID := testHelpers.CreateTestDeck(t, db).ID
	testHelpers.CreateTestGameResult(t, db, gameID, deckID, 1, 0)

	games, err := repo.GetAllByDeckWithResults(ctx, deckID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
	require.Len(t, games[0].Results, 1)
	assert.Equal(t, gameID, games[0].Results[0].GameID)
	assert.Equal(t, deckID, games[0].Results[0].DeckID)
}

func TestGetAllByPlayerWithResults(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	testDeck := testHelpers.CreateTestDeck(t, db)
	testHelpers.CreateTestGameResult(t, db, gameID, testDeck.ID, 1, 0)

	games, err := repo.GetAllByPlayerWithResults(ctx, testDeck.PlayerID)
	require.NoError(t, err)
	require.Len(t, games, 1)
	assert.Equal(t, gameID, games[0].ID)
	require.Len(t, games[0].Results, 1)
	assert.Equal(t, gameID, games[0].Results[0].GameID)
	assert.Equal(t, testDeck.ID, games[0].Results[0].DeckID)
}

func TestGetByIDWithResults_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)
	ctx := context.Background()
	formatID := testHelpers.GetCommanderFormatID(t, db)

	podID := testHelpers.CreateTestPod(t, db)
	gameID, err := repo.Add(ctx, "Friday Night", podID, formatID)
	require.NoError(t, err)
	deckID := testHelpers.CreateTestDeck(t, db).ID
	testHelpers.CreateTestGameResult(t, db, gameID, deckID, 2, 1)

	got, err := repo.GetByIDWithResults(ctx, gameID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, gameID, got.ID)
	assert.Equal(t, "Friday Night", got.Description)
	require.Len(t, got.Results, 1)
	assert.Equal(t, gameID, got.Results[0].GameID)
	assert.Equal(t, deckID, got.Results[0].DeckID)
}

func TestGetByIDWithResults_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameRepo(db)

	got, err := repo.GetByIDWithResults(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}
