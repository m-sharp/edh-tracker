package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/pod"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type PodRouter struct {
	log  *zap.Logger
	pods pod.Functions
}

func NewPodRouter(log *zap.Logger, biz *business.Business) *PodRouter {
	return &PodRouter{
		log:  log.Named("PodRouter"),
		pods: biz.Pods,
	}
}

func (p *PodRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
		{
			Path:        "/api/pod",
			Method:      http.MethodGet,
			Handler:     p.GetPod,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod",
			Method:      http.MethodPost,
			Handler:     p.PodCreate,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod",
			Method:      http.MethodPatch,
			Handler:     p.UpdatePod,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod",
			Method:      http.MethodDelete,
			Handler:     p.DeletePod,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod/player",
			Method:      http.MethodPost,
			Handler:     p.AddPlayer,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod/player",
			Method:      http.MethodPatch,
			Handler:     p.PromotePlayer,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod/player",
			Method:      http.MethodDelete,
			Handler:     p.KickPlayer,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod/invite",
			Method:      http.MethodPost,
			Handler:     p.GenerateInvite,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod/join",
			Method:      http.MethodPost,
			Handler:     p.JoinByInvite,
			RequireAuth: true,
		},
		{
			Path:        "/api/pod/leave",
			Method:      http.MethodPost,
			Handler:     p.LeavePod,
			RequireAuth: true,
		},
	}
}

// requireManager checks that the caller is a manager of the given pod.
// Returns false and writes a 403 if not.
func (p *PodRouter) requireManager(w http.ResponseWriter, r *http.Request, podID, callerPlayerID int) bool {
	role, err := p.pods.GetRole(r.Context(), podID, callerPlayerID)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to check pod role", "internal error")
		return false
	}
	if role != "manager" {
		http.Error(w, "Forbidden: pod manager role required", http.StatusForbidden)
		return false
	}
	return true
}

