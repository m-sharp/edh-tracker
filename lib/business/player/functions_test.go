package player

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	playerrepo "github.com/m-sharp/edh-tracker/lib/repositories/player"
	playerPodRolerepo "github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
)

func TestGetByID_Success(t *testing.T) {
	playerRepo := &testHelpers.MockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: 5}, Name: "Alice"}, nil
		},
	}
	gameResultRepo := &testHelpers.MockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return &gameresultrepo.Aggregate{Games: 3, Kills: 2, Record: map[int]int{1: 1}}, nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
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
	playerRepo := &testHelpers.MockPlayerRepo{
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
	playerRepo := &testHelpers.MockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: 5}, Name: "Alice"}, nil
		},
	}
	gameResultRepo := &testHelpers.MockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return nil, errors.New("stats error")
		},
	}
	fn := GetByID(playerRepo, gameResultRepo, nil)
	_, err := fn(context.Background(), 5)
	assert.Error(t, err)
}

func TestGetAll_Success(t *testing.T) {
	playerRepo := &testHelpers.MockPlayerRepo{
		GetAllFn: func(ctx context.Context) ([]playerrepo.Model, error) {
			return []playerrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Name: "Alice"},
				{GormModelBase: base.GormModelBase{ID: 2}, Name: "Bob"},
			}, nil
		},
	}
	gameResultRepo := &testHelpers.MockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return &gameresultrepo.Aggregate{Record: map[int]int{}}, nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
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
	playerRepo := &testHelpers.MockPlayerRepo{
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
	playerRepo := &testHelpers.MockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return nil, nil
		},
	}
	fn := GetPlayerName(playerRepo)
	_, err := fn(context.Background(), 99)
	assert.ErrorContains(t, err, "not found")
}

func TestUpdate_Success(t *testing.T) {
	updated := false
	playerRepo := &testHelpers.MockPlayerRepo{
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
	playerRepo := &testHelpers.MockPlayerRepo{
		UpdateFn: func(ctx context.Context, playerID int, name string) error {
			return errors.New("db error")
		},
	}
	fn := Update(playerRepo)
	err := fn(context.Background(), 1, "NewName")
	assert.Error(t, err)
}

func TestGetAllByPod_Success(t *testing.T) {
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRolerepo.Model, error) {
			return []playerPodRolerepo.Model{
				{PlayerID: 1, Role: "manager"},
				{PlayerID: 2, Role: "member"},
			}, nil
		},
	}
	playerRepo := &testHelpers.MockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{GormModelBase: base.GormModelBase{ID: playerID}, Name: "Player"}, nil
		},
	}
	gameResultRepo := &testHelpers.MockGameResultRepo{
		GetStatsForPlayerFn: func(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
			return &gameresultrepo.Aggregate{Record: map[int]int{}}, nil
		},
	}
	podRepo := &testHelpers.MockPodRepo{
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
	roleRepo := &testHelpers.MockPlayerPodRoleRepo{
		GetMembersWithRolesFn: func(ctx context.Context, podID int) ([]playerPodRolerepo.Model, error) {
			return []playerPodRolerepo.Model{
				{PlayerID: 99, Role: "member"},
			}, nil
		},
	}
	playerRepo := &testHelpers.MockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return nil, nil // not found → skip
		},
	}
	fn := GetAllByPod(playerRepo, nil, nil, roleRepo)
	got, err := fn(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, got, 0)
}
