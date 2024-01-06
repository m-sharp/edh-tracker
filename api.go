package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/routers"
)

type ApiRouter struct {
	cfg     *lib.Config
	log     *zap.Logger
	routers []lib.ApiRouter
}

func NewApiRouter(cfg *lib.Config, log *zap.Logger, client *lib.DBClient) *ApiRouter {
	inst := &ApiRouter{
		cfg: cfg,
		log: log.Named("ApiRoute"),
		routers: []lib.ApiRouter{
			routers.NewPlayerRouter(log, client),
			routers.NewDeckRouter(log, client),
			routers.NewGameRouter(log, client),
		},
	}

	return inst
}

func (a *ApiRouter) SetupRoutes(router *mux.Router) {
	for _, subRouter := range a.routers {
		for _, route := range subRouter.GetRoutes() {
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
