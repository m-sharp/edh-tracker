package format_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/format"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetAll(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewFormatRepo(db)
	ctx := context.Background()

	db.Create(&format.Model{Name: "FormatA"})
	db.Create(&format.Model{Name: "FormatB"})

	formats, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(formats), 2)
}

func TestGetById_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewFormatRepo(db)
	ctx := context.Background()

	m := format.Model{Name: "FormatC"}
	db.Create(&m)
	require.Greater(t, m.ID, 0)

	got, err := repo.GetById(ctx, m.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, "FormatC", got.Name)
}

func TestGetById_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewFormatRepo(db)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByName_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewFormatRepo(db)
	ctx := context.Background()

	m := format.Model{Name: "FormatD"}
	db.Create(&m)

	got, err := repo.GetByName(ctx, "FormatD")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, m.ID, got.ID)
}

func TestGetByName_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewFormatRepo(db)
	got, err := repo.GetByName(context.Background(), "NoSuchFormat")
	require.NoError(t, err)
	assert.Nil(t, got)
}
