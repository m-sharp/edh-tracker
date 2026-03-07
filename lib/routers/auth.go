package routers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/m-sharp/edh-tracker/lib/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business"
	"github.com/m-sharp/edh-tracker/lib/business/user"
	"github.com/m-sharp/edh-tracker/lib/trackerHttp"
)

const (
	providerGoogle = "google"
)

type AuthRouter struct {
	log         *zap.Logger
	oauthCfg    *oauth2.Config
	usersBiz    user.Functions
	jwtSecret   string
	secure      bool
	frontendURL string
}

func NewAuthRouter(log *zap.Logger, cfg *lib.Config, biz *business.Business) *AuthRouter {
	clientID, err := cfg.Get(lib.GoogleClientID)
	if err != nil {
		log.Fatal("Failed to get Google client ID", zap.Error(err))
	}

	clientSecret, err := cfg.Get(lib.GoogleClientSecret)
	if err != nil {
		log.Fatal("Failed to get Google client secret", zap.Error(err))
	}

	redirectURL, err := cfg.Get(lib.OAuthRedirectURL)
	if err != nil {
		log.Fatal("Failed to get OAuth Redirect URL", zap.Error(err))
	}

	jwtSecret, err := cfg.Get(lib.JWTSecret)
	if err != nil {
		log.Fatal("Failed to get JWT secret", zap.Error(err))
	}

	frontendURL, err := cfg.Get(lib.FrontendURL)
	if err != nil {
		log.Fatal("Failed to get frontend URL", zap.Error(err))
	}

	devMode, _ := cfg.Get(lib.Development)

	return &AuthRouter{
		log: log.Named("AuthRouter"),
		oauthCfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		usersBiz:    biz.Users,
		jwtSecret:   jwtSecret,
		secure:      devMode != "true",
		frontendURL: frontendURL,
	}
}

func (a *AuthRouter) GetRoutes() []*trackerHttp.Route {
	return []*trackerHttp.Route{
		{
			Path:    "/api/auth/google",
			Method:  http.MethodGet,
			Handler: a.Login,
		},
		{
			Path:    "/api/auth/google/callback",
			Method:  http.MethodGet,
			Handler: a.Callback,
		},
		{
			Path:    "/api/auth/logout",
			Method:  http.MethodPost,
			Handler: a.Logout,
			NoAuth:  true,
		},
		{
			Path:        "/api/auth/me",
			Method:      http.MethodGet,
			Handler:     a.Me,
			RequireAuth: true,
		},
	}
}

func (a *AuthRouter) Login(w http.ResponseWriter, r *http.Request) {
	redirectPath := r.URL.Query().Get("redirect")

	nonce, err := generateNonce()
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to generate CSRF nonce", "internal error")
		return
	}

	// Store nonce and redirect path in separate short-lived cookies
	trackerHttp.SetCookie(w, trackerHttp.CSRFCookieName, nonce, a.secure, trackerHttp.CookieMaxAge5m)

	if redirectPath != "" {
		trackerHttp.SetCookie(w, trackerHttp.RedirectCookieName, redirectPath, a.secure, trackerHttp.CookieMaxAge5m)
	}

	authURL := a.oauthCfg.AuthCodeURL(nonce, oauth2.AccessTypeOnline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func generateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (a *AuthRouter) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Validate CSRF state
	csrfCookie, err := r.Cookie(trackerHttp.CSRFCookieName)
	if err != nil {
		http.Error(w, "missing CSRF cookie", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != csrfCookie.Value {
		http.Error(w, "invalid state parameter", http.StatusBadRequest)
		return
	}

	// Read optional redirect path
	redirectPath := "/"
	if rc, err := r.Cookie(trackerHttp.RedirectCookieName); err == nil {
		redirectPath = rc.Value
	}

	// Clear CSRF and redirect cookies
	trackerHttp.ClearCookie(w, trackerHttp.CSRFCookieName, false)
	trackerHttp.ClearCookie(w, trackerHttp.RedirectCookieName, false)

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := a.oauthCfg.Exchange(ctx, code)
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to exchange OAuth code", "authentication failed")
		return
	}

	// Fetch Google user profile
	googleUser, err := fetchGoogleUserInfo(a.oauthCfg.Client(ctx, token))
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to fetch Google user info", "authentication failed")
		return
	}

	// Find or create user
	u, err := a.usersBiz.GetByOAuth(ctx, providerGoogle, googleUser.Sub)
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to look up user by OAuth", "authentication failed")
		return
	}
	if u == nil {
		u, err = a.usersBiz.CreateWithOAuth(
			ctx,
			googleUser.Name,
			providerGoogle,
			googleUser.Sub,
			googleUser.Email,
			googleUser.Name,
			googleUser.Picture,
		)
		if err != nil {
			trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to create user via OAuth", "authentication failed")
			return
		}
	}

	// Issue JWT and set session cookie
	tokenStr, err := trackerHttp.IssueJWT(u.ID, u.PlayerID, a.jwtSecret)
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to issue JWT", "authentication failed")
		return
	}

	trackerHttp.SetCookie(w, trackerHttp.SessionCookieName, tokenStr, a.secure, trackerHttp.CookieMaxAge24h)

	// Redirect to frontend
	dest, err := url.JoinPath(a.frontendURL, redirectPath)
	if err != nil {
		dest = a.frontendURL + "/"
	}
	http.Redirect(w, r, dest, http.StatusFound)
}

func (a *AuthRouter) Logout(w http.ResponseWriter, r *http.Request) {
	trackerHttp.ClearCookie(w, trackerHttp.SessionCookieName, true)
	w.WriteHeader(http.StatusNoContent)
}

func (a *AuthRouter) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, ok := utils.UserFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	u, err := a.usersBiz.GetByID(ctx, userID)
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to get user", "failed to get user")
		return
	}
	if u == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	marshalled, err := json.Marshal(u)
	if err != nil {
		trackerHttp.WriteError(a.log, w, http.StatusInternalServerError, err, "Failed to marshal user", "failed to get user")
		return
	}

	trackerHttp.WriteJson(a.log, w, marshalled)
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func fetchGoogleUserInfo(client *http.Client) (*googleUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to request Google userinfo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Google userinfo response: %w", err)
	}

	var info googleUserInfo
	if err = json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse Google userinfo: %w", err)
	}

	return &info, nil
}
