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
			Path:    "/api/pod",
			Method:  http.MethodGet,
			Handler: p.GetPod,
		},
		{
			Path:    "/api/pod",
			Method:  http.MethodPost,
			Handler: p.PodCreate,
		},
		{
			Path:    "/api/pod/player",
			Method:  http.MethodPost,
			Handler: p.AddPlayer,
		},
	}
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
	log := p.log.With(zap.String("PodName", e.Name))

	if err = e.Validate(); err != nil {
		trackerHttp.WriteError(log, w, http.StatusBadRequest, err, "Pod failed validation", err.Error())
		return
	}

	log.Info("Saving new Pod record")
	if _, err = p.pods.Create(ctx, e.Name); err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Pod record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *PodRouter) AddPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to add Player to Pod"

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

	log.Info("Adding Player to Pod")
	if err = p.pods.AddPlayer(ctx, input.PodID, input.PlayerID); err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Player to Pod", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
