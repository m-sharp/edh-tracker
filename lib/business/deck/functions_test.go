package deck

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckrepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	deckCommanderrepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	podrepo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
)

// mockDeckRepo implements repos.DeckRepository
type mockDeckRepo struct {
	GetByIdFn           func(ctx context.Context, deckID int) (*deckrepo.Model, error)
	AddFn               func(ctx context.Context, playerID int, name string, formatID int) (int, error)
	UpdateFn            func(ctx context.Context, deckID int, fields deckrepo.UpdateFields) error
	SoftDeleteFn        func(ctx context.Context, id int) error
	GetAllByPlayerIDsFn func(ctx context.Context, playerIDs []int) ([]deckrepo.Model, error)
}

func (m *mockDeckRepo) GetAll(ctx context.Context) ([]deckrepo.Model, error) {
	panic("unexpected call to GetAll")
}
func (m *mockDeckRepo) GetAllForPlayer(ctx context.Context, playerID int) ([]deckrepo.Model, error) {
	panic("unexpected call to GetAllForPlayer")
}
func (m *mockDeckRepo) GetAllByPlayerIDs(ctx context.Context, playerIDs []int) ([]deckrepo.Model, error) {
	if m.GetAllByPlayerIDsFn != nil {
		return m.GetAllByPlayerIDsFn(ctx, playerIDs)
	}
	panic("unexpected call to GetAllByPlayerIDs")
}
func (m *mockDeckRepo) Update(ctx context.Context, deckID int, fields deckrepo.UpdateFields) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, deckID, fields)
	}
	panic("unexpected call to Update")
}
func (m *mockDeckRepo) GetById(ctx context.Context, deckID int) (*deckrepo.Model, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, deckID)
	}
	panic("unexpected call to GetById")
}
func (m *mockDeckRepo) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, playerID, name, formatID)
	}
	panic("unexpected call to Add")
}
func (m *mockDeckRepo) BulkAdd(ctx context.Context, decks []deckrepo.Model) ([]deckrepo.Model, error) {
	panic("unexpected call to BulkAdd")
}
func (m *mockDeckRepo) Retire(ctx context.Context, deckID int) error {
	panic("unexpected call to Retire")
}
func (m *mockDeckRepo) SoftDelete(ctx context.Context, id int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	panic("unexpected call to SoftDelete")
}

// mockDeckCommanderRepo implements repos.DeckCommanderRepository
type mockDeckCommanderRepo struct {
	GetByDeckIdFn    func(ctx context.Context, deckID int) (*deckCommanderrepo.Model, error)
	AddFn            func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error)
	DeleteByDeckIDFn func(ctx context.Context, deckID int) error
	addCalled        bool
}

func (m *mockDeckCommanderRepo) GetByDeckId(ctx context.Context, deckID int) (*deckCommanderrepo.Model, error) {
	if m.GetByDeckIdFn != nil {
		return m.GetByDeckIdFn(ctx, deckID)
	}
	panic("unexpected call to GetByDeckId")
}
func (m *mockDeckCommanderRepo) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
	m.addCalled = true
	if m.AddFn != nil {
		return m.AddFn(ctx, deckID, commanderID, partnerCommanderID)
	}
	panic("unexpected call to Add")
}
func (m *mockDeckCommanderRepo) BulkAdd(ctx context.Context, entries []deckCommanderrepo.Model) error {
	panic("unexpected call to BulkAdd")
}
func (m *mockDeckCommanderRepo) DeleteByDeckID(ctx context.Context, deckID int) error {
	if m.DeleteByDeckIDFn != nil {
		return m.DeleteByDeckIDFn(ctx, deckID)
	}
	panic("unexpected call to DeleteByDeckID")
}

// mockGameResultRepoForDeck implements repos.GameResultRepository (used in GetAll tests)
type mockGameResultRepoForDeck struct {
	GetStatsForDeckFn func(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error)
}

func (m *mockGameResultRepoForDeck) GetByGameId(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
	if m.GetStatsForDeckFn != nil {
		return m.GetStatsForDeckFn(ctx, deckID)
	}
	panic("unexpected call to GetStatsForDeck")
}
func (m *mockGameResultRepoForDeck) GetByID(ctx context.Context, resultID int) (*gameresultrepo.Model, error) {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) Add(ctx context.Context, model gameresultrepo.Model) (int, error) {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) BulkAdd(ctx context.Context, results []gameresultrepo.Model) error {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call")
}

func makeFormat(id int, name string) *format.Entity {
	return &format.Entity{ID: id, Name: name}
}

func TestCreate_CommanderFormat_NoCommanderID(t *testing.T) {
	deckRepo := &mockDeckRepo{}
	deckCmdrRepo := &mockDeckCommanderRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return makeFormat(1, "commander"), nil
	}

	fn := Create(deckRepo, deckCmdrRepo, getFormat)
	_, err := fn(context.Background(), 1, "Test Deck", 1, nil, nil)
	assert.ErrorContains(t, err, "commander_id is required")
}

