package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/commander"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type CommanderRouter struct {
	log        *zap.Logger
	commanders commander.Functions
}

func NewCommanderRouter(log *zap.Logger, biz *business.Business) *CommanderRouter {
	return &CommanderRouter{
		log:        log.Named("CommanderRouter"),
		commanders: biz.Commanders,
	}
}

func (c *CommanderRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
		{
			Path:    "/api/commanders",
			Method:  http.MethodGet,
			Handler: c.GetAllCommanders,
		},
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

func (c *CommanderRouter) GetAllCommanders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Commander records"

	commanders, err := c.commanders.GetAll(ctx)
	if err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(commanders)
	if err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusInternalServerError, err, "Failed to marshal records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(c.log, w, marshalled)
}

func (c *CommanderRouter) GetCommanderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	commanderID, err := trackerHttp.GetQueryId(r, "commander_id")
	if err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusBadRequest, err, "Bad commander_id query string specified", err.Error())
		return
	}

	errMsg := "Failed to get Commander record"
	cmd, err := c.commanders.GetByID(ctx, commanderID)
	if err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}
	if cmd == nil {
		trackerHttp.WriteError(c.log, w, http.StatusNotFound, nil, "Commander not found", "Commander not found")
		return
	}

	marshalled, err := json.Marshal(cmd)
	if err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusInternalServerError, err, "Failed to marshal records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(c.log, w, marshalled)
}

func (c *CommanderRouter) CommanderCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to create new Commander"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusInternalServerError, err, "Failed to read Commander POST body", errMsg)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		trackerHttp.WriteError(c.log, w, http.StatusBadRequest, err, "Failed to unmarshal Commander body", errMsg)
		return
	}

	if req.Name == "" {
		trackerHttp.WriteError(c.log, w, http.StatusBadRequest, nil, "Missing commander name", "Commander name is required")
		return
	}

	log := c.log.With(zap.String("Name", req.Name))
	log.Info("Saving new Commander record")

	id, err := c.commanders.Create(ctx, req.Name)
	if err != nil {
		trackerHttp.WriteError(log, w, http.StatusInternalServerError, err, "Failed to add Commander record", errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{ID: id})
}
