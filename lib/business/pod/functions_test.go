package pod

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	podrepo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
	"github.com/m-sharp/edh-tracker/lib/repositories/podInvite"
)

// --- mock structs ---

type mockPodRepo struct {
	GetByIDFn          func(ctx context.Context, podID int) (*podrepo.Model, error)
	GetByPlayerIDFn    func(ctx context.Context, playerID int) ([]podrepo.Model, error)
	AddFn              func(ctx context.Context, name string) (int, error)
	AddPlayerToPodFn   func(ctx context.Context, podID, playerID int) error
	RemovePlayerFn     func(ctx context.Context, podID, playerID int) error
	SoftDeleteFn       func(ctx context.Context, podID int) error
	UpdateFn           func(ctx context.Context, podID int, name string) error
	GetIDsByPlayerIDFn func(ctx context.Context, playerID int) ([]int, error)
}

func (m *mockPodRepo) GetAll(ctx context.Context) ([]podrepo.Model, error) {
	panic("unexpected call to GetAll")
}
func (m *mockPodRepo) GetByID(ctx context.Context, podID int) (*podrepo.Model, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, podID)
	}
	panic("unexpected call to GetByID")
}
func (m *mockPodRepo) GetByPlayerID(ctx context.Context, playerID int) ([]podrepo.Model, error) {
	if m.GetByPlayerIDFn != nil {
		return m.GetByPlayerIDFn(ctx, playerID)
	}
	panic("unexpected call to GetByPlayerID")
}
func (m *mockPodRepo) GetByName(ctx context.Context, name string) (*podrepo.Model, error) {
	panic("unexpected call to GetByName")
}
func (m *mockPodRepo) GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error) {
	if m.GetIDsByPlayerIDFn != nil {
		return m.GetIDsByPlayerIDFn(ctx, playerID)
	}
	panic("unexpected call to GetIDsByPlayerID")
}
func (m *mockPodRepo) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
	panic("unexpected call to GetPlayerIDs")
}
func (m *mockPodRepo) Add(ctx context.Context, name string) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, name)
	}
	panic("unexpected call to Add")
}
func (m *mockPodRepo) BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error {
	panic("unexpected call to BulkAddPlayers")
}
func (m *mockPodRepo) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
	if m.AddPlayerToPodFn != nil {
		return m.AddPlayerToPodFn(ctx, podID, playerID)
	}
	panic("unexpected call to AddPlayerToPod")
}
func (m *mockPodRepo) SoftDelete(ctx context.Context, podID int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, podID)
	}
	panic("unexpected call to SoftDelete")
}
func (m *mockPodRepo) Update(ctx context.Context, podID int, name string) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, podID, name)
	}
	panic("unexpected call to Update")
}
func (m *mockPodRepo) RemovePlayer(ctx context.Context, podID, playerID int) error {
	if m.RemovePlayerFn != nil {
		return m.RemovePlayerFn(ctx, podID, playerID)
	}
	panic("unexpected call to RemovePlayer")
}

type mockRoleRepo struct {
	GetRoleFn             func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error)
	SetRoleFn             func(ctx context.Context, podID, playerID int, role string) error
	GetMembersWithRolesFn func(ctx context.Context, podID int) ([]playerPodRole.Model, error)
}

func (m *mockRoleRepo) GetRole(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
	if m.GetRoleFn != nil {
		return m.GetRoleFn(ctx, podID, playerID)
	}
	panic("unexpected call to GetRole")
}
func (m *mockRoleRepo) SetRole(ctx context.Context, podID, playerID int, role string) error {
	if m.SetRoleFn != nil {
		return m.SetRoleFn(ctx, podID, playerID, role)
	}
	panic("unexpected call to SetRole")
}
func (m *mockRoleRepo) GetMembersWithRoles(ctx context.Context, podID int) ([]playerPodRole.Model, error) {
	if m.GetMembersWithRolesFn != nil {
		return m.GetMembersWithRolesFn(ctx, podID)
	}
	panic("unexpected call to GetMembersWithRoles")
}
func (m *mockRoleRepo) BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error {
	panic("unexpected call to BulkAdd")
}

