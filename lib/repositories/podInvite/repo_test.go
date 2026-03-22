package podInvite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetByCode_Found_WithExpiry(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodInviteRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	future := time.Now().Add(7 * 24 * time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, repo.Add(ctx, podID, playerID, "code-with-expiry", &future))

	got, err := repo.GetByCode(ctx, "code-with-expiry")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, podID, got.PodID)
	assert.Equal(t, playerID, got.CreatedByPlayerID)
	assert.Equal(t, "code-with-expiry", got.InviteCode)
	assert.Equal(t, 0, got.UsedCount)
	require.NotNil(t, got.ExpiresAt)
	assert.Equal(t, future, got.ExpiresAt.UTC().Truncate(time.Second))
}

func TestGetByCode_Found_NoExpiry(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodInviteRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.Add(ctx, podID, playerID, "code-no-expiry", nil))

	got, err := repo.GetByCode(ctx, "code-no-expiry")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "code-no-expiry", got.InviteCode)
	assert.Nil(t, got.ExpiresAt)
}

func TestGetByCode_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodInviteRepo(db)

	got, err := repo.GetByCode(context.Background(), "does-not-exist")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestAdd_WithExpiry(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodInviteRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	future := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, repo.Add(ctx, podID, playerID, "add-with-expiry", &future))

	got, err := repo.GetByCode(ctx, "add-with-expiry")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.ExpiresAt)
	assert.Equal(t, future, got.ExpiresAt.UTC().Truncate(time.Second))
}

func TestAdd_NoExpiry(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodInviteRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.Add(ctx, podID, playerID, "add-no-expiry", nil))

	got, err := repo.GetByCode(ctx, "add-no-expiry")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Nil(t, got.ExpiresAt)
}

func TestIncrementUsedCount(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodInviteRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.Add(ctx, podID, playerID, "inc-code", nil))

	require.NoError(t, repo.IncrementUsedCount(ctx, "inc-code"))
	got, err := repo.GetByCode(ctx, "inc-code")
	require.NoError(t, err)
	assert.Equal(t, 1, got.UsedCount)

	require.NoError(t, repo.IncrementUsedCount(ctx, "inc-code"))
	got, err = repo.GetByCode(ctx, "inc-code")
	require.NoError(t, err)
	assert.Equal(t, 2, got.UsedCount)
}
