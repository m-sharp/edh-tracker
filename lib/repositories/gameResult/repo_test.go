package gameResult_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetByGameId(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID1 := testHelpers.CreateTestDeck(t, db).ID
	deckID2 := testHelpers.CreateTestDeck(t, db).ID

	_, err := repo.Add(ctx, gameResult.Model{GameID: gameID, DeckID: deckID1, Place: 1, KillCount: 2})
	require.NoError(t, err)
	_, err = repo.Add(ctx, gameResult.Model{GameID: gameID, DeckID: deckID2, Place: 2, KillCount: 0})
	require.NoError(t, err)

	got, err := repo.GetByGameId(ctx, gameID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestGetByGameId_Empty(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)

	got, err := repo.GetByGameId(context.Background(), 999999)
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 0)
}

func TestGetByID_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID := testHelpers.CreateTestDeck(t, db).ID

	id, err := repo.Add(ctx, gameResult.Model{GameID: gameID, DeckID: deckID, Place: 1, KillCount: 1})
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, gameID, got.GameID)
	assert.Equal(t, deckID, got.DeckID)
	assert.Equal(t, 1, got.Place)
	assert.Equal(t, 1, got.KillCount)
}

func TestGetByID_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)

	got, err := repo.GetByID(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID := testHelpers.CreateTestDeck(t, db).ID

	id, err := repo.Add(ctx, gameResult.Model{GameID: gameID, DeckID: deckID, Place: 2, KillCount: 1})
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID1 := testHelpers.CreateTestDeck(t, db).ID
	deckID2 := testHelpers.CreateTestDeck(t, db).ID

	results := []gameResult.Model{
		{GameID: gameID, DeckID: deckID1, Place: 1, KillCount: 2},
		{GameID: gameID, DeckID: deckID2, Place: 2, KillCount: 0},
	}
	err := repo.BulkAdd(ctx, results)
	require.NoError(t, err)

	got, err := repo.GetByGameId(ctx, gameID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestUpdate(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	gameID := testHelpers.CreateTestGame(t, db)
	deckID := testHelpers.CreateTestDeck(t, db).ID

	id, err := repo.Add(ctx, gameResult.Model{GameID: gameID, DeckID: deckID, Place: 2, KillCount: 0})
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, 1, 3, deckID))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 1, got.Place)
	assert.Equal(t, 3, got.KillCount)
}

func TestUpdate_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)

	err := repo.Update(context.Background(), 999999, 1, 0, 1)
	assert.ErrorContains(t, err, "unexpected number of rows")
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	id := testHelpers.CreateTestGameResult(t, db, testHelpers.CreateTestGame(t, db), testHelpers.CreateTestDeck(t, db).ID, 1, 0)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetStatsForPlayer(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	// Create a player with a deck
	testDeck := testHelpers.CreateTestDeck(t, db)
	playerID := testDeck.PlayerID

	// Three games for that player's deck
	game1 := testHelpers.CreateTestGame(t, db)
	game2 := testHelpers.CreateTestGame(t, db)
	game3 := testHelpers.CreateTestGame(t, db)

	_, err := repo.Add(ctx, gameResult.Model{GameID: game1, DeckID: testDeck.ID, Place: 1, KillCount: 2})
	require.NoError(t, err)
	_, err = repo.Add(ctx, gameResult.Model{GameID: game2, DeckID: testDeck.ID, Place: 2, KillCount: 0})
	require.NoError(t, err)
	_, err = repo.Add(ctx, gameResult.Model{GameID: game3, DeckID: testDeck.ID, Place: 1, KillCount: 1})
	require.NoError(t, err)

	got, err := repo.GetStatsForPlayer(ctx, playerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 3, got.Games)
	assert.Equal(t, 3, got.Kills)
	assert.Equal(t, map[int]int{1: 2, 2: 1}, got.Record)
}

func TestGetStatsForPlayer_Empty(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)

	got, err := repo.GetStatsForPlayer(context.Background(), 999999)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Games)
	assert.Len(t, got.Record, 0)
}

func TestGetStatsForDeck(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)
	ctx := context.Background()

	deckID := testHelpers.CreateTestDeck(t, db).ID
	gameID := testHelpers.CreateTestGame(t, db)

	_, err := repo.Add(ctx, gameResult.Model{GameID: gameID, DeckID: deckID, Place: 1, KillCount: 2})
	require.NoError(t, err)

	got, err := repo.GetStatsForDeck(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 1, got.Games)
	assert.Equal(t, 2, got.Kills)
}

func TestGetStatsForDeck_Empty(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewGameResultRepo(db)

	got, err := repo.GetStatsForDeck(context.Background(), 999999)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Games)
	assert.Len(t, got.Record, 0)
}
