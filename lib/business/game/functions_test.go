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
	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckrepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	gamerepo "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameresultrepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

func validInputs() []gameResult.InputEntity {
	return []gameResult.InputEntity{
		{DeckID: 10, Place: 1, Kills: 2},
		{DeckID: 11, Place: 2, Kills: 0},
	}
}

func TestCreate_OtherFormat_SkipsDeckFormatCheck(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		AddFn: func(ctx context.Context, description string, podID, formatID int) (int, error) {
			return 1, nil
		},
	}
	gameResultRepo := &testHelpers.MockGameResultRepo{
		BulkAddFn: func(ctx context.Context, results []gameresultrepo.Model) error {
			return nil
		},
	}
	deckRepo := &testHelpers.MockDeckRepo{}
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
	assert.False(t, deckRepo.GetByIdCalled, "deck repo should not be called for other format")
}

func TestCreate_MatchingFormat_Success(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		AddFn: func(ctx context.Context, description string, podID, formatID int) (int, error) {
			return 1, nil
		},
	}
	gameResultRepo := &testHelpers.MockGameResultRepo{
		BulkAddFn: func(ctx context.Context, results []gameresultrepo.Model) error {
			return nil
		},
	}
	deckRepo := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckrepo.Model, error) {
			return &deckrepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, FormatID: 1}, nil
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

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
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

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
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

	fn := Create(zap.NewNop(), gameRepo, gameResultRepo, deckRepo, getFormat)
	err := fn(context.Background(), "Game", 1, 1, inputs)
	assert.ErrorContains(t, err, "deck_id is required")
}

func TestGetByID_NotFound(t *testing.T) {
	gameRepo := &testHelpers.MockGameRepo{
		GetByIDWithResultsFn: func(ctx context.Context, gameID int) (*gamerepo.Model, error) {
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
		GetAllByPodWithResultsFn: func(ctx context.Context, podID int) ([]gamerepo.Model, error) {
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
		GetAllByPlayerWithResultsFn: func(ctx context.Context, playerID int) ([]gamerepo.Model, error) {
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
		GetAllByPlayerWithResultsFn: func(ctx context.Context, playerID int) ([]gamerepo.Model, error) {
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
