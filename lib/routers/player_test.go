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

	"github.com/m-sharp/edh-tracker/lib/business/player"
)

func newTestPlayerRouter(players player.Functions) *PlayerRouter {
	return &PlayerRouter{
		log:     zap.NewNop(),
		players: players,
	}
}

func TestPlayerRouter_GetAll_Success(t *testing.T) {
	players := []player.Entity{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	router := newTestPlayerRouter(player.Functions{
		GetAll: func(ctx context.Context) ([]player.Entity, error) { return players, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/players", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []player.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 2)
}

func TestPlayerRouter_GetAll_Error(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{
		GetAll: func(ctx context.Context) ([]player.Entity, error) {
			return nil, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/players", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestPlayerRouter_GetById_Success(t *testing.T) {
	p := &player.Entity{Name: "Alice"}
	router := newTestPlayerRouter(player.Functions{
		GetByID: func(ctx context.Context, id int) (*player.Entity, error) { return p, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/player?player_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetPlayerById(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got player.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "Alice", got.Name)
}

func TestPlayerRouter_GetById_MissingParam(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{})

	req := httptest.NewRequest(http.MethodGet, "/api/player", nil)
	rr := httptest.NewRecorder()
	router.GetPlayerById(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPlayerRouter_GetById_Error(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{
		GetByID: func(ctx context.Context, id int) (*player.Entity, error) {
			return nil, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/player?player_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetPlayerById(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestPlayerRouter_Add_Success(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{
		Create: func(ctx context.Context, name string) (int, error) { return 1, nil },
	})

	body, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: "Charlie"})
	req := httptest.NewRequest(http.MethodPost, "/api/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.PlayerCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestPlayerRouter_Add_MissingName(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{})

	body, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.PlayerCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
