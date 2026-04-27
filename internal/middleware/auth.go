package middleware

import (
	"context"
	"net/http"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/db"
	"github.com/alexedwards/scs/v2"
)

type ctxKey int

const userKey ctxKey = 0

func RequireAuth(sm *scs.SessionManager, q *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sm.GetString(r.Context(), "user_id")
			if userID == "" {
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
				return
			}

			user, err := q.GetUserById(r.Context(), userID)

			if err != nil {
				_ = sm.Destroy(r.Context())
				http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFrom(ctx context.Context) (db.User, bool) {
	u, ok := ctx.Value(userKey).(db.User)
	return u, ok
}
