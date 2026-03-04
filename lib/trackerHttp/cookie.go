package trackerHttp

import "net/http"

const (
	SessionCookieName  = "pod_tracker_session"
	CSRFCookieName     = "pod_tracker_csrf"
	RedirectCookieName = "pod_tracker_redirect"

	CookieMaxAge24h = 86400
	CookieMaxAge5m  = 300

	cookiePath = "/"
)

func SetCookie(
	w http.ResponseWriter,
	name, value string,
	secure bool,
	maxAge int,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     cookiePath,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		MaxAge:   maxAge,
	})
}

func ClearCookie(w http.ResponseWriter, name string, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     cookiePath,
		HttpOnly: httpOnly,
		MaxAge:   -1,
	})
}
