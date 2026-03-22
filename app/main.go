package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	absStatic, err := filepath.Abs(h.staticPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Strip the leading "/" so filepath.Join treats it as relative,
	// then clean to resolve any ".." sequences.
	rel := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	path := filepath.Join(absStatic, rel)

	// Guard: reject anything that escaped the static root
	if !strings.HasPrefix(path, absStatic+string(filepath.Separator)) && path != absStatic {
		http.ServeFile(w, r, filepath.Join(absStatic, h.indexPath))
		return
	}

	// If the file doesn't exist on disk, serve index.html (SPA fallback)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(absStatic, h.indexPath))
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	cfg, err := lib.NewConfig()
	if err != nil {
		log.Fatalf("Error creating Config: %s", err.Error())
	}

	logger := lib.GetLogger(cfg)

	server := lib.NewWebServer(cfg, logger, func(router *mux.Router) {
		router.PathPrefix("/").Handler(spaHandler{staticPath: "app", indexPath: "index.html"})
	})

	if err := server.Serve(); err != nil {
		logger.Fatal("Server stopped listening", zap.Error(err))
	}
}

// ToDo: Add sitemap.xml
// ToDo: Add robots.txt
// ToDo: Add a favicon.ico
// TODO: Do we even need this? Serve some other way?
