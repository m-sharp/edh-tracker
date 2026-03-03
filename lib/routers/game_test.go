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

	"github.com/m-sharp/edh-tracker/lib/business/game"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
)

func newTestGameRouter(games game.Functions) *GameRouter {
	return &GameRouter{
		log:   zap.NewNop(),
		games: games,
	}
}

// gameCreateBody builds a minimal valid createGameRequest JSON body.
func gameCreateBody(t *testing.T, formatID, deckID int) *bytes.Reader {
	t.Helper()
	req := createGameRequest{
		Description: "Test Game",
		PodID:       1,
		FormatID:    formatID,
		Results: []gameResult.InputEntity{
			{DeckID: deckID, Place: 1, Kills: 0},
		},
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	return bytes.NewReader(b)
}

func TestGameRouter_GetGames_ByPod(t *testing.T) {
	games := []game.Entity{{Description: "Game 1"}}
	router := newTestGameRouter(game.Functions{
		GetAllByPod: func(ctx context.Context, podId int) ([]game.Entity, error) {
			return games, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []game.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestGameRouter_GetGames_ByDeck(t *testing.T) {
	games := []game.Entity{{Description: "Game 1"}}
	router := newTestGameRouter(game.Functions{
		GetAllByDeck: func(ctx context.Context, deckId int) ([]game.Entity, error) {
			return games, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?deck_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []game.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestGameRouter_GetGames_MissingQueryParam(t *testing.T) {
	router := newTestGameRouter(game.Functions{})

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_GetGame_Success(t *testing.T) {
	g := &game.Entity{Description: "My Game"}
	router := newTestGameRouter(game.Functions{
		GetByID: func(ctx context.Context, gameId int) (*game.Entity, error) { return g, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/game?game_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetGameById(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got game.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "My Game", got.Description)
}

func TestGameRouter_GetGame_MissingParam(t *testing.T) {
	router := newTestGameRouter(game.Functions{})

	req := httptest.NewRequest(http.MethodGet, "/api/game", nil)
	rr := httptest.NewRecorder()
	router.GetGameById(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_GetGame_Error(t *testing.T) {
	router := newTestGameRouter(game.Functions{
		GetByID: func(ctx context.Context, gameId int) (*game.Entity, error) {
			return nil, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/game?game_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetGameById(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGameRouter_Add_Success(t *testing.T) {
	router := newTestGameRouter(game.Functions{
		Create: func(ctx context.Context, description string, podID, formatID int, results []gameResult.InputEntity) error {
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, 1, 10))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGameRouter_Add_CreateError(t *testing.T) {
	router := newTestGameRouter(game.Functions{
		Create: func(ctx context.Context, description string, podID, formatID int, results []gameResult.InputEntity) error {
			return errors.New("format not found")
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, 99, 10))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGameRouter_Add_EmptyResults(t *testing.T) {
	router := newTestGameRouter(game.Functions{})

	body, _ := json.Marshal(createGameRequest{
		Description: "Test",
		PodID:       1,
		FormatID:    1,
		Results:     []gameResult.InputEntity{},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_Add_InvalidResult(t *testing.T) {
	router := newTestGameRouter(game.Functions{})

	body, _ := json.Marshal(createGameRequest{
		Description: "Test",
		PodID:       1,
		FormatID:    1,
		Results: []gameResult.InputEntity{
			{DeckID: 0, Place: 1, Kills: 0}, // DeckID=0 is invalid
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/game", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
