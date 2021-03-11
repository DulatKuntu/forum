package app

import (
	"net/http"
	"strings"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//log middleware
func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s  %s  %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(rw, r)
	})
}

// authorized middleware
func (app *Application) isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(strings.Split(token, " ")) != 2 {
			app.clientError(rw, http.StatusUnauthorized, "")
			return
		}
		_, err := app.containsToken(strings.Split(token, " ")[1])
		if err != nil {
			app.clientError(rw, http.StatusUnauthorized, "")
			return
		}

		next.ServeHTTP(rw, r)

	})
}
