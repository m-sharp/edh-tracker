package trackerHttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/m-sharp/edh-tracker/lib/utils"
)

type jwtClaims struct {
	UserID   int `json:"user_id"`
	PlayerID int `json:"player_id"`
	jwt.RegisteredClaims
}

// RequireAuth validates the edh_session JWT cookie and injects userID/playerID into context.
// Re-issues a sliding 24h cookie on every valid request.
func RequireAuth(jwtSecret string, secure bool) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, err := sessionFromRequest(r, jwtSecret)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			reissueSession(w, claims.UserID, claims.PlayerID, jwtSecret, secure)

			ctx := utils.ContextWithUserInfo(r.Context(), claims.UserID, claims.PlayerID)
			next(w, r.WithContext(ctx))
		}
	}
}

// OptionalAuth reads the JWT cookie if present and injects context values, but never rejects.
func OptionalAuth(jwtSecret string, secure bool) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if claims, err := sessionFromRequest(r, jwtSecret); err == nil {
				reissueSession(w, claims.UserID, claims.PlayerID, jwtSecret, secure)
				ctx := utils.ContextWithUserInfo(r.Context(), claims.UserID, claims.PlayerID)
				r = r.WithContext(ctx)
			}

			next(w, r)
		}
	}
}

// IssueJWT creates a signed JWT with the given userID and playerID, expiring in 24 hours.
func IssueJWT(userID, playerID int, secret string) (string, error) {
	claims := jwtClaims{
		UserID:   userID,
		PlayerID: playerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signed, nil
}

func sessionFromRequest(r *http.Request, jwtSecret string) (*jwtClaims, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, fmt.Errorf("no session cookie: %w", err)
	}

	return validateJWT(cookie.Value, jwtSecret)
}

func validateJWT(tokenStr, secret string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func reissueSession(w http.ResponseWriter, userID, playerID int, secret string, secure bool) {
	tokenStr, err := IssueJWT(userID, playerID, secret)
	if err != nil {
		return // silently fail — existing session remains valid
	}

	SetCookie(w, SessionCookieName, tokenStr, secure, CookieMaxAge24h)
}
