package format

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

func newRepo(t *testing.T) (*Repository, *gorm.DB) {
	t.Helper()
	db := base.NewTestDB(t)
	return &Repository{db: db}, db
}

func TestGetAll(t *testing.T) {
	repo, db := newRepo(t)
	ctx := context.Background()

	db.Create(&Model{Name: "FormatA"})
	db.Create(&Model{Name: "FormatB"})

	formats, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(formats), 2)
}

func TestGetById_Found(t *testing.T) {
	repo, db := newRepo(t)
	ctx := context.Background()

	m := Model{Name: "FormatC"}
	db.Create(&m)
	require.Greater(t, m.ID, 0)

	got, err := repo.GetById(ctx, m.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, "FormatC", got.Name)
}

func TestGetById_NotFound(t *testing.T) {
	repo, _ := newRepo(t)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByName_Found(t *testing.T) {
	repo, db := newRepo(t)
	ctx := context.Background()

	m := Model{Name: "FormatD"}
	db.Create(&m)

	got, err := repo.GetByName(ctx, "FormatD")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, m.ID, got.ID)
}

func TestGetByName_NotFound(t *testing.T) {
	repo, _ := newRepo(t)
	got, err := repo.GetByName(context.Background(), "NoSuchFormat")
	require.NoError(t, err)
	assert.Nil(t, got)
}
