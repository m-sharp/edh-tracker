package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/player"
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

func (p *PlayerRouter) GetRoutes() []*lib.Route {
	return []*lib.Route{
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
			Method:  http.MethodPost,
			Handler: p.PlayerCreate,
		},
	}
}

func (p *PlayerRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Player records"

	players, err := p.players.GetAll(ctx)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(players)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(p.log, w, marshalled)
}

func (p *PlayerRouter) GetPlayerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: Use route param instead?
	playerID, err := lib.GetQueryId(r, "player_id")
	if err != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, err, "Bad player_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Player record"
	playerEntity, err := p.players.GetByID(ctx, playerID)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(playerEntity)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(p.log, w, marshalled)
}

func (p *PlayerRouter) PlayerCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create new Player"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read player POST body", errMsg)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal Player body", errMsg)
		return
	}
	log := p.log.With(zap.String("NewPlayer", req.Name))

	if req.Name == "" {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Missing player name", "name is required")
		return
	}

	log.Info("Saving new Player record")
	if _, err = p.players.Create(ctx, req.Name); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Player record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
