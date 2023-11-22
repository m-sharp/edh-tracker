package routers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
)

type DeckRouter struct {
	log      *zap.Logger
	provider *models.DeckProvider
}

func NewDeckRouter(log *zap.Logger, client *lib.DBClient) *DeckRouter {
	return &DeckRouter{
		log:      log.Named("DeckRouter"),
		provider: models.NewDeckProvider(client),
	}
}

func (d *DeckRouter) GetRoutes() []*lib.Route {
	return []*lib.Route{
		{
			Path:    "/api/decks",
			Method:  http.MethodGet,
			Handler: d.GetAll,
		},
		{
			Path:    "/api/decks-by-player",
			Method:  http.MethodGet,
			Handler: d.GetAllForPlayer,
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

	decks, err := d.provider.GetAll(ctx)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(decks)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(d.log, w, marshalled)
}

func (d *DeckRouter) GetAllForPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	playerIdStr := r.URL.Query().Get("player_id")
	if playerIdStr == "" {
		lib.WriteError(
			d.log, w, http.StatusBadRequest, nil,
			"No player_id query string specified", "missing player_id",
		)
		return
	}

	playerId, err := strconv.Atoi(playerIdStr)
	if err != nil {
		lib.WriteError(
			d.log, w, http.StatusBadRequest, err,
			"Failed to convert player_id to int", "bad player_id",
		)
		return
	}

	errMsg := "Failed to get Deck records"
	decks, err := d.provider.GetAllForPlayer(ctx, playerId)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(decks)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(d.log, w, marshalled)
}

func (d *DeckRouter) GetDeckById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deckIdStr := r.URL.Query().Get("deck_id")
	if deckIdStr == "" {
		lib.WriteError(
			d.log, w, http.StatusBadRequest, nil,
			"No deck_id query string specified", "missing deck_id",
		)
		return
	}

	deckId, err := strconv.Atoi(deckIdStr)
	if err != nil {
		lib.WriteError(
			d.log, w, http.StatusBadRequest, err,
			"Failed to convert deck_id to int", "bad deck_id",
		)
		return
	}

	errMsg := "Failed to get Deck records"
	decks, err := d.provider.GetById(ctx, deckId)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(decks)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(d.log, w, marshalled)
}

func (d *DeckRouter) DeckCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create new Deck"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to read deck POST body", errMsg)
		return
	}

	var deck models.Deck
	if err := json.Unmarshal(body, &deck); err != nil {
		lib.WriteError(d.log, w, http.StatusBadRequest, err, "Failed to unmarshal Deck body", errMsg)
		return
	}
	log := d.log.With(zap.Int("PlayerId", deck.PlayerId), zap.String("Commander", deck.Commander))

	if err := deck.Validate(); err != nil {
		lib.WriteError(log, w, http.StatusBadRequest, err, "New Deck failed validation", errMsg)
		return
	}

	log.Info("Saving new Deck record")
	if err := d.provider.Add(ctx, deck.PlayerId, deck.Commander); err != nil {
		if err == models.ErrDeckExists {
			lib.WriteError(log, w, http.StatusBadRequest, err, "Attempted to create preexisting deck", err.Error())
		} else {
			lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Deck record", errMsg)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (d *DeckRouter) RetireDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deckIdStr := r.URL.Query().Get("deck_id")
	if deckIdStr == "" {
		lib.WriteError(
			d.log, w, http.StatusBadRequest, nil,
			"No deck_id query string specified", "missing deck_id",
		)
		return
	}

	deckId, err := strconv.Atoi(deckIdStr)
	if err != nil {
		lib.WriteError(
			d.log, w, http.StatusBadRequest, err,
			"Failed to convert deck_id to int", "bad deck_id",
		)
		return
	}

	errMsg := "Failed to retire deck"
	if err := d.provider.Retire(ctx, deckId); err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
