package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

const (
	testProvider    = "google"
	testSubject     = "oauth-subject-123"
	testEmail       = "test@example.com"
	testDisplayName = "Test User"
	testAvatarURL   = "https://example.com/avatar.png"
)

func TestGetByID_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, 2)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, playerID, got.PlayerID)
	assert.Equal(t, 2, got.RoleID)
}

func TestGetByID_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)

	got, err := repo.GetByID(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByPlayerID_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, 2)
	require.NoError(t, err)

	got, err := repo.GetByPlayerID(ctx, playerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, playerID, got.PlayerID)
}

func TestGetByPlayerID_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)

	got, err := repo.GetByPlayerID(context.Background(), 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByOAuth_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.AddWithOAuth(ctx, playerID, 2, testProvider, testSubject, testEmail, testDisplayName, testAvatarURL)
	require.NoError(t, err)

	got, err := repo.GetByOAuth(ctx, testProvider, testSubject)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, testProvider, *got.OAuthProvider)
	assert.Equal(t, testSubject, *got.OAuthSubject)
}

func TestGetByOAuth_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)

	got, err := repo.GetByOAuth(context.Background(), "no-provider", "no-subject")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetRoleByName_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)

	got, err := repo.GetRoleByName(context.Background(), "player")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "player", got.Name)
}

func TestGetRoleByName_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)

	got, err := repo.GetRoleByName(context.Background(), "does-not-exist")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, 2)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, playerID, got.PlayerID)
	assert.Equal(t, 2, got.RoleID)
	assert.Nil(t, got.OAuthProvider)
}

func TestAddWithOAuth(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.AddWithOAuth(ctx, playerID, 2, testProvider, testSubject, testEmail, testDisplayName, testAvatarURL)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, playerID, got.PlayerID)
	require.NotNil(t, got.OAuthProvider)
	assert.Equal(t, testProvider, *got.OAuthProvider)
	require.NotNil(t, got.OAuthSubject)
	assert.Equal(t, testSubject, *got.OAuthSubject)
	require.NotNil(t, got.Email)
	assert.Equal(t, testEmail, *got.Email)
	require.NotNil(t, got.DisplayName)
	assert.Equal(t, testDisplayName, *got.DisplayName)
	require.NotNil(t, got.AvatarURL)
	assert.Equal(t, testAvatarURL, *got.AvatarURL)
}

func TestCreatePlayerAndUser(t *testing.T) {
	// MySQL doesn't support nested transactions, so we can't use NewTestDB's wrapping transaction here.
	// Use a plain connection and clean up explicitly.
	db := testHelpers.NewTestDBNoTx(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerName := fmt.Sprintf("CreatePU-Player-%s", uuid.NewString())

	got, err := repo.CreatePlayerAndUser(ctx, playerName, 2, testProvider, testSubject+"-cpu", testEmail, testDisplayName, testAvatarURL)
	require.NoError(t, err)
	require.NotNil(t, got)

	t.Cleanup(func() {
		db.Exec("DELETE FROM user WHERE id = ?", got.ID)
		db.Exec("DELETE FROM player WHERE id = ?", got.PlayerID)
	})

	assert.Greater(t, got.ID, 0)
	assert.Greater(t, got.PlayerID, 0)
	assert.Equal(t, 2, got.RoleID)
	require.NotNil(t, got.OAuthProvider)
	assert.Equal(t, testProvider, *got.OAuthProvider)

	// Verify user is readable
	fetched, err := repo.GetByID(ctx, got.ID)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, got.PlayerID, fetched.PlayerID)
}

func TestCreatePlayerAndUser_Rollback(t *testing.T) {
	db := testHelpers.NewTestDBNoTx(t)
	repo := testHelpers.NewUserRepo(db)
	playerRepo := testHelpers.NewPlayerRepo(db)
	ctx := context.Background()

	playerName := fmt.Sprintf("Rollback-Player-%s", uuid.NewString())

	// roleID 9999 doesn't exist → FK violation → transaction must roll back
	_, err := repo.CreatePlayerAndUser(ctx, playerName, 9999, testProvider, testSubject+"-rb", testEmail, testDisplayName, testAvatarURL)
	require.Error(t, err)

	// Verify player was not committed
	players, err := playerRepo.GetByNames(ctx, []string{playerName})
	require.NoError(t, err)
	assert.Empty(t, players)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)
	p3 := testHelpers.CreateTestPlayer(t, db)

	err := repo.BulkAdd(ctx, []int{p1, p2, p3}, 2)
	require.NoError(t, err)

	for _, pid := range []int{p1, p2, p3} {
		got, err := repo.GetByPlayerID(ctx, pid)
		require.NoError(t, err)
		require.NotNil(t, got, "expected user for player %d", pid)
		assert.Equal(t, 2, got.RoleID)
	}
}

func TestBulkAdd_Empty(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)

	err := repo.BulkAdd(context.Background(), []int{}, 2)
	require.NoError(t, err)
}

func TestSoftDelete(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewUserRepo(db)
	ctx := context.Background()

	playerID := testHelpers.CreateTestPlayer(t, db)
	id, err := repo.Add(ctx, playerID, 2)
	require.NoError(t, err)

	require.NoError(t, repo.SoftDelete(ctx, id))

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Nil(t, got, "soft-deleted user should not be returned by GetByID")
}
