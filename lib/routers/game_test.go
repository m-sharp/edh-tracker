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
	router := newFullGameRouter(
		game.Functions{
			Create: func(ctx context.Context, description string, podID, formatID int, results []gameResult.InputEntity) error {
				return nil
			},
		},
		gameResult.Functions{},
		func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	)

	req := withAuth(httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, 1, 10)), 42)
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGameRouter_Add_CreateError(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			Create: func(ctx context.Context, description string, podID, formatID int, results []gameResult.InputEntity) error {
				return errors.New("format not found")
			},
		},
		gameResult.Functions{},
		func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	)

	req := withAuth(httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, 99, 10)), 42)
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGameRouter_Add_NonMember_Forbidden(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{},
		gameResult.Functions{},
		func(ctx context.Context, podID, playerID int) (string, error) { return "", nil },
	)

	req := withAuth(httptest.NewRequest(http.MethodPost, "/api/game", gameCreateBody(t, 1, 10)), 42)
	rr := httptest.NewRecorder()
	router.GameCreate(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
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

func newFullGameRouter(games game.Functions, gameResults gameResult.Functions, getPodRole func(ctx context.Context, podID, playerID int) (string, error)) *GameRouter {
	return &GameRouter{
		log:         zap.NewNop(),
		games:       games,
		gameResults: gameResults,
		getPodRole:  getPodRole,
	}
}

// podManagerRole returns a getPodRole func that always grants manager.
func podManagerRole() func(ctx context.Context, podID, playerID int) (string, error) {
	return func(ctx context.Context, podID, playerID int) (string, error) {
		return "manager", nil
	}
}

// gameEntity returns a minimal game.Entity for a given podID.
func gameEntityForPod(podID int) *game.Entity {
	return &game.Entity{ID: 1, PodID: podID}
}

func TestGameRouter_GetGames_ByPlayer(t *testing.T) {
	games := []game.Entity{{Description: "Game 1"}}
	router := newTestGameRouter(game.Functions{
		GetAllByPlayer: func(ctx context.Context, playerID int) ([]game.Entity, error) {
			return games, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?player_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []game.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
}

func TestGameRouter_UpdateGame_Success(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			GetByID: func(ctx context.Context, gameID int) (*game.Entity, error) { return gameEntityForPod(1), nil },
			Update:  func(ctx context.Context, gameID int, description string) error { return nil },
		},
		gameResult.Functions{},
		podManagerRole(),
	)

	body, _ := json.Marshal(struct {
		Description string `json:"description"`
	}{Description: "Updated"})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/game?game_id=1", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdateGame(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGameRouter_UpdateGame_NotManager(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			GetByID: func(ctx context.Context, gameID int) (*game.Entity, error) { return gameEntityForPod(1), nil },
		},
		gameResult.Functions{},
		func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	)

	body, _ := json.Marshal(struct{ Description string }{Description: "x"})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/game?game_id=1", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdateGame(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestGameRouter_UpdateGame_MissingParam(t *testing.T) {
	router := newFullGameRouter(game.Functions{}, gameResult.Functions{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/game", nil), 42)
	rr := httptest.NewRecorder()
	router.UpdateGame(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_DeleteGame_Success(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			GetByID:    func(ctx context.Context, gameID int) (*game.Entity, error) { return gameEntityForPod(1), nil },
			SoftDelete: func(ctx context.Context, gameID int) error { return nil },
		},
		gameResult.Functions{},
		podManagerRole(),
	)

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/api/game?game_id=1", nil), 42)
	rr := httptest.NewRecorder()
	router.DeleteGame(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGameRouter_DeleteGame_MissingParam(t *testing.T) {
	router := newFullGameRouter(game.Functions{}, gameResult.Functions{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/api/game", nil), 42)
	rr := httptest.NewRecorder()
	router.DeleteGame(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_AddGameResult_Success(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			GetByID:   func(ctx context.Context, gameID int) (*game.Entity, error) { return gameEntityForPod(1), nil },
			AddResult: func(ctx context.Context, gameID, deckID, playerID, place, killCount int) (int, error) { return 1, nil },
		},
		gameResult.Functions{},
		podManagerRole(),
	)

	body, _ := json.Marshal(addGameResultRequest{GameID: 1, DeckID: 10, PlayerID: 42, Place: 2, KillCount: 1})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/api/game/result", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.AddGameResult(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGameRouter_AddGameResult_MissingGameID(t *testing.T) {
	router := newFullGameRouter(game.Functions{}, gameResult.Functions{}, podManagerRole())

	body, _ := json.Marshal(addGameResultRequest{GameID: 0, DeckID: 10, PlayerID: 42, Place: 2})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/api/game/result", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.AddGameResult(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_UpdateGameResult_Success(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			GetByID:      func(ctx context.Context, gameID int) (*game.Entity, error) { return gameEntityForPod(1), nil },
			UpdateResult: func(ctx context.Context, resultID, place, killCount, deckID int) error { return nil },
		},
		gameResult.Functions{
			GetGameIDForResult: func(ctx context.Context, resultID int) (int, error) { return 1, nil },
		},
		podManagerRole(),
	)

	body, _ := json.Marshal(updateGameResultRequest{Place: 2, KillCount: 1, DeckID: 10})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/game/result?result_id=5", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdateGameResult(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGameRouter_UpdateGameResult_MissingParam(t *testing.T) {
	router := newFullGameRouter(game.Functions{}, gameResult.Functions{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/game/result", nil), 42)
	rr := httptest.NewRecorder()
	router.UpdateGameResult(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_DeleteGameResult_Success(t *testing.T) {
	router := newFullGameRouter(
		game.Functions{
			GetByID:      func(ctx context.Context, gameID int) (*game.Entity, error) { return gameEntityForPod(1), nil },
			DeleteResult: func(ctx context.Context, resultID int) error { return nil },
		},
		gameResult.Functions{
			GetGameIDForResult: func(ctx context.Context, resultID int) (int, error) { return 1, nil },
		},
		podManagerRole(),
	)

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/api/game/result?result_id=5", nil), 42)
	rr := httptest.NewRecorder()
	router.DeleteGameResult(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGameRouter_DeleteGameResult_MissingParam(t *testing.T) {
	router := newFullGameRouter(game.Functions{}, gameResult.Functions{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/api/game/result", nil), 42)
	rr := httptest.NewRecorder()
	router.DeleteGameResult(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGameRouter_GetAll_Paginated_ByPod(t *testing.T) {
	entities := []game.Entity{{ID: 1, Description: "Game 1"}, {ID: 2, Description: "Game 2"}}
	router := newTestGameRouter(game.Functions{
		GetAllByPodPaginated: func(ctx context.Context, podID, limit, offset int) ([]game.Entity, int, error) {
			assert.Equal(t, 5, podID)
			assert.Equal(t, 10, limit)
			assert.Equal(t, 0, offset)
			return entities, 25, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?pod_id=5&limit=10&offset=0", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got business.PaginatedResponse[game.Entity]
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, 25, got.Total)
	assert.Equal(t, 10, got.Limit)
	assert.Equal(t, 0, got.Offset)
	assert.Len(t, got.Items, 2)
}

func TestGameRouter_GetAll_Paginated_ByDeck(t *testing.T) {
	entities := []game.Entity{{ID: 3}}
	router := newTestGameRouter(game.Functions{
		GetAllByDeckPaginated: func(ctx context.Context, deckID, limit, offset int) ([]game.Entity, int, error) {
			return entities, 1, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?deck_id=2&limit=5&offset=0", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got business.PaginatedResponse[game.Entity]
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, 1, got.Total)
	assert.Len(t, got.Items, 1)
}

func TestGameRouter_GetAll_Paginated_ByPlayer(t *testing.T) {
	entities := []game.Entity{{ID: 7}}
	router := newTestGameRouter(game.Functions{
		GetAllByPlayerIDPaginated: func(ctx context.Context, playerID, limit, offset int) ([]game.Entity, int, error) {
			return entities, 50, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?player_id=42&limit=25&offset=25", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got business.PaginatedResponse[game.Entity]
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, 50, got.Total)
	assert.Equal(t, 25, got.Limit)
	assert.Equal(t, 25, got.Offset)
}

func TestGameRouter_GetAll_Paginated_Error(t *testing.T) {
	router := newTestGameRouter(game.Functions{
		GetAllByPodPaginated: func(ctx context.Context, podID, limit, offset int) ([]game.Entity, int, error) {
			return nil, 0, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?pod_id=1&limit=10", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGameRouter_GetAll_NoLimit_PlainArray(t *testing.T) {
	games := []game.Entity{{ID: 1}, {ID: 2}, {ID: 3}}
	router := newTestGameRouter(game.Functions{
		GetAllByPod: func(ctx context.Context, podID int) ([]game.Entity, error) {
			return games, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/games?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []game.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 3)
}
