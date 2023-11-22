package lib

import (
	"net/http"

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
