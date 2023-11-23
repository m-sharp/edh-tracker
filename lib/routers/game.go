package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
)

type GameRouter struct {
	log      *zap.Logger
	provider *models.GameProvider
}

func NewGameRouter(log *zap.Logger, client *lib.DBClient) *GameRouter {
	return &GameRouter{
		log:      log.Named("GameRouter"),
		provider: models.NewGameProvider(log, client),
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

// ToDo: Eventually, this will probably need pagination
func (g *GameRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Game records"

	games, err := g.provider.GetAll(ctx)
	if err != nil {
		lib.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
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

	gameDetails, err := g.provider.GetGameById(ctx, gameId)
	if err != nil {
		lib.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(gameDetails)
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

	var gameDetails models.GameDetails
	if err := json.Unmarshal(body, &gameDetails); err != nil {
		lib.WriteError(g.log, w, http.StatusBadRequest, err, "Failed to unmarshal GameDetails body", errMsg)
		return
	}
	log := g.log.With(
		zap.String("GameDescription", gameDetails.Description),
		zap.Any("GameResults", gameDetails.Results),
	)

	for _, result := range gameDetails.Results {
		if err := result.Validate(); err != nil {
			lib.WriteError(
				log, w, http.StatusBadRequest, err,
				"Game result failed validation",
				fmt.Sprintf("Game result failed validation: %s", err.Error()),
			)
			return
		}
	}

	log.Info("Saving new Game record")
	if err := g.provider.Add(ctx, gameDetails.Description, gameDetails.Results...); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Game record", errMsg)
	}

	w.WriteHeader(http.StatusCreated)
}
