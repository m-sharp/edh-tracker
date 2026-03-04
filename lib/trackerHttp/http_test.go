package trackerHttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m-sharp/edh-tracker/lib/utils"
)

func TestCallerPlayerID_WithAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(utils.ContextWithUserInfo(req.Context(), 1, 42))
	rr := httptest.NewRecorder()

	playerID, ok := CallerPlayerID(rr, req)

	assert.True(t, ok)
	assert.Equal(t, 42, playerID)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCallerPlayerID_NoAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	playerID, ok := CallerPlayerID(rr, req)

	assert.False(t, ok)
	assert.Equal(t, 0, playerID)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
