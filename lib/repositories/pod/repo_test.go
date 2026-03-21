package pod_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetAll(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	id1, err := repo.Add(ctx, "Pod Alpha")
	require.NoError(t, err)
	id2, err := repo.Add(ctx, "Pod Beta")
	require.NoError(t, err)

	pods, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(pods), 2)

	// Soft-deleted pod should not appear
	require.NoError(t, repo.SoftDelete(ctx, id2))
	pods, err = repo.GetAll(ctx)
	require.NoError(t, err)
	for _, p := range pods {
		assert.NotEqual(t, id2, p.ID)
	}
	_ = id1
}

func TestGetByID_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Pod Alpha")
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "Pod Alpha", got.Name)
}

func TestGetByID_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)

	got, err := repo.GetByID(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByName_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Pod Alpha")
	require.NoError(t, err)

	got, err := repo.GetByName(ctx, "Pod Alpha")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
}

func TestGetByName_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)

	got, err := repo.GetByName(context.Background(), "NoSuchPod")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByPlayerID(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.AddPlayerToPod(ctx, podID, playerID))

	pods, err := repo.GetByPlayerID(ctx, playerID)
	require.NoError(t, err)
	require.Len(t, pods, 1)
	assert.Equal(t, podID, pods[0].ID)
}

func TestGetIDsByPlayerID(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.AddPlayerToPod(ctx, podID, playerID))

	ids, err := repo.GetIDsByPlayerID(ctx, playerID)
	require.NoError(t, err)
	assert.Contains(t, ids, podID)
}

func TestGetPlayerIDs(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.AddPlayerToPod(ctx, podID, playerID))

	ids, err := repo.GetPlayerIDs(ctx, podID)
	require.NoError(t, err)
	assert.Contains(t, ids, playerID)
}

func TestAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Pod Alpha")
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Pod Alpha", got.Name)
}

func TestBulkAddPlayers(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.BulkAddPlayers(ctx, podID, []int{p1, p2}))

	ids, err := repo.GetPlayerIDs(ctx, podID)
	require.NoError(t, err)
	assert.Contains(t, ids, p1)
	assert.Contains(t, ids, p2)
}

func TestAddPlayerToPod(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.AddPlayerToPod(ctx, podID, playerID))

	ids, err := repo.GetPlayerIDs(ctx, podID)
	require.NoError(t, err)
	assert.Contains(t, ids, playerID)
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Pod Alpha")
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestUpdate(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	id, err := repo.Add(ctx, "Pod Alpha")
	require.NoError(t, err)

	require.NoError(t, repo.Update(ctx, id, "Pod Beta"))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Pod Beta", got.Name)
}

func TestRemovePlayer(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.BulkAddPlayers(ctx, podID, []int{p1, p2}))

	require.NoError(t, repo.RemovePlayer(ctx, podID, p1))

	ids, err := repo.GetPlayerIDs(ctx, podID)
	require.NoError(t, err)
	assert.NotContains(t, ids, p1)
	assert.Contains(t, ids, p2)
}
