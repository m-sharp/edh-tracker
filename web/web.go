package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	Port = 8081
)

type Server struct {
	cfg    *lib.Config
	log    *zap.Logger
	router *mux.Router
}

func NewWebServer(cfg *lib.Config, log *zap.Logger, api *ApiRouter) *Server {
	inst := &Server{
		cfg:    cfg,
		log:    log.Named("WebServer"),
		router: mux.NewRouter(),
	}
	api.SetupRoutes(inst.router)
	inst.setupRoutes()
	return inst
}

func (s *Server) setupRoutes() {
	// ToDo: Try SPA handler? - https://github.com/gorilla/mux#serving-single-page-applications
	s.router.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir("app/"))))
}

func (s *Server) Serve() error {
	isDev := isDevelopment(s.cfg)
	csrfSecret, err := s.cfg.Get(lib.CSRFSecret)
	if err != nil {
		return fmt.Errorf("failed to get CSRF Secret from config: %w", err)
	}

	s.log.Info("Now listening!", zap.Int("Port", Port))
	return http.ListenAndServe(
		fmt.Sprintf(":%d", Port),
		// ToDo: Posts from react will need help getting to token - https://github.com/gorilla/csrf#javascript-applications
		csrf.Protect([]byte(csrfSecret), csrf.Secure(!isDev))(s.router),
	)
}

func isDevelopment(cfg *lib.Config) bool {
	isDev, err := cfg.Get(lib.Development)
	if err != nil {
		isDev = "false"
	}
	development, err := strconv.ParseBool(isDev)
	if err != nil {
		development = false
	}

	return development
}
