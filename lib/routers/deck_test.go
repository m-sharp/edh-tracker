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

	"github.com/m-sharp/edh-tracker/lib/business/deck"
)

func newTestDeckRouter(decks deck.Functions) *DeckRouter {
	return &DeckRouter{
		log:   zap.NewNop(),
		decks: decks,
	}
}

func TestDeckRouter_GetAll_Success(t *testing.T) {
	decks := []deck.EntityWithStats{
		{Entity: deck.Entity{Name: "Deck1"}},
	}
	router := newTestDeckRouter(deck.Functions{
		GetAll: func(ctx context.Context) ([]deck.EntityWithStats, error) { return decks, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []deck.EntityWithStats
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestDeckRouter_GetAll_Error(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		GetAll: func(ctx context.Context) ([]deck.EntityWithStats, error) {
			return nil, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeckRouter_Add_Success(t *testing.T) {
	commanderID := 5
	router := newTestDeckRouter(deck.Functions{
		Create: func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error) {
			return 10, nil
		},
	})

	body, _ := json.Marshal(newDeckRequest{PlayerID: 1, Name: "My Deck", FormatID: 1, CommanderID: &commanderID})
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestDeckRouter_Add_CreateError(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		Create: func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error) {
			return 0, errors.New("commander_id is required for commander format decks")
		},
	})

	body, _ := json.Marshal(newDeckRequest{PlayerID: 1, Name: "My Deck", FormatID: 1}) // no CommanderID
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeckRouter_Add_MissingPlayerID(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{})

	body, _ := json.Marshal(newDeckRequest{Name: "My Deck", FormatID: 1})
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeckRouter_Retire_Success(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		Retire: func(ctx context.Context, deckID int) error { return nil },
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", nil)
	rr := httptest.NewRecorder()
	router.RetireDeck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeckRouter_Retire_MissingParam(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{})

	req := httptest.NewRequest(http.MethodPatch, "/api/deck", nil)
	rr := httptest.NewRecorder()
	router.RetireDeck(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeckRouter_Retire_Error(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		Retire: func(ctx context.Context, deckID int) error { return errors.New("db error") },
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", nil)
	rr := httptest.NewRecorder()
	router.RetireDeck(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
