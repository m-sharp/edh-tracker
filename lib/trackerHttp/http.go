package trackerHttp

import (
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/utils"
)

type MiddlewareFunc func(nextHandler http.HandlerFunc) http.HandlerFunc

func CORSMiddleware(origin string) MiddlewareFunc {
	return func(nextHandler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

			nextHandler(w, r)
		}
	}
}

func CORSPreflightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GNUMiddleware adds the X-Clacks-Overhead header to keep names alive (https://wiki.lspace.org/GNU_Terry_Pratchett)
func GNUMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Clacks-Overhead", "GNU Steve Harp, GNU Terry Pratchett")

		nextHandler(w, r)
	}
}

type Route struct {
	Path        string
	Method      string
	Handler     http.HandlerFunc
	MiddleWare  MiddlewareFunc
	RequireAuth bool // explicitly require auth (for GET routes that need it)
	NoAuth      bool // opt out of automatic auth for state-changing routes
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

// CallerPlayerID extracts the authenticated player ID from the request context.
// Returns 0, false and writes a 401 response if not present or zero.
func CallerPlayerID(w http.ResponseWriter, r *http.Request) (int, bool) {
	_, playerID, ok := utils.UserFromContext(r.Context())
	if !ok || playerID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	return playerID, true
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
