package deckCommander_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetByDeckId_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	deckID := testHelpers.CreateTestDeck(t, db).ID
	got, err := repo.GetByDeckId(context.Background(), deckID)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByDeckId_Found_WithoutPartner(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	deckID := testHelpers.CreateTestDeck(t, db).ID
	commanderID := testHelpers.CreateTestCommander(t, db)

	id, err := repo.Add(ctx, deckID, commanderID, nil)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByDeckId(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, deckID, got.DeckID)
	assert.Equal(t, commanderID, got.CommanderID)
	assert.Nil(t, got.PartnerCommanderID)
}

func TestGetByDeckId_Found_WithPartner(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	deckID := testHelpers.CreateTestDeck(t, db).ID
	commanderID := testHelpers.CreateTestCommander(t, db)
	partnerID := testHelpers.CreateTestCommander(t, db)

	id, err := repo.Add(ctx, deckID, commanderID, &partnerID)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByDeckId(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, deckID, got.DeckID)
	assert.Equal(t, commanderID, got.CommanderID)
	require.NotNil(t, got.PartnerCommanderID)
	assert.Equal(t, partnerID, *got.PartnerCommanderID)
}

func TestAdd_WithoutPartner(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	deckID := testHelpers.CreateTestDeck(t, db).ID
	commanderID := testHelpers.CreateTestCommander(t, db)

	id, err := repo.Add(ctx, deckID, commanderID, nil)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByDeckId(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, deckID, got.DeckID)
	assert.Equal(t, commanderID, got.CommanderID)
	assert.Nil(t, got.PartnerCommanderID)
}

func TestAdd_WithPartner(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	deckID := testHelpers.CreateTestDeck(t, db).ID
	commanderID := testHelpers.CreateTestCommander(t, db)
	partnerID := testHelpers.CreateTestCommander(t, db)

	id, err := repo.Add(ctx, deckID, commanderID, &partnerID)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByDeckId(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, deckID, got.DeckID)
	assert.Equal(t, commanderID, got.CommanderID)
	require.NotNil(t, got.PartnerCommanderID)
	assert.Equal(t, partnerID, *got.PartnerCommanderID)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	entries := []deckCommander.Model{
		{DeckID: testHelpers.CreateTestDeck(t, db).ID, CommanderID: testHelpers.CreateTestCommander(t, db)},
		{DeckID: testHelpers.CreateTestDeck(t, db).ID, CommanderID: testHelpers.CreateTestCommander(t, db)},
		{DeckID: testHelpers.CreateTestDeck(t, db).ID, CommanderID: testHelpers.CreateTestCommander(t, db)},
	}
	require.NoError(t, repo.BulkAdd(ctx, entries))

	for _, e := range entries {
		got, err := repo.GetByDeckId(ctx, e.DeckID)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, e.CommanderID, got.CommanderID)
	}
}

func TestDeleteByDeckID(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckCommanderRepo(db)
	ctx := context.Background()

	deckID := testHelpers.CreateTestDeck(t, db).ID
	commanderID := testHelpers.CreateTestCommander(t, db)

	_, err := repo.Add(ctx, deckID, commanderID, nil)
	require.NoError(t, err)

	got, err := repo.GetByDeckId(ctx, deckID)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.NoError(t, repo.DeleteByDeckID(ctx, deckID))

	got, err = repo.GetByDeckId(ctx, deckID)
	require.NoError(t, err)
	assert.Nil(t, got)
}
