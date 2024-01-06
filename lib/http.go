package lib

import (
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type MiddlewareFunc func(nextHandler http.HandlerFunc) http.HandlerFunc

// ToDo: Unnecessary?
func CORSMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		nextHandler(w, r)
	}
}

// GNUMiddleware adds the X-Clacks-Overhead header to keep names alive (https://wiki.lspace.org/GNU_Terry_Pratchett)
func GNUMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Clacks-Overhead", "GNU Steve Harp, GNU Terry Pratchett")

		nextHandler(w, r)
	}
}

type Route struct {
	Path       string
	Method     string
	Handler    http.HandlerFunc
	MiddleWare MiddlewareFunc
}

type ApiRouter interface {
	GetRoutes() []*Route
}

func WriteError(log *zap.Logger, w http.ResponseWriter, statusCode int, err error, logMsg, errMsg string) {
	log.Error(logMsg, zap.Error(err))
	http.Error(w, errMsg, statusCode)
}

func WriteJson(log *zap.Logger, w http.ResponseWriter, marshalled []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(marshalled); err != nil {
		log.Error("Failed to return records", zap.Error(err))
		http.Error(w, "failed to return records", http.StatusInternalServerError)
	}
}

func GetQueryId(r *http.Request, key string) (int, error) {
	idStr := r.URL.Query().Get(key)
	if idStr == "" {
		return 0, fmt.Errorf("missing query string value for %q", key)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert id %q into an int: %w", idStr, err)
	}

	return id, nil
}
