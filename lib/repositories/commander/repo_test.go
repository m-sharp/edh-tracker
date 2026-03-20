package commander

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

func TestGetById_Found(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Atraxa")
	require.NoError(t, err)

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "Atraxa", got.Name)
}

func TestGetById_NotFound(t *testing.T) {
	repo := newRepo(t)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByName_Found(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Atraxa")
	require.NoError(t, err)

	got, err := repo.GetByName(ctx, "Atraxa")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
}

func TestGetByName_NotFound(t *testing.T) {
	repo := newRepo(t)
	got, err := repo.GetByName(context.Background(), "NoSuchCommander")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByNames(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	_, err := repo.Add(ctx, "Atraxa")
	require.NoError(t, err)
	_, err = repo.Add(ctx, "Najeela")
	require.NoError(t, err)

	// Both match
	got, err := repo.GetByNames(ctx, []string{"Atraxa", "Najeela"})
	require.NoError(t, err)
	assert.Len(t, got, 2)

	// Partial match
	got, err = repo.GetByNames(ctx, []string{"Atraxa"})
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "Atraxa", got[0].Name)
}

func TestAdd(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Atraxa")
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	got, err := repo.BulkAdd(ctx, []string{"Atraxa", "Najeela", "Edgar"})
	require.NoError(t, err)
	assert.Len(t, got, 3)
	for _, m := range got {
		assert.Greater(t, m.ID, 0)
	}
}
