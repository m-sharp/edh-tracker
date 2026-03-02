package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/game"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
)

type GameRouter struct {
	log   *zap.Logger
	games game.Functions
}

func NewGameRouter(log *zap.Logger, biz *business.Business) *GameRouter {
	return &GameRouter{
		log:   log.Named("GameRouter"),
		games: biz.Games,
	}
}

func (g *GameRouter) GetRoutes() []*lib.Route {
	return []*lib.Route{
		{
			Path:    "/api/games",
			Method:  http.MethodGet,
			Handler: g.GetAll,
		},
		{
			Path:    "/api/game",
			Method:  http.MethodGet,
			Handler: g.GetGameById,
		},
		{
			Path:    "/api/game",
			Method:  http.MethodPost,
			Handler: g.GameCreate,
		},
	}
}

type createGameRequest struct {
	Description string                   `json:"description"`
	PodID       int                      `json:"pod_id"`
	FormatID    int                      `json:"format_id"`
	Results     []gameResult.InputEntity `json:"results"`
}

// ToDo: Eventually, this will probably need pagination
func (g *GameRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Game records"

	podId, _ := lib.GetQueryId(r, "pod_id")
	deckId, _ := lib.GetQueryId(r, "deck_id")

	if podId == 0 && deckId == 0 {
		lib.WriteError(g.log, w, http.StatusBadRequest, fmt.Errorf("missing required query param"), "Missing pod_id or deck_id query param", "pod_id or deck_id query param is required")
		return
	}

	var (
		games []game.Entity
		err   error
	)
	if deckId != 0 {
		games, err = g.games.GetAllByDeck(ctx, deckId)
		if err != nil {
			lib.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	} else {
		games, err = g.games.GetAllByPod(ctx, podId)
		if err != nil {
			lib.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	}

	marshalled, err := json.Marshal(games)
	if err != nil {
		lib.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(g.log, w, marshalled)
}

func (g *GameRouter) GetGameById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Game Details"

	gameId, err := lib.GetQueryId(r, "game_id")
	if err != nil {
		lib.WriteError(g.log, w, http.StatusBadRequest, err, "Bad game_id query string specified", err.Error())
		return
	}

	gameEntity, err := g.games.GetByID(ctx, gameId)
	if err != nil {
		lib.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(gameEntity)
	if err != nil {
		lib.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(g.log, w, marshalled)
}

func (g *GameRouter) GameCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create Game record"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to read Game POST body", errMsg)
		return
	}

	var req createGameRequest
	if err = json.Unmarshal(body, &req); err != nil {
		lib.WriteError(g.log, w, http.StatusBadRequest, err, "Failed to unmarshal Game body", errMsg)
		return
	}
	log := g.log.With(
		zap.String("GameDescription", req.Description),
		zap.Any("GameResults", req.Results),
	)

	if len(req.Results) == 0 {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "No game results provided", "at least one game result is required")
		return
	}

	for _, result := range req.Results {
		if err = result.Validate(); err != nil {
			lib.WriteError(
				log, w, http.StatusBadRequest, err,
				"Game result failed validation",
				fmt.Sprintf("Game result failed validation: %s", err.Error()),
			)
			return
		}
	}

	log.Info("Saving new Game record")
	if err = g.games.Create(ctx, req.Description, req.PodID, req.FormatID, req.Results); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to create Game record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
