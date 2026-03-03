package game

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckrepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	gamerepo "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

// mockGameRepo implements repos.GameRepository
type mockGameRepo struct {
	GetAllByPodFn func(ctx context.Context, podID int) ([]gamerepo.Model, error)
	GetByIdFn     func(ctx context.Context, gameID int) (*gamerepo.Model, error)
	AddFn         func(ctx context.Context, description string, podID, formatID int) (int, error)
}

func (m *mockGameRepo) GetAllByPod(ctx context.Context, podID int) ([]gamerepo.Model, error) {
	if m.GetAllByPodFn != nil {
		return m.GetAllByPodFn(ctx, podID)
	}
	panic("unexpected call to GetAllByPod")
}
func (m *mockGameRepo) GetAllByDeck(ctx context.Context, deckID int) ([]gamerepo.Model, error) {
	panic("unexpected call to GetAllByDeck")
}
func (m *mockGameRepo) GetById(ctx context.Context, gameID int) (*gamerepo.Model, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, gameID)
	}
	panic("unexpected call to GetById")
}
func (m *mockGameRepo) Add(ctx context.Context, description string, podID, formatID int) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, description, podID, formatID)
	}
	panic("unexpected call to Add")
}
func (m *mockGameRepo) BulkAdd(ctx context.Context, games []gamerepo.Model) ([]int, error) {
	panic("unexpected call to BulkAdd")
}
func (m *mockGameRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}

// mockGameResultRepo implements repos.GameResultRepository
type mockGameResultRepo struct {
	BulkAddFn func(ctx context.Context, results []gameresultrepo.Model) error
}

func (m *mockGameResultRepo) GetByGameId(ctx context.Context, gameID int) ([]gameresultrepo.Model, error) {
	panic("unexpected call to GetByGameId")
}
func (m *mockGameResultRepo) GetStatsForPlayer(ctx context.Context, playerID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForPlayer")
}
func (m *mockGameResultRepo) GetStatsForDeck(ctx context.Context, deckID int) (*gameresultrepo.Aggregate, error) {
	panic("unexpected call to GetStatsForDeck")
}
func (m *mockGameResultRepo) BulkAdd(ctx context.Context, results []gameresultrepo.Model) error {
	if m.BulkAddFn != nil {
		return m.BulkAddFn(ctx, results)
	}
	panic("unexpected call to BulkAdd")
}
func (m *mockGameResultRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}

// mockDeckRepo implements repos.DeckRepository
type mockDeckRepo struct {
	GetByIdFn func(ctx context.Context, deckID int) (*deckrepo.Model, error)
	getCalled bool
}

func (m *mockDeckRepo) GetAll(ctx context.Context) ([]deckrepo.Model, error) {
	panic("unexpected call to GetAll")
}
func (m *mockDeckRepo) GetAllForPlayer(ctx context.Context, playerID int) ([]deckrepo.Model, error) {
	panic("unexpected call to GetAllForPlayer")
}
func (m *mockDeckRepo) GetById(ctx context.Context, deckID int) (*deckrepo.Model, error) {
	m.getCalled = true
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, deckID)
	}
	panic("unexpected call to GetById")
}
func (m *mockDeckRepo) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
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

func validInputs() []gameResult.InputEntity {
	return []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 2},
		{DeckID: 11, Place: 2, Kills: 0},
	}
}

func TestCreate_OtherFormat_SkipsDeckFormatCheck(t *testing.T) {
	gameRepo := &mockGameRepo{
		AddFn: func(ctx context.Context, description string, podID, formatID int) (int, error) {
			return 1, nil
		},
	}
	gameResultRepo := &mockGameResultRepo{
		BulkAddFn: func(ctx context.Context, results []gameresultrepo.Model) error {
			return nil
		},
	}
	deckRepo := &mockDeckRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: 2, Name: "other"}, nil
	}

	// DeckID=10 has FormatID=99 which would mismatch format 2, but "other" skips the check
	inputs := []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 0},
	}

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
	err := fn(context.Background(), "Game", 1, 2, inputs)
	require.NoError(t, err)
	assert.False(t, deckRepo.getCalled, "deck repo should not be called for other format")
}

func TestCreate_MatchingFormat_Success(t *testing.T) {
	gameRepo := &mockGameRepo{
		AddFn: func(ctx context.Context, description string, podID, formatID int) (int, error) {
			return 1, nil
		},
	}
	gameResultRepo := &mockGameResultRepo{
		BulkAddFn: func(ctx context.Context, results []gameresultrepo.Model) error {
			return nil
		},
	}
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{ModelBase: base.ModelBase{ID: deckID}, FormatID: 1}, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: 1, Name: "commander"}, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 2},
	}

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
	err := fn(context.Background(), "Game", 1, 1, inputs)
	require.NoError(t, err)
}

func TestCreate_FormatMismatch_Error(t *testing.T) {
	gameRepo := &mockGameRepo{}
	gameResultRepo := &mockGameResultRepo{}
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			// deck has format 2, game has format 1
			return &deckrepo.Model{ModelBase: base.ModelBase{ID: deckID}, FormatID: 2}, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: 1, Name: "commander"}, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 0},
	}

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
	err := fn(context.Background(), "Game", 1, 1, inputs)
	assert.ErrorContains(t, err, "format does not match")
}

func TestCreate_FormatNotFound_Error(t *testing.T) {
	gameRepo := &mockGameRepo{}
	gameResultRepo := &mockGameResultRepo{}
	deckRepo := &mockDeckRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return nil, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 0},
	}

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
	err := fn(context.Background(), "Game", 1, 99, inputs)
	assert.Error(t, err)
}

func TestCreate_InvalidInput_Error(t *testing.T) {
	gameRepo := &mockGameRepo{}
	gameResultRepo := &mockGameResultRepo{}
	deckRepo := &mockDeckRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		panic("should not be called — Validate runs first")
	}

	// DeckID=0 is invalid
	inputs := []gameResult.InputEntity{
		{DeckID: 0, Place: 1, Kills: 0},
	}

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
	err := fn(context.Background(), "Game", 1, 1, inputs)
	assert.ErrorContains(t, err, "deck_id is required")
}

func TestGetByID_NotFound(t *testing.T) {
	gameRepo := &mockGameRepo{
		GetByIdFn: func(ctx context.Context, gameID int) (*gamerepo.Model, error) {
			return nil, nil
		},
	}
	getGameResults := func(ctx context.Context, gameID int) ([]gameResult.Entity, error) {
		panic("should not be called")
	}

	fn := GetByID(zap.NewNop(), gameRepo, getGameResults)
	got, err := fn(context.Background(), 999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetAllByPod_ResultErrorDropsGame(t *testing.T) {
	gameRepo := &mockGameRepo{
		GetAllByPodFn: func(ctx context.Context, podID int) ([]gamerepo.Model, error) {
			return []gamerepo.Model{
				{ModelBase: base.ModelBase{ID: 1}, PodID: podID, FormatID: 1},
				{ModelBase: base.ModelBase{ID: 2}, PodID: podID, FormatID: 1},
			}, nil
		},
	}
	getGameResults := func(ctx context.Context, gameID int) ([]gameResult.Entity, error) {
		if gameID == 1 {
			return nil, errors.New("results error for game 1")
		}
		return []gameResult.Entity{{ID: 10, GameID: 2}}, nil
	}

	fn := GetAllByPod(zap.NewNop(), gameRepo, getGameResults)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	// game 1 dropped due to error, game 2 included
	assert.Len(t, got, 1)
	assert.Equal(t, 2, got[0].ID)
}
