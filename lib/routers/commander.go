package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
)

type CommanderRouter struct {
	log           *zap.Logger
	commanderRepo models.CommanderRepositoryInterface
}

func NewCommanderRouter(log *zap.Logger, repos *models.Repositories) *CommanderRouter {
	return &CommanderRouter{
		log:           log.Named("CommanderRouter"),
		commanderRepo: repos.Commanders,
	}
}

func (c *CommanderRouter) GetRoutes() []*lib.Route {
	return []*lib.Route{
		{
			Path:    "/api/commander",
			Method:  http.MethodGet,
			Handler: c.GetCommanderById,
		},
		{
			Path:    "/api/commander",
			Method:  http.MethodPost,
			Handler: c.CommanderCreate,
		},
	}
}

func (c *CommanderRouter) GetCommanderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	commanderID, err := lib.GetQueryId(r, "commander_id")
	if err != nil {
		lib.WriteError(c.log, w, http.StatusBadRequest, err, "Bad commander_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Commander record"
	commander, err := c.commanderRepo.GetById(ctx, commanderID)
	if err != nil {
		lib.WriteError(c.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}
	if commander == nil {
		lib.WriteError(c.log, w, http.StatusNotFound, nil, "Commander not found", "Commander not found")
		return
	}

	marshalled, err := json.Marshal(commander)
	if err != nil {
		lib.WriteError(c.log, w, http.StatusInternalServerError, err, "Failed to marshal records as JSON", errMsg)
		return
	}

	lib.WriteJson(c.log, w, marshalled)
}

func (c *CommanderRouter) CommanderCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create new Commander"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lib.WriteError(c.log, w, http.StatusInternalServerError, err, "Failed to read Commander POST body", errMsg)
		return
	}

	var commander models.Commander
	if err := json.Unmarshal(body, &commander); err != nil {
		lib.WriteError(c.log, w, http.StatusBadRequest, err, "Failed to unmarshal Commander body", errMsg)
		return
	}

	if commander.Name == "" {
		lib.WriteError(c.log, w, http.StatusBadRequest, nil, "Missing commander name", "Commander name is required")
		return
	}

	log := c.log.With(zap.String("Name", commander.Name))
	log.Info("Saving new Commander record")

	if _, err := c.commanderRepo.Add(ctx, commander.Name); err != nil {
		lib.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Commander record", errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
