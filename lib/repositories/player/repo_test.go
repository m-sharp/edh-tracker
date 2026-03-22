package player_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetAll(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
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
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
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
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	got, err := repo.GetById(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByName_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	got, err := repo.GetByName(ctx, "Alice")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
}

func TestGetByName_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	got, err := repo.GetByName(context.Background(), "NoSuchPlayer")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByNames(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
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
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	ctx := context.Background()

	got, err := repo.BulkAdd(ctx, []string{"Alice", "Bob", "Carol"})
	require.NoError(t, err)
	assert.Len(t, got, 3)
	for _, m := range got {
		assert.Greater(t, m.ID, 0)
	}
}

func TestUpdate_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, "Alicia"))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Alicia", got.Name)
}

func TestUpdate_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	err := repo.Update(context.Background(), 999999, "Ghost")
	assert.ErrorContains(t, err, "unexpected number of rows")
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Alice")
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetById(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSoftDelete_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	err := repo.SoftDelete(context.Background(), 999999)
	assert.ErrorContains(t, err, "unexpected number of rows")
}

func TestSoftDelete_CascadesToAllPlayerRows(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerRepo(db)
	podRepo := testHelpers.NewPodRepo(db)
	roleRepo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	playerID, err := repo.Add(ctx, "CascadeTestPlayer")
	require.NoError(t, err)

	testDeck := testHelpers.CreateTestDeckWithCommander(t, db)
	// Reassign the deck to our player by creating a new one directly
	deckRepo := testHelpers.NewDeckRepo(db)
	formatID := testHelpers.GetCommanderFormatID(t, db)
	deckID, err := deckRepo.Add(ctx, playerID, "Player Deck", formatID)
	require.NoError(t, err)

	dcRepo := testHelpers.NewDeckCommanderRepo(db)
	cmdID := testHelpers.CreateTestCommander(t, db)
	_, err = dcRepo.Add(ctx, deckID, cmdID, nil)
	require.NoError(t, err)

	podID := testHelpers.CreateTestPod(t, db)
	require.NoError(t, podRepo.AddPlayerToPod(ctx, podID, playerID))
	require.NoError(t, roleRepo.BulkAdd(ctx, podID, []int{playerID}, playerPodRole.RoleMember))

	require.NoError(t, repo.SoftDelete(ctx, playerID))

	var deckCount, dcCount, podCount, roleCount int64
	require.NoError(t, db.Unscoped().Table("deck").
		Where("player_id = ? AND deleted_at IS NOT NULL", playerID).Count(&deckCount).Error)
	require.NoError(t, db.Unscoped().Table("deck_commander").
		Where("deck_id = ? AND deleted_at IS NOT NULL", deckID).Count(&dcCount).Error)
	require.NoError(t, db.Unscoped().Table("player_pod").
		Where("player_id = ? AND deleted_at IS NOT NULL", playerID).Count(&podCount).Error)
	require.NoError(t, db.Unscoped().Table("player_pod_role").
		Where("player_id = ? AND deleted_at IS NOT NULL", playerID).Count(&roleCount).Error)

	assert.Equal(t, int64(1), deckCount)
	assert.Equal(t, int64(1), dcCount)
	assert.Equal(t, int64(1), podCount)
	assert.Equal(t, int64(1), roleCount)

	_ = testDeck // created to ensure no cross-contamination
}
