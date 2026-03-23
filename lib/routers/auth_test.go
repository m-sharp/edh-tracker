package routers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/m-sharp/edh-tracker/lib/business/user"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

// newTestAuthRouter builds an AuthRouter with stubbed token exchange and Google user fetcher.
func newTestAuthRouter(usersBiz user.Functions, fetcher func(*http.Client) (*googleUserInfo, error)) *AuthRouter {
	return &AuthRouter{
		log:          zap.NewNop(),
		oauthCfg:     &oauth2.Config{},
		usersBiz:     usersBiz,
		jwtSecret:    "test-secret",
		secure:       false,
		frontendURL:  "http://localhost:8081",
		tokenExchanger: func(_ context.Context, _ string) (*oauth2.Token, error) {
			return &oauth2.Token{}, nil
		},
		googleUserFetcher: fetcher,
	}
}

// stubGoogleUser returns a fetcher that always returns the given googleUserInfo.
func stubGoogleUser(info *googleUserInfo) func(*http.Client) (*googleUserInfo, error) {
	return func(_ *http.Client) (*googleUserInfo, error) {
		return info, nil
	}
}

// callbackRequest builds a GET /api/auth/google/callback request with the CSRF state cookie set.
func callbackRequest(nonce string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/google/callback?state="+nonce+"&code=fake-code", nil)
	req.AddCookie(&http.Cookie{Name: trackerHttp.CSRFCookieName, Value: nonce})
	return req
}

var fixedGoogleUser = &googleUserInfo{
	Sub:     "google-sub-123",
	Email:   "alice@example.com",
	Name:    "Alice",
	Picture: "https://example.com/pic.png",
}

func TestCallback_ExistingOAuthUser(t *testing.T) {
	existingUser := &user.Entity{ID: 1, PlayerID: 10}

	var gotProvider, gotSubject string
	biz := user.Functions{
		GetByOAuth: func(_ context.Context, provider, subject string) (*user.Entity, error) {
			gotProvider = provider
			gotSubject = subject
			return existingUser, nil
		},
	}
	router := newTestAuthRouter(biz, stubGoogleUser(fixedGoogleUser))

	rr := httptest.NewRecorder()
	router.Callback(rr, callbackRequest("test-nonce"))

	// Should redirect (no error) and set a session cookie
	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, providerGoogle, gotProvider)
	assert.Equal(t, fixedGoogleUser.Sub, gotSubject, "GetByOAuth must be called with the sub from Google")
	require.NotEmpty(t, rr.Result().Cookies())
	var sessionCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == trackerHttp.SessionCookieName {
			sessionCookie = c
		}
	}
	require.NotNil(t, sessionCookie, "session cookie should be set")
}

func TestCallback_EmailFallback_LinksExistingUser(t *testing.T) {
	seededUser := &user.Entity{ID: 5, PlayerID: 20}
	linkedUser := &user.Entity{ID: 5, PlayerID: 20}

	var linkOAuthSubject string
	biz := user.Functions{
		GetByOAuth: func(_ context.Context, _, _ string) (*user.Entity, error) {
			return nil, nil
		},
		GetByEmail: func(_ context.Context, _ string) (*user.Entity, error) {
			return seededUser, nil
		},
		LinkOAuth: func(_ context.Context, _ int, _, subject, _, _, _ string) (*user.Entity, error) {
			linkOAuthSubject = subject
			return linkedUser, nil
		},
	}
	router := newTestAuthRouter(biz, stubGoogleUser(fixedGoogleUser))

	rr := httptest.NewRecorder()
	router.Callback(rr, callbackRequest("test-nonce-2"))

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, fixedGoogleUser.Sub, linkOAuthSubject, "LinkOAuth must be called with the sub from Google")
	var sessionCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == trackerHttp.SessionCookieName {
			sessionCookie = c
		}
	}
	require.NotNil(t, sessionCookie, "session cookie should be set")
}

func TestCallback_NewUser_CreateWithOAuth(t *testing.T) {
	newUser := &user.Entity{ID: 99, PlayerID: 77}

	var createWithOAuthSubject string
	biz := user.Functions{
		GetByOAuth: func(_ context.Context, _, _ string) (*user.Entity, error) {
			return nil, nil
		},
		GetByEmail: func(_ context.Context, _ string) (*user.Entity, error) {
			return nil, nil
		},
		CreateWithOAuth: func(_ context.Context, _, _, subject, _, _, _ string) (*user.Entity, error) {
			createWithOAuthSubject = subject
			return newUser, nil
		},
	}
	router := newTestAuthRouter(biz, stubGoogleUser(fixedGoogleUser))

	rr := httptest.NewRecorder()
	router.Callback(rr, callbackRequest("test-nonce-3"))

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, fixedGoogleUser.Sub, createWithOAuthSubject, "CreateWithOAuth must be called with the sub from Google")
	var sessionCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == trackerHttp.SessionCookieName {
			sessionCookie = c
		}
	}
	require.NotNil(t, sessionCookie, "session cookie should be set")
}

func TestCallback_GetByEmail_Error(t *testing.T) {
	biz := user.Functions{
		GetByOAuth: func(_ context.Context, _, _ string) (*user.Entity, error) {
			return nil, nil
		},
		GetByEmail: func(_ context.Context, _ string) (*user.Entity, error) {
			return nil, errors.New("db error")
		},
	}
	router := newTestAuthRouter(biz, stubGoogleUser(fixedGoogleUser))

	rr := httptest.NewRecorder()
	router.Callback(rr, callbackRequest("test-nonce-4"))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
