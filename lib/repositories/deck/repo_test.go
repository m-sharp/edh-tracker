package deck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetAll(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id1, err := repo.Add(ctx, playerID, "Active Deck", 1)
	require.NoError(t, err)

	id2, err := repo.Add(ctx, playerID, "Retired Deck", 1)
	require.NoError(t, err)
	require.NoError(t, repo.Retire(ctx, id2))

	got, err := repo.GetAll(ctx)
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
	}
	assert.Contains(t, ids, id1)
	assert.NotContains(t, ids, id2)
}

func TestGetAllForPlayer(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	id1, err := repo.Add(ctx, p1, "Active", 1)
	require.NoError(t, err)
	id2, err := repo.Add(ctx, p1, "Retired", 1)
	require.NoError(t, err)
	require.NoError(t, repo.Retire(ctx, id2))

	_, err = repo.Add(ctx, p2, "Other Player Deck", 1)
	require.NoError(t, err)

	got, err := repo.GetAllForPlayer(ctx, p1)
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
	}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2) // includes retired
	for _, d := range got {
		assert.Equal(t, p1, d.PlayerID)
	}
}

func TestGetAllByPlayerIDs(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)
	p3 := testHelpers.CreateTestPlayer(t, db)

	id1, err := repo.Add(ctx, p1, "P1 Deck", 1)
	require.NoError(t, err)
	id2, err := repo.Add(ctx, p2, "P2 Deck", 1)
	require.NoError(t, err)
	id3, err := repo.Add(ctx, p3, "P3 Deck", 1)
	require.NoError(t, err)

	got, err := repo.GetAllByPlayerIDs(ctx, []int{p1, p2})
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
	}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
	assert.NotContains(t, ids, id3) // p3 not in query
}

func TestGetAllByPlayerIDs_Empty(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	got, err := repo.GetAllByPlayerIDs(context.Background(), []int{})
	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func TestGetById_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "Krenko Goblins", 1)
	require.NoError(t, err)

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, playerID, got.PlayerID)
	assert.Equal(t, "Krenko Goblins", got.Name)
	assert.Equal(t, 1, got.FormatID)
	assert.False(t, got.Retired)
}

func TestGetById_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(context.Background(), playerID, "New Deck", 1)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	decks := []deck.Model{
		{PlayerID: p1, Name: "Bulk Deck A", FormatID: 1},
		{PlayerID: p1, Name: "Bulk Deck B", FormatID: 1},
		{PlayerID: p2, Name: "Bulk Deck C", FormatID: 1},
	}
	got, err := repo.BulkAdd(ctx, decks)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	for _, d := range got {
		assert.Greater(t, d.ID, 0)
	}
}

func TestUpdate_PartialFields(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "Original Name", 1)
	require.NoError(t, err)

	newName := "Updated Name"
	require.NoError(t, repo.Update(ctx, id, deck.UpdateFields{Name: &newName}))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, 1, got.FormatID) // unchanged
}

func TestUpdate_MultipleFields(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "Original Name", 1)
	require.NoError(t, err)

	newName := "Updated Name"
	require.NoError(t, repo.Update(ctx, id, deck.UpdateFields{
		Name:    &newName,
		Retired: new(true),
	}))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
	assert.Equal(t, 1, got.FormatID) // unchanged
	assert.True(t, got.Retired)
}

func TestUpdate_NoFields(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "No Change Deck", 1)
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, deck.UpdateFields{}))
}

func TestUpdate_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	name := "Ghost"
	err := repo.Update(context.Background(), 999999, deck.UpdateFields{Name: &name})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected rows")
}

func TestRetire(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "To Retire", 1)
	require.NoError(t, err)

	require.NoError(t, repo.Retire(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.Retired)
}

func TestRetire_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	err := repo.Retire(context.Background(), 999999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected rows")
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewDeckRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, "To Delete", 1)
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}
