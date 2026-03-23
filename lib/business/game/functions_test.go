package game

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckrepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	gamerepo "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	repoTestHelpers "github.com/m-sharp/edh-tracker/lib/repositories/testHelpers"
)

func validInputs() []gameResult.InputEntity {
	return []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 2},
		{DeckID: 11, Place: 2, Kills: 0},
	}
}

// newTestClient creates a lib.DBClient from the integration test DB.
func newTestClient(t *testing.T) (*lib.DBClient, *gorm.DB) {
	t.Helper()
	db := repoTestHelpers.NewTestDB(t)
	return &lib.DBClient{GormDb: db}, db
}

// TestCreate_OtherFormat_SkipsDeckFormatCheck verifies that when format is "other",
// no deck format validation is performed, and the game + results are created atomically.
func TestCreate_OtherFormat_SkipsDeckFormatCheck(t *testing.T) {
	client, db := newTestClient(t)

	// Create a real deck and pod to satisfy FK constraints.
	testDeck := repoTestHelpers.CreateTestDeck(t, db)
	podID := repoTestHelpers.CreateTestPod(t, db)
	otherFormatID := repoTestHelpers.GetCommanderFormatID(t, db) // use commander format but name "other" via mock

	deckRepo := &testHelpers.MockDeckRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: otherFormatID, Name: "other"}, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: testDeck.ID, Place: 1, Kills: 0},
	}

	fn := Create(zap.NewNop(), nil, nil, deckRepo, getFormat, client)
	err := fn(context.Background(), "Game", podID, otherFormatID, inputs)
	require.NoError(t, err)
	assert.False(t, deckRepo.GetByIdCalled, "deck repo should not be called for other format")
}

// TestCreate_MatchingFormat_Success verifies that when deck format matches game format,
// the game and results are created successfully.
func TestCreate_MatchingFormat_Success(t *testing.T) {
	client, db := newTestClient(t)

	testDeck := repoTestHelpers.CreateTestDeck(t, db)
	podID := repoTestHelpers.CreateTestPod(t, db)
	formatID := repoTestHelpers.GetCommanderFormatID(t, db)

	deckRepo := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, FormatID: formatID}, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: formatID, Name: "commander"}, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: testDeck.ID, Place: 1, Kills: 2},
	}

	fn := Create(zap.NewNop(), nil, nil, deckRepo, getFormat, client)
	err := fn(context.Background(), "Game", podID, formatID, inputs)
	require.NoError(t, err)
}

func TestCreate_FormatMismatch_Error(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{}
	gameResultRepo := &testHelpers.MockGameResultRepo{}
	deckRepo := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			// deck has format 2, game has format 1
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, FormatID: 2}, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return &format.Entity{ID: 1, Name: "commander"}, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 0},
	}

	// nil client is safe here because the function returns before reaching the transaction.
	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat, nil)
	err := fn(context.Background(), "Game", 1, 1, inputs)
	assert.ErrorContains(t, err, "format does not match")
}

func TestCreate_FormatNotFound_Error(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{}
	gameResultRepo := &testHelpers.MockGameResultRepo{}
	deckRepo := &testHelpers.MockDeckRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return nil, nil
	}

	inputs := []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 0},
	}

	// nil client is safe here because the function returns before reaching the transaction.
	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat, nil)
	err := fn(context.Background(), "Game", 1, 99, inputs)
	assert.Error(t, err)
}

func TestCreate_InvalidInput_Error(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{}
	gameResultRepo := &testHelpers.MockGameResultRepo{}
	deckRepo := &testHelpers.MockDeckRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		panic("should not be called — Validate runs first")
	}

	// DeckID=0 is invalid
	inputs := []gameResult.InputEntity{
		{DeckID: 0, Place: 1, Kills: 0},
	}

	// nil client is safe here because the function returns before reaching the transaction.
	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat, nil)
	err := fn(context.Background(), "Game", 1, 1, inputs)
	assert.ErrorContains(t, err, "deck_id is required")
}

func TestGetByID_NotFound(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetByIDFn: func(ctx context.Context, gameID int) (*gamerepo.Model, error) {
			return nil, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		panic("should not be called")
	}

	fn := GetByID(zap.NewNop(), gameRepo, enrichGameResults)
	got, err := fn(context.Background(), 999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestGetAllByPod_ResultErrorDropsGame(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPodFn: func(ctx context.Context, podID int) ([]gamerepo.Model, error) {
			return []gamerepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, PodID: podID, FormatID: 1, Results: []gameresultrepo.Model{{GameID: 1, DeckID: 1}}},
				{GormModelBase: base.GormModelBase{ID: 2}, PodID: podID, FormatID: 1, Results: []gameresultrepo.Model{{GameID: 2, DeckID: 2}}},
			}, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		if len(models) > 0 && models[0].GameID == 1 {
			return nil, errors.New("results error for game 1")
		}
		return []gameResult.Entity{{ID: 10, GameID: 2}}, nil
	}

	fn := GetAllByPod(zap.NewNop(), gameRepo, enrichGameResults)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	// game 1 dropped due to error, game 2 included
	assert.Len(t, got, 1)
	assert.Equal(t, 2, got[0].ID)
}

func TestGetAllByPlayer_ResultErrorDropsGame(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPlayerIDFn: func(ctx context.Context, playerID int) ([]gamerepo.Model, error) {
			return []gamerepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Results: []gameresultrepo.Model{{GameID: 1, DeckID: 1}}},
				{GormModelBase: base.GormModelBase{ID: 2}, Results: []gameresultrepo.Model{{GameID: 2, DeckID: 2}}},
			}, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		if len(models) > 0 && models[0].GameID == 1 {
			return nil, errors.New("result error")
		}
		return []gameResult.Entity{{ID: 10, GameID: 2}}, nil
	}

	fn := GetAllByPlayer(zap.NewNop(), gameRepo, enrichGameResults)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, 2, got[0].ID)
}

