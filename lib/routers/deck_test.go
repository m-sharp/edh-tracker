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

	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/deck"
)

func newTestDeckRouter(decks deck.Functions) *DeckRouter {
	return &DeckRouter{
		log:   zap.NewNop(),
		decks: decks,
	}
}

func TestDeckRouter_GetAll_NoFilter_Returns400(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{})

	req := httptest.NewRequest(http.MethodGet, "/api/decks", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "pod_id or player_id query param is required")
}

func TestDeckRouter_GetAll_Paginated_NoFilter_Returns400(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{})

	req := httptest.NewRequest(http.MethodGet, "/api/decks?limit=10", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "pod_id or player_id query param is required")
}

func TestDeckRouter_Add_Success(t *testing.T) {
	commanderID := 5
	router := newTestDeckRouter(deck.Functions{
		Create: func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error) {
			assert.Equal(t, 42, playerID, "playerID should come from JWT, not body")
			return 10, nil
		},
	})

	body, _ := json.Marshal(newDeckRequest{Name: "My Deck", FormatID: 1, CommanderID: &commanderID})
	r := withAuth(httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestDeckRouter_Add_NoAuth(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{})

	body, _ := json.Marshal(newDeckRequest{Name: "My Deck", FormatID: 1})
	r := httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeckRouter_Add_CreateError(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		Create: func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error) {
			return 0, errors.New("commander_id is required for commander format decks")
		},
	})

	body, _ := json.Marshal(newDeckRequest{Name: "My Deck", FormatID: 1}) // no CommanderID
	r := withAuth(httptest.NewRequest(http.MethodPost, "/api/deck", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.DeckCreate(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeckRouter_Update_Success(t *testing.T) {
	retired := true
	callerID := 42
	router := newTestDeckRouter(deck.Functions{
		GetByID: func(ctx context.Context, deckID int) (*deck.EntityWithStats, error) {
			return &deck.EntityWithStats{Entity: deck.Entity{ID: deckID, PlayerID: callerID}}, nil
		},
		Update: func(ctx context.Context, deckID int, fields deck.UpdateFields) error {
			return nil
		},
	})

	body, _ := json.Marshal(updateDeckRequest{Retired: &retired})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", bytes.NewReader(body)), callerID)
	rr := httptest.NewRecorder()
	router.UpdateDeck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeckRouter_Update_MissingParam(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{})

	body, _ := json.Marshal(updateDeckRequest{})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/deck", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdateDeck(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeckRouter_Update_Forbidden(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		GetByID: func(ctx context.Context, deckID int) (*deck.EntityWithStats, error) {
			return &deck.EntityWithStats{Entity: deck.Entity{ID: deckID, PlayerID: 99}}, nil
		},
	})

	body, _ := json.Marshal(updateDeckRequest{})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdateDeck(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeckRouter_Update_Error(t *testing.T) {
	callerID := 42
	router := newTestDeckRouter(deck.Functions{
		GetByID: func(ctx context.Context, deckID int) (*deck.EntityWithStats, error) {
			return &deck.EntityWithStats{Entity: deck.Entity{ID: deckID, PlayerID: callerID}}, nil
		},
		Update: func(ctx context.Context, deckID int, fields deck.UpdateFields) error {
			return errors.New("db error")
		},
	})

	body, _ := json.Marshal(updateDeckRequest{})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/deck?deck_id=1", bytes.NewReader(body)), callerID)
	rr := httptest.NewRecorder()
	router.UpdateDeck(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeckRouter_GetAll_ByPod_Success(t *testing.T) {
	decks := []deck.EntityWithStats{
		{Entity: deck.Entity{Name: "Pod Deck"}},
	}
	router := newTestDeckRouter(deck.Functions{
		GetAllByPod: func(ctx context.Context, podID int) ([]deck.EntityWithStats, error) { return decks, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []deck.EntityWithStats
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
	assert.Equal(t, "Pod Deck", got[0].Name)
}

func TestDeckRouter_GetAll_Paginated_ByPod(t *testing.T) {
	entities := []deck.EntityWithStats{
		{Entity: deck.Entity{Name: "Deck A"}},
		{Entity: deck.Entity{Name: "Deck B"}},
	}
	router := newTestDeckRouter(deck.Functions{
		GetAllByPodPaginated: func(ctx context.Context, podID, limit, offset int) ([]deck.EntityWithStats, int, error) {
			assert.Equal(t, 3, podID)
			assert.Equal(t, 10, limit)
			assert.Equal(t, 0, offset)
			return entities, 30, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks?pod_id=3&limit=10&offset=0", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got business.PaginatedResponse[deck.EntityWithStats]
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, 30, got.Total)
	assert.Equal(t, 10, got.Limit)
	assert.Equal(t, 0, got.Offset)
	assert.Len(t, got.Items, 2)
}

func TestDeckRouter_GetAll_Paginated_ByPlayer(t *testing.T) {
	entities := []deck.EntityWithStats{
		{Entity: deck.Entity{Name: "Player Deck"}},
	}
	router := newTestDeckRouter(deck.Functions{
		GetAllByPlayerPaginated: func(ctx context.Context, playerID, limit, offset int) ([]deck.EntityWithStats, int, error) {
			return entities, 5, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks?player_id=7&limit=5&offset=0", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got business.PaginatedResponse[deck.EntityWithStats]
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, 5, got.Total)
	assert.Len(t, got.Items, 1)
}

func TestDeckRouter_GetAll_Paginated_Error(t *testing.T) {
	router := newTestDeckRouter(deck.Functions{
		GetAllByPodPaginated: func(ctx context.Context, podID, limit, offset int) ([]deck.EntityWithStats, int, error) {
			return nil, 0, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks?pod_id=1&limit=10", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeckRouter_GetAll_NoLimit_PlainArray(t *testing.T) {
	decks := []deck.EntityWithStats{
		{Entity: deck.Entity{Name: "D1"}},
		{Entity: deck.Entity{Name: "D2"}},
	}
	router := newTestDeckRouter(deck.Functions{
		GetAllByPod: func(ctx context.Context, podID int) ([]deck.EntityWithStats, error) { return decks, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/decks?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []deck.EntityWithStats
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 2)
}
