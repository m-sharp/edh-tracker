package pod

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	podrepo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
	"github.com/m-sharp/edh-tracker/lib/repositories/podInvite"
)

// --- Create ---

func TestCreate_Success(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		AddFn: func(ctx context.Context, name string) (int, error) { return 5, nil },
	}
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error { return nil },
	}

	fn := Create(podRepo, roleRepo)
	podID, err := fn(context.Background(), "My Pod", 1)
	require.NoError(t, err)
	assert.Equal(t, 5, podID)
}

func TestCreate_SetRoleError(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		AddFn: func(ctx context.Context, name string) (int, error) { return 5, nil },
	}
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error {
			return errors.New("db error")
		},
	}

	fn := Create(podRepo, roleRepo)
	_, err := fn(context.Background(), "My Pod", 1)
	assert.Error(t, err)
}

// --- AddPlayer ---

func TestAddPlayer_Success(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		AddPlayerToPodFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error { return nil },
	}

	fn := AddPlayer(podRepo, roleRepo)
	err := fn(context.Background(), 1, 2)
	require.NoError(t, err)
}

func TestAddPlayer_SetRoleError(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		AddPlayerToPodFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error {
			return errors.New("db error")
		},
	}

	fn := AddPlayer(podRepo, roleRepo)
	err := fn(context.Background(), 1, 2)
	assert.Error(t, err)
}

// --- GetRole ---

func TestGetRole_NilModel(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return nil, nil
		},
	}

	fn := GetRole(roleRepo)
	role, err := fn(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, "", role)
}

func TestGetRole_Found(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "manager"}, nil
		},
	}

	fn := GetRole(roleRepo)
	role, err := fn(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, "manager", role)
}

// --- PromoteToManager ---

func TestPromoteToManager_CallerNotInPod(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return nil, nil // nil means not in pod
		},
	}

	fn := PromoteToManager(roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	assert.ErrorContains(t, err, "forbidden")
}

func TestPromoteToManager_CallerIsMember(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "member"}, nil
		},
	}

	fn := PromoteToManager(roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	assert.ErrorContains(t, err, "forbidden")
}

func TestPromoteToManager_Success(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "manager"}, nil
		},
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error { return nil },
	}

	fn := PromoteToManager(roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	require.NoError(t, err)
}

// --- JoinByInvite ---

func TestJoinByInvite_InviteNotFound(t *testing.T) {
	inviteRepo := &testHelpers.MockPodInviteRepo{
		GetByCodeFn: func(ctx context.Context, code string) (*podInvite.Model, error) {
			return nil, nil
		},
	}

	fn := JoinByInvite(inviteRepo, nil, nil)
	_, err := fn(context.Background(), "bad-code", 5)
	assert.ErrorContains(t, err, "not found")
}

func TestJoinByInvite_InviteExpired(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	inviteRepo := &testHelpers.MockPodInviteRepo{
		GetByCodeFn: func(ctx context.Context, code string) (*podInvite.Model, error) {
			return &podInvite.Model{PodID: 1, InviteCode: code, ExpiresAt: &past}, nil
		},
	}

	fn := JoinByInvite(inviteRepo, nil, nil)
	_, err := fn(context.Background(), "old-code", 5)
	assert.ErrorContains(t, err, "expired")
}

func TestJoinByInvite_Success(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	setRoleCalled := false
	incrementCalled := false

	inviteRepo := &testHelpers.MockPodInviteRepo{
		GetByCodeFn: func(ctx context.Context, code string) (*podInvite.Model, error) {
			return &podInvite.Model{PodID: 3, InviteCode: code, ExpiresAt: &future}, nil
		},
		IncrementUsedCountFn: func(ctx context.Context, code string) error {
			incrementCalled = true
			return nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
		AddPlayerToPodFn: func(ctx context.Context, podID, playerID int) error { return nil },
		GetByIDFn: func(ctx context.Context, podID int) (*podrepo.Model, error) {
			return &podrepo.Model{Name: "Test Pod"}, nil
		},
	}
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error {
			setRoleCalled = true
			return nil
		},
	}

	fn := JoinByInvite(inviteRepo, podRepo, roleRepo)
	e, err := fn(context.Background(), "valid-code", 5)
	require.NoError(t, err)
	require.NotNil(t, e)
	assert.Equal(t, "Test Pod", e.Name)
	assert.True(t, setRoleCalled)
	assert.True(t, incrementCalled)
}

// --- Leave ---

