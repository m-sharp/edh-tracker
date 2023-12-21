package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
)

type PlayerRouter struct {
	log      *zap.Logger
	provider *models.PlayerProvider
}

func NewPlayerRouter(log *zap.Logger, client *lib.DBClient) *PlayerRouter {
	return &PlayerRouter{
		log:      log.Named("PlayerRouter"),
		provider: models.NewPlayerProvider(client),
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

	players, err := p.provider.GetAll(ctx)
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
	playerId, err := lib.GetQueryId(r, "player_id")
	if err != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, err, "Bad player_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Player record"
	players, err := p.provider.GetById(ctx, playerId)
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

func (p *PlayerRouter) PlayerCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create new Player"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(p.log, w, http.StatusInternalServerError, err, "Failed to read player POST body", errMsg)
		return
	}

	var player models.Player
	if err := json.Unmarshal(body, &player); err != nil {
		lib.WriteError(p.log, w, http.StatusBadRequest, err, "Failed to unmarshal Player body", errMsg)
		return
	}
	log := p.log.With(zap.String("NewPlayer", player.Name))

	if err := player.Validate(); err != nil {
		lib.WriteError(log, w, http.StatusBadRequest, err, "New Player failed validation", errMsg)
		return
	}

	log.Info("Saving new Player record")
	if err := p.provider.Add(ctx, player.Name); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Player record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
