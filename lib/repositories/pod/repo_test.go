package pod_test

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

func TestGetByIDWithMembers(t *testing.T) {
	t.Run("returns pod with members", func(t *testing.T) {
		db := testHelpers.NewTestDB(t)
		repo := testHelpers.NewPodRepo(db)
		roleRepo := testHelpers.NewPlayerPodRoleRepo(db)
		ctx := context.Background()

		podID := testHelpers.CreateTestPod(t, db)
		p1 := testHelpers.CreateTestPlayer(t, db)
		p2 := testHelpers.CreateTestPlayer(t, db)

		require.NoError(t, roleRepo.BulkAdd(ctx, podID, []int{p1, p2}, playerPodRole.RoleMember))

		got, err := repo.GetByIDWithMembers(ctx, podID)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, podID, got.ID)
		require.Len(t, got.Members, 2)

		memberIDs := []int{got.Members[0].PlayerID, got.Members[1].PlayerID}
		assert.Contains(t, memberIDs, p1)
		assert.Contains(t, memberIDs, p2)
		for _, m := range got.Members {
			assert.Equal(t, playerPodRole.RoleMember, m.Role)
		}
	})

	t.Run("excludes soft-deleted member", func(t *testing.T) {
		db := testHelpers.NewTestDB(t)
		repo := testHelpers.NewPodRepo(db)
		roleRepo := testHelpers.NewPlayerPodRoleRepo(db)
		ctx := context.Background()

		podID := testHelpers.CreateTestPod(t, db)
		p1 := testHelpers.CreateTestPlayer(t, db)
		p2 := testHelpers.CreateTestPlayer(t, db)
		p3 := testHelpers.CreateTestPlayer(t, db)

		require.NoError(t, roleRepo.BulkAdd(ctx, podID, []int{p1, p2, p3}, playerPodRole.RoleMember))

		// Soft-delete p3's role row directly via GORM
		p3Role, err := roleRepo.GetRole(ctx, podID, p3)
		require.NoError(t, err)
		require.NotNil(t, p3Role)
		require.NoError(t, db.Delete(&playerPodRole.Model{}, p3Role.ID).Error)

		got, err := repo.GetByIDWithMembers(ctx, podID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got.Members, 2)
		for _, m := range got.Members {
			assert.NotEqual(t, p3, m.PlayerID)
		}
	})

	t.Run("returns nil for unknown pod", func(t *testing.T) {
		db := testHelpers.NewTestDB(t)
		repo := testHelpers.NewPodRepo(db)

		got, err := repo.GetByIDWithMembers(context.Background(), 999999)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
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

func TestRemovePlayer_CascadesToPlayerPodRole(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	roleRepo := testHelpers.NewPlayerPodRoleRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	playerID := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.AddPlayerToPod(ctx, podID, playerID))
	require.NoError(t, roleRepo.BulkAdd(ctx, podID, []int{playerID}, playerPodRole.RoleMember))

	require.NoError(t, repo.RemovePlayer(ctx, podID, playerID))

	var count int64
	require.NoError(t, db.Unscoped().Table("player_pod_role").
		Where("pod_id = ? AND player_id = ? AND deleted_at IS NOT NULL", podID, playerID).
		Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestSoftDelete_CascadesToChildTables(t *testing.T) {
	db := testHelpers.NewTestDB(t)
	repo := testHelpers.NewPodRepo(db)
	roleRepo := testHelpers.NewPlayerPodRoleRepo(db)
	inviteRepo := testHelpers.NewPodInviteRepo(db)
	ctx := context.Background()

	podID := testHelpers.CreateTestPod(t, db)
	p1 := testHelpers.CreateTestPlayer(t, db)
	p2 := testHelpers.CreateTestPlayer(t, db)

	require.NoError(t, repo.BulkAddPlayers(ctx, podID, []int{p1, p2}))
	require.NoError(t, roleRepo.BulkAdd(ctx, podID, []int{p1, p2}, playerPodRole.RoleMember))
	require.NoError(t, inviteRepo.Add(ctx, podID, p1, "test-invite-code", nil))

	require.NoError(t, repo.SoftDelete(ctx, podID))

	var playerPodCount, roleCount, inviteCount int64
	require.NoError(t, db.Unscoped().Table("player_pod").
		Where("pod_id = ? AND deleted_at IS NOT NULL", podID).Count(&playerPodCount).Error)
	require.NoError(t, db.Unscoped().Table("player_pod_role").
		Where("pod_id = ? AND deleted_at IS NOT NULL", podID).Count(&roleCount).Error)
	require.NoError(t, db.Unscoped().Table("pod_invite").
		Where("pod_id = ? AND deleted_at IS NOT NULL", podID).Count(&inviteCount).Error)

	assert.Equal(t, int64(2), playerPodCount)
	assert.Equal(t, int64(2), roleCount)
	assert.Equal(t, int64(1), inviteCount)
}