func TestCreate_CommanderFormat_WithCommander(t *testing.T) {
	deckRepo := &mockDeckRepo{
		AddFn: func(ctx context.Context, playerID int, name string, formatID int) (int, error) {
			return 10, nil
		},
	}
	deckCmdrRepo := &mockDeckCommanderRepo{
		AddFn: func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
			return 1, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return makeFormat(1, "commander"), nil
	}

	commanderID := 5
	fn := Create(deckRepo, deckCmdrRepo, getFormat)
	id, err := fn(context.Background(), 1, "Krenko Goblins", 1, &commanderID, nil)
	require.NoError(t, err)
	assert.Equal(t, 10, id)
}

func TestCreate_OtherFormat_NoCommander(t *testing.T) {
	deckCmdrRepo := &mockDeckCommanderRepo{}
	deckRepo := &mockDeckRepo{
		AddFn: func(ctx context.Context, playerID int, name string, formatID int) (int, error) {
			return 11, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return makeFormat(2, "other"), nil
	}

	fn := Create(deckRepo, deckCmdrRepo, getFormat)
	id, err := fn(context.Background(), 1, "Casual Deck", 2, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 11, id)
	assert.False(t, deckCmdrRepo.addCalled, "deckCmdrRepo.Add should not be called for other format")
}

func TestCreate_FormatNotFound(t *testing.T) {
	deckRepo := &mockDeckRepo{}
	deckCmdrRepo := &mockDeckCommanderRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return nil, nil
	}

	fn := Create(deckRepo, deckCmdrRepo, getFormat)
	_, err := fn(context.Background(), 1, "Test Deck", 99, nil, nil)
	assert.Error(t, err)
}

func TestGetCommanderEntry_NoEntry(t *testing.T) {
	deckCmdrRepo := &mockDeckCommanderRepo{
		GetByDeckIdFn: func(ctx context.Context, deckID int) (*deckCommanderrepo.Model, error) {
			return nil, nil
		},
	}
	getCommanderName := func(ctx context.Context, id int) (string, error) {
		panic("should not be called")
	}

	fn := GetCommanderEntry(deckCmdrRepo, getCommanderName)
	got, err := fn(context.Background(), 7)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetCommanderEntry_WithCommander(t *testing.T) {
	deckCmdrRepo := &mockDeckCommanderRepo{
		GetByDeckIdFn: func(ctx context.Context, deckID int) (*deckCommanderrepo.Model, error) {
			return &deckCommanderrepo.Model{ModelBase: base.ModelBase{ID: 1}, DeckID: 7, CommanderID: 5}, nil
		},
	}
	getCommanderName := func(ctx context.Context, id int) (string, error) {
		if id == 5 {
			return "Krenko", nil
		}
		panic("unexpected commander id")
	}

	fn := GetCommanderEntry(deckCmdrRepo, getCommanderName)
	got, err := fn(context.Background(), 7)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Krenko", got.CommanderName)
	assert.Nil(t, got.PartnerCommanderID)
}

func TestGetCommanderEntry_WithPartner(t *testing.T) {
	partnerID := 6
	deckCmdrRepo := &mockDeckCommanderRepo{
		GetByDeckIdFn: func(ctx context.Context, deckID int) (*deckCommanderrepo.Model, error) {
			return &deckCommanderrepo.Model{
				ModelBase:          base.ModelBase{ID: 1},
				DeckID:             7,
				CommanderID:        5,
				PartnerCommanderID: &partnerID,
			}, nil
		},
	}
	getCommanderName := func(ctx context.Context, id int) (string, error) {
		switch id {
		case 5:
			return "Krenko", nil
		case 6:
			return "Goblin Chieftain", nil
		default:
			panic("unexpected commander id")
		}
	}

	fn := GetCommanderEntry(deckCmdrRepo, getCommanderName)
	got, err := fn(context.Background(), 7)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Krenko", got.CommanderName)
	require.NotNil(t, got.PartnerCommanderID)
	assert.Equal(t, 6, *got.PartnerCommanderID)
	require.NotNil(t, got.PartnerCommanderName)
	assert.Equal(t, "Goblin Chieftain", *got.PartnerCommanderName)
}

func TestGetDeckName_Success(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: 20}, Name: "Krenko Goblins"}, nil
		},
	}
	fn := GetDeckName(deckRepo)
	name, err := fn(context.Background(), 20)
	require.NoError(t, err)
	assert.Equal(t, "Krenko Goblins", name)
}

func TestGetDeckName_NotFound(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return nil, nil
		},
	}
	fn := GetDeckName(deckRepo)
	_, err := fn(context.Background(), 99)
	assert.Error(t, err)
}

// mockPodRepo implements repos.PodRepository (only GetPlayerIDs is exercised here)
type mockPodRepo struct {
	GetPlayerIDsFn func(ctx context.Context, podID int) ([]int, error)
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
	panic("unexpected call to GetIDsByPlayerID")
}
func (m *mockPodRepo) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
	if m.GetPlayerIDsFn != nil {
		return m.GetPlayerIDsFn(ctx, podID)
	}
	panic("unexpected call to GetPlayerIDs")
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
func (m *mockPodRepo) SoftDelete(ctx context.Context, podID int) error {
	panic("unexpected call to SoftDelete")
}
func (m *mockPodRepo) Update(ctx context.Context, podID int, name string) error {
	panic("unexpected call to Update")
}
func (m *mockPodRepo) RemovePlayer(ctx context.Context, podID, playerID int) error {
	panic("unexpected call to RemovePlayer")
}

