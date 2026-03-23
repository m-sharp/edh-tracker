package routers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business/pod"
	"github.com/m-sharp/edh-tracker/lib/errs"
	"github.com/m-sharp/edh-tracker/lib/utils"
)

func newTestPodRouter(pods pod.Functions) *PodRouter {
	return &PodRouter{
		log:  zap.NewNop(),
		pods: pods,
	}
}

// withAuth injects a playerID into the request context to simulate an authenticated user.
func withAuth(r *http.Request, playerID int) *http.Request {
	return r.WithContext(utils.ContextWithUserInfo(r.Context(), 1, playerID))
}

// --- GetPod ---

func TestPodRouter_GetPod_ByPodID(t *testing.T) {
	e := &pod.Entity{ID: 5, Name: "My Pod"}
	router := newTestPodRouter(pod.Functions{
		GetByID: func(ctx context.Context, podID int) (*pod.Entity, error) { return e, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/pod?pod_id=5", nil)
	rr := httptest.NewRecorder()
	router.GetPod(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got pod.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "My Pod", got.Name)
}

func TestPodRouter_GetPod_ByPlayerID(t *testing.T) {
	pods := []pod.Entity{{ID: 1, Name: "Pod A"}, {ID: 2, Name: "Pod B"}}
	router := newTestPodRouter(pod.Functions{
		GetByPlayerID: func(ctx context.Context, playerID int) ([]pod.Entity, error) { return pods, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/api/pod?player_id=10", nil)
	rr := httptest.NewRecorder()
	router.GetPod(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []pod.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 2)
}

func TestPodRouter_GetPod_BothParamsMissing(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	req := httptest.NewRequest(http.MethodGet, "/api/pod", nil)
	rr := httptest.NewRecorder()
	router.GetPod(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_GetPod_GetByIDError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetByID: func(ctx context.Context, podID int) (*pod.Entity, error) {
			return nil, errors.New("db error")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/pod?pod_id=5", nil)
	rr := httptest.NewRecorder()
	router.GetPod(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// --- PodCreate ---

func TestPodRouter_PodCreate_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		Create: func(ctx context.Context, name string, creatorPlayerID int) (int, error) { return 1, nil },
	})

	body, _ := json.Marshal(pod.Entity{Name: "My Pod"})
	req := httptest.NewRequest(http.MethodPost, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.PodCreate(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestPodRouter_PodCreate_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.Entity{Name: "My Pod"})
	req := httptest.NewRequest(http.MethodPost, "/api/pod", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.PodCreate(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_PodCreate_ValidationFailure(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.Entity{Name: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.PodCreate(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_PodCreate_CreateError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		Create: func(ctx context.Context, name string, creatorPlayerID int) (int, error) {
			return 0, errors.New("db error")
		},
	})

	body, _ := json.Marshal(pod.Entity{Name: "My Pod"})
	req := httptest.NewRequest(http.MethodPost, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.PodCreate(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// --- UpdatePod ---

func TestPodRouter_UpdatePod_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole: func(ctx context.Context, podID, playerID int) (string, error) { return "manager", nil },
		Update:  func(ctx context.Context, podID int, name string) error { return nil },
	})

	body, _ := json.Marshal(pod.UpdatePodInputEntity{PodID: 1, Name: "New Name"})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.UpdatePod(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPodRouter_UpdatePod_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.UpdatePodInputEntity{PodID: 1, Name: "New Name"})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.UpdatePod(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_UpdatePod_ValidationFailure(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.UpdatePodInputEntity{PodID: 0, Name: ""})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.UpdatePod(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_UpdatePod_NotManager(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole: func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	})

	body, _ := json.Marshal(pod.UpdatePodInputEntity{PodID: 1, Name: "New Name"})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.UpdatePod(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestPodRouter_UpdatePod_UpdateError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole: func(ctx context.Context, podID, playerID int) (string, error) { return "manager", nil },
		Update:  func(ctx context.Context, podID int, name string) error { return errors.New("db error") },
	})

	body, _ := json.Marshal(pod.UpdatePodInputEntity{PodID: 1, Name: "New Name"})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.UpdatePod(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// --- DeletePod ---

func TestPodRouter_DeletePod_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole:    func(ctx context.Context, podID, playerID int) (string, error) { return "manager", nil },
		SoftDelete: func(ctx context.Context, podID, callerPlayerID int) error { return nil },
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/pod?pod_id=1", nil)
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.DeletePod(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPodRouter_DeletePod_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	req := httptest.NewRequest(http.MethodDelete, "/api/pod?pod_id=1", nil)
	rr := httptest.NewRecorder()
	router.DeletePod(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_DeletePod_MissingPodID(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	req := httptest.NewRequest(http.MethodDelete, "/api/pod", nil)
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.DeletePod(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_DeletePod_NotManager(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole: func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/pod?pod_id=1", nil)
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.DeletePod(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

// --- AddPlayer ---

func TestPodRouter_AddPlayer_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole:   func(ctx context.Context, podID, playerID int) (string, error) { return "manager", nil },
		AddPlayer: func(ctx context.Context, podID, playerID int) error { return nil },
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.AddPlayer(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestPodRouter_AddPlayer_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.AddPlayer(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_AddPlayer_ValidationFailure(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 0, PlayerID: 0})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.AddPlayer(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_AddPlayer_NotManager(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole: func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.AddPlayer(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

// --- PromotePlayer ---

func TestPodRouter_PromotePlayer_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		PromoteToManager: func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error { return nil },
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.PromotePlayer(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPodRouter_PromotePlayer_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.PromotePlayer(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_PromotePlayer_PromoteError_Forbidden(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		PromoteToManager: func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error {
			return fmt.Errorf("forbidden: caller is not a manager: %w", errs.ErrForbidden)
		},
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.PromotePlayer(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestPodRouter_PromotePlayer_PromoteError_DBError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		PromoteToManager: func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error {
			return errors.New("db connection error")
		},
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodPatch, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.PromotePlayer(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// --- KickPlayer ---

func TestPodRouter_KickPlayer_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		RemovePlayer: func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error { return nil },
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodDelete, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.KickPlayer(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPodRouter_KickPlayer_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodDelete, "/api/pod/player", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.KickPlayer(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_KickPlayer_RemoveError_Forbidden(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		RemovePlayer: func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error {
			return fmt.Errorf("forbidden: caller is not a manager: %w", errs.ErrForbidden)
		},
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodDelete, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.KickPlayer(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestPodRouter_KickPlayer_RemoveError_DBError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		RemovePlayer: func(ctx context.Context, podID, callerPlayerID, targetPlayerID int) error {
			return errors.New("db connection error")
		},
	})

	body, _ := json.Marshal(pod.PlayerPodInputEntity{PodID: 1, PlayerID: 2})
	req := httptest.NewRequest(http.MethodDelete, "/api/pod/player", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.KickPlayer(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// --- GenerateInvite ---

func TestPodRouter_GenerateInvite_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole:        func(ctx context.Context, podID, playerID int) (string, error) { return "manager", nil },
		GenerateInvite: func(ctx context.Context, podID, callerPlayerID int) (string, error) { return "invite-code-abc", nil },
	})

	body, _ := json.Marshal(struct {
		PodID int `json:"pod_id"`
	}{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/invite", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.GenerateInvite(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got pod.InviteEntity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "invite-code-abc", got.InviteCode)
}

func TestPodRouter_GenerateInvite_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(struct {
		PodID int `json:"pod_id"`
	}{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/invite", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.GenerateInvite(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_GenerateInvite_NotManager(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		GetRole: func(ctx context.Context, podID, playerID int) (string, error) { return "member", nil },
	})

	body, _ := json.Marshal(struct {
		PodID int `json:"pod_id"`
	}{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/invite", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.GenerateInvite(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

// --- JoinByInvite ---

func TestPodRouter_JoinByInvite_Success(t *testing.T) {
	e := &pod.Entity{ID: 5, Name: "My Pod"}
	router := newTestPodRouter(pod.Functions{
		JoinByInvite: func(ctx context.Context, inviteCode string, playerID int) (*pod.Entity, error) { return e, nil },
	})

	body, _ := json.Marshal(pod.JoinInputEntity{InviteCode: "valid-code"})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/join", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.JoinByInvite(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got pod.Entity
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, "My Pod", got.Name)
}

func TestPodRouter_JoinByInvite_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.JoinInputEntity{InviteCode: "valid-code"})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/join", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.JoinByInvite(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_JoinByInvite_ValidationFailure(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.JoinInputEntity{InviteCode: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/join", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.JoinByInvite(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_JoinByInvite_JoinError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		JoinByInvite: func(ctx context.Context, inviteCode string, playerID int) (*pod.Entity, error) {
			return nil, errors.New("invite code not found or expired")
		},
	})

	body, _ := json.Marshal(pod.JoinInputEntity{InviteCode: "bad-code"})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/join", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.JoinByInvite(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// --- LeavePod ---

func TestPodRouter_LeavePod_Success(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		Leave: func(ctx context.Context, podID, playerID int) error { return nil },
	})

	body, _ := json.Marshal(pod.LeaveInputEntity{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/leave", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.LeavePod(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPodRouter_LeavePod_Unauthenticated(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.LeaveInputEntity{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/leave", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.LeavePod(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestPodRouter_LeavePod_ValidationFailure(t *testing.T) {
	router := newTestPodRouter(pod.Functions{})

	body, _ := json.Marshal(pod.LeaveInputEntity{PodID: 0})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/leave", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.LeavePod(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPodRouter_LeavePod_LeaveError_Forbidden(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		Leave: func(ctx context.Context, podID, playerID int) error {
			return fmt.Errorf("forbidden: cannot leave as sole manager: %w", errs.ErrForbidden)
		},
	})

	body, _ := json.Marshal(pod.LeaveInputEntity{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/leave", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.LeavePod(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestPodRouter_LeavePod_LeaveError_DBError(t *testing.T) {
	router := newTestPodRouter(pod.Functions{
		Leave: func(ctx context.Context, podID, playerID int) error {
			return errors.New("db connection error")
		},
	})

	body, _ := json.Marshal(pod.LeaveInputEntity{PodID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/pod/leave", bytes.NewReader(body))
	req = withAuth(req, 10)
	rr := httptest.NewRecorder()
	router.LeavePod(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
