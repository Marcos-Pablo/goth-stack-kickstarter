package session

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/config"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func New(db *sql.DB, cfg *config.Config) *scs.SessionManager {
	sm := scs.New()
	sm.Store = sqlite3store.New(db)
	sm.Lifetime = 24 * time.Hour
	sm.IdleTimeout = 30 * time.Minute
	sm.Cookie.Name = "session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = cfg.AppEnv == "production"
	return sm
}