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
	log        *zap.Logger
	gameRepo   *models.GameProvider
	formatRepo *models.FormatProvider
	deckRepo   *models.DeckProvider
}

func NewGameRouter(log *zap.Logger, repos *models.Repositories) *GameRouter {
	return &GameRouter{
		log:        log.Named("GameRouter"),
		gameRepo:   repos.Games,
		formatRepo: repos.Formats,
		deckRepo:   repos.Decks,
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

	podId, _ := lib.GetQueryId(r, "pod_id")
	deckId, _ := lib.GetQueryId(r, "deck_id")

	if podId == 0 && deckId == 0 {
		lib.WriteError(g.log, w, http.StatusBadRequest, fmt.Errorf("missing required query param"), "Missing pod_id or deck_id query param", "pod_id or deck_id query param is required")
		return
	}

	var (
		games []models.GameDetails
		err   error
	)
	if deckId != 0 {
		games, err = g.gameRepo.GetAllByDeck(ctx, deckId)
		if err != nil {
			lib.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	} else {
		games, err = g.gameRepo.GetAllByPod(ctx, podId)
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

	gameDetails, err := g.gameRepo.GetGameById(ctx, gameId)
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

	format, err := g.formatRepo.GetById(ctx, gameDetails.FormatID)
	if err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to look up format", errMsg)
		return
	}
	if format == nil {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Format not found", "Invalid format_id")
		return
	}

	// For non-"other" formats, verify all decks share the game's format
	if format.Name != "other" {
		for _, result := range gameDetails.Results {
			deck, err := g.deckRepo.GetById(ctx, result.DeckId)
			if err != nil {
				lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to look up deck", errMsg)
				return
			}
			if deck == nil {
				lib.WriteError(log, w, http.StatusBadRequest, nil,
					fmt.Sprintf("Deck %d not found", result.DeckId),
					fmt.Sprintf("Deck %d not found", result.DeckId),
				)
				return
			}
			if deck.FormatID != gameDetails.FormatID {
				lib.WriteError(log, w, http.StatusBadRequest, nil,
					fmt.Sprintf("Deck %d format does not match game format", result.DeckId),
					fmt.Sprintf("Deck %d is not in the correct format for this game", result.DeckId),
				)
				return
			}
		}
	}

	log.Info("Saving new Game record")
	if err := g.gameRepo.Add(ctx, gameDetails.Description, gameDetails.PodID, gameDetails.FormatID, gameDetails.Results...); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Game record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
