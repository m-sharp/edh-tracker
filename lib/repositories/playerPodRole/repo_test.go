package playerPodRole_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	"github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func TestGetRole_Found(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.SetRole(ctx, podID, playerID, playerPodRole.RoleManager))

	got, err := repo.GetRole(ctx, podID, playerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, playerPodRole.RoleManager, got.Role)
	assert.Equal(t, podID, got.PodID)
	assert.Equal(t, playerID, got.PlayerID)
}

func TestGetRole_NotFound(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)

	got, err := repo.GetRole(context.Background(), 999999, 999999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSetRole_Insert(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.SetRole(ctx, podID, playerID, playerPodRole.RoleMember))

	got, err := repo.GetRole(ctx, podID, playerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, playerPodRole.RoleMember, got.Role)
}

func TestSetRole_Update(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.SetRole(ctx, podID, playerID, playerPodRole.RoleMember))
	require.NoError(t, repo.SetRole(ctx, podID, playerID, playerPodRole.RoleManager))

	got, err := repo.GetRole(ctx, podID, playerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, playerPodRole.RoleManager, got.Role)
}

func TestSetRole_Restore(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.SetRole(ctx, podID, playerID, playerPodRole.RoleMember))

	// Soft-delete the row manually
	require.NoError(t, db.Where("pod_id = ? AND player_id = ?", podID, playerID).Delete(&playerPodRole.Model{}).Error)

	// SetRole should restore the soft-deleted row
	require.NoError(t, repo.SetRole(ctx, podID, playerID, playerPodRole.RoleMember))

	got, err := repo.GetRole(ctx, podID, playerID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, playerPodRole.RoleMember, got.Role)
}

func TestGetMembersWithRoles(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)
	p3 := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.SetRole(ctx, podID, p1, playerPodRole.RoleManager))
	require.NoError(t, repo.SetRole(ctx, podID, p2, playerPodRole.RoleMember))
	require.NoError(t, repo.SetRole(ctx, podID, p3, playerPodRole.RoleMember))

	// Soft-delete p3 — should not appear in results
	require.NoError(t, db.Where("pod_id = ? AND player_id = ?", podID, p3).Delete(&playerPodRole.Model{}).Error)

	members, err := repo.GetMembersWithRoles(ctx, podID)
	require.NoError(t, err)
	assert.Len(t, members, 2)

	playerIDs := []int{members[0].PlayerID, members[1].PlayerID}
	assert.Contains(t, playerIDs, p1)
	assert.Contains(t, playerIDs, p2)
}

func TestBulkAdd(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.BulkAdd(ctx, podID, []int{p1, p2}, playerPodRole.RoleMember))

	members, err := repo.GetMembersWithRoles(ctx, podID)
	require.NoError(t, err)
	assert.Len(t, members, 2)

	playerIDs := []int{members[0].PlayerID, members[1].PlayerID}
	assert.Contains(t, playerIDs, p1)
	assert.Contains(t, playerIDs, p2)
}
