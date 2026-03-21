package player

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	playerrepo "github.com/m-sharp/edh-tracker/lib/repositories/player"
	playerPodRolerepo "github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	podrepo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
)

// mockPlayerRepo implements repos.PlayerRepository
type mockPlayerRepo struct {
	GetAllFn    func(ctx context.Context) ([]playerrepo.Model, error)
	GetByIdFn   func(ctx context.Context, playerID int) (*playerrepo.Model, error)
	GetByNameFn func(ctx context.Context, name string) (*playerrepo.Model, error)
	UpdateFn    func(ctx context.Context, playerID int, name string) error
}

func (m *mockPlayerRepo) GetAll(ctx context.Context) ([]playerrepo.Model, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	panic("unexpected call to GetAll")
}
func (m *mockPlayerRepo) GetById(ctx context.Context, playerID int) (*playerrepo.Model, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, playerID)
	}
	panic("unexpected call to GetById")
}
func (m *mockPlayerRepo) GetByName(ctx context.Context, name string) (*playerrepo.Model, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	panic("unexpected call to GetByName")
}
func (m *mockPlayerRepo) GetByNames(ctx context.Context, names []string) ([]playerrepo.Model, error) {
	panic("unexpected call to GetByNames")
}
func (m *mockPlayerRepo) Add(ctx context.Context, name string) (int, error) {
	panic("unexpected call to Add")
}
func (m *mockPlayerRepo) BulkAdd(ctx context.Context, names []string) ([]playerrepo.Model, error) {
	panic("unexpected call to BulkAdd")
}
func (m *mockPlayerRepo) Update(ctx context.Context, playerID int, name string) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, playerID, name)
	}
	panic("unexpected call to Update")
}
func (m *mockPlayerRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}

// mockGameResultRepo implements repos.GameResultRepository
type mockGameResultRepo struct {
	GetStatsForPlayerFn func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error)
}

func (m *mockGameResultRepo) GetByGameId(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	panic("unexpected call to GetByGameId")
}
func (m *mockGameResultRepo) GetByGameIDWithDeckInfo(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	panic("unexpected call to GetByGameIDWithDeckInfo")
}
func (m *mockGameResultRepo) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	if m.GetStatsForPlayerFn != nil {
		return m.GetStatsForPlayerFn(ctx, playerID)
	}
	panic("unexpected call to GetStatsForPlayer")
}
func (m *mockGameResultRepo) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForDeck")
}
func (m *mockGameResultRepo) GetByID(ctx context.Context, resultID int) (*gameresultrepo.Model, error) {
	panic("unexpected call to GetByID")
}
func (m *mockGameResultRepo) Add(ctx context.Context, model gameresultrepo.Model) (int, error) {
	panic("unexpected call to Add")
}
func (m *mockGameResultRepo) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	panic("unexpected call to Update")
}
func (m *mockGameResultRepo) BulkAdd(ctx context.Context, results []gameresultrepo.Model) error {
	panic("unexpected call to BulkAdd")
}
func (m *mockGameResultRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}

// mockPodRepo implements repos.PodRepository
type mockPodRepo struct {
	GetIDsByPlayerIDFn func(ctx context.Context, playerID int) ([]int, error)
}

func (m *mockPodRepo) GetAll(ctx context.Context) ([]podrepo.Model, error) {
	panic("unexpected call to GetAll")
}
func (m *mockPodRepo) GetByIDWithMembers(ctx context.Context, podID int) (*podrepo.Model, error) {
	panic("unexpected call to GetByIDWithMembers")
}
func (m *mockPodRepo) GetByID(ctx context.Context, podID int) (*podrepo.Model, error) {
	panic("unexpected call to GetByID")
}
func (m *mockPodRepo) GetByPlayerID(ctx context.Context, playerID int) ([]podrepo.Model, error) {
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
func (m *mockPodRepo) Add(ctx context.Context, name string) (int, error) {
	panic("unexpected call to Add")
}
func (m *mockPodRepo) BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error {
	panic("unexpected call to BulkAddPlayers")
}
func (m *mockPodRepo) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
	panic("unexpected call to AddPlayerToPod")
}
func (m *mockPodRepo) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
	panic("unexpected call to GetPlayerIDs")
}
func (m *mockPodRepo) SoftDelete(ctx context.Context, podID int) error {
	panic("unexpected call to SoftDelete")
}
func (m *mockPodRepo) Update(ctx context.Context, podID int, name string) error {
	panic("unexpected call to Update")
}
func (m *mockPodRepo) RemovePlayer(ctx context.Context, podID, playerID int) error {
	panic("unexpected call to RemovePlayer")
}

