package deck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

func newRepo(t *testing.T) *Repository {
	t.Helper()
	db := base.NewTestDB(t)
	return &Repository{db: db}
}

func TestGetAll(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	formatID := 1
	// Add one active deck and one retired deck
	id1, err := repo.Add(ctx, 1, "Active Deck", formatID)
	require.NoError(t, err)

	id2, err := repo.Add(ctx, 1, "Retired Deck", formatID)
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
	repo := newRepo(t)
	ctx := context.Background()

	// Add active and retired decks for player 1
	id1, err := repo.Add(ctx, 1, "Active", 1)
	require.NoError(t, err)
	id2, err := repo.Add(ctx, 1, "Retired", 1)
	require.NoError(t, err)
	require.NoError(t, repo.Retire(ctx, id2))

	// Add deck for player 2
	_, err = repo.Add(ctx, 2, "Other Player Deck", 1)
	require.NoError(t, err)

	got, err := repo.GetAllForPlayer(ctx, 1)
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
	}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2) // includes retired
	for _, d := range got {
		assert.Equal(t, 1, d.PlayerID)
	}
}

func TestGetAllByPlayerIDs(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id1, err := repo.Add(ctx, 1, "P1 Deck", 1)
	require.NoError(t, err)
	id2, err := repo.Add(ctx, 2, "P2 Deck", 1)
	require.NoError(t, err)
	id3, err := repo.Add(ctx, 3, "P3 Deck", 1)
	require.NoError(t, err)

	got, err := repo.GetAllByPlayerIDs(ctx, []int{1, 2})
	require.NoError(t, err)

	var ids []int
	for _, d := range got {
		ids = append(ids, d.ID)
	}
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
	assert.NotContains(t, ids, id3) // player 3 not in query
}

func TestGetAllByPlayerIDs_Empty(t *testing.T) {
	repo := newRepo(t)
	got, err := repo.GetAllByPlayerIDs(context.Background(), []int{})
	require.NoError(t, err)
	assert.Len(t, got, 0)
}

func TestGetById_Found(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, 1, "Krenko Goblins", 1)
	require.NoError(t, err)

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, 1, got.PlayerID)
	assert.Equal(t, "Krenko Goblins", got.Name)
	assert.Equal(t, 1, got.FormatID)
	assert.False(t, got.Retired)
}

func TestGetById_NotFound(t *testing.T) {
	repo := newRepo(t)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestAdd(t *testing.T) {
	repo := newRepo(t)
	id, err := repo.Add(context.Background(), 1, "New Deck", 1)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	decks := []Model{
		{PlayerID: 1, Name: "Bulk Deck A", FormatID: 1},
		{PlayerID: 1, Name: "Bulk Deck B", FormatID: 1},
		{PlayerID: 2, Name: "Bulk Deck C", FormatID: 1},
	}
	got, err := repo.BulkAdd(ctx, decks)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	for _, d := range got {
		assert.Greater(t, d.ID, 0)
	}
}

func TestUpdate_PartialFields(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, 1, "Original Name", 1)
	require.NoError(t, err)

	newName := "Updated Name"
	require.NoError(t, repo.Update(ctx, id, UpdateFields{Name: &newName}))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, 1, got.FormatID) // unchanged
}

func TestUpdate_NoFields(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, 1, "No Change Deck", 1)
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, UpdateFields{}))
}

func TestUpdate_NotFound(t *testing.T) {
	repo := newRepo(t)
	name := "Ghost"
	err := repo.Update(context.Background(), 999999, UpdateFields{Name: &name})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected rows")
}

func TestRetire(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, 1, "To Retire", 1)
	require.NoError(t, err)

	require.NoError(t, repo.Retire(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.Retired)
}

func TestRetire_NotFound(t *testing.T) {
	repo := newRepo(t)
	err := repo.Retire(context.Background(), 999999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected rows")
}

func TestSoftDelete(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, 1, "To Delete", 1)
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}
