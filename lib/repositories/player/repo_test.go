package player

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

	_, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)
	id2, err := repo.Add(ctx, "Bob")
	require.NoError(t, err)

	players, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(players), 2)

	// Soft-deleted player should not appear
	require.NoError(t, repo.SoftDelete(ctx, id2))
	players, err = repo.GetAll(ctx)
	require.NoError(t, err)
	for _, p := range players {
		assert.NotEqual(t, id2, p.ID)
	}
}

func TestGetById_Found(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "Alice", got.Name)
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

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	got, err := repo.GetByName(ctx, "Alice")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
}

func TestGetByName_NotFound(t *testing.T) {
	repo := newRepo(t)
	got, err := repo.GetByName(context.Background(), "NoSuchPlayer")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByNames(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	_, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)
	_, err = repo.Add(ctx, "Bob")
	require.NoError(t, err)

	got, err := repo.GetByNames(ctx, []string{"Alice", "Bob"})
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestAdd(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	got, err := repo.BulkAdd(ctx, []string{"Alice", "Bob", "Carol"})
	require.NoError(t, err)
	assert.Len(t, got, 3)
	for _, m := range got {
		assert.Greater(t, m.ID, 0)
	}
}

func TestUpdate_Found(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, "Alicia"))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Alicia", got.Name)
}

func TestUpdate_NotFound(t *testing.T) {
	repo := newRepo(t)
	err := repo.Update(context.Background(), 999999, "Ghost")
	assert.ErrorContains(t, err, "unexpected number of rows")
}

func TestSoftDelete(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSoftDelete_NotFound(t *testing.T) {
	repo := newRepo(t)
	err := repo.SoftDelete(context.Background(), 999999)
	assert.ErrorContains(t, err, "unexpected number of rows")
}
