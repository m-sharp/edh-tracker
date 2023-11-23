package routers

import (
	"encoding/json"
	"io"
	"net/http"

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
	playerId, _ := lib.GetQueryId(r, "player_id")

	var (
		decks []models.Deck
		err   error
	)
	if playerId != 0 {
		decks, err = d.provider.GetAllForPlayer(ctx, playerId)
		if err != nil {
			lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	} else {
		decks, err = d.provider.GetAll(ctx)
		if err != nil {
			lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
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

	deckId, err := lib.GetQueryId(r, "deck_id")
	if err != nil {
		lib.WriteError(d.log, w, http.StatusBadRequest, err, "Bad bad_id query string specified", err.Error())
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

	deckId, err := lib.GetQueryId(r, "deck_id")
	if err != nil {
		lib.WriteError(d.log, w, http.StatusBadRequest, err, "Bad bad_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to retire deck"
	if err := d.provider.Retire(ctx, deckId); err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
