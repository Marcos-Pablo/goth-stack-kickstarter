package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sm.GetString(r.Context(), "user_id") == "" {
				http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
