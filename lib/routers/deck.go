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
	log               *zap.Logger
	deckRepo          *models.DeckProvider
	deckCommanderRepo *models.DeckCommanderProvider
	formatRepo        *models.FormatProvider
}

func NewDeckRouter(log *zap.Logger, repos *models.Repositories) *DeckRouter {
	return &DeckRouter{
		log:               log.Named("DeckRouter"),
		deckRepo:          repos.Decks,
		deckCommanderRepo: repos.DeckCommanders,
		formatRepo:        repos.Formats,
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
	playerID, _ := lib.GetQueryId(r, "player_id")

	var (
		decks []models.DeckWithStats
		err   error
	)
	if playerID != 0 {
		decks, err = d.deckRepo.GetAllForPlayer(ctx, playerID)
		if err != nil {
			lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
			return
		}
	} else {
		decks, err = d.deckRepo.GetAll(ctx)
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

	deckID, err := lib.GetQueryId(r, "deck_id")
	if err != nil {
		lib.WriteError(d.log, w, http.StatusBadRequest, err, "Bad deck_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Deck record"
	deck, err := d.deckRepo.GetById(ctx, deckID)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(deck)
	if err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to marshall records as JSON", errMsg)
		return
	}

	lib.WriteJson(d.log, w, marshalled)
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
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, "Failed to read deck POST body", errMsg)
		return
	}

	var req newDeckRequest
	if err = json.Unmarshal(body, &req); err != nil {
		lib.WriteError(d.log, w, http.StatusBadRequest, err, "Failed to unmarshal Deck body", errMsg)
		return
	}

	log := d.log.With(zap.Int("PlayerID", req.PlayerID), zap.String("Name", req.Name), zap.Int("FormatID", req.FormatID))

	if req.PlayerID == 0 {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Missing player_id", "player_id is required")
		return
	}
	if req.Name == "" {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Missing deck name", "deck name is required")
		return
	}
	if req.FormatID == 0 {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Missing format_id", "format_id is required")
		return
	}

	format, err := d.formatRepo.GetById(ctx, req.FormatID)
	if err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to look up format", errMsg)
		return
	}
	if format == nil {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Format not found", "Invalid format_id")
		return
	}

	if format.Name == "commander" && req.CommanderID == nil {
		lib.WriteError(log, w, http.StatusBadRequest, nil, "Missing commander_id for commander-format deck", "commander_id is required for commander format")
		return
	}

	log.Info("Saving new Deck record")
	deckID, err := d.deckRepo.Add(ctx, req.PlayerID, req.Name, req.FormatID)
	if err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Deck record", errMsg)
		return
	}

	if format.Name == "commander" {
		if _, err = d.deckCommanderRepo.Add(ctx, deckID, *req.CommanderID, req.PartnerCommanderID); err != nil {
			lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add DeckCommander record", errMsg)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (d *DeckRouter) RetireDeck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deckID, err := lib.GetQueryId(r, "deck_id")
	if err != nil {
		lib.WriteError(d.log, w, http.StatusBadRequest, err, "Bad deck_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to retire deck"
	if err = d.deckRepo.Retire(ctx, deckID); err != nil {
		lib.WriteError(d.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
