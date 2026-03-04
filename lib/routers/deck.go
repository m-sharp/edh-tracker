package routers

import (
	"encoding/json"
	"io"
	"net/http"

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
			Handler: d.RetireDeck,
		},
	}
}

// ToDo: Eventually, this will probably need paging
func (d *DeckRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Deck records"
	playerID, _ := trackerHttp.GetQueryId(r, "player_id")

	var (
		decks []deck.EntityWithStats
		err   error
	)
	if playerID != 0 {
		decks, err = d.decks.GetAllForPlayer(ctx, playerID)
		if err != nil {
			trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	} else {
		// TODO: Should probably not exist. Also, it's slowwwww
		decks, err = d.decks.GetAll(ctx)
		if err != nil {
			trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	}

	marshalled, err := json.Marshal(decks)
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

func (d *DeckRouter) RetireDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deckID, err := trackerHttp.GetQueryId(r, "deck_id")
	if err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusBadRequest, err, "Bad deck_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to retire deck"
	if err = d.decks.Retire(ctx, deckID); err != nil {
		trackerHttp.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
