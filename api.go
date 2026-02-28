package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/models"
	"github.com/m-sharp/edh-tracker/lib/routers"
)

type ApiRouter struct {
	cfg     *lib.Config
	log     *zap.Logger
	routers []lib.ApiRouter
}

func NewApiRouter(cfg *lib.Config, log *zap.Logger, repos *models.Repositories) *ApiRouter {
	inst := &ApiRouter{
		cfg: cfg,
		log: log.Named("ApiRoute"),
		routers: []lib.ApiRouter{
			routers.NewPlayerRouter(log, repos),
			routers.NewDeckRouter(log, repos),
			routers.NewGameRouter(log, repos),
			routers.NewPodRouter(log, repos),
		},
	}

	return inst
}

func (a *ApiRouter) SetupRoutes(router *mux.Router) {
	for _, subRouter := range a.routers {
		for _, route := range subRouter.GetRoutes() {
			// Handle CORS preflight requests
			if route.Method == http.MethodPost {
				router.HandleFunc(route.Path, lib.CORSMiddleware(lib.CORSPreflightHandler)).Methods(http.MethodOptions)
			}

			if route.MiddleWare != nil {
				router.HandleFunc(
					route.Path,
					lib.GNUMiddleware(lib.CORSMiddleware(route.MiddleWare(route.Handler))),
				).Methods(route.Method)
			} else {
				router.HandleFunc(
					route.Path,
					lib.GNUMiddleware(lib.CORSMiddleware(route.Handler)),
				).Methods(route.Method)
			}
		}
	}
}
