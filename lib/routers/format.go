package routers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
)

type FormatRouter struct {
	log        *zap.Logger
	formatRepo *models.FormatProvider
}

func NewFormatRouter(log *zap.Logger, repos *models.Repositories) *FormatRouter {
	return &FormatRouter{
		log:        log.Named("FormatRouter"),
		formatRepo: repos.Formats,
	}
}

func (f *FormatRouter) GetRoutes() []*lib.Route {
	return []*lib.Route{
		{
			Path:    "/api/formats",
			Method:  http.MethodGet,
			Handler: f.GetAll,
		},
	}
}

func (f *FormatRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errMsg := "Failed to get Format records"

	formats, err := f.formatRepo.GetAll(ctx)
	if err != nil {
		lib.WriteError(f.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(formats)
	if err != nil {
		lib.WriteError(f.log, w, http.StatusInternalServerError, err, "Failed to marshal records as JSON", errMsg)
		return
	}

	lib.WriteJson(f.log, w, marshalled)
}
