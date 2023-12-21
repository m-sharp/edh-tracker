package lib

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	Port = 8081
)

type Server struct {
	cfg    *Config
	log    *zap.Logger
	router *mux.Router
}

func NewWebServer(cfg *Config, log *zap.Logger, setupRoutes func(router *mux.Router)) *Server {
	inst := &Server{
		cfg:    cfg,
		log:    log.Named("WebServer"),
		router: mux.NewRouter(),
	}
	setupRoutes(inst.router)
	return inst
}

func (s *Server) Serve() error {
	s.log.Info("Now listening!", zap.Int("Port", Port))
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), s.router)
}
