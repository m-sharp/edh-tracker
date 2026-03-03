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
	podrepo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
)

// mockPlayerRepo implements repos.PlayerRepository
type mockPlayerRepo struct {
	GetAllFn    func(ctx context.Context) ([]playerrepo.Model, error)
	GetByIdFn   func(ctx context.Context, playerID int) (*playerrepo.Model, error)
	GetByNameFn func(ctx context.Context, name string) (*playerrepo.Model, error)
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
func (m *mockGameResultRepo) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	if m.GetStatsForPlayerFn != nil {
		return m.GetStatsForPlayerFn(ctx, playerID)
	}
	panic("unexpected call to GetStatsForPlayer")
}
func (m *mockGameResultRepo) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForDeck")
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

func TestGetByID_Success(t *testing.T) {
	playerRepo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, playerID int) (*playerrepo.Model, error) {
			return &playerrepo.Model{ModelBase: base.ModelBase{ID: 5}, Name: "Alice"}, nil
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
			return &playerrepo.Model{ModelBase: base.ModelBase{ID: 5}, Name: "Alice"}, nil
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
				{ModelBase: base.ModelBase{ID: 1}, Name: "Alice"},
				{ModelBase: base.ModelBase{ID: 2}, Name: "Bob"},
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
			return &playerrepo.Model{ModelBase: base.ModelBase{ID: 1}, Name: "Alice"}, nil
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
