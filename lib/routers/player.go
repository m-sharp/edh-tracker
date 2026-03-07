package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/player"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type PlayerRouter struct {
	log     *zap.Logger
	players player.Functions
}

func NewPlayerRouter(log *zap.Logger, biz *business.Business) *PlayerRouter {
	return &PlayerRouter{
		log:     log.Named("PlayerRouter"),
		players: biz.Players,
	}
}

func (p *PlayerRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
		{
			Path:    "/api/players",
			Method:  http.MethodGet,
			Handler: p.GetAll,
		},
		{
			Path:    "/api/player",
			Method:  http.MethodGet,
			Handler: p.GetPlayerById,
		},
		{
			Path:    "/api/player",
			Method:  http.MethodPatch,
			Handler: p.UpdatePlayer,
		},
	}
}

func (p *PlayerRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Player records"

	podID, _ := trackerHttp.GetQueryId(r, "pod_id")

	var (
		marshalled []byte
		marshalErr error
	)

	if podID != 0 {
		players, err := p.players.GetAllByPod(ctx, podID)
		if err != nil {
			trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
		marshalled, marshalErr = json.Marshal(players)
	} else {
		players, err := p.players.GetAll(ctx)
		if err != nil {
			trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
		marshalled, marshalErr = json.Marshal(players)
	}

	if marshalErr != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, marshalErr, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(p.log, w, marshalled)
}

func (p *PlayerRouter) GetPlayerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: Use route param instead?
	playerID, err := trackerHttp.GetQueryId(r, "player_id")
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Bad player_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Player record"
	playerEntity, err := p.players.GetByID(ctx, playerID)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(playerEntity)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(p.log, w, marshalled)
}

func (p *PlayerRouter) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to update Player"

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	playerID, err := trackerHttp.GetQueryId(r, "player_id")
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Bad player_id query string specified", err.Error())
		return
	}

	if callerPlayerID != playerID {
		http.Error(w, "Forbidden: you may only update your own player record", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read player PATCH body", errMsg)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal Player update body", errMsg)
		return
	}

	if req.Name == "" {
		trackerHttp.WriteError(p.log, w, http.StatusBadRequest, nil, "Missing player name", "name is required")
		return
	}

	if err = p.players.Update(ctx, playerID, req.Name); err != nil {
		trackerHttp.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to update Player record", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