func TestDeckUpdate_Success_NoCommander(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 42}, nil
		},
		UpdateFn: func(ctx context.Context, deckID int, fields deckrepo.UpdateFields) error {
			return nil
		},
	}
	deckCmdrRepo := &mockDeckCommanderRepo{}

	fn := Update(deckRepo, deckCmdrRepo)
	name := "New Name"
	err := fn(context.Background(), 1, 42, UpdateFields{Name: &name})
	require.NoError(t, err)
}

func TestDeckUpdate_WithCommander(t *testing.T) {
	deleteCalled := false
	addCalled := false
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 10}, nil
		},
		UpdateFn: func(ctx context.Context, deckID int, fields deckrepo.UpdateFields) error {
			return nil
		},
	}
	deckCmdrRepo := &mockDeckCommanderRepo{
		DeleteByDeckIDFn: func(ctx context.Context, deckID int) error {
			deleteCalled = true
			return nil
		},
		AddFn: func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
			addCalled = true
			return 1, nil
		},
	}

	commanderID := 5
	fn := Update(deckRepo, deckCmdrRepo)
	err := fn(context.Background(), 1, 10, UpdateFields{CommanderID: &commanderID})
	require.NoError(t, err)
	assert.True(t, deleteCalled, "DeleteByDeckID should be called when CommanderID is set")
	assert.True(t, addCalled, "Add should be called to set the new commander")
}

func TestDeckUpdate_NotFound(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return nil, nil
		},
	}
	fn := Update(deckRepo, &mockDeckCommanderRepo{})
	err := fn(context.Background(), 99, 42, UpdateFields{})
	assert.ErrorContains(t, err, "not found")
}

func TestDeckUpdate_Forbidden(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 10}, nil
		},
	}
	fn := Update(deckRepo, &mockDeckCommanderRepo{})
	err := fn(context.Background(), 1, 99, UpdateFields{}) // callerPlayerID=99 != deck.PlayerID=10
	assert.ErrorContains(t, err, "forbidden")
}

func TestDeckSoftDelete_Success(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 7}, nil
		},
		SoftDeleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	fn := SoftDelete(deckRepo)
	err := fn(context.Background(), 1, 7)
	require.NoError(t, err)
}

func TestDeckSoftDelete_NotFound(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return nil, nil
		},
	}
	fn := SoftDelete(deckRepo)
	err := fn(context.Background(), 99, 7)
	assert.ErrorContains(t, err, "not found")
}

func TestDeckSoftDelete_Forbidden(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 10}, nil
		},
	}
	fn := SoftDelete(deckRepo)
	err := fn(context.Background(), 1, 99) // caller=99, owner=10
	assert.ErrorContains(t, err, "forbidden")
}

func TestDeckGetAllByPod_Success(t *testing.T) {
	podRepo := &mockPodRepo{
		GetPlayerIDsFn: func(ctx context.Context, podID int) ([]int, error) {
			return []int{1, 2}, nil
		},
	}
	deckRepo := &mockDeckRepo{
		GetAllByPlayerIDsFn: func(ctx context.Context, playerIDs []int) ([]deckrepo.Model, error) {
			return []deckrepo.Model{
				{GormModelBase: base.GormModelBase{ID: 10}, PlayerID: 1, FormatID: 1, Name: "Deck A"},
			}, nil
		},
	}
	gameResultRepo := &mockGameResultRepoForDeck{
		GetStatsForDeckFn: func(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
			return nil, nil
		},
	}
	getPlayerName := func(ctx context.Context, id int) (string, error) { return "Alice", nil }
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: 1, Name: "commander"}, nil
	}
	getCommanderEntry := func(ctx context.Context, deckID int) (*CommanderInfo, error) {
		return nil, nil
	}

	fn := GetAllByPod(deckRepo, podRepo, gameResultRepo, getPlayerName, getFormat, getCommanderEntry)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "Deck A", got[0].Name)
}

func TestDeckGetAllByPod_EmptyPod(t *testing.T) {
	podRepo := &mockPodRepo{
		GetPlayerIDsFn: func(ctx context.Context, podID int) ([]int, error) {
			return []int{}, nil
		},
	}
	fn := GetAllByPod(&mockDeckRepo{}, podRepo, &mockGameResultRepoForDeck{}, nil, nil, nil)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestDeckGetAllByPod_Error(t *testing.T) {
	podRepo := &mockPodRepo{
		GetPlayerIDsFn: func(ctx context.Context, podID int) ([]int, error) {
			return nil, errors.New("db error")
		},
	}
	fn := GetAllByPod(&mockDeckRepo{}, podRepo, &mockGameResultRepoForDeck{}, nil, nil, nil)
	_, err := fn(context.Background(), 5)
	assert.Error(t, err)
}
