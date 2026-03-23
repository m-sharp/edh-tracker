package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	bizpkg "github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/game"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
	"github.com/m-sharp/edh-tracker/lib/business/pod"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type GameRouter struct {
	log         *zap.Logger
	games       game.Functions
	gameResults gameResult.Functions
	getPodRole  pod.GetRoleFunc
}

func NewGameRouter(log *zap.Logger, biz *bizpkg.Business) *GameRouter {
	return &GameRouter{
		log:         log.Named("GameRouter"),
		games:       biz.Games,
		gameResults: biz.GameResults,
		getPodRole:  biz.Pods.GetRole,
	}
}

func (g *GameRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
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
		{
			Path:    "/api/game",
			Method:  http.MethodPatch,
			Handler: g.UpdateGame,
		},
		{
			Path:    "/api/game",
			Method:  http.MethodDelete,
			Handler: g.DeleteGame,
		},
		{
			Path:    "/api/game/result",
			Method:  http.MethodPost,
			Handler: g.AddGameResult,
		},
		{
			Path:    "/api/game/result",
			Method:  http.MethodPatch,
			Handler: g.UpdateGameResult,
		},
		{
			Path:    "/api/game/result",
			Method:  http.MethodDelete,
			Handler: g.DeleteGameResult,
		},
	}
}

// requirePodManager fetches the game by gameID, then checks that callerPlayerID is a pod manager.
// Returns false and writes the appropriate error response if the check fails.
// TODO: For permissions, do we need a view for fetching a player that includes their permissions for pods & decks would give us lists/maps to just check ids against?
func (g *GameRouter) requirePodManager(w http.ResponseWriter, r *http.Request, gameID, callerPlayerID int) bool {
	gameEntity, err := g.games.GetByID(r.Context(), gameID)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to look up game", "internal error")
		return false
	}
	if gameEntity == nil {
		http.Error(w, "game not found", http.StatusNotFound)
		return false
	}

	role, err := g.getPodRole(r.Context(), gameEntity.PodID, callerPlayerID)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to check pod role", "internal error")
		return false
	}
	if role != "manager" {
		http.Error(w, "Forbidden: pod manager role required", http.StatusForbidden)
		return false
	}
	return true
}

type createGameRequest struct {
	Description string                   `json:"description"`
	PodID       int                      `json:"pod_id"`
	FormatID    int                      `json:"format_id"`
	Results     []gameResult.InputEntity `json:"results"`
}

func (g *GameRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Game records"

	podId, _ := trackerHttp.GetQueryId(r, "pod_id")
	deckId, _ := trackerHttp.GetQueryId(r, "deck_id")
	playerID, _ := trackerHttp.GetQueryId(r, "player_id")

	if podId == 0 && deckId == 0 && playerID == 0 {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, fmt.Errorf("missing required query param"), "Missing required query param", "pod_id, deck_id, or player_id query param is required")
		return
	}

	limit, _ := trackerHttp.GetQueryId(r, "limit")
	offset, _ := trackerHttp.GetQueryId(r, "offset")

	if limit > 0 {
		g.getAllPaginated(w, r, podId, deckId, playerID, limit, offset)
		return
	}

	var (
		games []game.Entity
		err   error
	)
	switch {
	case deckId != 0:
		games, err = g.games.GetAllByDeck(ctx, deckId)
	case playerID != 0:
		games, err = g.games.GetAllByPlayer(ctx, playerID)
	default:
		games, err = g.games.GetAllByPod(ctx, podId)
	}
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(games)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(g.log, w, marshalled)
}

func (g *GameRouter) getAllPaginated(w http.ResponseWriter, r *http.Request, podId, deckId, playerID, limit, offset int) {
	ctx := r.Context()
	errMsg := "Failed to get Game records"

	var (
		entities []game.Entity
		total    int
		err      error
	)
	switch {
	case deckId != 0:
		entities, total, err = g.games.GetAllByDeckPaginated(ctx, deckId, limit, offset)
	case playerID != 0:
		entities, total, err = g.games.GetAllByPlayerIDPaginated(ctx, playerID, limit, offset)
	default:
		entities, total, err = g.games.GetAllByPodPaginated(ctx, podId, limit, offset)
	}
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	resp := bizpkg.PaginatedResponse[game.Entity]{
		Items:  entities,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	marshalled, err := json.Marshal(resp)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}
	trackerHttp.WriteJson(g.log, w, marshalled)
}

func (g *GameRouter) GetGameById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Game Details"

	gameId, err := trackerHttp.GetQueryId(r, "game_id")
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Bad game_id query string specified", err.Error())
		return
	}

	gameEntity, err := g.games.GetByID(ctx, gameId)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(gameEntity)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(g.log, w, marshalled)
}