type mockInviteRepo struct {
	GetByCodeFn          func(ctx context.Context, code string) (*podInvite.Model, error)
	AddFn                func(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error
	IncrementUsedCountFn func(ctx context.Context, code string) error
}

func (m *mockInviteRepo) GetByCode(ctx context.Context, code string) (*podInvite.Model, error) {
	if m.GetByCodeFn != nil {
		return m.GetByCodeFn(ctx, code)
	}
	panic("unexpected call to GetByCode")
}
func (m *mockInviteRepo) Add(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
	if m.AddFn != nil {
		return m.AddFn(ctx, podID, createdByPlayerID, code, expiresAt)
	}
	panic("unexpected call to Add")
}
func (m *mockInviteRepo) IncrementUsedCount(ctx context.Context, code string) error {
	if m.IncrementUsedCountFn != nil {
		return m.IncrementUsedCountFn(ctx, code)
	}
	panic("unexpected call to IncrementUsedCount")
}

// --- Create ---

func TestCreate_Success(t *testing.T) {
	podRepo := &mockPodRepo{
		AddFn: func(ctx context.Context, name string) (int, error) { return 5, nil },
	}
	roleRepo := &mockRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error { return nil },
	}

	fn := Create(podRepo, roleRepo)
	podID, err := fn(context.Background(), "My Pod", 1)
	require.NoError(t, err)
	assert.Equal(t, 5, podID)
}

