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

func newTestGameRouter(game *mockGameRepo, format *mockFormatRepo, deck *mockDeckRepo) *GameRouter {
	return &GameRouter{
		log:        zap.NewNop(),
		gameRepo:   game,
		formatRepo: format,
		deckRepo:   deck,
	}
}

// gameCreateBody builds a minimal valid GameDetails JSON body.
func gameCreateBody(t *testing.T, formatID, deckID int) *bytes.Reader {
	t.Helper()
	details := models.GameDetails{
		Game: models.Game{
			Description: "Test Game",
			PodID:       1,
			FormatID:    formatID,
		},
		Results: []models.GameResult{
			{DeckId: deckID, Place: 1, Kills: 0},
		},
	}
	b, err := json.Marshal(details)
	require.NoError(t, err)
	return bytes.NewReader(b)
}

func TestGameRouter_GetGames_ByPod(t *testing.T) {
	games := []models.GameDetails{{Game: models.Game{Description: "Game 1"}}}
	gameRepo := &mockGameRepo{
		GetAllByPodFn: func(ctx context.Context, podId int) ([]models.GameDetails, error) {
			return games, nil
		},
	}
	router := newTestGameRouter(gameRepo, &mockFormatRepo{}, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/games?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []models.GameDetails
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestGameRouter_GetGames_ByDeck(t *testing.T) {
	games := []models.GameDetails{{Game: models.Game{Description: "Game 1"}}}
	gameRepo := &mockGameRepo{
		GetAllByDeckFn: func(ctx context.Context, deckId int) ([]models.GameDetails, error) {
			return games, nil
		},
	}
	router := newTestGameRouter(gameRepo, &mockFormatRepo{}, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/games?deck_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []models.GameDetails
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestGameRouter_GetGames_MissingQueryParam(t *testing.T) {
	router := newTestGameRouter(&mockGameRepo{}, &mockFormatRepo{}, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_GetGame_Success(t *testing.T) {
	game := &models.GameDetails{Game: models.Game{Description: "My Game"}}
	gameRepo := &mockGameRepo{
		GetGameByIdFn: func(ctx context.Context, gameId int) (*models.GameDetails, error) {
			return game, nil
		},
	}
	router := newTestGameRouter(gameRepo, &mockFormatRepo{}, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/game?game_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetGameById(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got models.GameDetails
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "My Game", got.Description)
}

func TestGameRouter_GetGame_MissingParam(t *testing.T) {
	router := newTestGameRouter(&mockGameRepo{}, &mockFormatRepo{}, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/game", nil)
	rr := httptest.NewRecorder()
	router.GetGameById(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_GetGame_Error(t *testing.T) {
	gameRepo := &mockGameRepo{
		GetGameByIdFn: func(ctx context.Context, gameId int) (*models.GameDetails, error) {
			return nil, errors.New("db error")
		},
	}
	router := newTestGameRouter(gameRepo, &mockFormatRepo{}, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/game?game_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetGameById(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGameRouter_Add_FormatMatch(t *testing.T) {
	const formatID = 1
	const deckID = 10
	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return &models.Format{Name: "commander"}, nil
		},
	}
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.DeckWithStats, error) {
			return &models.DeckWithStats{Deck: models.Deck{FormatID: formatID}}, nil
		},
	}
	gameRepo := &mockGameRepo{
		AddFn: func(ctx context.Context, description string, podID, fID int, results ...models.GameResult) error {
			return nil
		},
	}
	router := newTestGameRouter(gameRepo, formatRepo, deckRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, formatID, deckID))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGameRouter_Add_FormatMismatch(t *testing.T) {
	const gameFormatID = 1
	const deckFormatID = 2
	const deckID = 10

	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return &models.Format{Name: "commander"}, nil
		},
	}
	deckRepo := &mockDeckRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.DeckWithStats, error) {
			return &models.DeckWithStats{Deck: models.Deck{FormatID: deckFormatID}}, nil
		},
	}
	router := newTestGameRouter(&mockGameRepo{}, formatRepo, deckRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, gameFormatID, deckID))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_Add_OtherFormat_SkipsFormatCheck(t *testing.T) {
	const gameFormatID = 2
	const deckFormatID = 1 // differs from game, but "other" format skips the check
	const deckID = 10

	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return &models.Format{Name: "other"}, nil
		},
	}
	gameRepo := &mockGameRepo{
		AddFn: func(ctx context.Context, description string, podID, fID int, results ...models.GameResult) error {
			return nil
		},
	}
	// deckRepo should NOT be called for "other" format — leave it with nil fns so a call would panic
	router := newTestGameRouter(gameRepo, formatRepo, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, gameFormatID, deckID))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGameRouter_Add_InvalidFormat(t *testing.T) {
	formatRepo := &mockFormatRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.Format, error) {
			return nil, nil // format not found
		},
	}
	router := newTestGameRouter(&mockGameRepo{}, formatRepo, &mockDeckRepo{})

	req := httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, 99, 1))
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
