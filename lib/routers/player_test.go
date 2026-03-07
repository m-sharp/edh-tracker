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

func TestPlayerRouter_Update_Success(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{
		Update: func(ctx context.Context, playerID int, name string) error { return nil },
	})

	body, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: "Charlie"})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/player?player_id=42", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdatePlayer(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPlayerRouter_Update_MissingName(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{})

	body, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: ""})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/player?player_id=42", bytes.NewReader(body)), 42)
	rr := httptest.NewRecorder()
	router.UpdatePlayer(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPlayerRouter_Update_Forbidden(t *testing.T) {
	router := newTestPlayerRouter(player.Functions{})

	body, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: "NewName"})
	// callerID=99 != player_id=42 → forbidden
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/api/player?player_id=42", bytes.NewReader(body)), 99)
	rr := httptest.NewRecorder()
	router.UpdatePlayer(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestPlayerRouter_GetAll_ByPod_Success(t *testing.T) {
	players := []player.PlayerWithRoleEntity{
		{Entity: player.Entity{Name: "Alice"}, Role: "manager"},
	}
	router := newTestPlayerRouter(player.Functions{
		GetAllByPod: func(ctx context.Context, podID int) ([]player.PlayerWithRoleEntity, error) {
			return players, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/players?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []player.PlayerWithRoleEntity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
	assert.Equal(t, "manager", got[0].Role)
}
