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

func newTestPlayerRouter(repo *mockPlayerRepo) *PlayerRouter {
	return &PlayerRouter{
		log:        zap.NewNop(),
		playerRepo: repo,
	}
}

func TestPlayerRouter_GetAll_Success(t *testing.T) {
	players := []models.PlayerInfo{
		{Player: models.Player{Name: "Alice"}, PodIDs: []int{}},
		{Player: models.Player{Name: "Bob"}, PodIDs: []int{}},
	}
	repo := &mockPlayerRepo{
		GetAllFn: func(ctx context.Context) ([]models.PlayerInfo, error) { return players, nil },
	}
	router := newTestPlayerRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/players", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []models.PlayerInfo
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 2)
}

func TestPlayerRouter_GetAll_Error(t *testing.T) {
	repo := &mockPlayerRepo{
		GetAllFn: func(ctx context.Context) ([]models.PlayerInfo, error) {
			return nil, errors.New("db error")
		},
	}
	router := newTestPlayerRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/players", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestPlayerRouter_GetById_Success(t *testing.T) {
	player := &models.PlayerInfo{Player: models.Player{Name: "Alice"}, PodIDs: []int{}}
	repo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.PlayerInfo, error) {
			return player, nil
		},
	}
	router := newTestPlayerRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/player?player_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetPlayerById(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got models.PlayerInfo
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "Alice", got.Name)
}

func TestPlayerRouter_GetById_MissingParam(t *testing.T) {
	router := newTestPlayerRouter(&mockPlayerRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/player", nil)
	rr := httptest.NewRecorder()
	router.GetPlayerById(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPlayerRouter_GetById_Error(t *testing.T) {
	repo := &mockPlayerRepo{
		GetByIdFn: func(ctx context.Context, id int) (*models.PlayerInfo, error) {
			return nil, errors.New("db error")
		},
	}
	router := newTestPlayerRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/player?player_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetPlayerById(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestPlayerRouter_Add_Success(t *testing.T) {
	repo := &mockPlayerRepo{
		AddFn: func(ctx context.Context, name string) (int, error) { return 1, nil },
	}
	router := newTestPlayerRouter(repo)

	body, _ := json.Marshal(models.Player{Name: "Charlie"})
	req := httptest.NewRequest(http.MethodPost, "/api/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.PlayerCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestPlayerRouter_Add_MissingName(t *testing.T) {
	router := newTestPlayerRouter(&mockPlayerRepo{})

	body, _ := json.Marshal(models.Player{Name: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.PlayerCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