func TestGetAllByPlayer_RepoError(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPlayerIDFn: func(ctx context.Context, playerID int) ([]gamerepo.Model, error) {
			return nil, errors.New("db error")
		},
	}
	fn := GetAllByPlayer(zap.NewNop(), gameRepo, nil)
	_, err := fn(context.Background(), 5)
	assert.Error(t, err)
}

func TestGameUpdate_Success(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		UpdateFn: func(ctx context.Context, gameID int, description string) error {
			return nil
		},
	}
	fn := Update(gameRepo)
	err := fn(context.Background(), 1, "Updated description")
	require.NoError(t, err)
}

func TestGameUpdate_Error(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		UpdateFn: func(ctx context.Context, gameID int, description string) error {
			return errors.New("db error")
		},
	}
	fn := Update(gameRepo)
	err := fn(context.Background(), 1, "Updated description")
	assert.Error(t, err)
}

func TestGameSoftDelete_Success(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		SoftDeleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	fn := SoftDelete(gameRepo)
	err := fn(context.Background(), 1)
	require.NoError(t, err)
}

func TestAddResult_Success(t *testing.T) {
	gameResultRepo := &testHelpers.MockGameResultRepo{
		AddFn: func(ctx context.Context, model gameresultrepo.Model) (int, error) {
			return 99, nil
		},
	}
	fn := AddResult(gameResultRepo)
	id, err := fn(context.Background(), 1, 10, 42, 2, 1)
	require.NoError(t, err)
	assert.Equal(t, 99, id)
}

func TestUpdateResult_Success(t *testing.T) {
	gameResultRepo := &testHelpers.MockGameResultRepo{
		UpdateFn: func(ctx context.Context, resultID, place, killCount, deckID int) error {
			return nil
		},
	}
	fn := UpdateResult(gameResultRepo)
	err := fn(context.Background(), 1, 2, 1, 10)
	require.NoError(t, err)
}

func TestDeleteResult_Success(t *testing.T) {
	gameResultRepo := &testHelpers.MockGameResultRepo{
		SoftDeleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	fn := DeleteResult(gameResultRepo)
	err := fn(context.Background(), 1)
	require.NoError(t, err)
}

func TestGetAllByPodPaginated_Success(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPodPaginatedFn: func(ctx context.Context, podID, limit, offset int) ([]gamerepo.Model, int, error) {
			return []gamerepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, PodID: podID, FormatID: 1},
			}, 10, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		return []gameResult.Entity{}, nil
	}

	fn := GetAllByPodPaginated(zap.NewNop(), gameRepo, enrichGameResults)
	got, total, err := fn(context.Background(), 5, 1, 0)
	require.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Len(t, got, 1)
	assert.Equal(t, 1, got[0].ID)
}

func TestGetAllByPodPaginated_RepoError(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPodPaginatedFn: func(ctx context.Context, podID, limit, offset int) ([]gamerepo.Model, int, error) {
			return nil, 0, errors.New("db error")
		},
	}
	fn := GetAllByPodPaginated(zap.NewNop(), gameRepo, nil)
	_, _, err := fn(context.Background(), 5, 10, 0)
	assert.Error(t, err)
}

func TestGetAllByPodPaginated_ResultErrorDropsGame(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPodPaginatedFn: func(ctx context.Context, podID, limit, offset int) ([]gamerepo.Model, int, error) {
			return []gamerepo.Model{
				{GormModelBase: base.GormModelBase{ID: 1}, Results: []gameresultrepo.Model{{GameID: 1}}},
				{GormModelBase: base.GormModelBase{ID: 2}, Results: []gameresultrepo.Model{{GameID: 2}}},
			}, 5, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		if len(models) > 0 && models[0].GameID == 1 {
			return nil, errors.New("enrich error")
		}
		return []gameResult.Entity{}, nil
	}

	fn := GetAllByPodPaginated(zap.NewNop(), gameRepo, enrichGameResults)
	got, total, err := fn(context.Background(), 5, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, got, 1)
	assert.Equal(t, 2, got[0].ID)
}

func TestGetAllByDeckPaginated_Success(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByDeckPaginatedFn: func(ctx context.Context, deckID, limit, offset int) ([]gamerepo.Model, int, error) {
			return []gamerepo.Model{
				{GormModelBase: base.GormModelBase{ID: 3}, FormatID: 1},
			}, 7, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		return []gameResult.Entity{}, nil
	}

	fn := GetAllByDeckPaginated(zap.NewNop(), gameRepo, enrichGameResults)
	got, total, err := fn(context.Background(), 10, 5, 0)
	require.NoError(t, err)
	assert.Equal(t, 7, total)
	assert.Len(t, got, 1)
	assert.Equal(t, 3, got[0].ID)
}

func TestGetAllByPlayerIDPaginated_Success(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetAllByPlayerIDPaginatedFn: func(ctx context.Context, playerID, limit, offset int) ([]gamerepo.Model, int, error) {
			return []gamerepo.Model{
				{GormModelBase: base.GormModelBase{ID: 5}, FormatID: 1},
				{GormModelBase: base.GormModelBase{ID: 6}, FormatID: 1},
			}, 20, nil
		},
	}
	enrichGameResults := func(ctx context.Context, models []gameresultrepo.Model) ([]gameResult.Entity, error) {
		return []gameResult.Entity{}, nil
	}

	fn := GetAllByPlayerIDPaginated(zap.NewNop(), gameRepo, enrichGameResults)
	got, total, err := fn(context.Background(), 42, 10, 10)
	require.NoError(t, err)
	assert.Equal(t, 20, total)
	assert.Len(t, got, 2)
}
