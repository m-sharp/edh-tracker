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
			Path:    "/api/player",
			Method:  http.MethodGet,
			Handler: p.GetAll,
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

	players, err := p.provider.GetAll(ctx)
	if err != nil {
		p.log.Error("Failed to get Player records", zap.Error(err))
		http.Error(w, "failed to get Player records", http.StatusInternalServerError)
		return
	}

	marshalled, err := json.Marshal(players)
	if err != nil {
		p.log.Error("Failed to marshall Player records as JSON", zap.Error(err))
		http.Error(w, "failed to get Player records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(marshalled); err != nil {
		p.log.Error("Failed to return Player records", zap.Error(err))
		http.Error(w, "failed to get Player records", http.StatusInternalServerError)
	}
}

func (p *PlayerRouter) PlayerCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.log.Error("Failed to read Player post body", zap.Error(err))
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var player models.Player
	if err := json.Unmarshal(body, &player); err != nil {
		p.log.Error("Failed to unmarshal Player body", zap.Error(err))
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}
	log := p.log.With(zap.String("NewPlayer", player.Name))

	if err := player.Validate(); err != nil {
		log.Error("New Player failed validation", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Saving new Player record")
	if err := p.provider.Add(ctx, &player); err != nil {
		log.Error("Failed to add Player record", zap.Error(err))
		http.Error(w, "failed to save Player record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
