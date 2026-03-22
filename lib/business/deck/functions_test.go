package deck

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/testHelpers"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	commanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/commander"
	deckRepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	formatRepo "github.com/m-sharp/edh-tracker/lib/repositories/format"
	gameResultRepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	playerRepo "github.com/m-sharp/edh-tracker/lib/repositories/player"
)

func makeFormat(id int, name string) *format.Entity {
	return &format.Entity{ID: id, Name: name}
}

func TestCreate_CommanderFormat_NoCommanderID(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{}
	deckCmdrRepo := &testHelpers.MockDeckCommanderRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return makeFormat(1, "commander"), nil
	}

	fn := Create(dr, deckCmdrRepo, getFormat)
	_, err := fn(context.Background(), 1, "Test Deck", 1, nil, nil)
	assert.ErrorContains(t, err, "commander_id is required")
}

func TestCreate_CommanderFormat_WithCommander(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		AddFn: func(ctx context.Context, playerID int, name string, formatID int) (int, error) {
			return 10, nil
		},
	}
	deckCmdrRepo := &testHelpers.MockDeckCommanderRepo{
		AddFn: func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
			return 1, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return makeFormat(1, "commander"), nil
	}

	commanderID := 5
	fn := Create(dr, deckCmdrRepo, getFormat)
	id, err := fn(context.Background(), 1, "Krenko Goblins", 1, &commanderID, nil)
	require.NoError(t, err)
	assert.Equal(t, 10, id)
}

func TestCreate_OtherFormat_NoCommander(t *testing.T) {
	deckCmdrRepo := &testHelpers.MockDeckCommanderRepo{}
	dr := &testHelpers.MockDeckRepo{
		AddFn: func(ctx context.Context, playerID int, name string, formatID int) (int, error) {
			return 11, nil
		},
	}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return makeFormat(2, "other"), nil
	}

	fn := Create(dr, deckCmdrRepo, getFormat)
	id, err := fn(context.Background(), 1, "Casual Deck", 2, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 11, id)
	assert.False(t, deckCmdrRepo.AddCalled, "deckCmdrRepo.Add should not be called for other format")
}

func TestCreate_FormatNotFound(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{}
	deckCmdrRepo := &testHelpers.MockDeckCommanderRepo{}
	getFormat := func(ctx context.Context, id int) (*format.Entity, error) {
		return nil, nil
	}

	fn := Create(dr, deckCmdrRepo, getFormat)
	_, err := fn(context.Background(), 1, "Test Deck", 99, nil, nil)
	assert.Error(t, err)
}

func TestGetDeckName_Success(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return &deckRepo.Model{GormModelBase: base.GormModelBase{ID: 20}, Name: "Krenko Goblins"}, nil
		},
	}
	fn := GetDeckName(dr)
	name, err := fn(context.Background(), 20)
	require.NoError(t, err)
	assert.Equal(t, "Krenko Goblins", name)
}

func TestGetDeckName_NotFound(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return nil, nil
		},
	}
	fn := GetDeckName(dr)
	_, err := fn(context.Background(), 99)
	assert.Error(t, err)
}

func TestDeckUpdate_Success_NoCommander(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return &deckRepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 42}, nil
		},
		UpdateFn: func(ctx context.Context, deckID int, fields deckRepo.UpdateFields) error {
			return nil
		},
	}
	deckCmdrRepo := &testHelpers.MockDeckCommanderRepo{}

	fn := Update(dr, deckCmdrRepo)
	name := "New Name"
	err := fn(context.Background(), 1, 42, UpdateFields{Name: &name})
	require.NoError(t, err)
}

func TestDeckUpdate_WithCommander(t *testing.T) {
	deleteCalled := false
	addCalled := false
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return &deckRepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 10}, nil
		},
		UpdateFn: func(ctx context.Context, deckID int, fields deckRepo.UpdateFields) error {
			return nil
		},
	}
	deckCmdrRepo := &testHelpers.MockDeckCommanderRepo{
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
	fn := Update(dr, deckCmdrRepo)
	err := fn(context.Background(), 1, 10, UpdateFields{CommanderID: &commanderID})
	require.NoError(t, err)
	assert.True(t, deleteCalled, "DeleteByDeckID should be called when CommanderID is set")
	assert.True(t, addCalled, "Add should be called to set the new commander")
}

