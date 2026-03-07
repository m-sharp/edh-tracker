package main

import (
	gohttp "net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/routers"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

type ApiRouter struct {
	cfg     *lib.Config
	log     *zap.Logger
	routers []trackerHttp.ApiRouter
}

func NewApiRouter(cfg *lib.Config, log *zap.Logger, biz *business.Business) *ApiRouter {
	inst := &ApiRouter{
		cfg: cfg,
		log: log.Named("ApiRoute"),
		routers: []trackerHttp.ApiRouter{
			routers.NewAuthRouter(log, cfg, biz),
			routers.NewPlayerRouter(log, biz),
			routers.NewDeckRouter(log, biz),
			routers.NewGameRouter(log, biz),
			routers.NewPodRouter(log, biz),
			routers.NewFormatRouter(log, biz),
			routers.NewCommanderRouter(log, biz),
		},
	}

	return inst
}

func (a *ApiRouter) SetupRoutes(router *mux.Router) {
	jwtSecret, err := a.cfg.Get(lib.JWTSecret)
	if err != nil {
		a.log.Fatal("Failed to get JWT secret", zap.Error(err))
	}

	frontendURL, err := a.cfg.Get(lib.FrontendURL)
	if err != nil {
		a.log.Fatal("Failed to get frontend URL", zap.Error(err))
	}

	devMode, _ := a.cfg.Get(lib.Development)
	secure := devMode != "true"

	corsMW := trackerHttp.CORSMiddleware(frontendURL)
	requireAuth := trackerHttp.RequireAuth(jwtSecret, secure)

	for _, subRouter := range a.routers {
		for _, route := range subRouter.GetRoutes() {
			isMutating := route.Method == gohttp.MethodPost || route.Method == gohttp.MethodPatch || route.Method == gohttp.MethodDelete

			// Register CORS preflight for state-changing methods
			if isMutating {
				router.HandleFunc(route.Path, corsMW(trackerHttp.CORSPreflightHandler)).Methods(gohttp.MethodOptions)
			}

			handler := route.Handler
			if route.MiddleWare != nil {
				handler = route.MiddleWare(handler)
			}

			// Apply RequireAuth to flagged routes and all state-changing routes unless opted out
			needsAuth := route.RequireAuth || (!route.NoAuth && isMutating)
			if needsAuth {
				handler = requireAuth(handler)
			}

			router.HandleFunc(
				route.Path,
				trackerHttp.GNUMiddleware(corsMW(handler)),
			).Methods(route.Method)
		}
	}
}