func TestGetByID_Success(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: 5}, Name: "Alice"}, nil
		},
	}
	gameResultRepo := &mockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return &gameresultrepo.Aggregate{Games: 3, Kills: 2, Record: map[int]int{1: 1}}, nil
		},
	}
	podRepo := &mockPodRepo{
		GetIDsByPlayerIDFn: func(ctx context.Context, playerID int) ([]int, error) {
			return []int{1, 2}, nil
		},
	}

	fn := GetByID(playerRepo, gameResultRepo, podRepo)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 5, got.ID)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, 3, got.Stats.Games)
	assert.Equal(t, []int{1, 2}, got.PodIDs)
}

func TestGetByID_NotFound(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return nil, nil
		},
	}
	fn := GetByID(playerRepo, nil, nil)
	got, err := fn(context.Background(), 999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetByID_StatsError(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: 5}, Name: "Alice"}, nil
		},
	}
	gameResultRepo := &mockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return nil, errors.New("stats error")
		},
	}
	fn := GetByID(playerRepo, gameResultRepo, nil)
	_, err := fn(context.Background(), 5)
	assert.Error(t, err)
}

func TestGetAll_Success(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetAllFn: func(ctx context.Context) ([]playerrepo.Model, error) {
			return []playerrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Name: "Alice"},
				{GormModelBase: base.GormModelBase{ID: 2}, Name: "Bob"},
			}, nil
		},
	}
	gameResultRepo := &mockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return &gameresultrepo.Aggregate{Record: map[int]int{}}, nil
		},
	}
	podRepo := &mockPodRepo{
		GetIDsByPlayerIDFn: func(ctx context.Context, playerID int) ([]int, error) {
			return []int{}, nil
		},
	}

	fn := GetAll(playerRepo, gameResultRepo, podRepo)
	got, err := fn(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestGetPlayerName_Success(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: 1}, Name: "Alice"}, nil
		},
	}
	fn := GetPlayerName(playerRepo)
	name, err := fn(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "Alice", name)
}

func TestGetPlayerName_NotFound(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return nil, nil
		},
	}
	fn := GetPlayerName(playerRepo)
	_, err := fn(context.Background(), 99)
	assert.ErrorContains(t, err, "not found")
}

// mockPlayerPodRoleRepo implements repos.PlayerPodRoleRepository
type mockPlayerPodRoleRepo struct {
	GetMembersWithRolesFn func(ctx context.Context, podID int) ([]playerPodRolerepo.Model, error)
}

func (m *mockPlayerPodRoleRepo) GetRole(ctx context.Context, podID, playerID int) (*playerPodRolerepo.Model, error) {
	panic("unexpected call to GetRole")
}
func (m *mockPlayerPodRoleRepo) SetRole(ctx context.Context, podID, playerID int, role string) error {
	panic("unexpected call to SetRole")
}
func (m *mockPlayerPodRoleRepo) GetMembersWithRoles(ctx context.Context, podID int) ([]playerPodRolerepo.Model, error) {
	if m.GetMembersWithRolesFn != nil {
		return m.GetMembersWithRolesFn(ctx, podID)
	}
	panic("unexpected call to GetMembersWithRoles")
}
func (m *mockPlayerPodRoleRepo) BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error {
	panic("unexpected call to BulkAdd")
}

func TestUpdate_Success(t *testing.T) {
	updated := false
	playerRepo := &mockPlayerRepo{
		UpdateFn: func(ctx context.Context, playerID int, name string) error {
			updated = true
			return nil
		},
	}
	fn := Update(playerRepo)
	err := fn(context.Background(), 1, "NewName")
	require.NoError(t, err)
	assert.True(t, updated)
}

func TestUpdate_Error(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		UpdateFn: func(ctx context.Context, playerID int, name string) error {
			return errors.New("db error")
		},
	}
	fn := Update(playerRepo)
	err := fn(context.Background(), 1, "NewName")
	assert.Error(t, err)
}

func TestGetAllByPod_Success(t *testing.T) {
	roleRepo := &mockPlayerPodRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRolerepo.Model, error) {
			return []playerPodRolerepo.Model{
				{PlayerID: 1, Role: "manager"},
				{PlayerID: 2, Role: "member"},
			}, nil
		},
	}
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: playerID}, Name: "Player"}, nil
		},
	}
	gameResultRepo := &mockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return &gameresultrepo.Aggregate{Record: map[int]int{}}, nil
		},
	}
	podRepo := &mockPodRepo{
		GetIDsByPlayerIDFn: func(ctx context.Context, playerID int) ([]int, error) {
			return []int{10}, nil
		},
	}

	fn := GetAllByPod(playerRepo, gameResultRepo, podRepo, roleRepo)
	got, err := fn(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "manager", got[0].Role)
	assert.Equal(t, "member", got[1].Role)
}

func TestGetAllByPod_PlayerNotFound_Skipped(t *testing.T) {
	roleRepo := &mockPlayerPodRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRolerepo.Model, error) {
			return []playerPodRolerepo.Model{
				{PlayerID: 99, Role: "member"},
			}, nil
		},
	}
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return nil, nil // not found → skip
		},
	}
	fn := GetAllByPod(playerRepo, nil, nil, roleRepo)
	got, err := fn(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, got, 0)
}