func TestLeave_NonManagerLeavesOK(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "member"}, nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
		RemovePlayerFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}

	fn := Leave(podRepo, roleRepo)
	err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
}

func TestLeave_SoleManagerBlocked(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "manager"}, nil
		},
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRole.Model, error) {
			return []playerPodRole.Model{
				{PlayerID: 10, Role: "manager"},
				{PlayerID: 11, Role: "member"},
			}, nil
		},
	}

	fn := Leave(&testHelpers.MockPodRepo{}, roleRepo)
	err := fn(context.Background(), 1, 10)
	assert.ErrorContains(t, err, "forbidden")
}

func TestLeave_ManagerWithCoManagerOK(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "manager"}, nil
		},
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRole.Model, error) {
			return []playerPodRole.Model{
				{PlayerID: 10, Role: "manager"},
				{PlayerID: 11, Role: "manager"},
			}, nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
		RemovePlayerFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}

	fn := Leave(podRepo, roleRepo)
	err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
}

// --- RemovePlayer ---

func TestRemovePlayer_CallerNotManager(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return nil, nil
		},
	}

	fn := RemovePlayer(&testHelpers.MockPodRepo{}, roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	assert.ErrorContains(t, err, "forbidden")
}

func TestRemovePlayer_Success(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "manager"}, nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
		RemovePlayerFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}

	fn := RemovePlayer(podRepo, roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	require.NoError(t, err)
}

// --- Pass-through function minimal tests ---

func TestGetByID_Success(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		GetByIDFn: func(ctx context.Context, podID int) (*podrepo.Model, error) {
			return &podrepo.Model{Name: "My Pod"}, nil
		},
	}

	fn := GetByID(podRepo)
	e, err := fn(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, e)
	assert.Equal(t, "My Pod", e.Name)
}

func TestGetByID_Error(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		GetByIDFn: func(ctx context.Context, podID int) (*podrepo.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetByID(podRepo)
	_, err := fn(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetByPlayerID_Success(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		GetByPlayerIDFn: func(ctx context.Context, playerID int) ([]podrepo.Model, error) {
			return []podrepo.Model{{Name: "Pod A"}, {Name: "Pod B"}}, nil
		},
	}

	fn := GetByPlayerID(podRepo)
	entities, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, entities, 2)
}

func TestGetByPlayerID_Error(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		GetByPlayerIDFn: func(ctx context.Context, playerID int) ([]podrepo.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetByPlayerID(podRepo)
	_, err := fn(context.Background(), 1)
	assert.Error(t, err)
}

func TestSoftDelete_Success(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		SoftDeleteFn: func(ctx context.Context, podID int) error { return nil },
	}

	fn := SoftDelete(podRepo)
	err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
}

func TestSoftDelete_Error(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		SoftDeleteFn: func(ctx context.Context, podID int) error { return errors.New("db error") },
	}

	fn := SoftDelete(podRepo)
	err := fn(context.Background(), 1, 10)
	assert.Error(t, err)
}

func TestUpdate_Success(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		UpdateFn: func(ctx context.Context, podID int, name string) error { return nil },
	}

	fn := Update(podRepo)
	err := fn(context.Background(), 1, "New Name")
	require.NoError(t, err)
}

func TestUpdate_Error(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		UpdateFn: func(ctx context.Context, podID int, name string) error { return errors.New("db error") },
	}

	fn := Update(podRepo)
	err := fn(context.Background(), 1, "New Name")
	assert.Error(t, err)
}

func TestGetMembersWithRoles_Success(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRole.Model, error) {
			return []playerPodRole.Model{
				{PlayerID: 1, Role: "manager"},
				{PlayerID: 2, Role: "member"},
			}, nil
		},
	}

	fn := GetMembersWithRoles(roleRepo)
	result, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Manager", result[0].Role)
}

func TestGetMembersWithRoles_Error(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRole.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetMembersWithRoles(roleRepo)
	_, err := fn(context.Background(), 1)
	assert.Error(t, err)
}

func TestGenerateInvite_Success(t *testing.T) {
	inviteRepo := &testHelpers.MockPodInviteRepo{
		AddFn: func(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
			return nil
		},
	}

	fn := GenerateInvite(inviteRepo)
	code, err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, code)
}

func TestGenerateInvite_Error(t *testing.T) {
	inviteRepo := &testHelpers.MockPodInviteRepo{
		AddFn: func(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
			return errors.New("db error")
		},
	}

	fn := GenerateInvite(inviteRepo)
	_, err := fn(context.Background(), 1, 10)
	assert.Error(t, err)
}
