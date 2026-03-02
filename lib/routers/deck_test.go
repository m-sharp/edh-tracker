package routers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/models"
)

func newTestDeckRouter(deck *mockDeckRepo, dc *mockDeckCommanderRepo, format *mockFormatRepo) *DeckRouter {
	return &DeckRouter{
		log:               zap.NewNop(),
		deckRepo:          deck,
		deckCommanderRepo: dc,
		formatRepo:        format,
	}
}

func TestDeckRouter_GetAll_Success(t *testing.T) {
	decks := []models.DeckWithStats{
		{Deck: models.Deck{Name: "Deck1"}},
	}
	deckRepo := &mockDeckRepo{
		GetAllFn: func(ctx context.Context) ([]models.DeckWithStats, error) { return decks, nil },
	}
	router := newTestDeckRouter(deckRepo, &mockDeckCommanderRepo{}, &mockFormatRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/decks", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []models.DeckWithStats
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestDeckRouter_GetAll_Error(t *testing.T) {
	deckRepo := &mockDeckRepo{
		GetAllFn: func(ctx context.Context) ([]models.DeckWithStats, error) {
			return nil, errors.New("db error")
		},
	}
	router := newTestDeckRouter(deckRepo, &mockDeckCommanderRepo{}, &mockFormatRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/decks", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeckRouter_Add_CommanderFormat_WithCommanderID(t *testing.T) {
	commanderID := 5
	deckRepo := &mockDeckRepo{
		AddFn: func(ctx context.Context, playerID int, name string, formatID int) (int, error) { return 10, nil },
	}
	dcRepo := &mockDeckCommanderRepo{
		AddFn: func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
			return 1, nil
		},
	}
	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return &models.Format{Name: "commander"}, nil
		},
	}
	router := newTestDeckRouter(deckRepo, dcRepo, formatRepo)

	body, _ := json.Marshal(newDeckRequest{PlayerID: 1, Name: "My Deck", FormatID: 1, CommanderID: &commanderID})
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestDeckRouter_Add_CommanderFormat_NoCommanderID(t *testing.T) {
	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return &models.Format{Name: "commander"}, nil
		},
	}
	router := newTestDeckRouter(&mockDeckRepo{}, &mockDeckCommanderRepo{}, formatRepo)

	body, _ := json.Marshal(newDeckRequest{PlayerID: 1, Name: "My Deck", FormatID: 1}) // no CommanderID
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeckRouter_Add_OtherFormat_NoCommanderIDAllowed(t *testing.T) {
	deckRepo := &mockDeckRepo{
		AddFn: func(ctx context.Context, playerID int, name string, formatID int) (int, error) { return 10, nil },
	}
	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return &models.Format{Name: "other"}, nil
		},
	}
	router := newTestDeckRouter(deckRepo, &mockDeckCommanderRepo{}, formatRepo)

	body, _ := json.Marshal(newDeckRequest{PlayerID: 1, Name: "My Deck", FormatID: 2}) // no CommanderID
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestDeckRouter_Add_MissingPlayerID(t *testing.T) {
	router := newTestDeckRouter(&mockDeckRepo{}, &mockDeckCommanderRepo{}, &mockFormatRepo{})

	body, _ := json.Marshal(newDeckRequest{Name: "My Deck", FormatID: 1})
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeckRouter_Retire_Success(t *testing.T) {
	deckRepo := &mockDeckRepo{
		RetireFn: func(ctx context.Context, deckID int) error { return nil },
	}
	router := newTestDeckRouter(deckRepo, &mockDeckCommanderRepo{}, &mockFormatRepo{})

	req := httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", nil)
	rr := httptest.NewRecorder()
	router.RetireDeck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeckRouter_Retire_MissingParam(t *testing.T) {
	router := newTestDeckRouter(&mockDeckRepo{}, &mockDeckCommanderRepo{}, &mockFormatRepo{})

	req := httptest.NewRequest(http.MethodPatch, "/api/deck", nil)
	rr := httptest.NewRecorder()
	router.RetireDeck(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeckRouter_Retire_Error(t *testing.T) {
	deckRepo := &mockDeckRepo{
		RetireFn: func(ctx context.Context, deckID int) error { return errors.New("db error") },
	}
	router := newTestDeckRouter(deckRepo, &mockDeckCommanderRepo{}, &mockFormatRepo{})

	req := httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", nil)
	rr := httptest.NewRecorder()
	router.RetireDeck(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