func (p *PodRouter) GetPod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Pod records"

	podId, podErr := trackerHttp.GetQueryId(r, "pod_id")
	playerId, playerErr := trackerHttp.GetQueryId(r, "player_id")

	if podErr != nil && playerErr != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, podErr, "Missing pod_id or player_id query param", "pod_id or player_id query param is required")
		return
	}

	var (
		marshalled []byte
		err        error
	)

	if podId != 0 {
		var podEntity *pod.Entity
		podEntity, err = p.pods.GetByID(ctx, podId)
		if err != nil {
			trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
		marshalled, err = json.Marshal(podEntity)
	} else {
		var pods []pod.Entity
		pods, err = p.pods.GetByPlayerID(ctx, playerId)
		if err != nil {
			trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
		marshalled, err = json.Marshal(pods)
	}

	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(p.log, w, marshalled)
}

func (p *PodRouter) PodCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create Pod record"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read Pod POST body", errMsg)
		return
	}

	var e pod.Entity
	if err = json.Unmarshal(body, &e); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal Pod body", errMsg)
		return
	}
	log := p.log.With(zap.String("PodName", e.Name), zap.Int("PlayerID", callerID))

	if err = e.Validate(); err != nil {
		trackerHttp.WriteError(log, w, http.StatusBadRequest, err, "Pod failed validation", err.Error())
		return
	}

	log.Info("Saving new Pod record")
	if _, err = p.pods.Create(ctx, e.Name, callerID); err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Pod record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *PodRouter) UpdatePod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to update Pod"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read PATCH body", errMsg)
		return
	}

	var input pod.UpdatePodInputEntity
	if err = json.Unmarshal(body, &input); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal PATCH body", errMsg)
		return
	}
	if err = input.Validate(); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "UpdatePod failed validation", err.Error())
		return
	}

	if !p.requireManager(w, r, input.PodID, callerID) {
		return
	}

	if err = p.pods.Update(ctx, input.PodID, input.Name); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to update pod name", errMsg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PodRouter) DeletePod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to delete Pod"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	podID, err := trackerHttp.GetQueryId(r, "pod_id")
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Missing pod_id", "pod_id query param is required")
		return
	}

	if !p.requireManager(w, r, podID, callerID) {
		return
	}

	if err = p.pods.SoftDelete(ctx, podID, callerID); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to soft delete pod", errMsg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PodRouter) AddPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to add Player to Pod"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read PlayerPod POST body", errMsg)
		return
	}

	var input pod.PlayerPodInputEntity
	if err = json.Unmarshal(body, &input); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal PlayerPod body", errMsg)
		return
	}
	log := p.log.With(zap.Int("PodID", input.PodID), zap.Int("PlayerID", input.PlayerID))

	if err = input.Validate(); err != nil {
		trackerHttp.WriteError(log, w, http.StatusBadRequest, err, "PlayerPod failed validation", err.Error())
		return
	}

	if !p.requireManager(w, r, input.PodID, callerID) {
		return
	}

	log.Info("Adding Player to Pod")
	if err = p.pods.AddPlayer(ctx, input.PodID, input.PlayerID); err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Player to Pod", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *PodRouter) PromotePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to promote player"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read PATCH body", errMsg)
		return
	}

	var input pod.PlayerPodInputEntity
	if err = json.Unmarshal(body, &input); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal PATCH body", errMsg)
		return
	}
	if err = input.Validate(); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "PromotePlayer failed validation", err.Error())
		return
	}

	if err = p.pods.PromoteToManager(ctx, input.PodID, callerID, input.PlayerID); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusForbidden, err, "Failed to promote player to manager", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PodRouter) KickPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to remove player from pod"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read DELETE body", errMsg)
		return
	}

	var input pod.PlayerPodInputEntity
	if err = json.Unmarshal(body, &input); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal DELETE body", errMsg)
		return
	}
	if err = input.Validate(); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "KickPlayer failed validation", err.Error())
		return
	}

	if err = p.pods.RemovePlayer(ctx, input.PodID, callerID, input.PlayerID); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusForbidden, err, "Failed to remove player", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PodRouter) GenerateInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to generate invite"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read invite POST body", errMsg)
		return
	}

	var input struct {
		PodID int `json:"pod_id"`
	}
	if err = json.Unmarshal(body, &input); err != nil || input.PodID == 0 {
		http.Error(w, "pod_id is required", http.StatusBadRequest)
		return
	}

	if !p.requireManager(w, r, input.PodID, callerID) {
		return
	}

	code, err := p.pods.GenerateInvite(ctx, input.PodID, callerID)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to generate invite code", errMsg)
		return
	}

	marshalled, err := json.Marshal(pod.InviteEntity{InviteCode: code})
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshal invite response", errMsg)
		return
	}

	trackerHttp.WriteJson(p.log, w, marshalled)
}

func (p *PodRouter) JoinByInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to join pod"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read join POST body", errMsg)
		return
	}

	var input pod.JoinInputEntity
	if err = json.Unmarshal(body, &input); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal join body", errMsg)
		return
	}
	if err = input.Validate(); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "JoinByInvite failed validation", err.Error())
		return
	}

	podEntity, err := p.pods.JoinByInvite(ctx, input.InviteCode, callerID)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to join pod by invite", err.Error())
		return
	}

	marshalled, err := json.Marshal(podEntity)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshal pod response", errMsg)
		return
	}

	trackerHttp.WriteJson(p.log, w, marshalled)
}

func (p *PodRouter) LeavePod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to leave pod"

	callerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read leave POST body", errMsg)
		return
	}

	var input pod.LeaveInputEntity
	if err = json.Unmarshal(body, &input); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal leave body", errMsg)
		return
	}
	if err = input.Validate(); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "LeavePod failed validation", err.Error())
		return
	}

	if err = p.pods.Leave(ctx, input.PodID, callerID); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusForbidden, err, "Failed to leave pod", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
