package routers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type FormatRouter struct {
	log     *zap.Logger
	formats format.Functions
}

func NewFormatRouter(log *zap.Logger, biz *business.Business) *FormatRouter {
	return &FormatRouter{
		log:     log.Named("FormatRouter"),
		formats: biz.Formats,
	}
}

func (f *FormatRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
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

	formats, err := f.formats.GetAll(ctx)
	if err != nil {
		trackerHttp.WriteError(f.log, w, http.StatusInternalServerError, err, errMsg, errMsg)
		return
	}

	marshalled, err := json.Marshal(formats)
	if err != nil {
		trackerHttp.WriteError(f.log, w, http.StatusInternalServerError, err, "Failed to marshal records as JSON", errMsg)
		return
	}

	trackerHttp.WriteJson(f.log, w, marshalled)
}
