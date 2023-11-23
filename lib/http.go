package lib

import (
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type MiddlewareFunc func(nextHandler http.HandlerFunc) http.HandlerFunc

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
