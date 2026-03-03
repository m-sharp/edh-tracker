package deck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckrepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	deckCommanderrepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

// mockDeckRepo implements repos.DeckRepository
type mockDeckRepo struct {
	GetByIdFn func(ctx context.Context, deckID int) (*deckrepo.Model, error)
	AddFn     func(ctx context.Context, playerID int, name string, formatID int) (int, error)
}

func (m *mockDeckRepo) GetAll(ctx context.Context) ([]deckrepo.Model, error) {
	panic("unexpected call to GetAll")
}
func (m *mockDeckRepo) GetAllForPlayer(ctx context.Context, playerID int) ([]deckrepo.Model, error) {
	panic("unexpected call to GetAllForPlayer")
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
	panic("unexpected call to SoftDelete")
}

// mockDeckCommanderRepo implements repos.DeckCommanderRepository
type mockDeckCommanderRepo struct {
	GetByDeckIdFn func(ctx context.Context, deckID int) (*deckCommanderrepo.Model, error)
	AddFn         func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error)
	addCalled     bool
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

// mockGameResultRepoForDeck implements repos.GameResultRepository (used in GetAll tests)
type mockGameResultRepoForDeck struct{}

func (m *mockGameResultRepoForDeck) GetByGameId(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call")
}
func (m *mockGameResultRepoForDeck) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
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
			return &deckrepo.Model{ModelBase: base.ModelBase{ID: 20}, Name: "Krenko Goblins"}, nil
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