func TestDeckUpdate_NotFound(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return nil, nil
		},
	}
	fn := Update(dr, &testHelpers.MockDeckCommanderRepo{})
	err := fn(context.Background(), 99, 42, UpdateFields{})
	assert.ErrorContains(t, err, "not found")
}

func TestDeckUpdate_Forbidden(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return &deckRepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 10}, nil
		},
	}
	fn := Update(dr, &testHelpers.MockDeckCommanderRepo{})
	err := fn(context.Background(), 1, 99, UpdateFields{}) // callerPlayerID=99 != deck.PlayerID=10
	assert.ErrorContains(t, err, "forbidden")
}

func TestDeckSoftDelete_Success(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return &deckRepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 7}, nil
		},
		SoftDeleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	fn := SoftDelete(dr)
	err := fn(context.Background(), 1, 7)
	require.NoError(t, err)
}

func TestDeckSoftDelete_NotFound(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return nil, nil
		},
	}
	fn := SoftDelete(dr)
	err := fn(context.Background(), 99, 7)
	assert.ErrorContains(t, err, "not found")
}

func TestDeckSoftDelete_Forbidden(t *testing.T) {
	dr := &testHelpers.MockDeckRepo{
		GetByIdFn: func(ctx context.Context, deckID int) (*deckRepo.Model, error) {
			return &deckRepo.Model{GormModelBase: base.GormModelBase{ID: deckID}, PlayerID: 10}, nil
		},
	}
	fn := SoftDelete(dr)
	err := fn(context.Background(), 1, 99) // caller=99, owner=10
	assert.ErrorContains(t, err, "forbidden")
}

func TestDeckGetAllByPod_Success(t *testing.T) {
	partnerID := 6
	podRepo := &testHelpers.MockPodRepo{
		GetPlayerIDsFn: func(ctx context.Context, podID int) ([]int, error) {
			return []int{1, 2}, nil
		},
	}
	dr := &testHelpers.MockDeckRepo{
		GetAllByPlayerIDsFn: func(ctx context.Context, playerIDs []int) ([]deckRepo.Model, error) {
			return []deckRepo.Model{
				{
					GormModelBase: base.GormModelBase{ID: 10},
					PlayerID:      1,
					FormatID:      1,
					Name:          "Deck A",
					Player:        playerRepo.Model{Name: "Alice"},
					Format:        formatRepo.Model{Name: "commander"},
					Commander: &deckCommanderRepo.Model{
						CommanderID:        5,
						PartnerCommanderID: &partnerID,
						Commander:          commanderRepo.Model{Name: "Krenko"},
						PartnerCommander:   &commanderRepo.Model{Name: "Goblin Chieftain"},
					},
				},
			}, nil
		},
	}
	grRepo := &testHelpers.MockGameResultRepo{
		GetStatsForDeckFn: func(ctx context.Context, deckID int) (*gameResultRepo.Aggregate, error) {
			return nil, nil
		},
	}

	fn := GetAllByPod(dr, podRepo, grRepo)
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "Deck A", got[0].Name)
	assert.Equal(t, "Alice", got[0].PlayerName)
	assert.Equal(t, "Commander", got[0].FormatName)
	require.NotNil(t, got[0].Commanders)
	assert.Equal(t, "Krenko", got[0].Commanders.CommanderName)
	require.NotNil(t, got[0].Commanders.PartnerCommanderName)
	assert.Equal(t, "Goblin Chieftain", *got[0].Commanders.PartnerCommanderName)
}

func TestDeckGetAllByPod_EmptyPod(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		GetPlayerIDsFn: func(ctx context.Context, podID int) ([]int, error) {
			return []int{}, nil
		},
	}
	fn := GetAllByPod(&testHelpers.MockDeckRepo{}, podRepo, &testHelpers.MockGameResultRepo{})
	got, err := fn(context.Background(), 5)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestDeckGetAllByPod_Error(t *testing.T) {
	podRepo := &testHelpers.MockPodRepo{
		GetPlayerIDsFn: func(ctx context.Context, podID int) ([]int, error) {
			return nil, errors.New("db error")
		},
	}
	fn := GetAllByPod(&testHelpers.MockDeckRepo{}, podRepo, &testHelpers.MockGameResultRepo{})
	_, err := fn(context.Background(), 5)
	assert.Error(t, err)
}