func TestCreate_SetRoleError(t *testing.T) {
	podRepo := &mockPodRepo{
		AddFn: func(ctx context.Context, name string) (int, error) { return 5, nil },
	}
	roleRepo := &mockRoleRepo{
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
	podRepo := &mockPodRepo{
		AddPlayerToPodFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}
	roleRepo := &mockRoleRepo{
		SetRoleFn: func(ctx context.Context, podID, playerID int, role string) error { return nil },
	}

	fn := AddPlayer(podRepo, roleRepo)
	err := fn(context.Background(), 1, 2)
	require.NoError(t, err)
}

func TestAddPlayer_SetRoleError(t *testing.T) {
	podRepo := &mockPodRepo{
		AddPlayerToPodFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}
	roleRepo := &mockRoleRepo{
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
	roleRepo := &mockRoleRepo{
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
	roleRepo := &mockRoleRepo{
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
	roleRepo := &mockRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return nil, nil // nil means not in pod
		},
	}

	fn := PromoteToManager(roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	assert.ErrorContains(t, err, "forbidden")
}

func TestPromoteToManager_CallerIsMember(t *testing.T) {
	roleRepo := &mockRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "member"}, nil
		},
	}

	fn := PromoteToManager(roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	assert.ErrorContains(t, err, "forbidden")
}

func TestPromoteToManager_Success(t *testing.T) {
	roleRepo := &mockRoleRepo{
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
	inviteRepo := &mockInviteRepo{
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
	inviteRepo := &mockInviteRepo{
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

	inviteRepo := &mockInviteRepo{
		GetByCodeFn: func(ctx context.Context, code string) (*podInvite.Model, error) {
			return &podInvite.Model{PodID: 3, InviteCode: code, ExpiresAt: &future}, nil
		},
		IncrementUsedCountFn: func(ctx context.Context, code string) error {
			incrementCalled = true
			return nil
		},
	}
	podRepo := &mockPodRepo{
		AddPlayerToPodFn: func(ctx context.Context, podID, playerID int) error { return nil },
		GetByIDFn: func(ctx context.Context, podID int) (*podrepo.Model, error) {
			return &podrepo.Model{Name: "Test Pod"}, nil
		},
	}
	roleRepo := &mockRoleRepo{
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
	roleRepo := &mockRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "member"}, nil
		},
	}
	podRepo := &mockPodRepo{
		RemovePlayerFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}

	fn := Leave(podRepo, roleRepo)
	err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
}

func TestLeave_SoleManagerBlocked(t *testing.T) {
	roleRepo := &mockRoleRepo{
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

	fn := Leave(&mockPodRepo{}, roleRepo)
	err := fn(context.Background(), 1, 10)
	assert.ErrorContains(t, err, "forbidden")
}

func TestLeave_ManagerWithCoManagerOK(t *testing.T) {
	roleRepo := &mockRoleRepo{
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
	podRepo := &mockPodRepo{
		RemovePlayerFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}

	fn := Leave(podRepo, roleRepo)
	err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
}

// --- RemovePlayer ---

func TestRemovePlayer_CallerNotManager(t *testing.T) {
	roleRepo := &mockRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return nil, nil
		},
	}

	fn := RemovePlayer(&mockPodRepo{}, roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	assert.ErrorContains(t, err, "forbidden")
}

func TestRemovePlayer_Success(t *testing.T) {
	roleRepo := &mockRoleRepo{
		GetRoleFn: func(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error) {
			return &playerPodRole.Model{Role: "manager"}, nil
		},
	}
	podRepo := &mockPodRepo{
		RemovePlayerFn: func(ctx context.Context, podID, playerID int) error { return nil },
	}

	fn := RemovePlayer(podRepo, roleRepo)
	err := fn(context.Background(), 1, 10, 20)
	require.NoError(t, err)
}

// --- Pass-through function minimal tests ---

func TestGetByID_Success(t *testing.T) {
	podRepo := &mockPodRepo{
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
	podRepo := &mockPodRepo{
		GetByIDFn: func(ctx context.Context, podID int) (*podrepo.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetByID(podRepo)
	_, err := fn(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetByPlayerID_Success(t *testing.T) {
	podRepo := &mockPodRepo{
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
	podRepo := &mockPodRepo{
		GetByPlayerIDFn: func(ctx context.Context, playerID int) ([]podrepo.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetByPlayerID(podRepo)
	_, err := fn(context.Background(), 1)
	assert.Error(t, err)
}

func TestSoftDelete_Success(t *testing.T) {
	podRepo := &mockPodRepo{
		SoftDeleteFn: func(ctx context.Context, podID int) error { return nil },
	}

	fn := SoftDelete(podRepo)
	err := fn(context.Background(), 1, 10)
	require.NoError(t, err)
}

func TestSoftDelete_Error(t *testing.T) {
	podRepo := &mockPodRepo{
		SoftDeleteFn: func(ctx context.Context, podID int) error { return errors.New("db error") },
	}

	fn := SoftDelete(podRepo)
	err := fn(context.Background(), 1, 10)
	assert.Error(t, err)
}

func TestUpdate_Success(t *testing.T) {
	podRepo := &mockPodRepo{
		UpdateFn: func(ctx context.Context, podID int, name string) error { return nil },
	}

	fn := Update(podRepo)
	err := fn(context.Background(), 1, "New Name")
	require.NoError(t, err)
}

func TestUpdate_Error(t *testing.T) {
	podRepo := &mockPodRepo{
		UpdateFn: func(ctx context.Context, podID int, name string) error { return errors.New("db error") },
	}

	fn := Update(podRepo)
	err := fn(context.Background(), 1, "New Name")
	assert.Error(t, err)
}

func TestGetMembersWithRoles_Success(t *testing.T) {
	roleRepo := &mockRoleRepo{
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
	assert.Equal(t, "manager", result[0].Role)
}

func TestGetMembersWithRoles_Error(t *testing.T) {
	roleRepo := &mockRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRole.Model, error) {
			return nil, errors.New("db error")
		},
	}

	fn := GetMembersWithRoles(roleRepo)
	_, err := fn(context.Background(), 1)
	assert.Error(t, err)
}

func TestGenerateInvite_Success(t *testing.T) {
	inviteRepo := &mockInviteRepo{
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
	inviteRepo := &mockInviteRepo{
		AddFn: func(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
			return errors.New("db error")
		},
	}

	fn := GenerateInvite(inviteRepo)
	_, err := fn(context.Background(), 1, 10)
	assert.Error(t, err)
}
