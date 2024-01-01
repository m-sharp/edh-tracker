package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/m-sharp/edh-tracker/lib"
)

func main() {
	cfg, err := lib.NewConfig()
	if err != nil {
		log.Fatalf("Error creating Config: %s", err.Error())
	}

	logger := lib.GetLogger(cfg)

	server := lib.NewWebServer(cfg, logger, func(router *mux.Router) {
		// ToDo: Try SPA handler? - https://github.com/gorilla/mux#serving-single-page-applications
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("app/")))
	})

	if err := server.Serve(); err != nil {
		logger.Fatal("Server stopped listening", zap.Error(err))
	}
}

// ToDo: Add sitemap.xml
// ToDo: Add robots.txt
// ToDo: Add a favicon.ico
