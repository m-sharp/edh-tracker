package routers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/deck"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type DeckRouter struct {
	log   *zap.Logger
	decks deck.Functions
}

func NewDeckRouter(log *zap.Logger, biz *business.Business) *DeckRouter {
	return &DeckRouter{
		log:   log.Named("DeckRouter"),
		decks: biz.Decks,
	}
}

func (d *DeckRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
		{
			Path:    "/api/decks",
			Method:  http.MethodGet,
			Handler: d.GetAll,
		},
		{
			Path:    "/api/deck",
			Method:  http.MethodGet,
			Handler: d.GetDeckById,
		},
		{
			Path:    "/api/deck",
			Method:  http.MethodPost,
			Handler: d.DeckCreate,
		},
		{
			Path:    "/api/deck",
			Method:  http.MethodPatch,
			Handler: d.UpdateDeck,
		},
		{
			Path:    "/api/deck",
			Method:  http.MethodDelete,
			Handler: d.DeleteDeck,
		},
	}
}

func (d *DeckRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Deck records"
	playerID, _ := trackerHttp.GetQueryId(r, "player_id")
	podID, _ := trackerHttp.GetQueryId(r, "pod_id")
	limit, _ := trackerHttp.GetQueryId(r, "limit")
	offset, _ := trackerHttp.GetQueryId(r, "offset")

	if limit > 0 {
		d.getAllPaginated(w, r, podID, playerID, limit, offset)
		return
	}

	var (
		decks []deck.EntityWithStats
		err   error
	)
	switch {
	case podID != 0:
		decks, err = d.decks.GetAllByPod(ctx, podID)
	case playerID != 0:
		decks, err = d.decks.GetAllForPlayer(ctx, playerID)
	default:
		// TODO: Should probably not exist. Also, it's slowwwww
		decks, err = d.decks.GetAll(ctx)
	}
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(decks)
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(d.log, w, marshalled)
}

func (d *DeckRouter) getAllPaginated(w http.ResponseWriter, r *http.Request, podID, playerID, limit, offset int) {
	ctx := r.Context()
	errMsg := "Failed to get Deck records"

	var (
		entities []deck.EntityWithStats
		total    int
		err      error
	)
	switch {
	case podID != 0:
		entities, total, err = d.decks.GetAllByPodPaginated(ctx, podID, limit, offset)
	case playerID != 0:
		entities, total, err = d.decks.GetAllByPlayerPaginated(ctx, playerID, limit, offset)
	default:
		// TODO: Should probably not exist. Also, it's slowwwww
		decks, decksErr := d.decks.GetAll(ctx)
		if decksErr != nil {
			trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, decksErr, errMsg, errMsg)
			return
		}
		entities = decks
		total = len(decks)
	}
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	resp := business.PaginatedResponse[deck.EntityWithStats]{
		Items:  entities,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	marshalled, err := json.Marshal(resp)
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}
	trackerHttp.WriteJson(d.log, w, marshalled)
}

func (d *DeckRouter) GetDeckById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deckID, err := trackerHttp.GetQueryId(r, "deck_id")
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusBadRequest, err, "Bad deck_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Deck record"
	deckEntity, err := d.decks.GetByID(ctx, deckID)
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(deckEntity)
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(d.log, w, marshalled)
}

type newDeckRequest struct {
	PlayerID           int    `json:"player_id"`
	Name               string `json:"name"`
	FormatID           int    `json:"format_id"`
	CommanderID        *int   `json:"commander_id"`
	PartnerCommanderID *int   `json:"partner_commander_id"`
}

func (d *DeckRouter) DeckCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create new Deck"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to read deck POST body", errMsg)
		return
	}

	var req newDeckRequest
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusBadRequest, err, "Failed to unmarshal Deck body", errMsg)
		return
	}

	log := d.log.With(zap.Int("PlayerID", req.PlayerID), zap.String("Name", req.Name), zap.Int("FormatID", req.FormatID))

	if err = deck.ValidateCreate(req.PlayerID, req.Name, req.FormatID); err != nil {
		trackerHttp.WriteError(log, w, http.StatusBadRequest, err, "Deck create request failed validation", err.Error())
		return
	}

	log.Info("Saving new Deck record")
	if _, err = d.decks.Create(ctx, req.PlayerID, req.Name, req.FormatID, req.CommanderID, req.PartnerCommanderID); err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to create Deck record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type updateDeckRequest struct {
	Name               *string `json:"name"`
	FormatID           *int    `json:"format_id"`
	CommanderID        *int    `json:"commander_id"`
	PartnerCommanderID *int    `json:"partner_commander_id"`
	Retired            *bool   `json:"retired"`
}

func (d *DeckRouter) UpdateDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to update Deck"

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	deckID, err := trackerHttp.GetQueryId(r, "deck_id")
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusBadRequest, err, "Bad deck_id query string specified", err.Error())
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to read deck PATCH body", errMsg)
		return
	}

	var req updateDeckRequest
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusBadRequest, err, "Failed to unmarshal Deck update body", errMsg)
		return
	}

	fields := deck.UpdateFields{
		Name:               req.Name,
		FormatID:           req.FormatID,
		CommanderID:        req.CommanderID,
		PartnerCommanderID: req.PartnerCommanderID,
		Retired:            req.Retired,
	}

	if err = d.decks.Update(ctx, deckID, callerPlayerID, fields); err != nil {
		if strings.HasPrefix(err.Error(), "forbidden:") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to update Deck record", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (d *DeckRouter) DeleteDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to delete Deck"

	callerPlayerID, ok := trackerHttp.CallerPlayerID(w, r)
	if !ok {
		return
	}

	deckID, err := trackerHttp.GetQueryId(r, "deck_id")
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusBadRequest, err, "Bad deck_id query string specified", err.Error())
		return
	}

	if err = d.decks.SoftDelete(ctx, deckID, callerPlayerID); err != nil {
		if strings.HasPrefix(err.Error(), "forbidden:") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to delete Deck record", errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
