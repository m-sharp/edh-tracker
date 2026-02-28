package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
)

type PodRouter struct {
	log      *zap.Logger
	provider *models.PodProvider
}

func NewPodRouter(log *zap.Logger, client *lib.DBClient) *PodRouter {
	return &PodRouter{
		log:      log.Named("PodRouter"),
		provider: models.NewPodProvider(client),
	}
}

func (p *PodRouter) GetRoutes() []*lib.Route {
	return []*lib.Route{
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

	podId, podErr := lib.GetQueryId(r, "pod_id")
	playerId, playerErr := lib.GetQueryId(r, "player_id")

	if podErr != nil && playerErr != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, podErr, "Missing pod_id or player_id query param", "pod_id or player_id query param is required")
		return
	}

	var (
		marshalled []byte
		err        error
	)

	if podId != 0 {
		var pod *models.Pod
		pod, err = p.provider.GetByID(ctx, podId)
		if err != nil {
			lib.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
		marshalled, err = json.Marshal(pod)
	} else {
		var pods []models.Pod
		pods, err = p.provider.GetByPlayerID(ctx, playerId)
		if err != nil {
			lib.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
		marshalled, err = json.Marshal(pods)
	}

	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(p.log, w, marshalled)
}

func (p *PodRouter) PodCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create Pod record"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read Pod POST body", errMsg)
		return
	}

	var pod models.Pod
	if err := json.Unmarshal(body, &pod); err != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal Pod body", errMsg)
		return
	}
	log := p.log.With(zap.String("PodName", pod.Name))

	if err := pod.Validate(); err != nil {
		lib.WriteError(log, w, http.StatusBadRequest, err, "Pod failed validation", err.Error())
		return
	}

	log.Info("Saving new Pod record")
	if err := p.provider.Add(ctx, pod.Name); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Pod record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *PodRouter) AddPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to add Player to Pod"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read PlayerPod POST body", errMsg)
		return
	}

	var playerPod models.PlayerPod
	if err := json.Unmarshal(body, &playerPod); err != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal PlayerPod body", errMsg)
		return
	}
	log := p.log.With(zap.Int("PodID", playerPod.PodID), zap.Int("PlayerID", playerPod.PlayerID))

	if err := playerPod.Validate(); err != nil {
		lib.WriteError(log, w, http.StatusBadRequest, err, "PlayerPod failed validation", err.Error())
		return
	}

	log.Info("Adding Player to Pod")
	if err := p.provider.AddPlayerToPod(ctx, playerPod.PodID, playerPod.PlayerID); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Player to Pod", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