func (g *GameRouter) GameCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create Game record"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to read Game POST body", errMsg)
		return
	}

	var req createGameRequest
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Failed to unmarshal Game body", errMsg)
		return
	}
	log := g.log.With(
		zap.String("GameDescription", req.Description),
		zap.Any("GameResults", req.Results),
	)

	if len(req.Description) > 256 {
		trackerHttp.WriteError(log, w, http.StatusBadRequest, fmt.Errorf("description too long"), "Game description too long", "description must be 256 characters or fewer")
		return
	}

	if len(req.Results) == 0 {
		trackerHttp.WriteError(log, w, http.StatusBadRequest, nil, "No game results provided", "at least one game result is required")
		return
	}

	for _, result := range req.Results {
		if err = result.Validate(); err != nil {
			trackerHttp.WriteError(
				log, w, http.StatusBadRequest, err,
				"Game result failed validation",
				fmt.Sprintf("Game result failed validation: %s", err.Error()),
			)
			return
		}
	}

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	role, err := g.getPodRole(ctx, req.PodID, callerPlayerID)
	if err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to check pod membership", "internal error")
		return
	}
	if role == "" {
		http.Error(w, "Forbidden: must be a member of the pod", http.StatusForbidden)
		return
	}

	log.Info("Saving new Game record")
	if err = g.games.Create(ctx, req.Description, req.PodID, req.FormatID, req.Results); err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to create Game record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (g *GameRouter) UpdateGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to update Game"

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	gameID, err := trackerHttp.GetQueryId(r, "game_id")
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Bad game_id query string specified", err.Error())
		return
	}

	if !g.requirePodManager(w, r, gameID, callerPlayerID) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to read Game PATCH body", errMsg)
		return
	}

	var req struct {
		Description string `json:"description"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Failed to unmarshal Game update body", errMsg)
		return
	}

	if err = g.games.Update(ctx, gameID, req.Description); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to update Game record", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (g *GameRouter) DeleteGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	gameID, err := trackerHttp.GetQueryId(r, "game_id")
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Bad game_id query string specified", err.Error())
		return
	}

	if !g.requirePodManager(w, r, gameID, callerPlayerID) {
		return
	}

	if err = g.games.SoftDelete(ctx, gameID); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to delete Game record", "Failed to delete Game")
		return
	}

	w.WriteHeader(http.StatusOK)
}

type addGameResultRequest struct {
	GameID    int `json:"game_id"`
	DeckID    int `json:"deck_id"`
	PlayerID  int `json:"player_id"`
	Place     int `json:"place"`
	KillCount int `json:"kill_count"`
}

func (g *GameRouter) AddGameResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to add Game Result"

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to read Game Result POST body", errMsg)
		return
	}

	var req addGameResultRequest
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Failed to unmarshal Game Result body", errMsg)
		return
	}

	if req.GameID == 0 {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, nil, "Missing game_id", "game_id is required")
		return
	}

	if !g.requirePodManager(w, r, req.GameID, callerPlayerID) {
		return
	}

	if _, err = g.games.AddResult(ctx, req.GameID, req.DeckID, req.PlayerID, req.Place, req.KillCount); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to add Game Result record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type updateGameResultRequest struct {
	Place     int `json:"place"`
	KillCount int `json:"kill_count"`
	DeckID    int `json:"deck_id"`
}

func (g *GameRouter) UpdateGameResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to update Game Result"

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	resultID, err := trackerHttp.GetQueryId(r, "result_id")
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Bad result_id query string specified", err.Error())
		return
	}

	gameID, err := g.gameResults.GetGameIDForResult(ctx, resultID)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to look up result", errMsg)
		return
	}

	if !g.requirePodManager(w, r, gameID, callerPlayerID) {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to read Game Result PATCH body", errMsg)
		return
	}

	var req updateGameResultRequest
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Failed to unmarshal Game Result update body", errMsg)
		return
	}

	if err = g.games.UpdateResult(ctx, resultID, req.Place, req.KillCount, req.DeckID); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to update Game Result record", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (g *GameRouter) DeleteGameResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	resultID, err := trackerHttp.GetQueryId(r, "result_id")
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusBadRequest, err, "Bad result_id query string specified", err.Error())
		return
	}

	gameID, err := g.gameResults.GetGameIDForResult(ctx, resultID)
	if err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to look up result", "Failed to delete Game Result")
		return
	}

	if !g.requirePodManager(w, r, gameID, callerPlayerID) {
		return
	}

	if err = g.games.DeleteResult(ctx, resultID); err != nil {
		trackerHttp.WriteError(g.log, w, http.StatusInternalServerError, err, "Failed to delete Game Result record", "Failed to delete Game Result")
		return
	}

	w.WriteHeader(http.StatusOK)
}
